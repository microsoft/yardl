// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"fmt"
	"strings"

	"github.com/microsoft/yardl/tooling/internal/validation"
	"github.com/rs/zerolog/log"
)

func ValidateEvolution(latest *Environment, predecessors []*Environment, versionLabels []string) (*Environment, error) {

	// Initialize structures needed later for serialization codegen
	for _, ns := range latest.Namespaces {
		ns.Versions = append(ns.Versions, versionLabels...)
		ns.TypeDefChanges = make(map[string][]DefinitionChange)
		for _, p := range ns.Protocols {
			p.Versions = make(map[string]*ProtocolChange)
		}
	}

	// Compare each previous version with latest version
	for i, predecessor := range predecessors {
		log.Info().Msgf("Resolving changes from predecessor %s", versionLabels[i])

		definitionChanges, protocolChanges := resolveAllChanges(latest, predecessor)

		if err := validateChanges(definitionChanges, protocolChanges); err != nil {
			return nil, err
		}

		renameOldTypeDefinitions(predecessor, definitionChanges, versionLabels[i])

		for _, ch := range definitionChanges {
			if ch == nil {
				panic("NOPE")
			}
			log.Debug().Msgf("DefChange for OLD %s = NEW %s", defWithArgs(ch.PreviousDefinition()), defWithArgs(ch.LatestDefinition()))
		}

		// Save all TypeDefinition changes for codegen
		latest.GetTopLevelNamespace().TypeDefChanges[versionLabels[i]] = definitionChanges

		// Save all Protocol changes for codegen
		for _, p := range latest.GetTopLevelNamespace().Protocols {
			if protChange, ok := protocolChanges[p.GetQualifiedName()]; ok {
				p.Versions[versionLabels[i]] = protChange
			}
		}
	}

	return latest, nil
}

func renameOldTypeDefinitions(env *Environment, changes []DefinitionChange, versionLabel string) {
	oldNames := make(map[string]bool)
	for _, ch := range changes {
		td := ch.PreviousDefinition()
		oldNames[td.GetDefinitionMeta().GetQualifiedName()] = true

		log.Info().Msgf("Renaming OLD %s", td.GetDefinitionMeta().GetQualifiedName())
		td.GetDefinitionMeta().Name = td.GetDefinitionMeta().Name + "_" + versionLabel
	}

	Visit(env, func(self Visitor, node Node) {
		switch node := node.(type) {
		case TypeDefinition:
			oldName := node.GetDefinitionMeta().GetQualifiedName()
			if oldNames[oldName] {
				node.GetDefinitionMeta().Name = node.GetDefinitionMeta().Name + "_" + versionLabel
			}
			self.VisitChildren(node)
			// if nt, ok := node.(*NamedType); ok {
			// 	self.Visit(nt.Type)
			// }
		case *SimpleType:
			self.VisitChildren(node)
			self.Visit(node.ResolvedDefinition)
		}
		self.VisitChildren(node)
	})
}

// Emit User Warnings and aggregate Errors
func validateChanges(definitionChanges []DefinitionChange, protocolChanges map[string]*ProtocolChange) error {
	errorSink := &validation.ErrorSink{}
	validateTypeDefinitionChanges(definitionChanges, errorSink)
	validateProtocolChanges(protocolChanges, errorSink)
	return errorSink.AsError()
}

func validateTypeDefinitionChanges(changes []DefinitionChange, errorSink *validation.ErrorSink) {
	for _, ch := range changes {
		if ch == nil {
			panic("I don't want nil DefinitionChanges here")
		}

		td := ch.LatestDefinition()

		switch defChange := ch.(type) {

		case *CompatibilityChange:

		case *DefinitionChangeIncompatible:
			errorSink.Add(validationError(td, "changing '%s' is not backward compatible", td.GetDefinitionMeta().Name))

		case *RecordChange:
			oldRec := defChange.PreviousDefinition().(*RecordDefinition)
			newRec := defChange.LatestDefinition().(*RecordDefinition)

			for _, added := range defChange.FieldsAdded {
				if !TypeHasNullOption(added.Type) {
					log.Warn().Msgf("Adding non-Optional field '%s' may result in undefined behavior with previous versions", added.Name)
				}
			}

			for i, field := range oldRec.Fields {
				if defChange.FieldRemoved[i] {
					if !TypeHasNullOption(oldRec.Fields[i].Type) {
						log.Warn().Msgf("Removing non-Optional field '%s' may result in undefined behavior with previous versions", field.Name)
					}
					continue
				}

				if tc := defChange.FieldChanges[i]; tc != nil {
					if typeChangeIsError(tc) {
						newField := newRec.Fields[defChange.NewFieldIndex[i]]
						errorSink.Add(validationError(newField, "changing field '%s' from %s", newField.Name, typeChangeToError(tc)))
					}

					if warn := typeChangeToWarning(tc); warn != "" {
						log.Warn().Msgf("Changing field '%s' from %s", field.Name, warn)
					}
				}
			}

		case *NamedTypeChange:
			if tc := defChange.TypeChange; tc != nil {
				if typeChangeIsError(tc) {
					errorSink.Add(validationError(td, "changing type '%s' from %s", td.GetDefinitionMeta().Name, typeChangeToWarning(tc)))
				}
				if warn := typeChangeToWarning(tc); warn != "" {
					log.Warn().Msgf("Changing type '%s' from %s", td.GetDefinitionMeta().Name, warn)
				}
			}

		case *EnumChange:
			log.Warn().Msgf("EnumChange: %s", td.GetDefinitionMeta().Name)
			if tc := defChange.BaseTypeChange; tc != nil {
				errorSink.Add(validationError(td, "changing base type of '%s' is not backward compatible", td.GetDefinitionMeta().Name))
			}

		default:
			panic("Shouldn't get here")
		}
	}
}

func validateProtocolChanges(changes map[string]*ProtocolChange, errorSink *validation.ErrorSink) {
	for _, protChange := range changes {
		if protChange == nil {
			panic("I don't want nil ProtocolChanges here ")
		}

		pd := protChange.LatestDefinition().(*ProtocolDefinition)

		for _, reordered := range protChange.StepsReordered {
			errorSink.Add(validationError(reordered, "reordering step '%s' is not backward compatible", reordered.Name))
		}

		for _, removed := range protChange.StepsRemoved {
			errorSink.Add(validationError(pd, "removing step '%s' is not backward compatible", removed.Name))
		}

		for i, step := range pd.Sequence {
			if tc := protChange.StepChanges[i]; tc != nil {
				switch tc := tc.(type) {
				case *TypeChangeStepAdded:
					// A Step can be added to a Protocol if its Type can have an "empty" state
					typeCanBeEmpty := false
					switch t := GetUnderlyingType(step.Type).(type) {
					case *GeneralizedType:
						if t.Cases.HasNullOption() {
							typeCanBeEmpty = true
						} else if t.Dimensionality != nil {
							switch t.Dimensionality.(type) {
							case *Stream, *Vector, *Map:
								typeCanBeEmpty = true
							}
						}
					}
					if !typeCanBeEmpty {
						errorSink.Add(validationError(step, "adding step '%s' is not backward compatible", step.Name))
					}
				default:
					if typeChangeIsError(tc) {
						errorSink.Add(validationError(step.Type, "changing step '%s' from %s", step.Name, typeChangeToError(tc)))
					}

					if warn := typeChangeToWarning(tc); warn != "" {
						log.Warn().Msgf("Changing step '%s' from %s", step.Name, warn)
					}
				}
			}
		}
	}
}

func collectUserTypeDefsUsedInProtocols(env *Environment) []TypeDefinition {
	allTypeDefs := make(map[string]TypeDefinition)
	for _, ns := range env.Namespaces {
		for _, newTd := range ns.TypeDefinitions {
			allTypeDefs[newTd.GetDefinitionMeta().GetQualifiedName()] = newTd
		}
	}

	typeDefCollected := make(map[string]bool)
	var usedTypeDefs []TypeDefinition
	for _, ns := range env.Namespaces {
		for _, protocol := range ns.Protocols {
			Visit(protocol, func(self Visitor, node Node) {
				switch node := node.(type) {
				case TypeDefinition:
					self.VisitChildren(node)
					name := node.GetDefinitionMeta().GetQualifiedName()
					if !typeDefCollected[name] {
						if td, ok := allTypeDefs[name]; ok {
							typeDefCollected[name] = true
							usedTypeDefs = append(usedTypeDefs, td)
						}
					}
				case *SimpleType:
					self.VisitChildren(node)
					self.Visit(node.ResolvedDefinition)
				default:
					self.VisitChildren(node)
				}
			})
		}
	}
	return usedTypeDefs
}

// Returns the last name in the chain of resolved names
// i.e. the resolved Record, Enum, or NamedType name
func getResolvedName(td TypeDefinition) string {
	return getBaseDefinition(td).GetDefinitionMeta().GetQualifiedName()
}

func getParentDefinitions(td TypeDefinition, allTypeDefs []TypeDefinition) []TypeDefinition {
	parents := make([]TypeDefinition, 0)
	childName := td.GetDefinitionMeta().GetQualifiedName()

	for _, parentTd := range allTypeDefs {
		parentName := parentTd.GetDefinitionMeta().GetQualifiedName()
		if childName == parentName {
			continue
		}

		Visit(parentTd, func(self Visitor, node Node) {
			switch node := node.(type) {
			case PrimitiveDefinition:
				return

			case *NamedType:
				if node.GetDefinitionMeta().GetQualifiedName() == childName {
					parents = append(parents, parentTd)
				}
				self.Visit(node.Type)

			case TypeDefinition:
				if node.GetDefinitionMeta().GetQualifiedName() == childName {
					parents = append(parents, parentTd)
				}

			case *SimpleType:
				self.Visit(node.ResolvedDefinition)
			}
		})
	}

	return parents
}

// Only used for debugging
func defWithArgs(td TypeDefinition) string {
	// baseName := td.GetDefinitionMeta().GetQualifiedName()
	baseName := td.GetDefinitionMeta().Name

	args := make([]string, len(td.GetDefinitionMeta().TypeParameters))
	for i, param := range td.GetDefinitionMeta().TypeParameters {
		arg := "?"
		if i < len(td.GetDefinitionMeta().TypeArguments) {
			arg = TypeToShortSyntax(td.GetDefinitionMeta().TypeArguments[i], false)
		}
		args[i] = fmt.Sprintf("%s: %s", param.Name, arg)
	}
	argString := ""
	if len(args) > 0 {
		argString = fmt.Sprintf("<%s>", strings.Join(args, ", "))
	}
	return fmt.Sprintf("%s%s", baseName, argString)
}

// Returns the base definition for a NamedType
func getBaseDefinition(td TypeDefinition) TypeDefinition {
	tds := make([]TypeDefinition, 0)
	Visit(td, func(self Visitor, node Node) {
		switch node := node.(type) {
		case PrimitiveDefinition:
			return

		case *NamedType:
			tds = append(tds, node)
			self.Visit(node.Type)

		case TypeDefinition:
			tds = append(tds, node)

		case *SimpleType:
			self.Visit(node.ResolvedDefinition)
		}
	})

	return tds[len(tds)-1]
}

func getResolvedGenericDefinition(newTd, oldTd TypeDefinition) *DefinitionPair {
	// newBase := getBaseDefinition(newTd)
	// oldBase := getBaseDefinition(oldTd)
	// return &DefinitionPair{oldBase, newBase}

	// log.Debug().Msgf("Resolving %s <= %s", defWithArgs(newTd), defWithArgs(oldTd))

	resolveGenericDefinitions := func(out TypeDefinition, in TypeDefinition) TypeDefinition {
		typeArgs := make([]Type, 0)
		Visit(in, func(self Visitor, node Node) {
			switch node := node.(type) {
			case *SimpleType:
				// if len(node.TypeArguments) == len(out.GetDefinitionMeta().TypeParameters) {
				// 	typeArgs = node.TypeArguments
				// 	return
				// }
				self.Visit(node.ResolvedDefinition)

			case TypeDefinition:
				in = node
				if len(node.GetDefinitionMeta().TypeArguments) == len(out.GetDefinitionMeta().TypeParameters) {
					typeArgs = node.GetDefinitionMeta().TypeArguments
					// return
				}
				if nt, ok := node.(*NamedType); ok {
					self.Visit(nt.Type)
				}
			}
		})

		if len(typeArgs) > 0 {
			var err error
			out, err = MakeGenericType(out, typeArgs, false)
			if err != nil {
				log.Err(err).Msgf("NOPE")
				panic("HUH")
			}
		}

		return out
	}

	newDef := resolveGenericDefinitions(newTd, oldTd)
	oldDef := resolveGenericDefinitions(oldTd, newTd)

	// log.Debug().Msgf(" Resolved %s <= %s", defWithArgs(newDef), defWithArgs(oldDef))

	return &DefinitionPair{oldDef, newDef}
}

