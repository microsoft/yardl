% Abstract writer for protocol ProtocolWithComputedFields
classdef (Abstract) ProtocolWithComputedFieldsWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = ProtocolWithComputedFieldsWriterBase()
      obj.state_ = 0;
    end

    function close(obj)
      obj.close_();
      if obj.state_ ~= 2
        expected_method = obj.state_to_method_name_(bitand((int32(obj.state_) + 1), bitcmp(1, 'int8')));
        throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function write_record_with_computed_fields(obj, value)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      obj.write_record_with_computed_fields_(value);
      obj.state_ = 2;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"ProtocolWithComputedFields","sequence":[{"name":"recordWithComputedFields","type":"TestModel.RecordWithComputedFields"}]},"types":[{"name":"GenericRecordWithComputedFields","typeParameters":["T0","T1"],"fields":[{"name":"f1","type":[{"tag":"T0","type":"T0"},{"tag":"T1","type":"T1"}]}]},{"name":"MyTuple","typeParameters":["T1","T2"],"type":{"name":"Tuples.Tuple","typeArguments":["T1","T2"]}},{"name":"MyTuple","typeParameters":["T1","T2"],"type":{"name":"BasicTypes.MyTuple","typeArguments":["T1","T2"]}},{"name":"NamedNDArray","type":{"array":{"items":"int32","dimensions":[{"name":"dimA"},{"name":"dimB"}]}}},{"name":"RecordWithComputedFields","fields":[{"name":"arrayField","type":{"array":{"items":"int32","dimensions":[{"name":"x"},{"name":"y"}]}}},{"name":"arrayFieldMapDimensions","type":{"array":{"items":"int32","dimensions":[{"name":"x"},{"name":"y"}]}}},{"name":"dynamicArrayField","type":{"array":{"items":"int32"}}},{"name":"fixedArrayField","type":{"array":{"items":"int32","dimensions":[{"name":"x","length":3},{"name":"y","length":4}]}}},{"name":"intField","type":"int32"},{"name":"int8Field","type":"int8"},{"name":"uint8Field","type":"uint8"},{"name":"int16Field","type":"int16"},{"name":"uint16Field","type":"uint16"},{"name":"uint32Field","type":"uint32"},{"name":"int64Field","type":"int64"},{"name":"uint64Field","type":"uint64"},{"name":"sizeField","type":"size"},{"name":"float32Field","type":"float32"},{"name":"float64Field","type":"float64"},{"name":"complexfloat32Field","type":"complexfloat32"},{"name":"complexfloat64Field","type":"complexfloat64"},{"name":"stringField","type":"string"},{"name":"tupleField","type":{"name":"TestModel.MyTuple","typeArguments":["int32","int32"]}},{"name":"vectorField","type":{"vector":{"items":"int32"}}},{"name":"vectorOfVectorsField","type":{"vector":{"items":{"vector":{"items":"int32"}}}}},{"name":"fixedVectorField","type":{"vector":{"items":"int32","length":3}}},{"name":"optionalNamedArray","type":[null,"TestModel.NamedNDArray"]},{"name":"intFloatUnion","type":[{"tag":"int32","type":"int32"},{"tag":"float32","type":"float32"}]},{"name":"nullableIntFloatUnion","type":[null,{"tag":"int32","type":"int32"},{"tag":"float32","type":"float32"}]},{"name":"unionWithNestedGenericUnion","type":[{"tag":"int","explicitTag":true,"type":"int32"},{"tag":"genericRecordWithComputedFields","explicitTag":true,"type":{"name":"BasicTypes.GenericRecordWithComputedFields","typeArguments":["string","float32"]}}]},{"name":"mapField","type":{"map":{"keys":"string","values":"string"}}}]},{"name":"Tuple","typeParameters":["T1","T2"],"fields":[{"name":"v1","type":"T1"},{"name":"v2","type":"T2"}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_record_with_computed_fields_(obj, value)

    end_stream_(obj)
    close_(obj)
  end

  methods (Access=private)
    function raise_unexpected_state_(obj, actual)
      expected_method = obj.state_to_method_name_(obj.state_);
      actual_method = obj.state_to_method_name_(actual);
      throw(yardl.ProtocolError("Expected call to '%s' but received call to '%s'", expected_method, actual_method));
    end

    function name = state_to_method_name_(obj, state)
      if state == 0
        name = 'write_record_with_computed_fields';
      else
        name = '<unknown>';
      end
    end
  end
end