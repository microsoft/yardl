// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
	"fmt"

	"github.com/microsoft/yardl/tooling/internal/validation"
	"github.com/rs/zerolog/log"
)

const (
	// Annotations referenced in serialization codegen
	VersionAnnotationKey = "version"

	// Annotations used only for validation model evolution (local to this file)
	schemaAnnotationKey = "schema"
)

type ChangeTable map[string]DefinitionChange

func (ct ChangeTable) Add(change DefinitionChange) {
	log.Debug().Msgf("Adding Change for %s", change.LatestDefinition().GetDefinitionMeta().GetQualifiedName())
	if _, ok := ct[change.LatestDefinition().GetDefinitionMeta().GetQualifiedName()]; ok {
		// panic(fmt.Sprintf("Change already exists for %s", change.LatestDefinition().GetDefinitionMeta().GetQualifiedName()))
		log.Error().Msg(fmt.Sprintf("Change already exists for %s", change.LatestDefinition().GetDefinitionMeta().GetQualifiedName()))
	}
	ct[change.LatestDefinition().GetDefinitionMeta().GetQualifiedName()] = change
}

func ValidateEvolution(env *Environment, predecessors []*Environment, versionLabels []string) (*Environment, error) {

	initializeChangeAnnotations(env)

	for i, predecessor := range predecessors {
		log.Info().Msgf("Resolving changes from predecessor %s", versionLabels[i])

		annotatePredecessorSchemas(predecessor)

		definitionChanges, protocolChanges := doSomethingElse(env, predecessor)

		if err := validateChanges(definitionChanges, protocolChanges); err != nil {
			return nil, err
		}

		saveChangedDefinitions(env, definitionChanges, protocolChanges, versionLabels[i])
		continue

		// changeTable := make(ChangeTable)

		// if err := annotateAllChanges(env, predecessor, changeTable, versionLabels[i]); err != nil {
		// 	return nil, err
		// }

		// if err := validateChanges(env, changeTable); err != nil {
		// 	return nil, err
		// }

		// saveChangedDefinitions(env, changeTable, versionLabels[i])
	}

	return env, nil
}

// func validateChanges(env *Environment, changeTable ChangeTable) error {
func validateChanges(definitionChanges []DefinitionChange, protocolChanges []*ProtocolChange) error {
	// Emit User Warnings and aggregate Errors
	errorSink := &validation.ErrorSink{}
	validateTypeDefinitionChanges(definitionChanges, errorSink)
	validateProtocolChanges(protocolChanges, errorSink)
	return errorSink.AsError()
}