func resolveAllChanges(newEnv, oldEnv *Environment) ([]DefinitionChange, map[string]*ProtocolChange) {

	getAllTypeDefs := func(env *Environment) []TypeDefinition {
		allTypeDefs := make([]TypeDefinition, 0)
		for _, ns := range env.Namespaces {
			for _, td := range ns.TypeDefinitions {
				allTypeDefs = append(allTypeDefs, td)
			}
		}
		return allTypeDefs
	}

	getNameMappingWithinModel := func(allTypeDefs []TypeDefinition) map[string]map[string]bool {
		resolvesTo := make(map[string]map[string]bool)
		for _, tdA := range allTypeDefs {
			table := make(map[string]bool)
			for _, tdB := range allTypeDefs {
				if getResolvedName(tdA) == getResolvedName(tdB) {
					table[tdB.GetDefinitionMeta().GetQualifiedName()] = true
				}
			}
			resolvesTo[tdA.GetDefinitionMeta().GetQualifiedName()] = table
		}
		return resolvesTo
	}

	// We only need to compare new TypeDefinitions used within a Protocol
	// allNewTypeDefs := collectUserTypeDefsUsedInProtocols(newEnv)
	allNewTypeDefs := getAllTypeDefs(newEnv)
	newNameMapping := getNameMappingWithinModel(allNewTypeDefs)

	// But we need to examine all old TypeDefinitions for semantic equivalence
	allOldTypeDefs := getAllTypeDefs(oldEnv)
	oldNameMapping := getNameMappingWithinModel(allOldTypeDefs)

	semanticallyEqual := make(map[string]map[string]bool)
	for _, newTd := range allNewTypeDefs {
		table := make(map[string]bool)
		semanticallyEqual[newTd.GetDefinitionMeta().GetQualifiedName()] = table
	}

	for _, newTd := range allNewTypeDefs {
		newName := newTd.GetDefinitionMeta().GetQualifiedName()
		for _, oldTd := range allOldTypeDefs {
			oldName := oldTd.GetDefinitionMeta().GetQualifiedName()
			if newName == oldName {
				// We found a mapping between versions...
				// log.Debug().Msgf("Found a mapping between %s and %s", newName, oldName)

				for n := range newNameMapping[newName] {
					for o := range oldNameMapping[oldName] {
						semanticallyEqual[n][o] = true
						// log.Debug().Msgf("New %s === Old %s", n, o)
					}
				}
			}
		}
	}

	// Resolve all Generic TypeDefinitions across versions
	resolvedGenericDefinitions := make(map[string]map[string]*DefinitionPair)
	for _, newTd := range allNewTypeDefs {
		newName := newTd.GetDefinitionMeta().GetQualifiedName()
		resolvedGenericDefinitions[newName] = make(map[string]*DefinitionPair)

		for _, oldTd := range allOldTypeDefs {
			oldName := oldTd.GetDefinitionMeta().GetQualifiedName()
			if semanticallyEqual[newName][oldName] {
				// log.Info().Msgf("Semantically Equal %s <= %s", newName, oldName)

				// if len(newTd.GetDefinitionMeta().TypeParameters) > 0 || len(oldTd.GetDefinitionMeta().TypeParameters) > 0 {
				// One of these definitions is a Generic TypeDefinition

				resolved := getResolvedGenericDefinition(newTd, oldTd)
				resolvedGenericDefinitions[newName][oldName] = resolved
				// log.Info().Msgf("Resolved %s <= %s: %s <= %s", defWithArgs(newTd), defWithArgs(oldTd), defWithArgs(resolved.LatestDefinition()), defWithArgs(resolved.PreviousDefinition()))
				// }
			}

		}
	}

	// Filter all semantically equivalent TypeDefinition pairs to only include "base" TypeDefinitions
	// These the TypeDefinitions that will be directly compared
	bases := make(map[string]map[string]bool)
	for _, newTd := range allNewTypeDefs {
		newName := newTd.GetDefinitionMeta().GetQualifiedName()
		bases[newName] = make(map[string]bool)
		for _, oldTd := range allOldTypeDefs {
			oldName := oldTd.GetDefinitionMeta().GetQualifiedName()
			if semanticallyEqual[newName][oldName] {
				newBase := getBaseDefinition(newTd)
				oldBase := getBaseDefinition(oldTd)
				bases[newBase.GetDefinitionMeta().GetQualifiedName()][oldBase.GetDefinitionMeta().GetQualifiedName()] = true

				if newName == oldName {
					bases[newName][oldName] = true
				}
			}
		}
	}

	// For Debugging:
	// for _, newTd := range allNewTypeDefs {
	// 	newName := newTd.GetDefinitionMeta().GetQualifiedName()
	// 	for _, oldTd := range allOldTypeDefs {
	// 		oldName := oldTd.GetDefinitionMeta().GetQualifiedName()

	// 		if bases[newName][oldName] {
	// 			log.Info().Msgf("Base %s <= %s", defWithArgs(newTd), defWithArgs(oldTd))

	// 			if resolved, ok := resolvedGenericDefinitions[newName][oldName]; ok {
	// 				log.Info().Msgf("  Resolved %s <= %s", defWithArgs(resolved.LatestDefinition()), defWithArgs(resolved.PreviousDefinition()))
	// 			}
	// 		}
	// 	}
	// }

	compared := make(map[string]map[string]bool)
	changes := make(map[string]map[string]DefinitionChange)

	context := &EvolutionContext{
		SemanticMapping:  semanticallyEqual,
		Compared:         compared,
		Changes:          changes,
		ResolvedGenerics: resolvedGenericDefinitions,
	}

	for _, newTd := range allNewTypeDefs {
		newName := newTd.GetDefinitionMeta().GetQualifiedName()
		compared[newName] = make(map[string]bool)
		changes[newName] = make(map[string]DefinitionChange)

		for _, oldTd := range allOldTypeDefs {
			oldName := oldTd.GetDefinitionMeta().GetQualifiedName()

			// if semanticallyEqual[newName][oldName] {
			if bases[newName][oldName] {

				// log.Debug().Msgf("New %s semantically equals Old %s", newName, oldName)

				// TODO: Delete:
				// if len(newTd.GetDefinitionMeta().TypeParameters) > 0 || len(oldTd.GetDefinitionMeta().TypeParameters) > 0 {
				// 	// Skip Generic TypeDefinitions for now. We'll compare them when we compare ProtocolStep Types
				// 	continue
				// }

				resolved := context.ResolvedGenerics[newName][oldName]
				newDef := resolved.LatestDefinition()
				oldDef := resolved.PreviousDefinition()

				if ch := compareTypeDefinitions(newDef, oldDef, context); ch != nil {
					// if ch := compareTypeDefinitions(newTd, oldTd, context); ch != nil {
					if _, ok := ch.(*DefinitionChangeIncompatible); ok {
						log.Error().Msgf("New %s is incompatible with Old %s", newName, oldName)
					} else {
						log.Debug().Msgf("CHANGE for %s <= %s", newName, oldName)
					}
					changes[newName][oldName] = ch
				}
				compared[newName][oldName] = true
			}
		}
	}

	compatibility := make(map[string]map[string]DefinitionChange)
	for _, newTd := range allNewTypeDefs {
		newName := newTd.GetDefinitionMeta().GetQualifiedName()
		for _, oldTd := range allOldTypeDefs {
			oldName := oldTd.GetDefinitionMeta().GetQualifiedName()

			if ch, ok := changes[newName][oldName]; ok {

				switch ch.(type) {
				case *DefinitionChangeIncompatible:
				default:
					// Found a compatible DefinitionChange between two "base" types
					// Now we want to ensure all NamedTypes that REFERENCE this changed Definition are also marked as "changed"

					// log.Debug().Msgf("Finding compatibility for %s <= %s", newName, oldName)
					for _, newParent := range getParentDefinitions(newTd, allNewTypeDefs) {
						newParentName := newParent.GetDefinitionMeta().GetQualifiedName()

						compatibility[newParentName] = make(map[string]DefinitionChange)

						for _, oldParent := range getParentDefinitions(oldTd, allOldTypeDefs) {
							oldParentName := oldParent.GetDefinitionMeta().GetQualifiedName()

							if _, ok := changes[newParentName][oldParentName]; ok {
								// Already have this Change stored, so no Compatibility needed
								continue
							}
							if semanticallyEqual[newParentName][oldParentName] {
								// log.Debug().Msgf("Need Compatibility for %s <= %s", newParentName, oldParentName)
								// compatibility[newParentName][oldParentName] = &CompatibilityChange{DefinitionPair{oldParent, newParent}}
								// resolved := getResolvedGenericDefinition(newParent, oldParent, resolvedGenericDefinitions)
								resolved := context.ResolvedGenerics[newParentName][oldParentName]
								compatibility[newParentName][oldParentName] = &CompatibilityChange{*resolved}
							} else {
								panic("UH OH")
							}
						}
					}

				}
			}
		}
	}

	allProtocolChanges := resolveAllProtocolChanges(newEnv, oldEnv, context)

	allDefinitionChanges := make([]DefinitionChange, 0)
	for _, newTd := range allNewTypeDefs {
		newName := newTd.GetDefinitionMeta().GetQualifiedName()

		// _, newChanged := newDefsChanged[newName]
		// _, newReferenced := newDefsReferenced[newName]
		// needNewCompatibility := !newChanged && newReferenced

		for _, oldTd := range allOldTypeDefs {
			oldName := oldTd.GetDefinitionMeta().GetQualifiedName()

			// _, oldChanged := oldDefsChanged[oldName]
			// _, oldReferenced := oldDefsReferenced[oldName]
			// needOldCompatibility := !oldChanged && oldReferenced

			if semanticallyEqual[newName][oldName] {

				var compat DefinitionChange
				var change DefinitionChange
				var alias DefinitionChange

				if ch, ok := changes[newName][oldName]; ok {
					if compat != nil {
						panic("HUH")
					}

					kind := "Gud"
					if _, bad := ch.(*DefinitionChangeIncompatible); bad {
						kind = "Bad"
					}
					change = ch
					log.Debug().Msgf("%s DefChange: %s <= %s", kind, defWithArgs(ch.LatestDefinition()), defWithArgs(ch.PreviousDefinition()))
					allDefinitionChanges = append(allDefinitionChanges, ch)
				} else if compat, ok := compatibility[newName][oldName]; ok {
					log.Debug().Msgf("Compatibility: %s <= %s", defWithArgs(compat.LatestDefinition()), defWithArgs(compat.PreviousDefinition()))
					allDefinitionChanges = append(allDefinitionChanges, compat)
				} else if _, ok := newNameMapping[oldName]; !ok {
					if change != nil {
						panic("HUH")
					}

					// This old TypeDefinition was removed in the new model
					resolvedPair := context.ResolvedGenerics[newName][oldName]
					// newDef := resolvedPair.LatestDefinition()

					// alias = &CompatibilityChange{DefinitionPair{oldDef, newDef}}
					alias = &CompatibilityChange{*resolvedPair}
					log.Debug().Msgf("Alias Removed: %s = %s", defWithArgs(alias.PreviousDefinition()), defWithArgs(alias.LatestDefinition()))
					allDefinitionChanges = append(allDefinitionChanges, alias)
				}

				// allDefinitionChanges = append(allDefinitionChanges, compat, change, alias)
			}
		}
	}

	log.Debug().Msgf("Finished comparing TypeDefinitions")

	// Finally, filter out changes with duplicate "old" TypeDefinitions because we only need the first one for codegen
	neededDefinitionChanges := make([]DefinitionChange, 0)
	uniqueOldNames := make(map[string]bool)
	for _, change := range allDefinitionChanges {
		if change == nil {
			continue
		}
		oldName := change.PreviousDefinition().GetDefinitionMeta().GetQualifiedName()

		if uniqueOldNames[oldName] {
			continue
		}

		neededDefinitionChanges = append(neededDefinitionChanges, change)
		uniqueOldNames[oldName] = true

		// newName := change.LatestDefinition().GetDefinitionMeta().GetQualifiedName()
		// log.Debug().Msgf("Change for %s <= %s", newName, oldName)
	}

	return neededDefinitionChanges, allProtocolChanges
}

func resolveAllProtocolChanges(newEnv, oldEnv *Environment, context *EvolutionContext) map[string]*ProtocolChange {
	// Protocols may be reordered, added, or removed
	// We only care about pre-existing Protocols that CHANGED
	oldProts := make(map[string]*ProtocolDefinition)
	for _, ns := range oldEnv.Namespaces {
		for _, oldProt := range ns.Protocols {
			oldProts[oldProt.GetQualifiedName()] = oldProt
		}
	}

	allProtocolChanges := make(map[string]*ProtocolChange)
	for _, ns := range newEnv.Namespaces {
		for _, newProt := range ns.Protocols {
			oldProt, ok := oldProts[newProt.GetDefinitionMeta().GetQualifiedName()]
			if !ok {
				// Skip new ProtocolDefinition
				continue
			}

			// Annotate this ProtocolDefinition with any changes from previous version.
			protocolChange := compareProtocolDefinitions(newProt, oldProt, context)
			if protocolChange != nil {
				// Annotate the ProtocolChange with the Old ProtocolDefinition schema string
				protocolChange.PreviousSchema = GetProtocolSchemaString(oldProt, oldEnv.SymbolTable)

				// log.Debug().Msgf("Protocol %s changed", newProt.GetQualifiedName())
				allProtocolChanges[newProt.GetQualifiedName()] = protocolChange
			}
		}
	}

	return allProtocolChanges
}

type EvolutionContext struct {
	SemanticMapping  map[string]map[string]bool
	Compared         map[string]map[string]bool
	Changes          map[string]map[string]DefinitionChange
	AliasesRemoved   []DefinitionChange
	ResolvedGenerics map[string]map[string]*DefinitionPair
}

