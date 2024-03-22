% Abstract writer for protocol Unions
classdef (Abstract) UnionsWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = UnionsWriterBase()
      obj.state_ = 0;
    end

    function close(obj)
      obj.close_();
      if obj.state_ ~= 8
        expected_method = obj.state_to_method_name_(bitand((int32(obj.state_) + 1), bitcmp(1, 'int8')));
        throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function write_int_or_simple_record(obj, value)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      obj.write_int_or_simple_record_(value);
      obj.state_ = 2;
    end

    % Ordinal 1
    function write_int_or_record_with_vlens(obj, value)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      obj.write_int_or_record_with_vlens_(value);
      obj.state_ = 4;
    end

    % Ordinal 2
    function write_monosotate_or_int_or_simple_record(obj, value)
      if obj.state_ ~= 4
        obj.raise_unexpected_state_(4);
      end

      obj.write_monosotate_or_int_or_simple_record_(value);
      obj.state_ = 6;
    end

    % Ordinal 3
    function write_record_with_unions(obj, value)
      if obj.state_ ~= 6
        obj.raise_unexpected_state_(6);
      end

      obj.write_record_with_unions_(value);
      obj.state_ = 8;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"Unions","sequence":[{"name":"intOrSimpleRecord","type":[{"tag":"int32","type":"int32"},{"tag":"SimpleRecord","type":"TestModel.SimpleRecord"}]},{"name":"intOrRecordWithVlens","type":[{"tag":"int32","type":"int32"},{"tag":"RecordWithVlens","type":"TestModel.RecordWithVlens"}]},{"name":"monosotateOrIntOrSimpleRecord","type":[null,{"tag":"int32","type":"int32"},{"tag":"SimpleRecord","type":"TestModel.SimpleRecord"}]},{"name":"recordWithUnions","type":"BasicTypes.RecordWithUnions"}]},"types":[{"name":"DaysOfWeek","values":[{"symbol":"monday","value":1},{"symbol":"tuesday","value":2},{"symbol":"wednesday","value":4},{"symbol":"thursday","value":8},{"symbol":"friday","value":16},{"symbol":"saturday","value":32},{"symbol":"sunday","value":64}]},{"name":"Fruits","values":[{"symbol":"apple","value":0},{"symbol":"banana","value":1},{"symbol":"pear","value":2}]},{"name":"GenericNullableUnion2","typeParameters":["T1","T2"],"type":[null,{"tag":"T1","type":"T1"},{"tag":"T2","type":"T2"}]},{"name":"RecordWithUnions","fields":[{"name":"nullOrIntOrString","type":[null,{"tag":"int32","type":"int32"},{"tag":"string","type":"string"}]},{"name":"dateOrDatetime","type":[{"tag":"time","type":"time"},{"tag":"datetime","type":"datetime"}]},{"name":"nullOrFruitsOrDaysOfWeek","type":{"name":"BasicTypes.GenericNullableUnion2","typeArguments":["BasicTypes.Fruits","BasicTypes.DaysOfWeek"]}}]},{"name":"RecordWithVlens","fields":[{"name":"a","type":{"vector":{"items":"TestModel.SimpleRecord"}}},{"name":"b","type":"int32"},{"name":"c","type":"int32"}]},{"name":"SimpleRecord","fields":[{"name":"x","type":"int32"},{"name":"y","type":"int32"},{"name":"z","type":"int32"}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_int_or_simple_record_(obj, value)
    write_int_or_record_with_vlens_(obj, value)
    write_monosotate_or_int_or_simple_record_(obj, value)
    write_record_with_unions_(obj, value)

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
        name = 'write_int_or_simple_record';
      elseif state == 2
        name = 'write_int_or_record_with_vlens';
      elseif state == 4
        name = 'write_monosotate_or_int_or_simple_record';
      elseif state == 6
        name = 'write_record_with_unions';
      else
        name = '<unknown>';
      end
    end
  end
end