// func validateTypeDefinitionChanges(typeDefs []TypeDefinition, changeTable ChangeTable, errorSink *validation.ErrorSink) {
func validateTypeDefinitionChanges(changes []DefinitionChange, errorSink *validation.ErrorSink) {
	for _, ch := range changes {
		if ch == nil {
			panic("I don't want nil DefinitionChanges here ")
		}

		// TODO: Fix these user messages
		td := ch.LatestDefinition()

		switch defChange := ch.(type) {

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

// func validateProtocolChanges(protocols []*ProtocolDefinition, changeTable ChangeTable, errorSink *validation.ErrorSink) {
func validateProtocolChanges(changes []*ProtocolChange, errorSink *validation.ErrorSink) {
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

// Annotate the previous model with Protocol Schema strings for later
func annotatePredecessorSchemas(predecessor *Environment) {
	Visit(predecessor, func(self Visitor, node Node) {
		switch node := node.(type) {
		case *ProtocolDefinition:
			if node.GetDefinitionMeta().Annotations == nil {
				node.GetDefinitionMeta().Annotations = make(map[string]any)
			}
			node.GetDefinitionMeta().Annotations[schemaAnnotationKey] = GetProtocolSchemaString(node, predecessor.SymbolTable)

		default:
			self.VisitChildren(node)
		}
	})
}

// Prepare Annotations on the new model
func initializeChangeAnnotations(env *Environment) {
	Visit(env, func(self Visitor, node Node) {
		switch node := node.(type) {
		case *Namespace:
			node.TypeDefChanges = make(map[string][]DefinitionChange)
			self.VisitChildren(node)

		case *ProtocolDefinition:
			node.Versions = make(map[string]*ProtocolChange)

		default:
			self.VisitChildren(node)
		}
	})
}

// func saveChangedDefinitions(env *Environment, changeTable ChangeTable, versionLabel string) {
func saveChangedDefinitions(env *Environment, definitionChanges []DefinitionChange, protocolChanges []*ProtocolChange, versionLabel string) {
	Visit(env, func(self Visitor, node Node) {
		switch node := node.(type) {
		case *Namespace:
			node.Versions = append(node.Versions, versionLabel)

			for _, ch := range definitionChanges {
				if ch.PreviousDefinition().GetDefinitionMeta().Annotations == nil {
					ch.PreviousDefinition().GetDefinitionMeta().Annotations = make(map[string]any)
				}
				ch.PreviousDefinition().GetDefinitionMeta().Annotations[VersionAnnotationKey] = versionLabel
			}

			// namespaceDefChanges := make([]DefinitionChange, 0)
			// for _, td := range node.TypeDefinitions {
			// 	for _, ch := range definitionChanges {
			// 		if ch.LatestDefinition().GetDefinitionMeta().GetQualifiedName() == td.GetDefinitionMeta().GetQualifiedName() {
			// 			namespaceDefChanges = append(namespaceDefChanges, ch)
			// 		}
			// 	}
			// }
			// node.TypeDefChanges[versionLabel] = namespaceDefChanges
			node.TypeDefChanges[versionLabel] = definitionChanges

			// for _, td := range node.TypeDefinitions {
			// 	if change, ok := changeTable[td.GetDefinitionMeta().GetQualifiedName()]; ok {
			// 		log.Debug().Msgf("CHANGE: %s -> %s", change.PreviousDefinition().GetDefinitionMeta().GetQualifiedName(), change.LatestDefinition().GetDefinitionMeta().GetQualifiedName())
			// 		node.TypeDefChanges[versionLabel] = append(node.TypeDefChanges[versionLabel], change)
			// 		if change.PreviousDefinition().GetDefinitionMeta().Annotations == nil {
			// 			change.PreviousDefinition().GetDefinitionMeta().Annotations = make(map[string]any)
			// 		}
			// 		change.PreviousDefinition().GetDefinitionMeta().Annotations[VersionAnnotationKey] = versionLabel
			// 	}

			// }

			// for _, change := range node.TypeDefChanges[versionLabel] {
			// if change.PreviousDefinition().GetDefinitionMeta().Annotations == nil {
			// 	change.PreviousDefinition().GetDefinitionMeta().Annotations = make(map[string]any)
			// }
			// change.PreviousDefinition().GetDefinitionMeta().Annotations[VersionAnnotationKey] = versionLabel
			// }

			self.VisitChildren(node)

		case *ProtocolDefinition:
			for _, protChange := range protocolChanges {
				if protChange.LatestDefinition().GetDefinitionMeta().GetQualifiedName() == node.GetDefinitionMeta().GetQualifiedName() {
					node.Versions[versionLabel] = protChange
				}
			}
			// var changed *ProtocolChange
			// if ch, ok := changeTable[node.GetQualifiedName()].(*ProtocolChange); ok {
			// 	changed = ch
			// }
			// node.Versions[versionLabel] = changed

		default:
			self.VisitChildren(node)
		}
	})
}

// func annotateAllChanges(newNode, oldNode *Environment, changeTable ChangeTable, versionLabel string) error {
// 	oldNamespaces := make(map[string]*Namespace)
// 	for _, oldNs := range oldNode.Namespaces {
// 		oldNamespaces[oldNs.Name] = oldNs
// 	}

// 	for _, newNs := range newNode.Namespaces {
// 		if oldNs, ok := oldNamespaces[newNs.Name]; ok {
// 			annotateNamespaceChanges(newNs, oldNs, changeTable, versionLabel)
// 		} else {
// 			return fmt.Errorf("Namespace '%s' does not exist in previous version", newNs.Name)
// 		}
// 	}

// 	return nil
// }

func collectUserTypeDefsUsedInProtocols(ns *Namespace) []TypeDefinition {
	isUserTypeDef := make(map[string]bool)
	for _, newTd := range ns.TypeDefinitions {
		isUserTypeDef[newTd.GetDefinitionMeta().GetQualifiedName()] = true
	}
	typeDefCollected := make(map[string]bool)
	var usedTypeDefs []TypeDefinition
	for _, protocol := range ns.Protocols {
		Visit(protocol, func(self Visitor, node Node) {
			switch node := node.(type) {
			case TypeDefinition:
				self.VisitChildren(node)
				name := node.GetDefinitionMeta().GetQualifiedName()
				if isUserTypeDef[name] && !typeDefCollected[name] {
					typeDefCollected[name] = true
					usedTypeDefs = append(usedTypeDefs, node)
				}
			case *SimpleType:
				self.VisitChildren(node)
				self.Visit(node.ResolvedDefinition)
			default:
				self.VisitChildren(node)
			}
		})
	}
	return usedTypeDefs
}

func getResolvedName(td TypeDefinition) string {
	collectNames := func(td TypeDefinition) []string {
		names := make([]string, 0)
		Visit(td, func(self Visitor, node Node) {
			switch node := node.(type) {
			case PrimitiveDefinition:
				return

			case *NamedType:
				names = append(names, node.GetDefinitionMeta().GetQualifiedName())
				self.Visit(node.Type)

			case TypeDefinition:
				names = append(names, node.GetDefinitionMeta().GetQualifiedName())

			case *SimpleType:
				self.Visit(node.ResolvedDefinition)
			}
		})
		return names
	}

	names := collectNames(td)
	return names[len(names)-1]
}

func definitionsResolve(tdA, tdB TypeDefinition) bool {
	collectNames := func(td TypeDefinition) []string {
		names := make([]string, 0)
		Visit(td, func(self Visitor, node Node) {
			switch node := node.(type) {

			case PrimitiveDefinition:
				return
			case *NamedType:
				name := node.GetDefinitionMeta().GetQualifiedName()
				names = append(names, name)
				self.Visit(node.Type)
			case TypeDefinition:
				name := node.GetDefinitionMeta().GetQualifiedName()
				names = append(names, name)

			case *SimpleType:
				self.Visit(node.ResolvedDefinition)

			default:
				// self.VisitChildren(node)
			}
		})
		return names
	}

	namesA := collectNames(tdA)
	namesB := collectNames(tdB)
	// log.Debug().Msgf("%s resolves to %s", tdA.GetDefinitionMeta().GetQualifiedName(), namesA[len(namesA)-1])

	for _, nameA := range namesA {
		for _, nameB := range namesB {
			if nameA == nameB {
				return true
			}
		}
	}
	return false
}

func doSomethingElse(newEnv, oldEnv *Environment) ([]DefinitionChange, []*ProtocolChange) {

	getAllTypeDefs := func(env *Environment) []TypeDefinition {
		allTypeDefs := make([]TypeDefinition, 0)
		for _, ns := range env.Namespaces {
			for _, td := range ns.TypeDefinitions {
				allTypeDefs = append(allTypeDefs, td)
			}
		}
		return allTypeDefs
	}

	// getAllTypeDefsByName := func(env *Environment) map[string]TypeDefinition {
	// 	allTypeDefs := make(map[string]TypeDefinition)
	// 	for _, ns := range env.Namespaces {
	// 		for _, td := range ns.TypeDefinitions {
	// 			allTypeDefs[td.GetDefinitionMeta().GetQualifiedName()] = td
	// 		}
	// 	}
	// 	return allTypeDefs
	// }

	// getResolvedNames := func(allTypeDefs []TypeDefinition) map[string]string {
	// 	resolvedNames := make(map[string]string)
	// 	for _, td := range allTypeDefs {
	// 		resolvedNames[td.GetDefinitionMeta().GetQualifiedName()] = getResolvedName(td)
	// 	}
	// 	return resolvedNames
	// }

	getNameEquivalence := func(allTypeDefs []TypeDefinition) map[string]map[string]bool {
		// allTypeDefs := getAllTypeDefs(env)
		resolvesTo := make(map[string]map[string]bool)
		for _, tdA := range allTypeDefs {
			table := make(map[string]bool)
			for _, tdB := range allTypeDefs {
				// if definitionsResolve(tdA, tdB) {
				if getResolvedName(tdA) == getResolvedName(tdB) {
					table[tdB.GetDefinitionMeta().GetQualifiedName()] = true
				}
			}
			resolvesTo[tdA.GetDefinitionMeta().GetQualifiedName()] = table
		}
		return resolvesTo
	}

	allNewTypeDefs := getAllTypeDefs(newEnv)
	// allNewTypeDefsByName := getAllTypeDefsByName(newEnv)
	newNameMapping := getNameEquivalence(allNewTypeDefs)
	// newNameResolution := getResolvedNames(allNewTypeDefs)

	// For Debugging:
	// for _, newTd := range allNewTypeDefs {
	// 	for name := range newNameMapping[newTd.GetDefinitionMeta().GetQualifiedName()] {
	// 		log.Debug().Msgf("[New] %s resolves to %s", newTd.GetDefinitionMeta().GetQualifiedName(), name)
	// 	}
	// }

	allOldTypeDefs := getAllTypeDefs(oldEnv)
	// allOldTypeDefsByName := getAllTypeDefsByName(oldEnv)
	oldNameMapping := getNameEquivalence(allOldTypeDefs)
	// oldNameResolution := getResolvedNames(allOldTypeDefs)

	// For Debugging:
	// for _, oldTd := range allOldTypeDefs {
	// 	for name := range oldNameMapping[oldTd.GetDefinitionMeta().GetQualifiedName()] {
	// 		log.Debug().Msgf("[Old] %s resolves to %s", oldTd.GetDefinitionMeta().GetQualifiedName(), name)
	// 	}
	// }

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
						// log.Debug().Msgf("New %s == Old %s", n, o)
					}
				}
			}
		}
	}

	// For Debugging:
	// for _, newTd := range allNewTypeDefs {
	// 	for oldName := range semanticallyEqual[newTd.GetDefinitionMeta().GetQualifiedName()] {
	// 		log.Debug().Msgf("%s semantically equals %s", newTd.GetDefinitionMeta().GetQualifiedName(), oldName)
	// 	}
	// }

	// typeDefKind := func(td TypeDefinition) string {
	// 	switch td.(type) {
	// 	case *RecordDefinition:
	// 		return "Record"
	// 	case *NamedType:
	// 		return "Alias"
	// 	case *EnumDefinition:
	// 		return "Enum"
	// 	default:
	// 		panic("Shouldn't get here")
	// 	}
	// }

	compared := make(map[string]map[string]bool)
	changes := make(map[string]map[string]DefinitionChange)

	context := &EvolutionContext{
		SemanticMapping: semanticallyEqual,
		Compared:        compared,
		Changes:         changes,
	}

	for _, newTd := range allNewTypeDefs {
		newName := newTd.GetDefinitionMeta().GetQualifiedName()
		compared[newName] = make(map[string]bool)
		changes[newName] = make(map[string]DefinitionChange)

		for _, oldTd := range allOldTypeDefs {
			oldName := oldTd.GetDefinitionMeta().GetQualifiedName()
			if semanticallyEqual[newName][oldName] {

				// newResolvedName := newNameResolution[newName]
				// oldResolvedName := oldNameResolution[oldName]
				// if compared[newResolvedName][oldResolvedName] {
				// 	log.Debug().Msgf("Skip %s <= %s, already compared %s <= %s", newName, oldName, newResolvedName, oldResolvedName)
				// } else {
				// 	log.Debug().Msgf("Comparing new %s %s with old %s %s", typeDefKind(newTd), newName, typeDefKind(oldTd), oldName)

				// 	if ch := compareTypeDefinitions(newTd, oldTd, context); ch != nil {
				// 		changes[newName][oldName] = ch
				// 	}
				// }

				if ch := compareTypeDefinitions(newTd, oldTd, context); ch != nil {
					changes[newName][oldName] = ch
				}
				compared[newName][oldName] = true
			}
		}
	}

	// Protocols may be reordered, added, or removed
	// We only care about pre-existing Protocols that CHANGED
	oldProts := make(map[string]*ProtocolDefinition)
	for _, ns := range oldEnv.Namespaces {
		for _, oldProt := range ns.Protocols {
			oldProts[oldProt.GetQualifiedName()] = oldProt
		}
	}

	allProtocolChanges := make([]*ProtocolChange, 0)
	for _, ns := range newEnv.Namespaces {
		for _, newProt := range ns.Protocols {
			oldProt, ok := oldProts[newProt.GetDefinitionMeta().GetQualifiedName()]
			if !ok {
				// Skip new ProtocolDefinition
				continue
			}

			// Annotate this ProtocolDefinition with any changes from previous version.
			protocolChange := detectProtocolDefinitionChanges(newProt, oldProt, context)
			// changeTable[newProt.GetQualifiedName()] = protocolChange
			if protocolChange != nil {
				// log.Debug().Msgf("Protocol %s changed", newProt.GetQualifiedName())
				// changeTable.Add(protocolChange)
				allProtocolChanges = append(allProtocolChanges, protocolChange)
			}
		}
	}

	// Now we're finished detecting all changes between models... Time to collect/filter TypeDefinition changes
	allDefinitionChanges := make([]DefinitionChange, 0)
	// For Debugging:
	for _, newTd := range allNewTypeDefs {
		newName := newTd.GetDefinitionMeta().GetQualifiedName()
		for _, oldTd := range allOldTypeDefs {
			oldName := oldTd.GetDefinitionMeta().GetQualifiedName()

			if ch, ok := changes[newName][oldName]; ok {
				if ch == nil {
					panic("I don't want nil DefinitionChanges in my context")
				}
				// newChangeName := ch.LatestDefinition().GetDefinitionMeta().GetQualifiedName()
				// oldChangeName := ch.PreviousDefinition().GetDefinitionMeta().GetQualifiedName()
				// log.Debug().Msgf("Change for %s <= %s:  %s <= %s", newName, oldName, newChangeName, oldChangeName)

				allDefinitionChanges = append(allDefinitionChanges, ch)
			}
		}
	}
	for _, ch := range context.AliasesRemoved {
		// if ch.LatestDefinition().GetDefinitionMeta().GetQualifiedName() == newName {
		log.Debug().Msgf("Appending AliasRemoved %s => %s", ch.PreviousDefinition().GetDefinitionMeta().GetQualifiedName(), ch.LatestDefinition().GetDefinitionMeta().GetQualifiedName())
		allDefinitionChanges = append(allDefinitionChanges, ch)
		// }
	}

	neededDefinitionChanges := make([]DefinitionChange, 0)
	uniqueOldNames := make(map[string]bool)
	for _, change := range allDefinitionChanges {
		// newName := change.LatestDefinition().GetDefinitionMeta().GetQualifiedName()
		oldName := change.PreviousDefinition().GetDefinitionMeta().GetQualifiedName()

		if uniqueOldNames[oldName] {
			continue
		}

		// newResolvesTo := newNameResolution[newName]
		// oldResolvesTo := oldNameResolution[oldName]
		// log.Debug().Msgf("OLD %s resolve to %s | NEW %s resolves to %s", oldName, oldResolvesTo, newName, newResolvesTo)

		// if newName == newResolvesTo && oldName == oldResolvesTo {
		// neededDefinitionChanges = append(neededDefinitionChanges, change)
		// log.Debug().Msgf("Change for %s <= %s", newName, oldName)
		// }
		neededDefinitionChanges = append(neededDefinitionChanges, change)
		uniqueOldNames[oldName] = true
		// log.Debug().Msgf("Change for %s <= %s", newName, oldName)
	}

	// Now we need to "rename" all old TypeDefinitions to their new names...
	// for _, change := range allDefinitionChanges {
	// 	oldName := change.PreviousDefinition().GetDefinitionMeta().Name
	// 	oldQualifiedName := change.PreviousDefinition().GetDefinitionMeta().GetQualifiedName()

	// 	Visit(change.PreviousDefinition(), 	)
	// }

	return neededDefinitionChanges, allProtocolChanges
}