func compareTypeDefinitions(newTd, oldTd TypeDefinition, context *EvolutionContext) DefinitionChange {
	// log.Debug().Msgf("Comparing TypeDefinitions %s <= %s", newTd.GetDefinitionMeta().GetQualifiedName(), oldTd.GetDefinitionMeta().GetQualifiedName())

	newName := newTd.GetDefinitionMeta().GetQualifiedName()
	oldName := oldTd.GetDefinitionMeta().GetQualifiedName()

	// If these definitions aren't semantically equal, then they aren't compatible between versions
	if !context.SemanticMapping[newName][oldName] {
		log.Warn().Msgf("TypeDefinitions %s and %s aren't semantically equal", newName, oldName)
		return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
	}

	if context.Compared[newName][oldName] {
		// If we already compared these definitions return any pre-existing DefinitionChange
		if ch, ok := context.Changes[newName][oldName]; ok {
			return ch
		}
		return nil
	}
	log.Debug().Msgf("Comparing Defs %s <= %s", defWithArgs(newTd), defWithArgs(oldTd))

	switch newTd := newTd.(type) {
	case *NamedType:
		switch oldTd := oldTd.(type) {
		case *NamedType:
			// log.Warn().Msgf("  Comparing NamedTypes %s <= %s", TypeToShortSyntax(newTd.Type, true), TypeToShortSyntax(oldTd.Type, true))
			// log.Warn().Msgf("  Which is also: Types %s <= %s", TypeToShortSyntax(GetUnderlyingType(newTd.Type), true), TypeToShortSyntax(GetUnderlyingType(oldTd.Type), true))
			// resolved, ok := context.ResolvedGenerics[newName][oldName]
			// if ok {
			// 	log.Warn().Msgf("  Could try using rslv %s <= %s", defWithArgs(resolved.LatestDefinition()), defWithArgs(resolved.PreviousDefinition()))
			// }

			ch := compareTypes(newTd.Type, oldTd.Type, context)
			switch ch.(type) {
			case nil:
				return nil
			case *TypeChangeIncompatible:
				return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
			}
			return &NamedTypeChange{DefinitionPair{oldTd, newTd}, ch}

		case *RecordDefinition:
			newType, ok := newTd.Type.(*SimpleType)
			if !ok {
				// Old TypeDefinition is a RecordDefinition, but new TypeDefinition is not a SimpleType
				return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
			}

			ch := compareTypeDefinitions(newType.ResolvedDefinition, oldTd, context)
			if ch == nil {
				return nil
			}
			switch ch := ch.(type) {
			case *DefinitionChangeIncompatible:
				return ch
			default:
				// Alias Added
				return &NamedTypeChange{DefinitionPair{oldTd, newTd}, nil}
			}

		case *EnumDefinition:
			newType, ok := newTd.Type.(*SimpleType)
			if !ok {
				// Old TypeDefinition is an EnumDefinition, but new TypeDefinition is not a SimpleType
				return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
			}

			ch := compareTypeDefinitions(newType.ResolvedDefinition, oldTd, context)
			if ch == nil {
				return nil
			}
			switch ch := ch.(type) {
			case *DefinitionChangeIncompatible:
				return ch
			default:
				// Alias Added
				return &NamedTypeChange{DefinitionPair{oldTd, newTd}, nil}
			}

		default:
			return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
		}

	case *RecordDefinition:
		switch oldTd := oldTd.(type) {
		case *RecordDefinition:

			// Enforce that TypeParameters must be IDENTICAL (for now)
			if len(newTd.GetDefinitionMeta().TypeParameters) != len(oldTd.GetDefinitionMeta().TypeParameters) {
				return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
			}
			for i, newTypeParam := range newTd.GetDefinitionMeta().TypeParameters {
				oldTypeParam := oldTd.GetDefinitionMeta().TypeParameters[i]
				if newTypeParam.Name != oldTypeParam.Name {
					return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
				}
			}

			ch := compareRecordDefinitions(newTd, oldTd, context)
			if ch == nil {
				return nil
			}
			return ch

		case *NamedType:
			oldType, ok := oldTd.Type.(*SimpleType)
			if !ok {
				// New TypeDefinition is a RecordDefinition, but old TypeDefinition is not a SimpleType
				return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
			}

			// return compareTypeDefinitions(newTd, oldType.ResolvedDefinition, context)
			ch := compareTypeDefinitions(newTd, oldType.ResolvedDefinition, context)
			if ch == nil {
				return nil
			}
			switch ch := ch.(type) {
			case *DefinitionChangeIncompatible:
				return ch
			default:
				// Alias Removed
				// context.AliasesRemoved = append(context.AliasesRemoved, &NamedTypeChange{DefinitionPair{oldTd, newTd}, nil})
				return &NamedTypeChange{DefinitionPair{oldTd, newTd}, nil}
			}
		default:
			return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
		}

	case *EnumDefinition:
		switch oldTd := oldTd.(type) {
		case *EnumDefinition:
			ch := compareEnumDefinitions(newTd, oldTd, context)
			if ch == nil {
				return nil
			}
			return ch

		case *NamedType:
			oldType, ok := oldTd.Type.(*SimpleType)
			if !ok {
				// New TypeDefinition is an EnumDefinition, but old TypeDefinition is not a SimpleType
				return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
			}

			ch := compareTypeDefinitions(newTd, oldType.ResolvedDefinition, context)
			// if ch == nil {
			// 	return nil
			// }
			switch ch := ch.(type) {
			case nil:
				return nil
			case *DefinitionChangeIncompatible:
				return ch
			default:
				return &NamedTypeChange{DefinitionPair{oldTd, newTd}, nil}
			}

		default:
			return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
		}

	default:
		panic("Expected a TypeDefinition...")
	}

	panic("Shouldn't get here")
	return nil

}

// func OLD_compareTypeDefinitions(newTd, oldTd TypeDefinition, context *EvolutionContext) DefinitionChange {
// 	// log.Debug().Msgf("Comparing TypeDefinitions %s <= %s", newTd.GetDefinitionMeta().GetQualifiedName(), oldTd.GetDefinitionMeta().GetQualifiedName())

// 	newName := newTd.GetDefinitionMeta().GetQualifiedName()
// 	oldName := oldTd.GetDefinitionMeta().GetQualifiedName()

// 	// If these definitions aren't semantically equal, then they aren't compatible between versions
// 	if !context.SemanticMapping[newName][oldName] {
// 		log.Warn().Msgf("TypeDefinitions %s and %s aren't semantically equal", newName, oldName)
// 		return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 	}

// 	if context.Compared[newName][oldName] {
// 		// If we already compared these definitions return any pre-existing DefinitionChange
// 		if ch, ok := context.Changes[newName][oldName]; ok {
// 			return ch
// 		}
// 		return nil
// 	}
// 	log.Debug().Msgf("Comparing Defs %s <= %s", defWithArgs(newTd), defWithArgs(oldTd))

// 	// if len(newTd.GetDefinitionMeta().TypeArguments) > 0 {
// 	// 	generic, err := MakeGenericType(oldTd, newTd.GetDefinitionMeta().TypeArguments, false)
// 	// 	if err != nil {
// 	// 		panic("Can't do that")
// 	// 		return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 	// 	}
// 	// 	oldTd = generic
// 	// 	context.OldResolvedGenericDefs[newName][oldName] = generic
// 	// 	log.Debug().Msgf("Saved old resolved generic def %s", defWithArgs(generic))
// 	// }
// 	context.OldResolvedGenericDefs[newName][oldName] = oldTd

// 	// if len(oldTd.GetDefinitionMeta().TypeArguments) > 0 {
// 	// 	generic, err := MakeGenericType(newTd, oldTd.GetDefinitionMeta().TypeArguments, false)
// 	// 	if err != nil {
// 	// 		panic("Can't do that")
// 	// 		return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 	// 	}
// 	// 	newTd = generic
// 	// 	context.NewResolvedGenericDefs[newName][oldName] = generic
// 	// 	log.Debug().Msgf("Saved new resolved generic def %s", defWithArgs(generic))
// 	// }
// 	context.NewResolvedGenericDefs[newName][oldName] = newTd

// 	switch newTd := newTd.(type) {
// 	case *NamedType:
// 		switch oldTd := oldTd.(type) {
// 		case *NamedType:
// 			// newDef := newTd
// 			// oldDef := oldTd

// 			// log.Debug().Msgf("  Comparing NamedTypes %s <= %s", TypeToShortSyntax(newTd.Type, true), TypeToShortSyntax(oldTd.Type, true))

// 			// if oldSt, ok := oldTd.Type.(*SimpleType); ok {
// 			// 	if len(oldSt.TypeArguments) > 0 {

// 			// 		def := oldSt.ResolvedDefinition
// 			// 		log.Debug().Msgf("  Args: %d, Params: %d", len(def.GetDefinitionMeta().TypeArguments), len(newTd.GetDefinitionMeta().TypeParameters))
// 			// 		for len(def.GetDefinitionMeta().TypeArguments) < len(newTd.GetDefinitionMeta().TypeParameters) {
// 			// 			if nt, ok := def.(*NamedType); ok {
// 			// 				if st, ok := nt.Type.(*SimpleType); ok {
// 			// 					def = st.ResolvedDefinition
// 			// 				} else {
// 			// 					panic("HUH")
// 			// 				}
// 			// 			} else {
// 			// 				break
// 			// 				// panic("NOPE")
// 			// 			}
// 			// 		}

// 			// 		// genericNewTd, err := MakeGenericType(newTd, st.TypeArguments, false)
// 			// 		genericNewTd, err := MakeGenericType(newTd, def.GetDefinitionMeta().TypeArguments, false)
// 			// 		if err != nil {
// 			// 			log.Err(err).Msgf(" Failed to make generic type %s using %s", defWithArgs(newTd), TypeToShortSyntax(oldSt, true))
// 			// 			// These TypeDefinitions aren't compatible
// 			// 			// return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 			// 			panic("ARGH")
// 			// 		} else {
// 			// 			newDef = genericNewTd.(*NamedType)
// 			// 			context.NewResolvedGenericDefs[newName][oldName] = genericNewTd
// 			// 		}
// 			// 		log.Debug().Msgf("  oldTd was %s, now it is %s", defWithArgs(oldTd), defWithArgs(def))
// 			// 	}
// 			// }
// 			// log.Debug().Msgf("  newTd was %s, now it is %s", defWithArgs(newTd), defWithArgs(newDef))

// 			// if newSt, ok := newTd.Type.(*SimpleType); ok {
// 			// 	if len(newSt.TypeArguments) > 0 {

// 			// 		def := newSt.ResolvedDefinition
// 			// 		log.Debug().Msgf("  Args: %d, Params: %d", len(def.GetDefinitionMeta().TypeArguments), len(oldTd.GetDefinitionMeta().TypeParameters))
// 			// 		for len(def.GetDefinitionMeta().TypeArguments) < len(oldTd.GetDefinitionMeta().TypeParameters) {
// 			// 			log.Warn().Msgf("  Trying to resolve %s <= %s", defWithArgs(def), defWithArgs(oldTd))
// 			// 			if nt, ok := def.(*NamedType); ok {
// 			// 				if st, ok := nt.Type.(*SimpleType); ok {
// 			// 					def = st.ResolvedDefinition
// 			// 				} else {
// 			// 					panic("HUH")
// 			// 				}
// 			// 			} else {
// 			// 				log.Error().Msgf("Expected NamedType but %s is a %T", def.GetDefinitionMeta().GetQualifiedName(), def)

// 			// 				break
// 			// 				// panic("NOPE")
// 			// 			}
// 			// 		}

// 			// 		// genericOldTd, err := MakeGenericType(oldTd, st.TypeArguments, false)
// 			// 		genericOldTd, err := MakeGenericType(oldTd, def.GetDefinitionMeta().TypeArguments, false)
// 			// 		if err != nil {
// 			// 			log.Err(err).Msgf(" Failed to make generic type %s using %s", defWithArgs(oldTd), TypeToShortSyntax(newSt, true))
// 			// 			// These TypeDefinitions aren't compatible
// 			// 			// return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 			// 			panic("ARGHY")
// 			// 		} else {
// 			// 			oldDef = genericOldTd.(*NamedType)
// 			// 			context.OldResolvedGenericDefs[newName][oldName] = genericOldTd
// 			// 		}
// 			// 	}
// 			// }

// 			// newTd = newDef
// 			// oldTd = oldDef

// 			// newDef := newTd
// 			// oldDef := oldTd

// 			if len(newTd.TypeParameters) != len(oldTd.TypeParameters) {
// 				if len(newTd.TypeParameters) < len(oldTd.TypeParameters) {
// 					log.Debug().Msgf("Need to unpack new %s", defWithArgs(newTd))

// 					// st, ok := newDef.Type.(*SimpleType)
// 					// if !ok {
// 					// 	panic("HUH")
// 					// }

// 					// for len(st.TypeArguments) < len(oldDef.TypeParameters) {
// 					// 	if nt, ok := st.ResolvedDefinition.(*NamedType); ok {

// 					// 		newDef = nt

// 					// 		if it, ok := nt.Type.(*SimpleType); ok {
// 					// 			st = it
// 					// 		} else {
// 					// 			panic("HUH")
// 					// 		}
// 					// 	} else {
// 					// 		panic("HUH")
// 					// 	}
// 					// }
// 					var typeArgs []Type

// 					var def TypeDefinition = newTd
// 					for len(typeArgs) < len(oldTd.TypeParameters) {
// 						if nt, ok := def.(*NamedType); ok {
// 							if st, ok := nt.Type.(*SimpleType); ok {
// 								typeArgs = st.TypeArguments
// 								def = st.ResolvedDefinition
// 							} else {
// 								panic("HUH")
// 							}
// 						} else {
// 							panic("HUH")
// 						}
// 					}

// 					log.Debug().Msgf("Unpacked it to %s", defWithArgs(def))

// 					genericOldTd, err := MakeGenericType(oldTd, def.GetDefinitionMeta().TypeArguments, false)
// 					if err != nil {
// 						// log.Err(err).Msgf(" Failed to make generic type %s using %s", defWithArgs(newTd), TypeToShortSyntax(oldSt, true))
// 						// These TypeDefinitions aren't compatible
// 						// return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 						panic("ARGH")
// 					}

