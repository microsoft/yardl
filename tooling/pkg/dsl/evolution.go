// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package dsl

import (
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

		// Need to "rename" Old Definitions to their semantically equivalent New Definition Name
		for _, ch := range definitionChanges {
			// TODO: I don't think this will work for Generic TypeDefinitions with TypeArgs (needs to **Visit** each TypeDefinition)
			// log.Debug().Msgf("Will write alias (old) %s === (new) %s", ch.PreviousDefinition().GetDefinitionMeta().GetQualifiedName(), ch.LatestDefinition().GetDefinitionMeta().GetQualifiedName())
			ch.PreviousDefinition().GetDefinitionMeta().Name = ch.PreviousDefinition().GetDefinitionMeta().Name + "_" + versionLabels[i]
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
			panic("I don't want nil DefinitionChanges here ")
		}

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
	isUserTypeDef := make(map[string]bool)
	for _, ns := range env.Namespaces {
		for _, newTd := range ns.TypeDefinitions {
			isUserTypeDef[newTd.GetDefinitionMeta().GetQualifiedName()] = true
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
	}
	return usedTypeDefs
}

// Returns the last name in the chain of resolved names
// i.e. the resolved Record, Enum, or NamedType name
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
	allNewTypeDefs := collectUserTypeDefsUsedInProtocols(newEnv)
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
				// log.Debug().Msgf("New %s semantically equals Old %s", newName, oldName)
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

	// Now we're finished detecting all changes between models... Time to collect/filter TypeDefinition changes
	allDefinitionChanges := make([]DefinitionChange, 0)
	for _, newTd := range allNewTypeDefs {
		newName := newTd.GetDefinitionMeta().GetQualifiedName()
		for _, oldTd := range allOldTypeDefs {
			oldName := oldTd.GetDefinitionMeta().GetQualifiedName()

			if ch, ok := changes[newName][oldName]; ok {
				allDefinitionChanges = append(allDefinitionChanges, ch)
			}
		}
	}

	// Also save all "added aliases" for later codegen
	allDefinitionChanges = append(allDefinitionChanges, context.AliasesRemoved...)

	// Finally, filter out changes with duplicate "old" TypeDefinitions because we only need the first one for codegen
	neededDefinitionChanges := make([]DefinitionChange, 0)
	uniqueOldNames := make(map[string]bool)
	for _, change := range allDefinitionChanges {
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
			switch ch := ch.(type) {
			case *TypeChangeIncompatible:
				return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
			default:
				return &NamedTypeChange{DefinitionPair{oldTd, newTd}, ch}
			}

		case *RecordDefinition:
			oldType := &SimpleType{NodeMeta: *oldTd.GetNodeMeta(), Name: oldTd.GetDefinitionMeta().GetQualifiedName(), ResolvedDefinition: oldTd}
			ch := compareTypes(newTd.Type, oldType, context)
			if ch == nil {
				return nil
			}
			switch ch := ch.(type) {
			case *TypeChangeIncompatible:
				return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
			default:
				return &NamedTypeChange{DefinitionPair{oldTd, newTd}, ch}
			}

		case *EnumDefinition:
			oldType := &SimpleType{NodeMeta: *oldTd.GetNodeMeta(), Name: oldTd.GetDefinitionMeta().GetQualifiedName(), ResolvedDefinition: oldTd}
			ch := compareTypes(newTd.Type, oldType, context)
			if ch == nil {
				return nil
			}
			switch ch := ch.(type) {
			case *TypeChangeIncompatible:
				return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
			default:
				return &NamedTypeChange{DefinitionPair{oldTd, newTd}, ch}
			}

		default:
			return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
		}

	case *RecordDefinition:
		switch oldTd := oldTd.(type) {
		case *RecordDefinition:
			ch := compareRecordDefinitions(newTd, oldTd, context)
			if ch == nil {
				return nil
			}
			return ch

		case *NamedType:
			newType := &SimpleType{NodeMeta: *newTd.GetNodeMeta(), Name: newTd.GetDefinitionMeta().GetQualifiedName(), ResolvedDefinition: newTd}
			ch := compareTypes(newType, oldTd.Type, context)
			if ch == nil {
				return nil
			}
			switch ch := ch.(type) {
			case *TypeChangeIncompatible:
				return &DefinitionChangeIncompatible{DefinitionPair{oldTd, newTd}}
			default:
				return &NamedTypeChange{DefinitionPair{oldTd, newTd}, ch}
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
}

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

	newDef := newType.ResolvedDefinition
	oldDef := oldType.ResolvedDefinition

	newName := newDef.GetDefinitionMeta().GetQualifiedName()
	oldName := oldDef.GetDefinitionMeta().GetQualifiedName()

	// Here is where we do a lookup on the two ResolvedDefinitions
	// 1. Are they semantically equal?
	// 1. Did we already compare them?

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

	// The two TypeDefinitions are semantically equivalent, so context.Compared[newName][oldName] should be true
	// Did the TypeDefinition change between versions?
	if ch, ok := context.Changes[newName][oldName]; ok {
		if _, ok := ch.(*DefinitionChangeIncompatible); ok {
			return &TypeChangeIncompatible{TypePair{oldType, newType}}
		} else {
			return &TypeChangeDefinitionChanged{TypePair{oldType, newType}, ch}
		}
	}

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