type EvolutionContext struct {
	SemanticMapping map[string]map[string]bool
	Compared        map[string]map[string]bool
	Changes         map[string]map[string]DefinitionChange
	AliasesRemoved  []DefinitionChange
}

func compareTypeDefinitions(newTd, oldTd TypeDefinition, context *EvolutionContext) DefinitionChange {

	switch newTd := newTd.(type) {
	case *NamedType:
		switch oldTd := oldTd.(type) {
		case *NamedType:
			ch := compareTypes(newTd.Type, oldTd.Type, context)
			if ch == nil {
				return nil
			}
			return &NamedTypeChange{DefinitionPair{oldTd, newTd}, ch}

		case *RecordDefinition:
			oldType := &SimpleType{NodeMeta: *oldTd.GetNodeMeta(), Name: oldTd.GetDefinitionMeta().GetQualifiedName(), ResolvedDefinition: oldTd}
			ch := compareTypes(newTd.Type, oldType, context)
			if ch != nil {
				return &NamedTypeChange{DefinitionPair{oldTd, newTd}, ch}
			}
			return nil

		case *EnumDefinition:
			oldType := &SimpleType{NodeMeta: *oldTd.GetNodeMeta(), Name: oldTd.GetDefinitionMeta().GetQualifiedName(), ResolvedDefinition: oldTd}
			ch := compareTypes(newTd.Type, oldType, context)
			if ch != nil {
				return &NamedTypeChange{DefinitionPair{oldTd, newTd}, ch}
			}
			return nil

		default:
			return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
		}

	case *RecordDefinition:
		switch oldTd := oldTd.(type) {
		case *RecordDefinition:
			ch := detectRecordDefinitionChanges(newTd, oldTd, context)
			if ch == nil {
				return nil
			}
			return ch

		case *NamedType:
			newType := &SimpleType{NodeMeta: *newTd.GetNodeMeta(), Name: newTd.GetDefinitionMeta().GetQualifiedName(), ResolvedDefinition: newTd}
			ch := compareTypes(newType, oldTd.Type, context)
			if ch != nil {
				return &NamedTypeChange{DefinitionPair{oldTd, newTd}, ch}
			}
			return nil

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
			newType := &SimpleType{NodeMeta: *newTd.GetNodeMeta(), Name: newTd.GetDefinitionMeta().GetQualifiedName(), ResolvedDefinition: newTd}
			ch := compareTypes(newType, oldTd.Type, context)
			if ch != nil {
				return &NamedTypeChange{DefinitionPair{oldTd, newTd}, ch}
			}
			return nil
		default:
			return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
		}

	default:
		panic("Expected a TypeDefinition...")
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
			return compareGeneralizedToSimpleTypes(newType, oldType, context)
		default:
			panic("Shouldn't get here")
		}

	case *GeneralizedType:
		switch oldType := oldType.(type) {
		case *GeneralizedType:
			return compareGeneralizedTypes(newType, oldType, context)
		case *SimpleType:
			return compareSimpleToGeneralizedTypes(newType, oldType, context)
		default:
			panic("Shouldn't get here")
		}

	default:
		panic("Expected a type")
	}
}