// 					log.Debug().Msgf("    Resolved oldTd to %s", defWithArgs(genericOldTd))
// 					context.OldResolvedGenericDefs[newName][oldName] = genericOldTd

// 					oldTd = genericOldTd.(*NamedType)

// 				} else {
// 					log.Debug().Msgf("Need to unpack old %s", defWithArgs(oldTd))

// 					var typeArgs []Type

// 					var def TypeDefinition = oldTd
// 					for len(typeArgs) < len(newTd.TypeParameters) {
// 						if nt, ok := def.(*NamedType); ok {
// 							if st, ok := nt.Type.(*SimpleType); ok {
// 								typeArgs = st.TypeArguments
// 								def = st.ResolvedDefinition
// 							} else {
// 								panic("HUH")
// 							}
// 						} else {
// 							panic("HUH")
// 						}
// 					}

// 					log.Debug().Msgf("Unpacked it to %s", defWithArgs(def))

// 					genericNewTd, err := MakeGenericType(newTd, def.GetDefinitionMeta().TypeArguments, false)
// 					if err != nil {
// 						// log.Err(err).Msgf(" Failed to make generic type %s using %s", defWithArgs(newTd), TypeToShortSyntax(oldSt, true))
// 						// These TypeDefinitions aren't compatible
// 						// return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 						panic("ARGH")
// 					}

// 					log.Debug().Msgf("    Resolved newTd to %s", defWithArgs(genericNewTd))
// 					context.NewResolvedGenericDefs[newName][oldName] = genericNewTd

// 					newTd = genericNewTd.(*NamedType)
// 					// ch := compareTypes(genericNewTd.(*NamedType).Type, oldTd.Type, context)
// 					// log.Debug().Msgf("Compared the types to get %T", ch)
// 				}
// 			} else {
// 				// if st, ok := oldTd.Type.(*SimpleType); ok {
// 				// 	log.Debug().Msgf("    %d <- %d", len(newTd.TypeParameters), len(st.TypeArguments))
// 				// 	newDef, err := MakeGenericType(newTd, st.TypeArguments, false)
// 				// 	if err != nil {
// 				// 		panic("Can't do that")
// 				// 	}
// 				// 	newTd = newDef.(*NamedType)
// 				// }

// 				// if st, ok := newTd.Type.(*SimpleType); ok {
// 				// 	log.Debug().Msgf("    %d <- %d", len(oldTd.TypeParameters), len(st.TypeArguments))
// 				// 	oldDef, err := MakeGenericType(oldTd, st.TypeArguments, false)
// 				// 	if err != nil {
// 				// 		panic("Can't do that")
// 				// 	}
// 				// 	oldTd = oldDef.(*NamedType)
// 				// }
// 			}

// 			log.Debug().Msgf("At this point: %s <= %s", defWithArgs(newTd), defWithArgs(oldTd))

// 			// log.Debug().Msgf("  Comparing NamedTypes %s <= %s", TypeToShortSyntax(newTd.Type, true), TypeToShortSyntax(oldTd.Type, true))

// 			// if oldSt, ok := oldTd.Type.(*SimpleType); ok {
// 			// 	oldName := oldSt.ResolvedDefinition.GetDefinitionMeta().GetQualifiedName()
// 			// 	log.Debug().Msgf("    New Resolved Def: %s", defWithArgs(context.NewResolvedGenericDefs[newTd.GetQualifiedName()][oldName]))
// 			// }

// 			// if newSt, ok := newTd.Type.(*SimpleType); ok {
// 			// 	newName := newSt.ResolvedDefinition.GetDefinitionMeta().GetQualifiedName()
// 			// 	log.Debug().Msgf("    Old Resolved Def: %s", defWithArgs(context.OldResolvedGenericDefs[newName][oldTd.GetQualifiedName()]))
// 			// }

// 			// if oldSt, ok := oldTd.Type.(*SimpleType); ok {
// 			// 	if len(oldSt.TypeArguments) > 0 {

// 			// 		def := oldSt.ResolvedDefinition
// 			// 		log.Debug().Msgf("  Args: %d, Params: %d", len(def.GetDefinitionMeta().TypeArguments), len(newTd.GetDefinitionMeta().TypeParameters))
// 			// 		for len(def.GetDefinitionMeta().TypeArguments) < len(newTd.GetDefinitionMeta().TypeParameters) {
// 			// 			if nt, ok := def.(*NamedType); ok {
// 			// 				if st, ok := nt.Type.(*SimpleType); ok {
// 			// 					def = st.ResolvedDefinition
// 			// 				} else {
// 			// 					panic("HUH")
// 			// 				}
// 			// 			} else {
// 			// 				break
// 			// 				// panic("NOPE")
// 			// 			}
// 			// 		}

// 			// 		// genericNewTd, err := MakeGenericType(newTd, st.TypeArguments, false)
// 			// 		genericNewTd, err := MakeGenericType(newTd, def.GetDefinitionMeta().TypeArguments, false)
// 			// 		if err != nil {
// 			// 			log.Err(err).Msgf(" Failed to make generic type %s using %s", defWithArgs(newTd), TypeToShortSyntax(oldSt, true))
// 			// 			// These TypeDefinitions aren't compatible
// 			// 			// return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 			// 			panic("ARGH")
// 			// 		} else {
// 			// 			newDef = genericNewTd.(*NamedType)
// 			// 			context.NewResolvedGenericDefs[newName][oldName] = genericNewTd
// 			// 		}
// 			// 		log.Debug().Msgf("  oldTd was %s, now it is %s", defWithArgs(oldTd), defWithArgs(def))
// 			// 	}
// 			// }
// 			// log.Debug().Msgf("  newTd was %s, now it is %s", defWithArgs(newTd), defWithArgs(newDef))

// 			// if newSt, ok := newTd.Type.(*SimpleType); ok {
// 			// 	if len(newSt.TypeArguments) > 0 {

// 			// 		def := newSt.ResolvedDefinition
// 			// 		log.Debug().Msgf("  Args: %d, Params: %d", len(def.GetDefinitionMeta().TypeArguments), len(oldTd.GetDefinitionMeta().TypeParameters))
// 			// 		for len(def.GetDefinitionMeta().TypeArguments) < len(oldTd.GetDefinitionMeta().TypeParameters) {
// 			// 			log.Warn().Msgf("  Trying to resolve %s <= %s", defWithArgs(def), defWithArgs(oldTd))
// 			// 			if nt, ok := def.(*NamedType); ok {
// 			// 				if st, ok := nt.Type.(*SimpleType); ok {
// 			// 					def = st.ResolvedDefinition
// 			// 				} else {
// 			// 					panic("HUH")
// 			// 				}
// 			// 			} else {
// 			// 				log.Error().Msgf("Expected NamedType but %s is a %T", def.GetDefinitionMeta().GetQualifiedName(), def)

// 			// 				break
// 			// 				// panic("NOPE")
// 			// 			}
// 			// 		}

// 			// 		// genericOldTd, err := MakeGenericType(oldTd, st.TypeArguments, false)
// 			// 		genericOldTd, err := MakeGenericType(oldTd, def.GetDefinitionMeta().TypeArguments, false)
// 			// 		if err != nil {
// 			// 			log.Err(err).Msgf(" Failed to make generic type %s using %s", defWithArgs(oldTd), TypeToShortSyntax(newSt, true))
// 			// 			// These TypeDefinitions aren't compatible
// 			// 			// return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 			// 			panic("ARGHY")
// 			// 		} else {
// 			// 			oldDef = genericOldTd.(*NamedType)
// 			// 			context.OldResolvedGenericDefs[newName][oldName] = genericOldTd
// 			// 		}
// 			// 	}
// 			// }

// 			// newTd = newDef
// 			// oldTd = oldDef

// 			// log.Debug().Msgf("  Base     Definitions %s <= %s", defWithArgs(newTd), defWithArgs(oldTd))
// 			// log.Debug().Msgf("  Resolved Definitions %s <= %s", defWithArgs(newDef), defWithArgs(oldDef))

// 			log.Warn().Msgf("  Comparing NamedTypes %s <= %s", TypeToShortSyntax(newTd.Type, true), TypeToShortSyntax(oldTd.Type, true))
// 			ch := compareTypes(newTd.Type, oldTd.Type, context)
// 			// log.Warn().Msgf("  Comparing NamedTypes %s <= %s", TypeToShortSyntax(newDef.Type, true), TypeToShortSyntax(oldDef.Type, true))
// 			// ch := compareTypes(newDef.Type, oldDef.Type, context)
// 			if ch == nil {
// 				return nil
// 			}
// 			switch ch := ch.(type) {
// 			case *TypeChangeIncompatible:
// 				return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 			// case *TypeChangeDefinitionChanged:
// 			// 	return &NamedTypeChange{DefinitionPair{oldTd, newTd}, ch}
// 			default:
// 				// // When comparing two NamedTypes, either of which is Generic, the underlying
// 				// // Types cannot "change", unless it is just a definition change.
// 				// if IsGeneric(newTd) || IsGeneric(oldTd) {
// 				// 	log.Warn().Msgf("   %T", ch)
// 				// 	log.Warn().Msgf("Generic NamedTypes %s and %s cannot change types", newName, oldName)
// 				// 	return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 				// }
// 				return &NamedTypeChange{DefinitionPair{oldTd, newTd}, ch}
// 			}

// 		case *RecordDefinition:
// 			newType, ok := newTd.Type.(*SimpleType)
// 			if !ok {
// 				// Old TypeDefinition is a RecordDefinition, but new TypeDefinition is not a SimpleType
// 				return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 			}

// 			var oldDef TypeDefinition = oldTd
// 			if len(newType.TypeArguments) > 0 {

// 				// if len(newType.TypeArguments) != len(oldTd.GetDefinitionMeta().TypeParameters) {
// 				// 	newName := newType.ResolvedDefinition.GetDefinitionMeta().GetQualifiedName()
// 				// 	newTdResolved := context.NewResolvedGenericDefs[newName][oldName]

// 				// 	newDef, err := MakeGenericType(newTdResolved, newType.TypeArguments, false)
// 				// 	if err != nil {
// 				// 		panic(err)
// 				// 	}
// 				// 	log.Debug().Msgf("    Should try using %s <= %s", defWithArgs(newTdResolved), defWithArgs(oldTd))
// 				// 	log.Debug().Msgf("    which is now     %s", defWithArgs(newDef))
// 				// 	log.Debug().Msgf("    new resolved has %d type args", len(newDef.GetDefinitionMeta().TypeArguments))
// 				// }

// 				newDef := newType.ResolvedDefinition
// 				log.Debug().Msgf("  New Definition %s", defWithArgs(newDef))
// 				if len(newType.TypeArguments) != len(newDef.GetDefinitionMeta().TypeArguments) {
// 					panic("HUH")
// 				}

// 				// if len(newType.TypeArguments) != len(oldTd.GetDefinitionMeta().TypeParameters) {
// 				for len(newDef.GetDefinitionMeta().TypeArguments) != len(oldTd.GetDefinitionMeta().TypeParameters) {
// 					log.Debug().Msgf("  Cannot resolve %s using %s", defWithArgs(oldTd), defWithArgs(newDef))
// 					if nt, ok := newDef.(*NamedType); ok {
// 						if st, ok := nt.Type.(*SimpleType); ok {
// 							// var err error
// 							// newDef, err = MakeGenericType(st.ResolvedDefinition, st.TypeArguments, false)
// 							log.Debug().Msgf("    Resolving %s with %s", defWithArgs(st.ResolvedDefinition), defWithArgs(newDef))
// 							// newDef, err = MakeGenericType(st.ResolvedDefinition, newDef.GetDefinitionMeta().TypeArguments, false)
// 							// if err != nil {
// 							// 	panic(err)
// 							// }
// 							// newType = st
// 							newDef = st.ResolvedDefinition
// 						} else {
// 							panic("HUH")
// 						}
// 					} else {
// 						panic("NOPE")
// 					}
// 				}
// 				// log.Debug().Msgf("  Should resolve %s using %s", defWithArgs(oldTd), TypeToShortSyntax(newType, true))
// 				log.Debug().Msgf("  Should resolve %s using %s", defWithArgs(oldTd), defWithArgs(newDef))

// 				// generic, err := MakeGenericType(oldDef, newType.TypeArguments, false)
// 				generic, err := MakeGenericType(oldDef, newDef.GetDefinitionMeta().TypeArguments, false)
// 				if err != nil {
// 					log.Err(err).Msgf(" Failed to make generic type %s using %s", defWithArgs(oldTd), TypeToShortSyntax(newType, true))
// 					return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 				}
// 				oldDef = generic
// 				context.OldResolvedGenericDefs[newName][oldName] = generic
// 			}

// 			ch := compareTypeDefinitions(newType.ResolvedDefinition, oldDef, context)
// 			if ch == nil {
// 				return nil
// 			}
// 			switch ch := ch.(type) {
// 			case *DefinitionChangeIncompatible:
// 				return ch
// 			default:
// 				// Alias Added
// 				return &NamedTypeChange{DefinitionPair{oldDef, newTd}, nil}
// 			}

// 			// oldType := &SimpleType{NodeMeta: *oldTd.GetNodeMeta(), Name: oldTd.GetDefinitionMeta().GetQualifiedName(), ResolvedDefinition: oldTd}
// 			// ch := compareTypes(newTd.Type, oldType, context)
// 			// if ch == nil {
// 			// 	return nil
// 			// }
// 			// switch ch := ch.(type) {
// 			// case *TypeChangeIncompatible:
// 			// 	return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 			// default:
// 			// 	return &NamedTypeChange{DefinitionPair{oldTd, newTd}, ch}
// 			// }

