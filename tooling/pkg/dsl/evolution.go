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
		ns.DefinitionChanges = make(map[string][]DefinitionChange)
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
		// for _, ch := range definitionChanges {
		// 	log.Debug().Msgf("DefinitionChange OLD %s = NEW %s", defWithArgs(ch.PreviousDefinition()), defWithArgs(ch.LatestDefinition()))
		// }

		// Save all TypeDefinition changes for codegen
		latest.GetTopLevelNamespace().DefinitionChanges[versionLabels[i]] = definitionChanges

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

		// log.Debug().Msgf("Renaming OLD %s", td.GetDefinitionMeta().GetQualifiedName())
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

func getAllTypeDefinitions(env *Environment) []TypeDefinition {
	allTypeDefs := make([]TypeDefinition, 0)
	for _, ns := range env.Namespaces {
		for _, td := range ns.TypeDefinitions {
			allTypeDefs = append(allTypeDefs, td)
		}
	}
	return allTypeDefs
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

// Returns a map of semantically eqivalent TypeDefinitions within a model
func getIntraModelSemanticEquivalence(tds []TypeDefinition) map[string]map[string]bool {
	resolvesTo := make(map[string]map[string]bool)
	for _, tdA := range tds {
		tdAName := tdA.GetDefinitionMeta().GetQualifiedName()
		table := make(map[string]bool)
		for _, tdB := range tds {
			tdBName := tdB.GetDefinitionMeta().GetQualifiedName()
			if getResolvedName(tdA) == getResolvedName(tdB) {
				table[tdBName] = true
			}
		}
		resolvesTo[tdAName] = table
	}
	return resolvesTo
}

// Returns all NamedTypes that reference the given TypeDefinition, plus itself
func getReferencingDefinitions(td TypeDefinition, allTypeDefs []TypeDefinition) []TypeDefinition {
	parents := make([]TypeDefinition, 0)
	childName := td.GetDefinitionMeta().GetQualifiedName()

	for _, parentTd := range allTypeDefs {
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

// Useful for debugging Generic TypeDefinition resolution
func defWithArgs(td TypeDefinition) string {
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
		case PrimitiveDefinition, *GenericTypeParameter:
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
	// log.Debug().Msgf("Resolving %s <= %s", defWithArgs(newTd), defWithArgs(oldTd))
	resolveGenericDefinitions := func(out TypeDefinition, in TypeDefinition) TypeDefinition {
		typeArgs := make([]Type, 0)
		Visit(in, func(self Visitor, node Node) {
			switch node := node.(type) {
			case *SimpleType:
				self.Visit(node.ResolvedDefinition)

			case TypeDefinition:
				in = node
				if len(node.GetDefinitionMeta().TypeArguments) == len(out.GetDefinitionMeta().TypeParameters) {
					typeArgs = node.GetDefinitionMeta().TypeArguments
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
				log.Panic().Msgf("Unable to resolve generic TypeDefinition %s <= %s", defWithArgs(out), defWithArgs(in))
			}
		}

		return out
	}

	newDef := resolveGenericDefinitions(newTd, oldTd)
	oldDef := resolveGenericDefinitions(oldTd, newTd)
	// log.Debug().Msgf(" Resolved %s <= %s", defWithArgs(newDef), defWithArgs(oldDef))

	return &DefinitionPair{oldDef, newDef}
}

// Determine semantic equivalence between Old and New TypeDefinitions
func getCrossModelSemanticEquivalence(newTds, oldTds []TypeDefinition) map[string]map[string]bool {
	newNameMapping := getIntraModelSemanticEquivalence(newTds)
	oldNameMapping := getIntraModelSemanticEquivalence(oldTds)

	semanticallyEqual := make(map[string]map[string]bool)
	for _, newTd := range newTds {
		newName := newTd.GetDefinitionMeta().GetQualifiedName()
		semanticallyEqual[newName] = make(map[string]bool)
	}

	for _, newTd := range newTds {
		newName := newTd.GetDefinitionMeta().GetQualifiedName()
		for _, oldTd := range oldTds {
			oldName := oldTd.GetDefinitionMeta().GetQualifiedName()
			if newName == oldName {
				// We found a mapping between versions...
				for n := range newNameMapping[newName] {
					for o := range oldNameMapping[oldName] {
						// n and o are semantically equal iff
						// 1. n is not in the old model, OR
						if _, found := oldNameMapping[n]; !found {
							semanticallyEqual[n][o] = true
							continue
						}
						// 2. o is not in the new model, OR
						if _, found := newNameMapping[o]; !found {
							semanticallyEqual[n][o] = true
							continue
						}

						// 3. o and n are equivalent names in both versions
						if oldNameMapping[n][o] && newNameMapping[o][n] {
							semanticallyEqual[n][o] = true
							continue
						}
					}
				}
			}
		}
	}

	// For debugging:
	// for _, newTd := range newTds {
	// 	newName := newTd.GetDefinitionMeta().GetQualifiedName()
	// 	for _, oldTd := range oldTds {
	// 		oldName := oldTd.GetDefinitionMeta().GetQualifiedName()
	// 		if semanticallyEqual[newName][oldName] {
	// 			log.Debug().Msgf("New %s === Old %s", newName, oldName)
	// 		}
	// 	}
	// }

	return semanticallyEqual
}

func resolveAllChanges(newEnv, oldEnv *Environment) ([]DefinitionChange, map[string]*ProtocolChange) {
	allNewTypeDefs := getAllTypeDefinitions(newEnv)
	allOldTypeDefs := getAllTypeDefinitions(oldEnv)
	// semanticallyEqual := getCrossModelSemanticEquivalence(allNewTypeDefs, allOldTypeDefs)

	// Find all "base" TypeDefinitions that are semantically equivalent.
	oldDefToCompareWith := make(map[string]string)
	newDefToCompareWith := make(map[string]string)
	for _, newTd := range allNewTypeDefs {
		newName := newTd.GetDefinitionMeta().GetQualifiedName()
		for _, oldTd := range allOldTypeDefs {
			oldName := oldTd.GetDefinitionMeta().GetQualifiedName()

			savePairIfBothUnique := func(newDef, oldDef TypeDefinition) {
				newName := newDef.GetDefinitionMeta().GetQualifiedName()
				oldName := oldDef.GetDefinitionMeta().GetQualifiedName()

				_, haveOldDef := oldDefToCompareWith[newName]
				_, haveNewDef := newDefToCompareWith[oldName]

				if !haveOldDef && !haveNewDef {
					log.Debug().Msgf("Base Pair %s <= %s", defWithArgs(newTd), defWithArgs(oldTd))
					oldDefToCompareWith[newName] = oldName
					newDefToCompareWith[oldName] = newName
				}

			}

			// First match by Definition Name
			if newName == oldName {
				// We want to compare their "base" TypeDefinitions
				savePairIfBothUnique(getBaseDefinition(newTd), getBaseDefinition(oldTd))
				// We also want to compare them if they are the same NamedType
				savePairIfBothUnique(newTd, oldTd)
			}
		}
	}

	// SANITY check
	for _, newTd := range allNewTypeDefs {
		newName := newTd.GetDefinitionMeta().GetQualifiedName()
		for _, oldTd := range allOldTypeDefs {
			oldName := oldTd.GetDefinitionMeta().GetQualifiedName()

			if oldDefToCompareWith[newName] == oldName {
				if newDefToCompareWith[oldName] != newName {
					log.Panic().Msgf("Mismatched mapping %s <= %s", newName, oldName)
				}
			}
		}
	}

	// Generate a mapping of all TypeDefinitions that are semantically equivalent
	// across versions, i.e. the {New, Old} pair that should be compared for changes
	// Resolve all Generic TypeDefinitions across versions along the way
	semanticDefinitionPairs := make(map[string]map[string]*DefinitionPair)
	for _, newTd := range allNewTypeDefs {
		newName := newTd.GetDefinitionMeta().GetQualifiedName()
		semanticDefinitionPairs[newName] = make(map[string]*DefinitionPair)
	}
	for _, newTd := range allNewTypeDefs {
		newName := newTd.GetDefinitionMeta().GetQualifiedName()
		for _, oldTd := range allOldTypeDefs {
			oldName := oldTd.GetDefinitionMeta().GetQualifiedName()

			if oldDefToCompareWith[newName] == oldName {
				// These TypeDefinitions are semantically equivalent (cross-version) "base" Definitions
				// Resolve all definitions that reference this TypeDefinition
				for _, newParent := range getReferencingDefinitions(newTd, allNewTypeDefs) {
					newParentName := newParent.GetDefinitionMeta().GetQualifiedName()
					for _, oldParent := range getReferencingDefinitions(oldTd, allOldTypeDefs) {
						oldParentName := oldParent.GetDefinitionMeta().GetQualifiedName()

						if _, ok := semanticDefinitionPairs[newParentName][oldParentName]; !ok {
							resolved := getResolvedGenericDefinition(newParent, oldParent)
							semanticDefinitionPairs[newParentName][oldParentName] = resolved

							log.Debug().Msgf("Resolved %s <= %s", defWithArgs(resolved.LatestDefinition()), defWithArgs(resolved.PreviousDefinition()))
						}
					}
				}
			}
		}
	}

	changes := make(map[string]map[string]DefinitionChange)
	for _, newTd := range allNewTypeDefs {
		newName := newTd.GetDefinitionMeta().GetQualifiedName()
		changes[newName] = make(map[string]DefinitionChange)
	}

	context := &EvolutionContext{
		Changes:       changes,
		SemanticPairs: semanticDefinitionPairs,
	}

	for _, newTd := range allNewTypeDefs {
		newName := newTd.GetDefinitionMeta().GetQualifiedName()
		for _, oldTd := range allOldTypeDefs {
			oldName := oldTd.GetDefinitionMeta().GetQualifiedName()

			if oldDefToCompareWith[newName] == oldName {
				// These are semantically equivalent "base" TypeDefinitions

				resolvedPair := context.SemanticPairs[newName][oldName]
				if resolvedPair == nil {
					log.Panic().Msgf("Should have already compared %s <= %s", newName, oldName)
				}
				newDef := resolvedPair.LatestDefinition()
				oldDef := resolvedPair.PreviousDefinition()

				ch := compareTypeDefinitions(newDef, oldDef, context)
				log.Info().Msgf("Saving %T for %s <= %s", ch, newName, oldName)
				changes[newName][oldName] = ch

				// Now we want to ensure all NamedTypes that REFERENCE this changed Definition are also marked as "changed"
				// 		If base change is nil, save nil for pair
				// 		If base change was valid, save NamedTypeChange for pair
				// 		If base change was not valid, save DefinitionChangeIncompatible for pair
				for _, newParent := range getReferencingDefinitions(newDef, allNewTypeDefs) {
					newParentName := newParent.GetDefinitionMeta().GetQualifiedName()
					for _, oldParent := range getReferencingDefinitions(oldDef, allOldTypeDefs) {
						oldParentName := oldParent.GetDefinitionMeta().GetQualifiedName()

						// Skip if these are meant to be directly compared
						if oldDefToCompareWith[newParentName] == oldParentName {
							continue
						}

						// Skip if we already saved a Change for these Definitions
						if _, ok := changes[newParentName][oldParentName]; ok {
							continue
						}

						resolvedPair := context.SemanticPairs[newParentName][oldParentName]
						if resolvedPair == nil {
							log.Panic().Msgf("Should have already compared %s <= %s", newParentName, oldParentName)
						}

						switch ch.(type) {
						case nil:
							changes[newParentName][oldParentName] = nil
						case *DefinitionChangeIncompatible:
							// log.Debug().Msgf("Saving Incompatible Change for %s <= %s", newParentName, oldParentName)
							changes[newParentName][oldParentName] = &DefinitionChangeIncompatible{*resolvedPair}
						default:
							// log.Debug().Msgf("Saving NamedType Change for %s <= %s", newParentName, oldParentName)
							changes[newParentName][oldParentName] = &NamedTypeChange{*resolvedPair, nil}
						}
					}
				}
			}
		}
	}

	// Collect all DefinitionChanges in OLD Definition order - the order in which they'll be referenced by codegen
	// While simultaneously filtering so we only produce one DefinitionChange per each OLD TypeDefinition
	finalDefinitionChanges := make([]DefinitionChange, 0)
	uniqueOldNames := make(map[string]bool)
	newNameMapping := getIntraModelSemanticEquivalence(allNewTypeDefs)
	for _, oldTd := range allOldTypeDefs {
		oldName := oldTd.GetDefinitionMeta().GetQualifiedName()

		// First, look for its base semantic pair and emit that DefinitionChange
		for _, newTd := range allNewTypeDefs {
			newName := newTd.GetDefinitionMeta().GetQualifiedName()

			if uniqueOldNames[oldName] {
				// continue
				break
			}

			if newDefToCompareWith[oldName] == newName {

				change, ok := changes[newName][oldName]
				if !ok {
					log.Panic().Msgf("Should have already compared %s <= %s", newName, oldName)
				}
				if change == nil {
					continue
				}

				log.Debug().Msgf("Emitting %T for %s <= %s", change, defWithArgs(change.LatestDefinition()), defWithArgs(change.PreviousDefinition()))
				finalDefinitionChanges = append(finalDefinitionChanges, change)
				uniqueOldNames[oldName] = true

			}
		}

		// Then emit all Changes for NamedTypes that reference these TypeDefinitions
		for _, newTd := range allNewTypeDefs {
			newName := newTd.GetDefinitionMeta().GetQualifiedName()

			if uniqueOldNames[oldName] {
				// continue
				break
			}

			if resolvedPair, compared := context.SemanticPairs[newName][oldName]; compared {
				change, ok := changes[newName][oldName]
				if !ok {
					log.Panic().Msgf("Should have already compared %s <= %s", newName, oldName)
				}
				if change == nil {
					if _, ok := newNameMapping[oldName]; !ok {
						// This old TypeDefinition was removed or renamed using aliases in the new model
						change = &CompatibilityChange{*resolvedPair}
					} else {
						continue
					}
				}

				log.Debug().Msgf("Emitting %T for %s <= %s", change, defWithArgs(change.LatestDefinition()), defWithArgs(change.PreviousDefinition()))
				finalDefinitionChanges = append(finalDefinitionChanges, change)
				uniqueOldNames[oldName] = true
			}
		}
	}

	// Now we're finished comparing all TypeDefinitions and we can finally compare Protocols
	allProtocolChanges := resolveAllProtocolChanges(newEnv, oldEnv, context)

	return finalDefinitionChanges, allProtocolChanges
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
			if protocolChange := compareProtocolDefinitions(newProt, oldProt, context); protocolChange != nil {
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
	Changes        map[string]map[string]DefinitionChange
	AliasesRemoved []DefinitionChange
	SemanticPairs  map[string]map[string]*DefinitionPair
}

func compareTypeDefinitions(newTd, oldTd TypeDefinition, context *EvolutionContext) DefinitionChange {
	// log.Debug().Msgf("Comparing TypeDefinitions %s <= %s", newTd.GetDefinitionMeta().GetQualifiedName(), oldTd.GetDefinitionMeta().GetQualifiedName())

	newName := newTd.GetDefinitionMeta().GetQualifiedName()
	oldName := oldTd.GetDefinitionMeta().GetQualifiedName()

	if _, valid := context.SemanticPairs[newName][oldName]; !valid {
		return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
	}

	// If we already compared these definitions, return early
	if ch, ok := context.Changes[newName][oldName]; ok {
		return ch
	}

	// TypeParameters must effectively be identical.
	// Any TypeDefinitions with mismatched TypeParameters should be NamedTypes handled in the calling function.
	if len(newTd.GetDefinitionMeta().TypeParameters) != len(oldTd.GetDefinitionMeta().TypeParameters) {
		return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
	}

	switch newTd := newTd.(type) {
	case *NamedType:
		switch oldTd := oldTd.(type) {
		case *NamedType:
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

			ch := compareTypeDefinitions(newTd, oldType.ResolvedDefinition, context)
			if ch == nil {
				return nil
			}
			switch ch := ch.(type) {
			case *DefinitionChangeIncompatible:
				return ch
			default:
				// Alias Removed
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
			switch ch := ch.(type) {
			case nil:
				return nil
			case *DefinitionChangeIncompatible:
				return ch
			default:
				// Alias Removed
				return &NamedTypeChange{DefinitionPair{oldTd, newTd}, nil}
			}

		default:
			return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
		}

	default:
		log.Panic().Msgf("Expected a TypeDefinition... Not a %T", newTd)
	}

	return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
}

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
	switch newType := newType.(type) {

	case *SimpleType:
		switch oldType := oldType.(type) {
		case *SimpleType:
			return compareSimpleTypes(newType, oldType, context)
		case *GeneralizedType:
			newDef := newType.ResolvedDefinition
			if nt, ok := newDef.(*NamedType); ok {
				return compareTypes(nt.Type, oldType, context)
			}
			return compareGeneralizedToSimpleTypes(newType, oldType, context)
		default:
			return &TypeChangeIncompatible{TypePair{oldType, newType}}
		}

	case *GeneralizedType:
		switch oldType := oldType.(type) {
		case *GeneralizedType:
			return compareGeneralizedTypes(newType, oldType, context)
		case *SimpleType:
			oldDef := oldType.ResolvedDefinition
			if nt, ok := oldDef.(*NamedType); ok {
				return compareTypes(newType, nt.Type, context)
			}
			return compareSimpleToGeneralizedTypes(newType, oldType, context)
		default:
			return &TypeChangeIncompatible{TypePair{oldType, newType}}
		}

	case nil:
		switch oldType.(type) {
		case nil:
			return nil
		default:
			return &TypeChangeIncompatible{TypePair{oldType, newType}}
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

	// Compare the ResolvedDefinitions
	// 1. Are they semantically equal?
	// 2. Did we already compare them?
	// 3. Are TypeArguments compatible?
	if _, ok := context.SemanticPairs[newName][oldName]; ok {
		return compareSemanticallyEquivalentTypes(newType, oldType, context)
	} else {
		return compareOtherSimpleTypes(newType, oldType, context)
	}
}

func compareSemanticallyEquivalentTypes(newType, oldType *SimpleType, context *EvolutionContext) TypeChange {
	newDef := newType.ResolvedDefinition
	oldDef := oldType.ResolvedDefinition
	newName := newDef.GetDefinitionMeta().GetQualifiedName()
	oldName := oldDef.GetDefinitionMeta().GetQualifiedName()
	// log.Debug().Msgf("Checking: %s <= %s", defWithArgs(newDef), defWithArgs(oldDef))

	ch, ok := context.Changes[newName][oldName]
	if !ok {
		log.Panic().Msgf("Haven't yet compared %s <= %s", newName, oldName)
	}
	switch ch.(type) {
	case *DefinitionChangeIncompatible:
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	typeArgDefinitionChanged := false
	if len(newType.TypeArguments) > 0 || len(oldType.TypeArguments) > 0 {
		if len(newType.TypeArguments) == len(oldType.TypeArguments) {
			// We can just compare TypeArguments
			for i := range newType.TypeArguments {
				newTypeArg := newDef.GetDefinitionMeta().TypeArguments[i]
				oldTypeArg := oldDef.GetDefinitionMeta().TypeArguments[i]

				if ch := compareTypes(newTypeArg, oldTypeArg, context); ch != nil {
					switch ch.(type) {
					case *TypeChangeDefinitionChanged:
						typeArgDefinitionChanged = true
					default:
						// log.Error().Msgf("TypeArgs aren't compatible %s <= %s", TypeToShortSyntax(newTypeArg, true), TypeToShortSyntax(oldTypeArg, true))
						return &TypeChangeIncompatible{TypePair{oldType, newType}}
					}
				}
			}
		} else {
			// TypeArguments are not directly compatible
			// We have to resolve the Definition with fewer TypeArguments
			// Then compare TypeArguments by name

			resolvedGenericDefinitions := context.SemanticPairs[newName][oldName]
			newDefReference := resolvedGenericDefinitions.LatestDefinition()
			oldDefReference := resolvedGenericDefinitions.PreviousDefinition()
			// log.Debug().Msgf("Against: %s <= %s", defWithArgs(newDefReference), defWithArgs(oldDefReference))

			newTypeArgsByName := make(map[string]Type)
			for i, tp := range newDef.GetDefinitionMeta().TypeParameters {
				ta := newDef.GetDefinitionMeta().TypeArguments[i]
				newTypeArgsByName[tp.Name] = ta
				// log.Debug().Msgf("  New TA %s: %s", tp.Name, TypeToShortSyntax(ta, true))
			}

			oldTypeArgsByName := make(map[string]Type)
			for i, tp := range oldDef.GetDefinitionMeta().TypeParameters {
				ta := oldDef.GetDefinitionMeta().TypeArguments[i]
				oldTypeArgsByName[tp.Name] = ta
				// log.Debug().Msgf("  Old TA %s: %s", tp.Name, TypeToShortSyntax(ta, true))
			}

			newRefTypeArgs := make(map[string]Type)
			for i, tp := range newDefReference.GetDefinitionMeta().TypeParameters {
				if i < len(newDefReference.GetDefinitionMeta().TypeArguments) {
					ta := newDefReference.GetDefinitionMeta().TypeArguments[i]
					if st, ok := ta.(*SimpleType); ok {
						if innerTp, ok := st.ResolvedDefinition.(*GenericTypeParameter); ok {
							newRefTypeArgs[tp.Name] = oldTypeArgsByName[innerTp.Name]
							continue
						}
					}
					newRefTypeArgs[tp.Name] = ta
				} else {
					// Missing TypeArg, need to get it from the Type
					newRefTypeArgs[tp.Name] = newTypeArgsByName[tp.Name]
				}
			}

			oldRefTypeArgs := make(map[string]Type)
			for i, tp := range oldDefReference.GetDefinitionMeta().TypeParameters {
				if i < len(oldDefReference.GetDefinitionMeta().TypeArguments) {
					ta := oldDefReference.GetDefinitionMeta().TypeArguments[i]
					if st, ok := ta.(*SimpleType); ok {
						if innerTp, ok := st.ResolvedDefinition.(*GenericTypeParameter); ok {
							oldRefTypeArgs[tp.Name] = newTypeArgsByName[innerTp.Name]
							continue
						}
					}
					oldRefTypeArgs[tp.Name] = ta
				} else {
					// Missing TypeArg, need to get it from the Type
					oldRefTypeArgs[tp.Name] = oldTypeArgsByName[tp.Name]
				}
			}

			// for n, t := range newRefTypeArgs {
			// 	log.Debug().Msgf("  New Ref TypeArg %s: %s", n, TypeToShortSyntax(t, true))
			// }
			// for n, t := range oldRefTypeArgs {
			// 	log.Debug().Msgf("  Old Ref TypeArg %s: %s", n, TypeToShortSyntax(t, true))
			// }

			for _, tp := range newDefReference.GetDefinitionMeta().TypeParameters {
				expectedTypeArg := newRefTypeArgs[tp.Name]
				actualTypeArg := newTypeArgsByName[tp.Name]

				if ch := compareTypes(expectedTypeArg, actualTypeArg, context); ch != nil {
					switch ch.(type) {
					case *TypeChangeIncompatible:
						// log.Error().Msgf("TypeArgs aren't compatible %s <= %s", TypeToShortSyntax(expectedTypeArg, true), TypeToShortSyntax(actualTypeArg, true))
						return &TypeChangeIncompatible{TypePair{oldType, newType}}
					case *TypeChangeDefinitionChanged:
						typeArgDefinitionChanged = true
					}
				}
			}

			for _, tp := range oldDefReference.GetDefinitionMeta().TypeParameters {
				expectedTypeArg := oldRefTypeArgs[tp.Name]
				actualTypeArg := oldTypeArgsByName[tp.Name]

				if ch := compareTypes(expectedTypeArg, actualTypeArg, context); ch != nil {
					switch ch.(type) {
					case *TypeChangeIncompatible:
						// log.Error().Msgf("TypeArgs aren't compatible %s <= %s", TypeToShortSyntax(expectedTypeArg, true), TypeToShortSyntax(actualTypeArg, true))
						return &TypeChangeIncompatible{TypePair{oldType, newType}}
					case *TypeChangeDefinitionChanged:
						typeArgDefinitionChanged = true
					}
				}
			}

		}
	}

	if ch != nil || typeArgDefinitionChanged {
		return &TypeChangeDefinitionChanged{TypePair{oldType, newType}, ch}
	}
	// log.Debug().Msgf("These types are the same: %s <= %s", defWithArgs(newDef), defWithArgs(oldDef))
	return nil
}

func compareOtherSimpleTypes(newType, oldType *SimpleType, context *EvolutionContext) TypeChange {
	newDef := newType.ResolvedDefinition
	oldDef := oldType.ResolvedDefinition

	// Unwind old NamedType to compare underlying Types
	switch oldDef := oldDef.(type) {
	case *NamedType:
		ch := compareTypes(newType, oldDef.Type, context)
		if ch == nil {
			return nil
		}

		// TODO: I think I can test this with a change between two generalized types:
		// 	- OLD: contains a NamedType: Primitive
		// 	- NEW: contains a different Primitive?
		// Then, what happens if I need to "reference" the old parent type? Do I need that missing Alias?

		// switch ch := ch.(type) {
		// case *TypeChangeIncompatible:
		// default:
		// 	tdChange := &NamedTypeChange{DefinitionPair{oldDef, newDef}, ch}
		// 	context.AliasesRemoved = append(context.AliasesRemoved, tdChange)
		// 	// return &TypeChangeDefinitionChanged{TypePair{oldType, newType}, tdChange}
		// }
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
				log.Info().Msgf("all good")
				return &TypeChangeScalarToUnion{TypePair{oldType, newType}, i}
			}
		}
	}

	log.Error().Msgf("Can't do that!")

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