func compareSimpleTypes(newType, oldType *SimpleType, context *EvolutionContext) TypeChange {
	// TODO: Handle TypeArgs
	//

	// Both newType and oldType are SimpleTypes
	// Thus, the possible type changes here are:
	//  - Primitive to Primitive (possibly valid)
	//  - TypeDefinition to TypeDefinition (possibly valid)
	//  - Primitive to TypeDefinition (invalid)
	//  - TypeDefinition to Primitive (invalid)

	newDef := newType.ResolvedDefinition
	oldDef := oldType.ResolvedDefinition

	// Here is where we do a lookup on the two ResolvedDefinitions
	// 1. Are they semantically equal?
	// 1. Did we already compare them?

	if !context.SemanticMapping[newDef.GetDefinitionMeta().GetQualifiedName()][oldDef.GetDefinitionMeta().GetQualifiedName()] {
		// log.Debug().Msgf("SimpleType %s != %s", newDef.GetDefinitionMeta().GetQualifiedName(), oldDef.GetDefinitionMeta().GetQualifiedName())
		switch oldDef := oldDef.(type) {
		case *NamedType:
			ch := compareTypes(newType, oldDef.Type, context)
			if ch == nil {
				return nil
			}
			tdChange := &NamedTypeChange{DefinitionPair{oldDef, newDef}, ch}
			context.AliasesRemoved = append(context.AliasesRemoved, tdChange)
			log.Debug().Msgf("Alias removed %s <= %s", newDef.GetDefinitionMeta().GetQualifiedName(), oldDef.GetDefinitionMeta().GetQualifiedName())
			// return ch
			return &TypeChangeDefinitionChanged{TypePair{oldType, newType}, tdChange}
		}
		switch newDef := newDef.(type) {
		case *NamedType:
			ch := compareTypes(newDef.Type, oldType, context)
			if ch == nil {
				return nil
			}
			log.Debug().Msgf("Alias added %s <= %s", newDef.GetDefinitionMeta().GetQualifiedName(), oldDef.GetDefinitionMeta().GetQualifiedName())
			return ch
		}

		if _, ok := oldDef.(PrimitiveDefinition); ok {
			if _, ok := newDef.(PrimitiveDefinition); ok {
				return detectPrimitiveTypeChange(newType, oldType)
			}
			return &TypeChangeIncompatible{TypePair{oldType, newType}}
		}

		if _, ok := newDef.(PrimitiveDefinition); ok {
			return &TypeChangeIncompatible{TypePair{oldType, newType}}
		}

		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	if !context.Compared[newDef.GetDefinitionMeta().GetQualifiedName()][oldDef.GetDefinitionMeta().GetQualifiedName()] {
		panic(fmt.Sprintf("Why haven't we compared new %s with old %s", newDef.GetDefinitionMeta().GetQualifiedName(), oldDef.GetDefinitionMeta().GetQualifiedName()))
	}

	if ch, ok := context.Changes[newDef.GetDefinitionMeta().GetQualifiedName()][oldDef.GetDefinitionMeta().GetQualifiedName()]; ok {
		if ch == nil {
			panic("I don't want nil DefinitionChanges in my context")
		}

		if _, ok := ch.(*DefinitionChangeIncompatible); ok {
			return &TypeChangeIncompatible{TypePair{oldType, newType}}
		} else {
			return &TypeChangeDefinitionChanged{TypePair{oldType, newType}, ch}
		}
	}

	return nil
}

// func annotateNamespaceChanges(newNs, oldNs *Namespace, changeTable ChangeTable, versionLabel string) {
// 	// TypeDefinitions may be reordered, added, or removed, so we compare them by name
// 	oldTds := make(map[string]TypeDefinition)
// 	for _, oldTd := range oldNs.TypeDefinitions {
// 		oldTds[oldTd.GetDefinitionMeta().GetQualifiedName()] = oldTd
// 	}

// 	// newUsedTypeDefs := collectUserTypeDefsUsedInProtocols(newNs)
// 	newTypeDefs := newNs.TypeDefinitions

// 	// typeDefChanges := make([]DefinitionChange, 0)
// 	alreadyCompared := make(map[string]bool)
// 	for _, newTd := range newTypeDefs {
// 		oldTd, ok := oldTds[newTd.GetDefinitionMeta().GetQualifiedName()]
// 		if !ok {
// 			// Skip new TypeDefinition
// 			continue
// 		}

// 		// type NamedTypeUnwinder = func(TypeDefinition) TypeDefinition
// 		// removedAliases := make([]DefinitionChange, 0)
// 		// var unwindOldAlias, unwindNewAlias NamedTypeUnwinder

// 		// unwindOldAlias = func(oldTd TypeDefinition) TypeDefinition {
// 		// 	switch old := oldTd.(type) {
// 		// 	case *NamedType:
// 		// 		if _, isNamedType := newTd.(*NamedType); !isNamedType {
// 		// 			// Alias removed and we need to generate its compatibility serializers.
// 		// 			if oldType, ok := old.Type.(*SimpleType); ok {
// 		// 				compat := &NamedTypeChange{DefinitionPair{oldTd, newTd}, nil}
// 		// 				removedAliases = append([]DefinitionChange{compat}, removedAliases...)
// 		// 				oldTd = oldType.ResolvedDefinition
// 		// 				return unwindOldAlias(oldTd)
// 		// 			}
// 		// 		}
// 		// 	}
// 		// 	return oldTd
// 		// }

// 		// unwindNewAlias = func(newTd TypeDefinition) TypeDefinition {
// 		// 	switch new := newTd.(type) {
// 		// 	case *NamedType:
// 		// 		if _, isNamedType := oldTd.(*NamedType); !isNamedType {
// 		// 			// Alias added and we can safely ignore it.
// 		// 			if newType, ok := new.Type.(*SimpleType); ok {
// 		// 				newTd = newType.ResolvedDefinition
// 		// 				return unwindNewAlias(newTd)
// 		// 			}
// 		// 		}
// 		// 	}
// 		// 	return newTd
// 		// }

// 		// // TODO: The unwind helpers might modify oldTd/newTd out-of-order. Make them NOT capture locals

// 		// // "Unwind" any NamedTypes so we only compare underlying TypeDefinitions
// 		// oldTd = unwindOldAlias(oldTd)
// 		// newTd = unwindNewAlias(newTd)

// 		if alreadyCompared[newTd.GetDefinitionMeta().GetQualifiedName()] {
// 			// TODO: Remove this check if not needed, once integration tests are "complete"
// 			panic(fmt.Sprintf("Already Compared %s", newTd.GetDefinitionMeta().GetQualifiedName()))
// 			continue
// 		}

// 		// log.Debug().Msgf("Comparing TypeDefinitions %s and %s", newTd.GetDefinitionMeta().GetQualifiedName(), oldTd.GetDefinitionMeta().GetQualifiedName())
// 		defChange := detectTypeDefinitionChanges(newTd, oldTd, changeTable, true)
// 		if defChange != nil {
// 			// typeDefChanges = append(typeDefChanges, defChange)
// 			// typeDefChanges = append(typeDefChanges, removedAliases...)
// 			// for _, ch := range removedAliases {
// 			// 	changeTable.Add(ch)
// 			// }

// 			// Save this DefinitionChange so that, later, detectSimpleTypeChanges can determine if an underlying TypeDefinition changed
// 			changeTable.Add(defChange)
// 		}

// 		alreadyCompared[newTd.GetDefinitionMeta().GetQualifiedName()] = true
// 	}

// 	// Save all TypeDefinition changes for generating of compatibility serializers
// 	// newNs.TypeDefChanges[versionLabel] = typeDefChanges

// 	// Protocols may be reordered, added, or removed
// 	// We only care about pre-existing Protocols that CHANGED
// 	// oldProts := make(map[string]*ProtocolDefinition)
// 	// for _, oldProt := range oldNs.Protocols {
// 	// 	oldProts[oldProt.GetQualifiedName()] = oldProt
// 	// }

// 	// for _, newProt := range newNs.Protocols {
// 	// 	oldProt, ok := oldProts[newProt.GetDefinitionMeta().GetQualifiedName()]
// 	// 	if !ok {
// 	// 		// Skip new ProtocolDefinition
// 	// 		continue
// 	// 	}

// 	// 	// Annotate this ProtocolDefinition with any changes from previous version.
// 	// 	protocolChange := detectProtocolDefinitionChanges(newProt, oldProt, changeTable)
// 	// 	// changeTable[newProt.GetQualifiedName()] = protocolChange
// 	// 	if protocolChange != nil {
// 	// 		changeTable.Add(protocolChange)
// 	// 	}
// 	// }
// }

// Compares two TypeDefinitions with matching names
// func detectTypeDefinitionChanges(newTd, oldTd TypeDefinition, changeTable ChangeTable, isTopLevel bool) DefinitionChange {

// 	log.Debug().Msgf("Comparing TypeDefinitions %s and %s", newTd.GetDefinitionMeta().GetQualifiedName(), oldTd.GetDefinitionMeta().GetQualifiedName())

// 	// if ch, ok := changeTable[newTd.GetDefinitionMeta().GetQualifiedName()]; ok {
// 	// 	// We already compared these two TypeDefinitions and
// 	// 	if ch == nil {
// 	// 		// These two TypeDefinitions did not change between versions
// 	// 		return nil
// 	// 	}

// 	// 	if ch.PreviousDefinition().GetDefinitionMeta().GetQualifiedName() == oldTd.GetDefinitionMeta().GetQualifiedName() {
// 	// 		// This TypeDefinition changed from the previous version
// 	// 		return ch
// 	// 	}
// 	// }

// 	if oldNT, oldIsNamedType := oldTd.(*NamedType); oldIsNamedType {
// 		if _, newIsNamedType := newTd.(*NamedType); !newIsNamedType {
// 			// Alias removed in new version
// 			// Recurse into Alias
// 			newType := &SimpleType{NodeMeta: *newTd.GetNodeMeta(), Name: newTd.GetDefinitionMeta().GetQualifiedName(), ResolvedDefinition: newTd}
// 			ch := detectTypeChanges(newType, oldNT.Type, changeTable)
// 			if ch == nil {
// 				return nil
// 			}
// 			if _, ok := ch.(*TypeChangeIncompatible); ok {
// 				return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 			}
// 			return &NamedTypeChange{DefinitionPair{oldTd, newTd}, ch}
// 		}
// 	}

// 	if newNt, newIsNamedType := newTd.(*NamedType); newIsNamedType {
// 		if _, oldIsNamedType := oldTd.(*NamedType); !oldIsNamedType {
// 			// Alias added in new version
// 			// Recurse into Alias
// 			oldType := &SimpleType{NodeMeta: *oldTd.GetNodeMeta(), Name: oldTd.GetDefinitionMeta().GetQualifiedName(), ResolvedDefinition: oldTd}
// 			ch := detectTypeChanges(newNt.Type, oldType, changeTable)
// 			if ch == nil {
// 				return nil
// 			}
// 			if _, ok := ch.(*TypeChangeIncompatible); ok {
// 				return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 			}
// 			return &NamedTypeChange{DefinitionPair{oldTd, newTd}, ch}
// 		}
// 	}

// 	if newNt, newIsNamedType := newTd.(*NamedType); newIsNamedType {
// 		if oldNt, oldIsNamedType := oldTd.(*NamedType); oldIsNamedType {
// 			if typeChange := detectTypeChanges(newNt.Type, oldNt.Type, changeTable); typeChange != nil {
// 				return &NamedTypeChange{DefinitionPair{oldTd, newTd}, typeChange}
// 			}
// 			return nil
// 		}
// 	}

// 	if oldPrim, ok := oldTd.(PrimitiveDefinition); ok {
// 		if newPrim, ok := newTd.(PrimitiveDefinition); ok {
// 			return detectPrimitiveChange(newPrim, oldPrim)
// 		}
// 		return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 	}

// 	if _, ok := newTd.(PrimitiveDefinition); ok {
// 		return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 	}

// 	// if !isTopLevel {
// 	// 	return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 	// }

// 	// Haven't compared these two TypeDefinitions yet, so let's do it:
// 	switch newNode := newTd.(type) {
// 	case *RecordDefinition:
// 		switch oldTd := oldTd.(type) {
// 		case *RecordDefinition:
// 			// if ch := detectRecordDefinitionChanges(newNode, oldTd, changeTable); ch != nil {
// 			// 	return ch
// 			// }
// 			// return nil

// 		case *NamedType:
// 			panic("New TD is a RecordDefinition, Old TD should not be a NamedType")
// 			// if st, ok := oldTd.Type.(*SimpleType); ok {
// 			// 	// Alias removed in new version
// 			// 	if ch := detectTypeDefinitionChanges(newNode, st.ResolvedDefinition, changeTable); ch != nil {
// 			// 		changeTable.Add(ch)
// 			// 		return &NamedTypeChange{DefinitionPair{oldTd, newTd}, nil}
// 			// 	}
// 			// 	return nil
// 			// }
// 			// return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}

// 		default:
// 			// Changing a non-Record to a Record is not backward compatible
// 			return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 		}

// 	case *NamedType:
// 		panic("Shouldn't even get here")
// 		// switch oldTd := oldTd.(type) {
// 		// case *NamedType:
// 		// 	if typeChange := detectTypeChanges(newNode.Type, oldTd.Type, changeTable); typeChange != nil {
// 		// 		return &NamedTypeChange{DefinitionPair{oldTd, newTd}, typeChange}
// 		// 	}
// 		// 	return nil

// 		// default:
// 		// 	panic("New TD is a NamedType, Old TD should only be a NamedType")
// 		// if st, ok := newNode.Type.(*SimpleType); ok {
// 		// 	// Alias added in new version
// 		// 	if ch := detectTypeDefinitionChanges(st.ResolvedDefinition, oldTd, changeTable); ch != nil {
// 		// 		changeTable.Add(ch)
// 		// 		return &NamedTypeChange{DefinitionPair{oldTd, newTd}, nil}
// 		// 	}
// 		// 	return nil
// 		// }
// 		// return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 		// }

// 	case *EnumDefinition:
// 		switch oldTd := oldTd.(type) {
// 		case *EnumDefinition:
// 			if ch := detectEnumDefinitionChanges(newNode, oldTd, changeTable); ch != nil {
// 				return ch
// 			}
// 			return nil
// 		case *NamedType:
// 			panic("New TD is an EnumDefinition, Old TD should not be a NamedType")
// 			// if st, ok := oldTd.Type.(*SimpleType); ok {
// 			// 	if ch := detectTypeDefinitionChanges(newNode, st.ResolvedDefinition, changeTable); ch != nil {
// 			// 		changeTable.Add(ch)
// 			// 		return &NamedTypeChange{DefinitionPair{oldTd, newTd}, nil}
// 			// 	}
// 			// 	return nil
// 			// }
// 			// return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}

// 		default:
// 			return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 		}

// 	case PrimitiveDefinition:
// 		panic("Should have already handled PrimitiveDefinition")
// 		// switch oldTd := oldTd.(type) {
// 		// case PrimitiveDefinition:
// 		// 	// TODO
// 		// 	return nil
// 		// case *NamedType:
// 		// 	panic("New TD is a PrimitiveDefinition, Old TD should not be a NamedType")
// 		// 	if st, ok := oldTd.Type.(*SimpleType); ok {
// 		// 		if ch := detectTypeDefinitionChanges(newNode, st.ResolvedDefinition, changeTable); ch != nil {
// 		// 			return &NamedTypeChange{DefinitionPair{oldTd, newTd}, nil}
// 		// 		}
// 		// 		return nil
// 		// 	}
// 		// 	return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}

// 		// default:
// 		// 	return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
// 		// }

// 	default:
// 		// log.Debug().Msgf("What is this? %s was %s", newNode, oldTd)
// 		panic("Expected a TypeDefinition")
// 	}

// 	return nil
// }

// Compares two ProtocolDefinitions with matching names
// func detectProtocolDefinitionChanges(newProtocol, oldProtocol *ProtocolDefinition, changeTable ChangeTable) *ProtocolChange {
func detectProtocolDefinitionChanges(newProtocol, oldProtocol *ProtocolDefinition, context *EvolutionContext) *ProtocolChange {
	change := &ProtocolChange{
		DefinitionPair: DefinitionPair{oldProtocol, newProtocol},
		PreviousSchema: oldProtocol.GetDefinitionMeta().Annotations[schemaAnnotationKey].(string),
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

		// if typeChange := detectTypeChanges(newStep.Type, oldStep.Type, changeTable); typeChange != nil {
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
// func detectRecordDefinitionChanges(newRecord, oldRecord *RecordDefinition, changeTable ChangeTable) *RecordChange {
func detectRecordDefinitionChanges(newRecord, oldRecord *RecordDefinition, context *EvolutionContext) *RecordChange {
	change := &RecordChange{
		DefinitionPair: DefinitionPair{oldRecord, newRecord},
		FieldRemoved:   make([]bool, len(oldRecord.Fields)),
		FieldChanges:   make([]TypeChange, len(oldRecord.Fields)),
		NewFieldIndex:  make([]int, len(oldRecord.Fields)),
	}

	// Fields may be reordered
	// If they are, we want result to represent the old Record, for Serialization compatibility
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

		// if typeChange := detectTypeChanges(newField.Type, oldField.Type, changeTable); typeChange != nil {
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

// func detectEnumDefinitionChanges(newNode, oldEnum *EnumDefinition, changeTable ChangeTable) DefinitionChange {
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

	// if ch := detectTypeChanges(newBaseType, oldBaseType, changeTable); ch != nil {
	if ch := compareTypes(newBaseType, oldBaseType, context); ch != nil {
		// CHANGE: Changed Enum base type
		return &EnumChange{DefinitionPair{oldEnum, newNode}, ch}
	}

	return nil
}

// Compares two Types to determine if and how they changed
// func detectTypeChanges(newType, oldType Type, changeTable ChangeTable) TypeChange {
// 	// newType = GetUnderlyingType(newType)
// 	// oldType = GetUnderlyingType(oldType)

// 	switch newType := newType.(type) {

// 	case *SimpleType:
// 		switch oldType := oldType.(type) {
// 		case *SimpleType:
// 			return detectSimpleTypeChanges(newType, oldType, changeTable)
// 		case *GeneralizedType:
// 			return detectGeneralizedToSimpleTypeChanges(newType, oldType, changeTable)
// 		default:
// 			panic("Shouldn't get here")
// 		}

// 	case *GeneralizedType:
// 		switch oldType := oldType.(type) {
// 		case *GeneralizedType:
// 			return detectGeneralizedTypeChanges(newType, oldType, changeTable)
// 		case *SimpleType:
// 			return detectSimpleToGeneralizedTypeChanges(newType, oldType, changeTable)
// 		default:
// 			panic("Shouldn't get here")
// 		}

// 	default:
// 		panic("Expected a type")
// 	}
// }

// func detectSimpleTypeChanges(newType, oldType *SimpleType, changeTable ChangeTable) TypeChange {
// 	// TODO: Compare TypeArguments
// 	// This comparison depends on whether the ResolvedDefinition changed!
// 	if len(newType.TypeArguments) != len(oldType.TypeArguments) {
// 		// CHANGE: Changed number of TypeArguments
// 		return &TypeChangeIncompatible{TypePair{oldType, newType}}
// 	} else {
// 		for i := range newType.TypeArguments {
// 			if ch := detectTypeChanges(newType.TypeArguments[i], oldType.TypeArguments[i], changeTable); ch != nil {
// 				// CHANGE: Changed TypeArgument
// 				// TODO: Returning early skips other possible changes to the Type
// 				return ch
// 			}
// 		}
// 	}

// 	// Both newType and oldType are SimpleTypes
// 	// Thus, the possible type changes here are:
// 	//  - Primitive to Primitive (possibly valid)
// 	//  - TypeDefinition to TypeDefinition (possibly valid)
// 	//  - Primitive to TypeDefinition (invalid)
// 	//  - TypeDefinition to Primitive (invalid)

// 	newDef := newType.ResolvedDefinition
// 	oldDef := oldType.ResolvedDefinition

// 	// if _, ok := oldDef.(PrimitiveDefinition); ok {
// 	// 	if _, ok := newDef.(PrimitiveDefinition); ok {
// 	// 		return detectPrimitiveTypeChange(newType, oldType)
// 	// 	}
// 	// 	return &TypeChangeIncompatible{TypePair{oldType, newType}}
// 	// }

// 	// if _, ok := newDef.(PrimitiveDefinition); ok {
// 	// 	return &TypeChangeIncompatible{TypePair{oldType, newType}}
// 	// }

// 	// if ch, ok := changeTable[newDef.GetDefinitionMeta().GetQualifiedName()]; ok {
// 	// 	if ch != nil {
// 	// 		if ch.PreviousDefinition().GetDefinitionMeta().GetQualifiedName() == oldDef.GetDefinitionMeta().GetQualifiedName() {
// 	// 			return &TypeChangeDefinitionChanged{TypePair{oldType, newType}}
// 	// 		} else {
// 	// 			panic(fmt.Sprintf("Found change for %s but it doesn't match %s", newDef.GetDefinitionMeta().GetQualifiedName(), oldDef.GetDefinitionMeta().GetQualifiedName()))
// 	// 		}
// 	// 	}
// 	// }

// 	// if newDef.GetDefinitionMeta().GetQualifiedName() != oldDef.GetDefinitionMeta().GetQualifiedName() {
// 	// 	// CHANGE: Not the same underlying TypeDefinition
// 	// 	return &TypeChangeIncompatible{TypePair{oldType, newType}}
// 	// }

// 	if tdChange := detectTypeDefinitionChanges(newDef, oldDef, changeTable, false); tdChange != nil {
// 		if _, ok := tdChange.(*DefinitionChangeIncompatible); ok {
// 			return &TypeChangeIncompatible{TypePair{oldType, newType}}
// 		}
// 		return &TypeChangeDefinitionChanged{TypePair{oldType, newType}, tdChange}
// 	}
// 	return nil
// }

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

func detectPrimitiveChange(newPrimitive, oldPrimitive PrimitiveDefinition) DefinitionChange {
	if newPrimitive == oldPrimitive {
		return nil
	}

	// CHANGE: Changed Primitive type
	if oldPrimitive == PrimitiveString {
		if GetPrimitiveKind(newPrimitive) == PrimitiveKindInteger || GetPrimitiveKind(newPrimitive) == PrimitiveKindFloatingPoint {
			return &PrimitiveChangeStringToNumber{DefinitionPair{oldPrimitive, newPrimitive}}
		}
	}

	if GetPrimitiveKind(oldPrimitive) == PrimitiveKindInteger || GetPrimitiveKind(oldPrimitive) == PrimitiveKindFloatingPoint {
		if newPrimitive == PrimitiveString {
			return &PrimitiveChangeNumberToString{DefinitionPair{oldPrimitive, newPrimitive}}
		}

		if GetPrimitiveKind(newPrimitive) == PrimitiveKindInteger || GetPrimitiveKind(newPrimitive) == PrimitiveKindFloatingPoint {
			return &PrimitiveChangeNumberToNumber{DefinitionPair{oldPrimitive, newPrimitive}}
		}
	}

	return &DefinitionChangeIncompatible{DefinitionPair{oldPrimitive, newPrimitive}}
}

// func detectGeneralizedToSimpleTypeChanges(newType *SimpleType, oldType *GeneralizedType, changeTable ChangeTable) TypeChange {
func compareGeneralizedToSimpleTypes(newType *SimpleType, oldType *GeneralizedType, context *EvolutionContext) TypeChange {
	// Is it a change from Optional<T> to T (partially compatible)
	if oldType.Cases.IsOptional() {
		// switch detectTypeChanges(newType, oldType.Cases[1].Type, changeTable).(type) {
		switch compareTypes(newType, oldType.Cases[1].Type, context).(type) {
		case nil, *TypeChangeDefinitionChanged:
			return &TypeChangeOptionalToScalar{TypePair{oldType, newType}}
		}
	}

	// Is it a change from Union<T, ...> to T (partially compatible)
	if oldType.Cases.IsUnion() {
		for i, tc := range oldType.Cases {
			// switch detectTypeChanges(newType, tc.Type, changeTable).(type) {
			switch compareTypes(newType, tc.Type, context).(type) {
			case nil, *TypeChangeDefinitionChanged:
				return &TypeChangeUnionToScalar{TypePair{oldType, newType}, i}
			}
		}
	}

	// CHANGE: Incompatible change from Generalized to Simple
	return &TypeChangeIncompatible{TypePair{oldType, newType}}
}

// func detectSimpleToGeneralizedTypeChanges(newType *GeneralizedType, oldType *SimpleType, changeTable ChangeTable) TypeChange {
func compareSimpleToGeneralizedTypes(newType *GeneralizedType, oldType *SimpleType, context *EvolutionContext) TypeChange {
	// Is it a change from T to Optional<T> (partially compatible)
	if newType.Cases.IsOptional() {
		// switch detectTypeChanges(newType.Cases[1].Type, oldType, changeTable).(type) {
		switch compareTypes(newType.Cases[1].Type, oldType, context).(type) {
		case nil, *TypeChangeDefinitionChanged:
			return &TypeChangeScalarToOptional{TypePair{oldType, newType}}
		}
	}

	// Is it a change from T to Union<T, ...> (partially compatible)
	if newType.Cases.IsUnion() {
		for i, tc := range newType.Cases {
			// switch detectTypeChanges(tc.Type, oldType, changeTable).(type) {
			switch compareTypes(tc.Type, oldType, context).(type) {
			case nil, *TypeChangeDefinitionChanged:
				return &TypeChangeScalarToUnion{TypePair{oldType, newType}, i}
			}
		}
	}

	// CHANGE: Incompatible change from Simple to Generalized
	return &TypeChangeIncompatible{TypePair{oldType, newType}}
}

// func detectGeneralizedTypeChanges(newType, oldType *GeneralizedType, changeTable ChangeTable) TypeChange {
func compareGeneralizedTypes(newType, oldType *GeneralizedType, context *EvolutionContext) TypeChange {
	// A GeneralizedType can change in many ways...
	if newType.Cases.IsOptional() {
		return detectOptionalChanges(newType, oldType, context)
	} else if newType.Cases.IsUnion() {
		return detectUnionChanges(newType, oldType, context)
	} else {
		switch newType.Dimensionality.(type) {
		case nil:
			// TODO: Not an Optional, Union, Stream, Vector, Array, Map...
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

// func detectOptionalChanges(newType, oldType *GeneralizedType, changeTable ChangeTable) TypeChange {
func detectOptionalChanges(newType, oldType *GeneralizedType, context *EvolutionContext) TypeChange {
	if !oldType.Cases.IsOptional() {
		if oldType.Cases.IsUnion() && oldType.Cases.HasNullOption() {
			// An Optional<T> can become a Union<null, T, ...> ONLY if
			// 	1. type T does not change, or
			// 	2. type T's TypeDefinition changed

			// Look for a matching type in the old Union
			for i, c := range oldType.Cases[1:] {
				// switch detectTypeChanges(newType.Cases[1].Type, c.Type, changeTable).(type) {
				switch compareTypes(newType.Cases[1].Type, c.Type, context).(type) {
				case nil, *TypeChangeDefinitionChanged:
					return &TypeChangeUnionToOptional{TypePair{oldType, newType}, i + 1}
				}
			}
		}

		// CHANGE: Changed a non-Optional/Union to an Optional
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	// if ch := detectTypeChanges(newType.Cases[1].Type, oldType.Cases[1].Type, changeTable); ch != nil {
	if ch := compareTypes(newType.Cases[1].Type, oldType.Cases[1].Type, context); ch != nil {
		// CHANGE: Changed Optional type
		return &TypeChangeOptionalTypeChanged{TypePair{oldType, newType}, ch}
	}
	return nil
}

// func detectUnionChanges(newType, oldType *GeneralizedType, changeTable ChangeTable) TypeChange {
func detectUnionChanges(newType, oldType *GeneralizedType, context *EvolutionContext) TypeChange {
	if !oldType.Cases.IsUnion() {
		if oldType.Cases.IsOptional() && newType.Cases.HasNullOption() {
			for i, c := range newType.Cases[1:] {
				// switch detectTypeChanges(c.Type, oldType.Cases[1].Type, changeTable).(type) {
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

			// switch detectTypeChanges(newCase.Type, oldCase.Type, changeTable).(type) {
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

// func detectStreamChanges(newType, oldType *GeneralizedType, changeTable ChangeTable) TypeChange {
func detectStreamChanges(newType, oldType *GeneralizedType, context *EvolutionContext) TypeChange {
	if _, ok := oldType.Dimensionality.(*Stream); !ok {
		// CHANGE: Changed a non-Stream to a Stream
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	// if ch := detectTypeChanges(newType.Cases[0].Type, oldType.Cases[0].Type, changeTable); ch != nil {
	if ch := compareTypes(newType.Cases[0].Type, oldType.Cases[0].Type, context); ch != nil {
		// CHANGE: Changed Stream type
		return &TypeChangeStreamTypeChanged{TypePair{oldType, newType}, ch}
	}
	return nil
}

// func detectVectorChanges(newType, oldType *GeneralizedType, changeTable ChangeTable) TypeChange {
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

	// if ch := detectTypeChanges(newType.Cases[0].Type, oldType.Cases[0].Type, changeTable); ch != nil {
	if ch := compareTypes(newType.Cases[0].Type, oldType.Cases[0].Type, context); ch != nil {
		// CHANGE: Changed Vector type
		return &TypeChangeVectorTypeChanged{TypePair{oldType, newType}, ch}
	}

	return nil
}

// func detectArrayChanges(newType, oldType *GeneralizedType, changeTable ChangeTable) TypeChange {
func detectArrayChanges(newType, oldType *GeneralizedType, context *EvolutionContext) TypeChange {
	newDim := newType.Dimensionality.(*Array)
	oldDim, ok := oldType.Dimensionality.(*Array)
	if !ok {
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	// if ch := detectTypeChanges(newType.Cases[0].Type, oldType.Cases[0].Type, changeTable); ch != nil {
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

// func detectMapChanges(newType, oldType *GeneralizedType, changeTable ChangeTable) TypeChange {
func detectMapChanges(newType, oldType *GeneralizedType, context *EvolutionContext) TypeChange {
	newDim := newType.Dimensionality.(*Map)
	oldDim, ok := oldType.Dimensionality.(*Map)
	if !ok {
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}

	// if ch := detectTypeChanges(newDim.KeyType, oldDim.KeyType, changeTable); ch != nil {
	if ch := compareTypes(newDim.KeyType, oldDim.KeyType, context); ch != nil {
		// CHANGE: Changed Map key type
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}
	// if ch := detectTypeChanges(newType.Cases[0].Type, oldType.Cases[0].Type, changeTable); ch != nil {
	if ch := compareTypes(newType.Cases[0].Type, oldType.Cases[0].Type, context); ch != nil {
		// CHANGE: Changed Map value type
		return &TypeChangeIncompatible{TypePair{oldType, newType}}
	}
	return nil
}