// 			// if newSt, ok := newTd.Type.(*SimpleType); ok {
// 			// 	if len(newSt.TypeArguments) > 0 {
// 			// 		genericOldTd, err := MakeGenericType(oldTd, newSt.TypeArguments, false)
// 			// 		if err != nil {
// 			// 			panic("These TypeDefinitions aren't compatible")
// 			// 		}
// 			// 		ch := compareTypeDefinitions(newSt.ResolvedDefinition, genericOldTd, context)
// 			// 		if ch == nil {
// 			// 			return nil
// 			// 		}
// 			// 		switch ch := ch.(type) {
// 			// 		case *DefinitionChangeIncompatible:
// 			// 			return ch
// 			// 		default:
// 			// 			log.Error().Msgf("This is where I generate the NamedTypeChange %s <= %s", newTd.GetDefinitionMeta().GetQualifiedName(), genericOldTd.GetDefinitionMeta().GetQualifiedName())
// 			// 			return &NamedTypeChange{DefinitionPair{genericOldTd, newTd}, nil}
// 			// 		}
// 			// 	}

// 			// 	ch := compareTypeDefinitions(newTd, newSt.ResolvedDefinition, context)
// 			// 	if ch == nil {
// 			// 		return nil
// 			// 	}
// 			// 	switch ch := ch.(type) {
// 			// 	case *DefinitionChangeIncompatible:
// 			// 		return ch
// 			// 	default:
// 			// 		log.Error().Msgf("This is where I generate the NamedTypeChange %s <= %s", newTd.GetDefinitionMeta().GetQualifiedName(), oldTd.GetDefinitionMeta().GetQualifiedName())
// 			// 		return &NamedTypeChange{DefinitionPair{oldTd, newTd}, nil}
// 			// 	}
// 			// } else {
// 			// 	panic("HUH")
// 			// }

// 			// oldType := &SimpleType{NodeMeta: *oldTd.GetNodeMeta(), Name: oldTd.GetDefinitionMeta().GetQualifiedName(), ResolvedDefinition: oldTd}
// 			// ch := compareTypes(newTd.Type, oldType, context)
// 			// if ch == nil {
// 			// 	return nil
// 			// }
// 			// switch ch := ch.(type) {
// 			// case *TypeChangeIncompatible:
// 			// 	return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 			// default:
// 			// 	return &NamedTypeChange{DefinitionPair{oldTd, newTd}, ch}
// 			// }

// 		case *EnumDefinition:
// 			newType, ok := newTd.Type.(*SimpleType)
// 			if !ok {
// 				// Old TypeDefinition is an EnumDefinition, but new TypeDefinition is not a SimpleType
// 				return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 			}

// 			ch := compareTypeDefinitions(newType.ResolvedDefinition, oldTd, context)
// 			if ch == nil {
// 				return nil
// 			}
// 			switch ch := ch.(type) {
// 			case *DefinitionChangeIncompatible:
// 				return ch
// 			default:
// 				// Alias Added
// 				return &NamedTypeChange{DefinitionPair{oldTd, newTd}, nil}
// 			}

// 		default:
// 			return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 		}

// 	case *RecordDefinition:
// 		switch oldTd := oldTd.(type) {
// 		case *RecordDefinition:

// 			// Enforce that TypeParameters must be IDENTICAL (for now)
// 			if len(newTd.GetDefinitionMeta().TypeParameters) != len(oldTd.GetDefinitionMeta().TypeParameters) {
// 				return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 			}
// 			for i, newTypeParam := range newTd.GetDefinitionMeta().TypeParameters {
// 				oldTypeParam := oldTd.GetDefinitionMeta().TypeParameters[i]
// 				if newTypeParam.Name != oldTypeParam.Name {
// 					return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 				}
// 			}

// 			ch := compareRecordDefinitions(newTd, oldTd, context)
// 			if ch == nil {
// 				return nil
// 			}
// 			return ch

// 		case *NamedType:
// 			oldType, ok := oldTd.Type.(*SimpleType)
// 			if !ok {
// 				// New TypeDefinition is a RecordDefinition, but old TypeDefinition is not a SimpleType
// 				return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 			}

// 			var newDef TypeDefinition = newTd
// 			if len(oldType.TypeArguments) > 0 {

// 				oldDef := oldType.ResolvedDefinition
// 				for len(oldDef.GetDefinitionMeta().TypeArguments) != len(newDef.GetDefinitionMeta().TypeParameters) {
// 					if nt, ok := oldDef.(*NamedType); ok {
// 						if st, ok := nt.Type.(*SimpleType); ok {
// 							oldDef = st.ResolvedDefinition
// 						} else {
// 							panic("HUH")
// 						}
// 					} else {
// 						panic("NOPE")
// 					}
// 				}

// 				// generic, err := MakeGenericType(newDef, oldType.TypeArguments, false)
// 				generic, err := MakeGenericType(newDef, oldDef.GetDefinitionMeta().TypeArguments, false)
// 				if err != nil {
// 					log.Err(err).Msgf(" Failed to make generic type %s using %s", defWithArgs(oldTd), TypeToShortSyntax(oldType, true))
// 					return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 				}
// 				newDef = generic
// 				context.NewResolvedGenericDefs[newName][oldName] = generic
// 			}

// 			ch := compareTypeDefinitions(newDef, oldType.ResolvedDefinition, context)
// 			if ch == nil {
// 				return nil
// 			}
// 			switch ch := ch.(type) {
// 			case *DefinitionChangeIncompatible:
// 				return ch
// 			default:
// 				// Alias Removed
// 				// context.AliasesRemoved = append(context.AliasesRemoved, &NamedTypeChange{DefinitionPair{oldTd, newTd}, nil})
// 				return &NamedTypeChange{DefinitionPair{oldTd, newDef}, nil}
// 			}
// 			// if oldSt, ok := oldTd.Type.(*SimpleType); ok {
// 			// 	if len(oldSt.TypeArguments) > 0 {
// 			// 		genericNewTd, err := MakeGenericType(newTd, oldSt.TypeArguments, false)
// 			// 		if err != nil {
// 			// 			panic("These TypeDefinitions aren't compatible")
// 			// 		}
// 			// 		ch := compareTypeDefinitions(genericNewTd, oldSt.ResolvedDefinition, context)
// 			// 		if ch == nil {
// 			// 			return nil
// 			// 		}
// 			// 		switch ch := ch.(type) {
// 			// 		case *DefinitionChangeIncompatible:
// 			// 			return ch
// 			// 		default:
// 			// 			log.Error().Msgf("This is where I generate the NamedTypeChange %s <= %s", genericNewTd.GetDefinitionMeta().GetQualifiedName(), oldTd.GetDefinitionMeta().GetQualifiedName())
// 			// 			return &NamedTypeChange{DefinitionPair{oldTd, genericNewTd}, nil}
// 			// 		}
// 			// 	}

// 			// 	ch := compareTypeDefinitions(newTd, oldSt.ResolvedDefinition, context)
// 			// 	if ch == nil {
// 			// 		return nil
// 			// 	}
// 			// 	switch ch := ch.(type) {
// 			// 	case *DefinitionChangeIncompatible:
// 			// 		return ch
// 			// 	default:
// 			// 		log.Error().Msgf("This is where I generate the NamedTypeChange %s <= %s", newTd.GetDefinitionMeta().GetQualifiedName(), oldTd.GetDefinitionMeta().GetQualifiedName())
// 			// 		return &NamedTypeChange{DefinitionPair{oldTd, newTd}, nil}
// 			// 	}
// 			// } else {
// 			// 	panic("Can't change a ")
// 			// }

// 			// newType := &SimpleType{NodeMeta: *newTd.GetNodeMeta(), Name: newTd.GetDefinitionMeta().GetQualifiedName(), ResolvedDefinition: newTd}
// 			// ch := compareTypes(newType, oldTd.Type, context)
// 			// if ch == nil {
// 			// 	return nil
// 			// }
// 			// switch ch := ch.(type) {
// 			// case *TypeChangeIncompatible:
// 			// 	return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 			// default:
// 			// 	log.Error().Msgf("This is where I generate the NamedTypeChange %s <= %s", newTd.GetDefinitionMeta().GetQualifiedName(), oldTd.GetDefinitionMeta().GetQualifiedName())
// 			// 	return &NamedTypeChange{DefinitionPair{oldTd, newTd}, ch}
// 			// }

// 		default:
// 			return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 		}

// 	case *EnumDefinition:
// 		switch oldTd := oldTd.(type) {
// 		case *EnumDefinition:
// 			ch := compareEnumDefinitions(newTd, oldTd, context)
// 			if ch == nil {
// 				return nil
// 			}
// 			return ch

// 		case *NamedType:
// 			oldType, ok := oldTd.Type.(*SimpleType)
// 			if !ok {
// 				// New TypeDefinition is an EnumDefinition, but old TypeDefinition is not a SimpleType
// 				return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 			}

// 			ch := compareTypeDefinitions(newTd, oldType.ResolvedDefinition, context)
// 			if ch == nil {
// 				return nil
// 			}
// 			switch ch := ch.(type) {
// 			case *DefinitionChangeIncompatible:
// 				return ch
// 			default:
// 				// Alias Removed
// 				// context.AliasesRemoved = append(context.AliasesRemoved, &NamedTypeChange{DefinitionPair{oldTd, newTd}, nil})
// 				return &NamedTypeChange{DefinitionPair{oldTd, newTd}, nil}
// 			}

// 		default:
// 			return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 		}

// 	default:
// 		panic("Expected a TypeDefinition...")
// 	}

// 	panic("Shouldn't get here")
// 	return nil
// }

// Compares two ProtocolDefinitions with matching names
func compareProtocolDefinitions(newProtocol, oldProtocol *ProtocolDefinition, context *EvolutionContext) *ProtocolChange {
	change := &ProtocolChange{
		DefinitionPair: DefinitionPair{oldProtocol, newProtocol},
		StepChanges:    make([]TypeChange, len(newProtocol.Sequence)),
	}

	newSteps := make(map[string]*ProtocolStep)
	for _, newStep := range newProtocol.Sequence {
		newSteps[newStep.Name] = newStep
	}

	oldSteps := make(map[string]*ProtocolStep)
	oldStepIndices := make(map[string]int)
	for i, oldStep := range oldProtocol.Sequence {
		oldSteps[oldStep.Name] = oldStep
		oldStepIndices[oldStep.Name] = i

		if _, ok := newSteps[oldStep.Name]; !ok {
			// CHANGE: Removed this ProtocolStep
			change.StepsRemoved = append(change.StepsRemoved, oldStep)
		}
	}

	expectedIndex := 0
	for i, newStep := range newProtocol.Sequence {
		oldStep, ok := oldSteps[newStep.Name]
		if !ok {
			// CHANGE: Added this ProtocolStep
			change.StepChanges[i] = &TypeChangeStepAdded{TypePair{nil, newStep.Type}}
			continue
		}

		if oldStepIndices[newStep.Name] != expectedIndex {
			// CHANGE: Reordered this ProtocolStep
			change.StepsReordered = append(change.StepsReordered, newStep)
		}
		expectedIndex++

		log.Debug().Msgf("Comparing ProtocolStep %s <= %s", newStep.Name, oldStep.Name)
		if typeChange := compareTypes(newStep.Type, oldStep.Type, context); typeChange != nil {
			// CHANGE: ProtocolStep type changed
			change.StepChanges[i] = typeChange
		}
	}

	if len(change.StepsReordered) > 0 || len(change.StepsRemoved) > 0 {
		return change
	}
	for _, ch := range change.StepChanges {
		if ch != nil {
			return change
		}
	}

	return nil
}

// Compares two RecordDefinitions with matching names
func compareRecordDefinitions(newRecord, oldRecord *RecordDefinition, context *EvolutionContext) *RecordChange {
	change := &RecordChange{
		DefinitionPair: DefinitionPair{oldRecord, newRecord},
		FieldRemoved:   make([]bool, len(oldRecord.Fields)),
		FieldChanges:   make([]TypeChange, len(oldRecord.Fields)),
		NewFieldIndex:  make([]int, len(oldRecord.Fields)),
	}

	// Fields may be reordered
	oldFields := make(map[string]*Field)
	for _, f := range oldRecord.Fields {
		oldFields[f.Name] = f
	}

	newFields := make(map[string]*Field)
	newFieldIndices := make(map[string]int)
	for i, newField := range newRecord.Fields {
		newFields[newField.Name] = newField
		newFieldIndices[newField.Name] = i

		if _, ok := oldFields[newField.Name]; !ok {
			// CHANGE: New field
			change.FieldsAdded = append(change.FieldsAdded, newField)
		}
	}

	fieldsReordered := false
	for i, oldField := range oldRecord.Fields {
		newField, ok := newFields[oldField.Name]
		if !ok {
			// CHANGE: Removed field
			change.FieldRemoved[i] = true
			change.NewFieldIndex[i] = -1
			continue
		}

		change.NewFieldIndex[i] = newFieldIndices[oldField.Name]
		if change.NewFieldIndex[i] != i {
			fieldsReordered = true
		}

		if typeChange := compareTypes(newField.Type, oldField.Type, context); typeChange != nil {
			// CHANGE: Field type changed
			change.FieldChanges[i] = typeChange
		}
	}

	if fieldsReordered || len(change.FieldsAdded) > 0 {
		return change
	}
	for i := range oldRecord.Fields {
		if change.FieldRemoved[i] || change.FieldChanges[i] != nil {
			return change
		}
	}
	return nil
}

func compareEnumDefinitions(newNode, oldEnum *EnumDefinition, context *EvolutionContext) DefinitionChange {
	if newNode.IsFlags != oldEnum.IsFlags {
		// CHANGE: Changed Enum to Flags or vice versa
		return &DefinitionChangeIncompatible{DefinitionPair{oldEnum, newNode}}
	}

	oldBaseType := oldEnum.BaseType
	if oldBaseType == nil {
		oldBaseType = &SimpleType{ResolvedDefinition: PrimitiveInt32}
	}

	newBaseType := newNode.BaseType
	if newBaseType == nil {
		newBaseType = &SimpleType{ResolvedDefinition: PrimitiveInt32}
	}

	if ch := compareTypes(newBaseType, oldBaseType, context); ch != nil {
		// CHANGE: Changed Enum base type
		return &EnumChange{DefinitionPair{oldEnum, newNode}, ch}
	}

	return nil
}

func compareTypes(newType, oldType Type, context *EvolutionContext) TypeChange {

	newType = GetUnderlyingType(newType)
	oldType = GetUnderlyingType(oldType)

	switch newType := newType.(type) {

	case *SimpleType:
		switch oldType := oldType.(type) {
		case *SimpleType:
			return compareSimpleTypes(newType, oldType, context)
		case *GeneralizedType:
			newDef := newType.ResolvedDefinition
			if nt, ok := newDef.(*NamedType); ok {
				log.Debug().Msgf("  Unwinding new NT %s", defWithArgs(nt))
				return compareTypes(nt.Type, oldType, context)
			}
			return compareGeneralizedToSimpleTypes(newType, oldType, context)
		default:
			panic("Shouldn't get here")
		}

	case *GeneralizedType:
		switch oldType := oldType.(type) {
		case *GeneralizedType:
			return compareGeneralizedTypes(newType, oldType, context)
		case *SimpleType:
			oldDef := oldType.ResolvedDefinition
			if nt, ok := oldDef.(*NamedType); ok {
				log.Debug().Msgf("  Unwinding old NT %s", defWithArgs(nt))
				return compareTypes(newType, nt.Type, context)
			}
			return compareSimpleToGeneralizedTypes(newType, oldType, context)
		default:
			panic("Shouldn't get here")
		}

	default:
		panic("Expected a type")
	}
}

func compareSimpleTypes(newType, oldType *SimpleType, context *EvolutionContext) TypeChange {
	newDef := newType.ResolvedDefinition
	oldDef := oldType.ResolvedDefinition

	newName := newDef.GetDefinitionMeta().GetQualifiedName()
	oldName := oldDef.GetDefinitionMeta().GetQualifiedName()
	// log.Debug().Msgf("Comparing SimpleTypes %s <= %s", newName, oldName)

	// Here is where we do a lookup on the two ResolvedDefinitions
	// 1. Are they semantically equal?
	// 2. Did we already compare them?
	// 3. Are TypeArguments compatible?
	if context.SemanticMapping[newName][oldName] {
		log.Debug().Msgf("Checking: %s <= %s", defWithArgs(newDef), defWithArgs(oldDef))

		if !context.Compared[newName][oldName] {
			panic("Shouldn't get here")
		}

		ch, ok := context.Changes[newName][oldName]
		if ok {
			// log.Debug().Msgf("Found change for %s <= %s", defWithArgs(newDef), defWithArgs(oldDef))
			switch ch.(type) {
			case *DefinitionChangeIncompatible:
				return &TypeChangeIncompatible{TypePair{oldType, newType}}
			}
		} else {
			ch = nil
		}

		typeArgDefinitionChanged := false
		if len(newType.TypeArguments) > 0 || len(oldType.TypeArguments) > 0 {

			if len(newType.TypeArguments) != len(oldType.TypeArguments) {
				panic("HUH")
			}

			// resolvedGenericDefinitions := context.ResolvedGenerics[newName][oldName]
			// newDefReference := resolvedGenericDefinitions.LatestDefinition()
			// oldDefReference := resolvedGenericDefinitions.PreviousDefinition()
			// log.Debug().Msgf("Against: %s <= %s", defWithArgs(newDefReference), defWithArgs(oldDefReference))

			for i := range newType.TypeArguments {
				// newTypeArg := newType.TypeArguments[i]
				// oldTypeArg := oldType.TypeArguments[i]
				newTypeArg := newDef.GetDefinitionMeta().TypeArguments[i]
				oldTypeArg := oldDef.GetDefinitionMeta().TypeArguments[i]

				if ch := compareTypes(newTypeArg, oldTypeArg, context); ch != nil {
					switch ch.(type) {
					case *TypeChangeIncompatible:
						log.Error().Msgf("TypeArgs aren't compatible %s <= %s", TypeToShortSyntax(newTypeArg, true), TypeToShortSyntax(oldTypeArg, true))
						return &TypeChangeIncompatible{TypePair{oldType, newType}}
					case *TypeChangeDefinitionChanged:
						typeArgDefinitionChanged = true
					}
				}
			}
			// Plug them in, then compare...

			////////

			// log.Debug().Msgf("Against: %s <= %s", defWithArgs(ch.LatestDefinition()), defWithArgs(ch.PreviousDefinition()))
			// log.Debug().Msgf("AKA    : %s <= %s", defWithArgs(context.NewResolvedGenericDefs[newName][oldName]),
			// 	defWithArgs(context.OldResolvedGenericDefs[newName][oldName]))

			////
			// newResolved, err := MakeGenericType(newDefReference, newType.TypeArguments, false)
			// if err != nil {
			// 	log.Err(err).Msgf("These TypeDefinitions aren't compatible")
			// }

			// oldResolved, err := MakeGenericType(oldDefReference, oldType.TypeArguments, false)
			// if err != nil {
			// 	log.Err(err).Msgf("These TypeDefinitions aren't compatible")
			// }

			// log.Debug().Msgf("Or... : %s <= %s", defWithArgs(newResolved), defWithArgs(oldResolved))

			// newDefReference = newResolved
			// oldDefReference = oldResolved

			////

			// newTypeArgsByName := make(map[string]Type)
			// for i, tp := range newDef.GetDefinitionMeta().TypeParameters {
			// 	ta := newDef.GetDefinitionMeta().TypeArguments[i]
			// 	newTypeArgsByName[tp.Name] = ta
			// }
			// for i, ta := range newDefReference.GetDefinitionMeta().TypeArguments {
			// 	tp := newDefReference.GetDefinitionMeta().TypeParameters[i]
			// 	if st, ok := ta.(*SimpleType); ok {
			// 		if innerTp, ok := st.ResolvedDefinition.(*GenericTypeParameter); ok {
			// 			newTypeArgsByName[innerTp.Name] = newTypeArgsByName[tp.Name]
			// 		}
			// 	}
			// }
			// // for n, t := range newTypeArgsByName {
			// // 	log.Debug().Msgf("New Type TypeArg %s: %s", n, TypeToShortSyntax(t, true))
			// // }

			// oldTypeArgsByName := make(map[string]Type)
			// for i, tp := range oldDef.GetDefinitionMeta().TypeParameters {
			// 	ta := oldDef.GetDefinitionMeta().TypeArguments[i]
			// 	oldTypeArgsByName[tp.Name] = ta
			// }
			// for i, ta := range oldDefReference.GetDefinitionMeta().TypeArguments {
			// 	tp := oldDefReference.GetDefinitionMeta().TypeParameters[i]
			// 	if st, ok := ta.(*SimpleType); ok {
			// 		if innerTp, ok := st.ResolvedDefinition.(*GenericTypeParameter); ok {
			// 			oldTypeArgsByName[innerTp.Name] = oldTypeArgsByName[tp.Name]
			// 		}
			// 	}
			// }
			// // for n, t := range oldTypeArgsByName {
			// // 	log.Debug().Msgf("Old Type TypeArg %s: %s", n, TypeToShortSyntax(t, true))
			// // }

			// newResolvedTypeArgs := make([]Type, len(newType.TypeArguments))
			// for i, tp := range newDefReference.GetDefinitionMeta().TypeParameters {
			// 	var ta Type
			// 	if i >= len(newDefReference.GetDefinitionMeta().TypeArguments) {
			// 		// log.Debug().Msgf("Looking in OLD TYPE for TypeParam %s", tp.Name)
			// 		ta = oldTypeArgsByName[tp.Name]
			// 	} else {
			// 		ta = newDefReference.GetDefinitionMeta().TypeArguments[i]
			// 		if st, ok := ta.(*SimpleType); ok {
			// 			if innerTp, ok := st.ResolvedDefinition.(*GenericTypeParameter); ok {
			// 				// log.Debug().Msgf("Looking in OLD TYPE for TypeParam %s", innerTp.Name)
			// 				ta = oldTypeArgsByName[innerTp.Name]
			// 			}
			// 		}
			// 	}

			// 	if ta == nil {
			// 		// panic(fmt.Sprintf("Can't find old Type TypeArg for new TypeParam %s", tp.Name))
			// 		log.Error().Msgf("Can't find old Type TypeArg for new TypeParam %s", tp.Name)
			// 	} else {
			// 		newResolvedTypeArgs[i] = ta
			// 	}

			// }
			// // for i, t := range newResolvedTypeArgs {
			// // 	log.Debug().Msgf("New Resolved TypeArg %s: %s", ch.LatestDefinition().GetDefinitionMeta().TypeParameters[i].Name, TypeToShortSyntax(t, true))
			// // }

			// oldResolvedTypeArgs := make([]Type, len(oldType.TypeArguments))
			// for i, tp := range oldDefReference.GetDefinitionMeta().TypeParameters {
			// 	var ta Type
			// 	if i >= len(oldDefReference.GetDefinitionMeta().TypeArguments) {
			// 		// log.Debug().Msgf("Looking in NEW TYPE for TypeParam %s", tp.Name)
			// 		ta = newTypeArgsByName[tp.Name]
			// 	} else {
			// 		ta = oldDefReference.GetDefinitionMeta().TypeArguments[i]
			// 		if st, ok := ta.(*SimpleType); ok {
			// 			if innerTp, ok := st.ResolvedDefinition.(*GenericTypeParameter); ok {
			// 				// log.Debug().Msgf("Looking in NEW TYPE for TypeParam %s", innerTp.Name)
			// 				ta = newTypeArgsByName[innerTp.Name]
			// 			}
			// 		}
			// 	}

			// 	if ta == nil {
			// 		// panic(fmt.Sprintf("Can't find new Type TypeArg for old TypeParam %s", tp.Name))
			// 		log.Error().Msgf("Can't find new Type TypeArg for old TypeParam %s", tp.Name)
			// 	} else {
			// 		oldResolvedTypeArgs[i] = ta
			// 	}

			// }
			// // for i, t := range oldResolvedTypeArgs {
			// // 	log.Debug().Msgf("Old Resolved TypeArg %s: %s", ch.PreviousDefinition().GetDefinitionMeta().TypeParameters[i].Name, TypeToShortSyntax(t, true))
			// // }

			// // newDefArgsByName := make(map[string]Type)
			// // for i, tp := range ch.LatestDefinition().GetDefinitionMeta().TypeParameters {
			// // 	var ta Type
			// // 	if i < len(ch.LatestDefinition().GetDefinitionMeta().TypeArguments) {
			// // 		// TypeArg already exists in the Definition
			// // 		ta = ch.LatestDefinition().GetDefinitionMeta().TypeArguments[i]

			// // 		if st, ok := ta.(*SimpleType); ok {
			// // 			if innerTp, ok := st.ResolvedDefinition.(*GenericTypeParameter); ok {
			// // 				ta = newTypeArgsByName[tp.Name]
			// // 				newDefArgsByName[innerTp.Name] = ta
			// // 			}
			// // 		}
			// // 	} else {
			// // 		// TypeArg is missing in the Definition, so get it from The Type
			// // 		ta = newType.TypeArguments[i]
			// // 	}

			// // 	newDefArgsByName[tp.Name] = ta
			// // }
			// // for n, t := range newDefArgsByName {
			// // 	log.Debug().Msgf("New Def TypeArg %s: %s", n, TypeToShortSyntax(t, true))
			// // }

			// // oldDefArgsByName := make(map[string]Type)
			// // for i, tp := range ch.PreviousDefinition().GetDefinitionMeta().TypeParameters {
			// // 	var ta Type
			// // 	if i < len(ch.PreviousDefinition().GetDefinitionMeta().TypeArguments) {
			// // 		// TypeArg already exists in the Definition
			// // 		ta = ch.PreviousDefinition().GetDefinitionMeta().TypeArguments[i]

			// // 		if st, ok := ta.(*SimpleType); ok {
			// // 			if innerTp, ok := st.ResolvedDefinition.(*GenericTypeParameter); ok {
			// // 				ta = oldTypeArgsByName[tp.Name]
			// // 				oldDefArgsByName[innerTp.Name] = ta
			// // 			}
			// // 		}
			// // 	} else {
			// // 		// TypeArg is missing in the Definition, so get it from The Type
			// // 		ta = oldType.TypeArguments[i]
			// // 	}

			// // 	oldDefArgsByName[tp.Name] = ta
			// // }
			// // for n, t := range oldDefArgsByName {
			// // 	log.Debug().Msgf("Old Def TypeArg %s: %s", n, TypeToShortSyntax(t, true))
			// // }
			// for i, _ := range newDef.GetDefinitionMeta().TypeParameters {
			// 	newTypeArg := newDef.GetDefinitionMeta().TypeArguments[i]
			// 	expectedTypeArg := newResolvedTypeArgs[i]

			// 	// log.Debug().Msgf("Comparing TypeArgs %s <= %s", TypeToShortSyntax(newTypeArg, true), TypeToShortSyntax(expectedTypeArg, true))
			// 	ch := compareTypes(newTypeArg, expectedTypeArg, context)
			// 	switch ch := ch.(type) {
			// 	case nil:
			// 		continue
			// 	case *TypeChangeDefinitionChanged:
			// 		typeArgDefinitionChanged = true
			// 		continue
			// 	default:
			// 		// ERROR: GenericTypeArguments don't match between versions
			// 		log.Error().Msgf("GenericTypeArguments don't match between versions")
			// 		return &TypeChangeIncompatible{TypePair{ch.OldType(), ch.NewType()}}
			// 	}
			// }

			// for i, _ := range oldDef.GetDefinitionMeta().TypeParameters {
			// 	oldTypeArg := oldDef.GetDefinitionMeta().TypeArguments[i]
			// 	expectedTypeArg := oldResolvedTypeArgs[i]

			// 	// log.Debug().Msgf("Comparing TypeArgs %s <= %s", TypeToShortSyntax(oldTypeArg, true), TypeToShortSyntax(expectedTypeArg, true))
			// 	ch := compareTypes(expectedTypeArg, oldTypeArg, context)
			// 	switch ch := ch.(type) {
			// 	case nil:
			// 		continue
			// 	case *TypeChangeDefinitionChanged:
			// 		typeArgDefinitionChanged = true
			// 		continue
			// 	default:
			// 		// ERROR: GenericTypeArguments don't match between versions
			// 		log.Error().Msgf("GenericTypeArguments don't match between versions")
			// 		return &TypeChangeIncompatible{TypePair{ch.OldType(), ch.NewType()}}
			// 	}
			// }
		}

		// for newName, newType := range newDefArgsByName {
		// 	if _, ok := oldDefArgsByName[newName]; ok {
		// 		log.Debug().Msgf("Comparing TypeArgs %s", newName)
		// 	} else {

		// 		if st, ok := newType.(*SimpleType); ok {
		// 			if tp, ok := st.ResolvedDefinition.(*GenericTypeParameter); ok {
		// 				_, ok := oldDefArgsByName[tp.Name]
		// 				if !ok {
		// 					log.Error().Msgf("Didn't find %s", tp.Name)
		// 				}

		// 				log.Debug().Msgf("Comparing new TypeArgs %s and with old %s", newName, tp.Name)
		// 			}
		// 		}
		// 	}
		// }

		// newTypeArgs := make(map[string]Type)
		// for _, ta := range ch.LatestDefinition().GetDefinitionMeta().TypeArguments {
		// 	// TODO: This needs to be recursive...

		// 	// log.Debug().Msgf("-----")
		// 	// Print(ta)
		// 	// log.Debug().Msgf("-----")

		// 	if st, ok := ta.(*SimpleType); ok {
		// 		if tp, ok := st.ResolvedDefinition.(*GenericTypeParameter); ok {
		// 			newTypeArgs[tp.Name] = newTypeArgsByName[tp.Name]
		// 			log.Debug().Msgf("TypeArg %s: %s", tp.Name, TypeToShortSyntax(newTypeArgs[tp.Name], true))
		// 		}
		// 	}
		// }

		// oldTypeArgs := make(map[string]Type)
		// for _, ta := range ch.PreviousDefinition().GetDefinitionMeta().TypeArguments {
		// 	// TODO: This needs to be recursive...

		// 	log.Debug().Msgf("-----")
		// 	Print(ta)
		// 	log.Debug().Msgf("-----")

		// 	if st, ok := ta.(*SimpleType); ok {
		// 		if tp, ok := st.ResolvedDefinition.(*GenericTypeParameter); ok {
		// 			oldTypeArgs[tp.Name] = oldTypeArgsByName[tp.Name]
		// 			log.Debug().Msgf("TypeArg %s: %s", tp.Name, TypeToShortSyntax(oldTypeArgs[tp.Name], true))
		// 		}
		// 	}
		// }

		// compared := make(map[string]bool)
		// for newName, newTypeArg := range newTypeArgs {
		// 	if oldTypeArg, ok := oldTypeArgsByName[newName]; ok {
		// 		log.Debug().Msgf("Comparing TypeArgs %s", newName)
		// 		ch := compareTypes(newTypeArg, oldTypeArg, context)
		// 		compared[newName] = true
		// 		switch ch := ch.(type) {
		// 		case nil, *TypeChangeDefinitionChanged:
		// 			continue
		// 		default:
		// 			// ERROR: GenericTypeArguments don't match between versions
		// 			log.Error().Msgf("GenericTypeArguments don't match between versions")
		// 			return &TypeChangeIncompatible{TypePair{ch.OldType(), ch.NewType()}}
		// 		}
		// 	}
		// }
		// for oldName, oldTypeArg := range oldTypeArgs {
		// 	if !compared[newName] {
		// 		if newTypeArg, ok := newTypeArgs[oldName]; ok {
		// 			log.Debug().Msgf("Comparing TypeArgs %s", oldName)
		// 			ch := compareTypes(newTypeArg, oldTypeArg, context)
		// 			if ch != nil {
		// 				continue
		// 			}
		// 			switch ch := ch.(type) {
		// 			case nil, *TypeChangeDefinitionChanged:
		// 				continue
		// 			default:
		// 				// ERROR: GenericTypeArguments don't match between versions
		// 				log.Error().Msgf("GenericTypeArguments don't match between versions")
		// 				return &TypeChangeIncompatible{TypePair{ch.OldType(), ch.NewType()}}
		// 			}
		// 		}
		// 	}
		// }

		// oldResolved, err := MakeGenericType(ch.PreviousDefinition(), newDef.GetDefinitionMeta().TypeArguments, false)
		// if err != nil {
		// 	log.Error().Msgf("%s", err)
		// 	// panic("These TypeDefinitions aren't compatible")
		// }

		// for i, ta := range oldDef.GetDefinitionMeta().TypeArguments {
		// 	ch := compareTypes(ta, oldResolved.GetDefinitionMeta().TypeArguments[i], context)
		// 	if ch == nil {
		// 		continue
		// 	}
		// 	switch ch := ch.(type) {
		// 	case *TypeChangeDefinitionChanged:
		// 		continue
		// 	default:
		// 		// ERROR: GenericTypeArguments don't match between versions
		// 		// return &TypeChangeIncompatible{TypePair{oldType, newType}}
		// 		return &TypeChangeIncompatible{TypePair{ch.OldType(), ch.NewType()}}
		// 	}
		// }
		// for _, ta := range ch.PreviousDefinition().GetDefinitionMeta().TypeArguments {
		// 	if ta == nil {
		// 		// ??
		// 		continue
		// 	}

		// 	if TypeContainsGenericTypeParameter(ta) {
		// 		log.Debug().Msgf("TypeArg: %s", TypeToShortSyntax(ta, true))
		// 		// ??
		// 		continue
		// 	}
		// }

		// if _, ok := ch.(*DefinitionChangeIncompatible); ok {
		// 	return &TypeChangeIncompatible{TypePair{oldType, newType}}
		// } else {
		// 	return &TypeChangeDefinitionChanged{TypePair{oldType, newType}, ch}
		// }
		if ch != nil || typeArgDefinitionChanged {
			return &TypeChangeDefinitionChanged{TypePair{oldType, newType}, ch}
		}
		log.Debug().Msgf("These types are the same: %s <= %s", defWithArgs(newDef), defWithArgs(oldDef))
		return nil

		// } else {
		// 	// We've already compared these underlying definitions and there's no change
		// 	// So we just need to compare TypeArguments

		// 	// assert
		// 	if len(newDef.GetDefinitionMeta().TypeParameters) != len(oldDef.GetDefinitionMeta().TypeParameters) {
		// 		panic("These TypeDefinitions aren't compatible")
		// 	}
		// 	if len(newDef.GetDefinitionMeta().TypeArguments) != len(oldDef.GetDefinitionMeta().TypeParameters) {
		// 		panic("These TypeDefinitions aren't compatible")
		// 	}

		// 	for i, tpNew := range newDef.GetDefinitionMeta().TypeParameters {
		// 		taNew := newDef.GetDefinitionMeta().TypeArguments[i]

		// 		tpOld := oldDef.GetDefinitionMeta().TypeParameters[i]
		// 		taOld := oldDef.GetDefinitionMeta().TypeArguments[i]

		// 		if tpNew.Name != tpOld.Name {
		// 			panic("These TypeDefinitions aren't compatible")
		// 		}

		// 		ch := compareTypes(taNew, taOld, context)
		// 		switch ch := ch.(type) {
		// 		case nil, *TypeChangeDefinitionChanged:
		// 			continue
		// 		default:
		// 			// ERROR: GenericTypeArguments don't match between versions
		// 			log.Error().Msgf("GenericTypeArguments don't match between versions")
		// 			return &TypeChangeIncompatible{TypePair{ch.OldType(), ch.NewType()}}
		// 		}
		// 	}

		// 	return nil
		// }
	}

	if !context.SemanticMapping[newName][oldName] {
		// Unwind old NamedType to compare underlying Types
		switch oldDef := oldDef.(type) {
		case *NamedType:
			ch := compareTypes(newType, oldDef.Type, context)
			if ch == nil {
				return nil
			}
			switch ch := ch.(type) {
			case *TypeChangeIncompatible:
			default:
				tdChange := &NamedTypeChange{DefinitionPair{oldDef, newDef}, ch}
				context.AliasesRemoved = append(context.AliasesRemoved, tdChange)
				// return &TypeChangeDefinitionChanged{TypePair{oldType, newType}, tdChange}
			}
			return ch
		}

		// Unwind new NamedType to compare underlying types
		switch newDef := newDef.(type) {
		case *NamedType:
			ch := compareTypes(newDef.Type, oldType, context)
			if ch == nil {
				return nil
			}
			return ch
		}

		// Check if we're comparing Primitives
		if _, ok := oldDef.(PrimitiveDefinition); ok {
			if _, ok := newDef.(PrimitiveDefinition); ok {
				return detectPrimitiveTypeChange(newType, oldType)
			}
			return &TypeChangeIncompatible{TypePair{oldType, newType}}
		}
		if _, ok := newDef.(PrimitiveDefinition); ok {
			return &TypeChangeIncompatible{TypePair{oldType, newType}}
		}

		// Check if we're comparing GenericTypeParameters
		if oldTp, ok := oldDef.(*GenericTypeParameter); ok {
			if newTp, ok := newDef.(*GenericTypeParameter); ok {
				if oldTp.Name == newTp.Name {
					return nil
				}
			}
			return &TypeChangeIncompatible{TypePair{oldType, newType}}
		}
		if _, ok := newDef.(*GenericTypeParameter); ok {
			return &TypeChangeIncompatible{TypePair{oldType, newType}}
		}

		// These types aren't compatible
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	// For Debugging:
	// for i, newTypeArg := range newType.TypeArguments {
	// 	newDefTypeArg := newDef.GetDefinitionMeta().TypeArguments[i]
	// 	if newTypeArg != newDefTypeArg {
	// 		panic("TypeArg doesn't match")
	// 	}
	// }
	// for i, oldTypeArg := range oldType.TypeArguments {
	// 	oldDefTypeArg := oldDef.GetDefinitionMeta().TypeArguments[i]
	// 	if oldTypeArg != oldDefTypeArg {
	// 		panic("TypeArg doesn't match")
	// 	}
	// }

	// {
	// 	// For more debugging:
	// 	getParams := func(ps []*GenericTypeParameter) string {
	// 		s := ""
	// 		for _, p := range ps {
	// 			s += p.Name + ","
	// 		}
	// 		return s
	// 	}
	// 	getArgs := func(ts []Type) string {
	// 		s := ""
	// 		for _, t := range ts {
	// 			s += TypeToShortSyntax(t, true) + ","
	// 		}
	// 		return s
	// 	}
	// 	// newType := GetUnderlyingType(newType).(*SimpleType)
	// 	// oldType := GetUnderlyingType(oldType).(*SimpleType)
	// 	newTypeParams := getParams(newType.ResolvedDefinition.GetDefinitionMeta().TypeParameters)
	// 	newTypeArgs := getArgs(newType.ResolvedDefinition.GetDefinitionMeta().TypeArguments)
	// 	oldTypeParams := getParams(oldType.ResolvedDefinition.GetDefinitionMeta().TypeParameters)
	// 	oldTypeArgs := getArgs(oldType.ResolvedDefinition.GetDefinitionMeta().TypeArguments)

	// 	log.Debug().Msgf("Old %s has params %s and args %s", oldName, oldTypeParams, oldTypeArgs)
	// 	log.Debug().Msgf("New %s has params %s and args %s", newName, newTypeParams, newTypeArgs)
	// }

	// Unwind old NamedType to compare underlying Types
	// switch oldDef := oldDef.(type) {
	// case *NamedType:
	// 	ch := compareTypes(newType, oldDef.Type, context)
	// 	if ch == nil {
	// 		return nil
	// 	}
	// 	switch ch := ch.(type) {
	// 	case *TypeChangeIncompatible:
	// 	default:
	// 		tdChange := &NamedTypeChange{DefinitionPair{oldDef, newDef}, ch}
	// 		context.AliasesRemoved = append(context.AliasesRemoved, tdChange)
	// 		// return &TypeChangeDefinitionChanged{TypePair{oldType, newType}, tdChange}
	// 	}
	// 	return ch
	// }

	// // Unwind new NamedType to compare underlying types
	// switch newDef := newDef.(type) {
	// case *NamedType:
	// 	ch := compareTypes(newDef.Type, oldType, context)
	// 	if ch == nil {
	// 		return nil
	// 	}
	// 	return ch
	// }

	// if len(newType.TypeArguments) > 0 && len(oldType.TypeArguments) > 0 {
	// 	// Both types have TypeArguments, so we're comparing two Generic Types
	// 	// For now, the TypeArguments must be identical
	// 	// To determine whether the TypeArguments match, we need to fully resolve the Generic Type here

	// 	newType := GetUnderlyingType(newType).(*SimpleType)
	// 	oldType := GetUnderlyingType(oldType).(*SimpleType)

	// 	newTypeArgs := newType.ResolvedDefinition.GetDefinitionMeta().TypeArguments
	// 	oldTypeArgs := oldType.ResolvedDefinition.GetDefinitionMeta().TypeArguments
	// 	if len(newTypeArgs) != len(oldTypeArgs) {
	// 		log.Warn().Msgf("TypeArguments changed: %s <= %s", TypeToShortSyntax(newType, true), TypeToShortSyntax(oldType, true))
	// 		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	// 	}

	// 	for i, newTypeArg := range newType.ResolvedDefinition.GetDefinitionMeta().TypeArguments {
	// 		oldTypeArg := oldType.ResolvedDefinition.GetDefinitionMeta().TypeArguments[i]
	// 		if ch := compareTypes(newTypeArg, oldTypeArg, context); ch != nil {
	// 			// log.Warn().Msgf("TypeArgument %d changed: %s <= %s", i, newName, oldName)
	// 			log.Warn().Msgf("TypeArgument %d changed: %s <= %s", i, TypeToShortSyntax(newType, true), TypeToShortSyntax(oldType, true))
	// 			return &TypeChangeIncompatible{TypePair{oldType, newType}}
	// 		}
	// 	}
	// }

	// At this point, we know the TypeDefinitions are semantically equivalent
	// TODO: This check may need to occur EARLIER because we could be comparing a GenericType with a non-generic, but EQUAL type
	// if len(newDef.GetDefinitionMeta().TypeArguments) > 0 || len(oldDef.GetDefinitionMeta().TypeArguments) > 0 {
	// 	if context.Compared[newName][oldName] {
	// 		panic("Shouldn't get here")
	// 	}

	// 	typeDefChange := compareTypeDefinitions(newDef, oldDef, context)
	// 	if typeDefChange == nil {
	// 		return nil
	// 	}

	// 	log.Error().Msgf("Generic Type changed: %s <= %s", typeDefChange.LatestDefinition().GetDefinitionMeta().GetQualifiedName(), typeDefChange.PreviousDefinition().GetDefinitionMeta().GetQualifiedName())

	// 	switch typeDefChange.(type) {
	// 	case *DefinitionChangeIncompatible:
	// 		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	// 	default:
	// 		return &TypeChangeDefinitionChanged{TypePair{oldType, newType}, typeDefChange}
	// 	}
	// }

	// The two TypeDefinitions are semantically equivalent, so context.Compared[newName][oldName] should be true
	// Did the TypeDefinition change between versions?
	// if ch, ok := context.Changes[newName][oldName]; ok {
	// 	if _, ok := ch.(*DefinitionChangeIncompatible); ok {
	// 		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	// 	} else {
	// 		return &TypeChangeDefinitionChanged{TypePair{oldType, newType}, ch}
	// 	}
	// }

	panic("Shouldn't get here")
	return nil
}

func compareGeneralizedToSimpleTypes(newType *SimpleType, oldType *GeneralizedType, context *EvolutionContext) TypeChange {
	// Is it a change from Optional<T> to T (partially compatible)
	if oldType.Cases.IsOptional() {
		switch compareTypes(newType, oldType.Cases[1].Type, context).(type) {
		case nil, *TypeChangeDefinitionChanged:
			return &TypeChangeOptionalToScalar{TypePair{oldType, newType}}
		}
	}

	// Is it a change from Union<T, ...> to T (partially compatible)
	if oldType.Cases.IsUnion() {
		for i, tc := range oldType.Cases {
			switch compareTypes(newType, tc.Type, context).(type) {
			case nil, *TypeChangeDefinitionChanged:
				return &TypeChangeUnionToScalar{TypePair{oldType, newType}, i}
			}
		}
	}

	// CHANGE: Incompatible change from Generalized to Simple
	return &TypeChangeIncompatible{TypePair{oldType, newType}}
}

func compareSimpleToGeneralizedTypes(newType *GeneralizedType, oldType *SimpleType, context *EvolutionContext) TypeChange {
	// Is it a change from T to Optional<T> (partially compatible)
	if newType.Cases.IsOptional() {
		switch compareTypes(newType.Cases[1].Type, oldType, context).(type) {
		case nil, *TypeChangeDefinitionChanged:
			return &TypeChangeScalarToOptional{TypePair{oldType, newType}}
		}
	}

	// Is it a change from T to Union<T, ...> (partially compatible)
	if newType.Cases.IsUnion() {
		for i, tc := range newType.Cases {
			switch compareTypes(tc.Type, oldType, context).(type) {
			case nil, *TypeChangeDefinitionChanged:
				return &TypeChangeScalarToUnion{TypePair{oldType, newType}, i}
			}
		}
	}

	// CHANGE: Incompatible change from Simple to Generalized
	return &TypeChangeIncompatible{TypePair{oldType, newType}}
}

func compareGeneralizedTypes(newType, oldType *GeneralizedType, context *EvolutionContext) TypeChange {
	// A GeneralizedType can change in many ways...
	if newType.Cases.IsOptional() {
		return detectOptionalChanges(newType, oldType, context)
	} else if newType.Cases.IsUnion() {
		return detectUnionChanges(newType, oldType, context)
	} else {
		switch newType.Dimensionality.(type) {
		case nil:
			// Not an Optional, Union, Stream, Vector, Array, Map...
		case *Stream:
			return detectStreamChanges(newType, oldType, context)
		case *Vector:
			return detectVectorChanges(newType, oldType, context)
		case *Array:
			return detectArrayChanges(newType, oldType, context)
		case *Map:
			return detectMapChanges(newType, oldType, context)
		default:
			panic("Shouldn't get here")
		}
	}

	return nil
}

func detectPrimitiveTypeChange(newType, oldType *SimpleType) TypeChange {
	newPrimitive := newType.ResolvedDefinition.(PrimitiveDefinition)
	oldPrimitive := oldType.ResolvedDefinition.(PrimitiveDefinition)

	if newPrimitive == oldPrimitive {
		return nil
	}

	// CHANGE: Changed Primitive type
	if oldPrimitive == PrimitiveString {
		if GetPrimitiveKind(newPrimitive) == PrimitiveKindInteger || GetPrimitiveKind(newPrimitive) == PrimitiveKindFloatingPoint {
			return &TypeChangeStringToNumber{TypePair{oldType, newType}}
		}
	}

	if GetPrimitiveKind(oldPrimitive) == PrimitiveKindInteger || GetPrimitiveKind(oldPrimitive) == PrimitiveKindFloatingPoint {
		if newPrimitive == PrimitiveString {
			return &TypeChangeNumberToString{TypePair{oldType, newType}}
		}

		if GetPrimitiveKind(newPrimitive) == PrimitiveKindInteger || GetPrimitiveKind(newPrimitive) == PrimitiveKindFloatingPoint {
			return &TypeChangeNumberToNumber{TypePair{oldType, newType}}
		}
	}

	return &TypeChangeIncompatible{TypePair{oldType, newType}}
}

func detectOptionalChanges(newType, oldType *GeneralizedType, context *EvolutionContext) TypeChange {
	if !oldType.Cases.IsOptional() {
		if oldType.Cases.IsUnion() && oldType.Cases.HasNullOption() {
			// An Optional<T> can become a Union<null, T, ...> ONLY if
			// 	1. type T does not change, or
			// 	2. type T's TypeDefinition changed

			// Look for a matching type in the old Union
			for i, c := range oldType.Cases[1:] {
				switch compareTypes(newType.Cases[1].Type, c.Type, context).(type) {
				case nil, *TypeChangeDefinitionChanged:
					return &TypeChangeUnionToOptional{TypePair{oldType, newType}, i + 1}
				}
			}
		}

		// CHANGE: Changed a non-Optional/Union to an Optional
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if ch := compareTypes(newType.Cases[1].Type, oldType.Cases[1].Type, context); ch != nil {
		// CHANGE: Changed Optional type
		return &TypeChangeOptionalTypeChanged{TypePair{oldType, newType}, ch}
	}
	return nil
}

func detectUnionChanges(newType, oldType *GeneralizedType, context *EvolutionContext) TypeChange {
	if !oldType.Cases.IsUnion() {
		if oldType.Cases.IsOptional() && newType.Cases.HasNullOption() {
			for i, c := range newType.Cases[1:] {
				switch compareTypes(c.Type, oldType.Cases[1].Type, context).(type) {
				case nil, *TypeChangeDefinitionChanged:
					return &TypeChangeOptionalToUnion{TypePair{oldType, newType}, i + 1}
				}
			}
		}
		// CHANGE: Changed a non-Union/Optional to a Union
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	oldMatches := make([]bool, len(oldType.Cases))
	newMatches := make([]bool, len(newType.Cases))

	innerTypeDefsChanged := false
	// Search for a match for each Type in the new Union
	for i, newCase := range newType.Cases {
		for j, oldCase := range oldType.Cases {
			if oldMatches[j] {
				continue
			}

			log.Debug().Msgf("  Comparing Union cases %s and %s", TypeToShortSyntax(newCase.Type, true), TypeToShortSyntax(oldCase.Type, true))
			switch compareTypes(newCase.Type, oldCase.Type, context).(type) {
			case nil:
				// Found matching type
				newMatches[i] = true
				oldMatches[j] = true
			case *TypeChangeDefinitionChanged:
				// Found matching type with underlying definition changed
				newMatches[i] = true
				oldMatches[j] = true
				innerTypeDefsChanged = true
			}

			if newMatches[i] {
				break
			}
		}
	}

	// If newMatches is all False, then this isn't a valid Union type change
	// If newMatches is not all True, then type(s) were added to the Union
	// If oldMatches is not all True, then type(s) were removed from the Union
	// If newMatches and oldMatches are all true, then the Union types are the same, but possibly reordered
	anyMatch := false
	allMatch := true
	for _, m := range newMatches {
		if !m {
			allMatch = false
		} else {
			anyMatch = true
		}
	}
	for _, m := range oldMatches {
		if !m {
			allMatch = false
		}
	}

	if !anyMatch {
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if innerTypeDefsChanged || !allMatch {
		return &TypeChangeUnionTypesetChanged{TypePair{oldType, newType}, oldMatches, newMatches}
	}

	return nil
}

func detectStreamChanges(newType, oldType *GeneralizedType, context *EvolutionContext) TypeChange {
	if _, ok := oldType.Dimensionality.(*Stream); !ok {
		// CHANGE: Changed a non-Stream to a Stream
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if ch := compareTypes(newType.Cases[0].Type, oldType.Cases[0].Type, context); ch != nil {
		// CHANGE: Changed Stream type
		return &TypeChangeStreamTypeChanged{TypePair{oldType, newType}, ch}
	}
	return nil
}

func detectVectorChanges(newType, oldType *GeneralizedType, context *EvolutionContext) TypeChange {
	newDim := newType.Dimensionality.(*Vector)
	oldDim, ok := oldType.Dimensionality.(*Vector)
	if !ok {
		// CHANGE: Changed a non-Vector to a Vector
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if (oldDim.Length == nil) != (newDim.Length == nil) {
		// CHANGE: Changed from a fixed-length Vector to a variable-length Vector or vice versa
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}
	if newDim.Length != nil && *newDim.Length != *oldDim.Length {
		// CHANGE: Changed vector length
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if ch := compareTypes(newType.Cases[0].Type, oldType.Cases[0].Type, context); ch != nil {
		// CHANGE: Changed Vector type
		return &TypeChangeVectorTypeChanged{TypePair{oldType, newType}, ch}
	}

	return nil
}

func detectArrayChanges(newType, oldType *GeneralizedType, context *EvolutionContext) TypeChange {
	newDim := newType.Dimensionality.(*Array)
	oldDim, ok := oldType.Dimensionality.(*Array)
	if !ok {
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if ch := compareTypes(newType.Cases[0].Type, oldType.Cases[0].Type, context); ch != nil {
		// CHANGE: Changed Array type
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if (newDim.Dimensions == nil) != (oldDim.Dimensions == nil) {
		// CHANGE: Added or removed array dimensions
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if newDim.Dimensions != nil {
		newDimensions := *newDim.Dimensions
		oldDimensions := *oldDim.Dimensions

		if len(newDimensions) != len(oldDimensions) {
			// CHANGE: Mismatch in number of array dimensions
			return &TypeChangeIncompatible{TypePair{oldType, newType}}
		}

		for i, newDimension := range newDimensions {
			oldDimension := oldDimensions[i]

			if (newDimension.Length == nil) != (oldDimension.Length == nil) {
				// CHANGE: Added or removed array dimension length
				return &TypeChangeIncompatible{TypePair{oldType, newType}}
			}

			if newDimension.Length != nil && *newDimension.Length != *oldDimension.Length {
				// CHANGE: Changed array dimension length
				return &TypeChangeIncompatible{TypePair{oldType, newType}}
			}
		}
	}
	return nil
}

func detectMapChanges(newType, oldType *GeneralizedType, context *EvolutionContext) TypeChange {
	newDim := newType.Dimensionality.(*Map)
	oldDim, ok := oldType.Dimensionality.(*Map)
	if !ok {
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if ch := compareTypes(newDim.KeyType, oldDim.KeyType, context); ch != nil {
		// CHANGE: Changed Map key type
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}
	if ch := compareTypes(newType.Cases[0].Type, oldType.Cases[0].Type, context); ch != nil {
		// CHANGE: Changed Map value type
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}
	return nil
}
