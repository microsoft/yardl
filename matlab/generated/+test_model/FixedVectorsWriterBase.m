% Abstract writer for protocol FixedVectors
classdef (Abstract) FixedVectorsWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = FixedVectorsWriterBase()
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
    function write_fixed_int_vector(obj, value)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      obj.write_fixed_int_vector_(value);
      obj.state_ = 2;
    end

    % Ordinal 1
    function write_fixed_simple_record_vector(obj, value)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      obj.write_fixed_simple_record_vector_(value);
      obj.state_ = 4;
    end

    % Ordinal 2
    function write_fixed_record_with_vlens_vector(obj, value)
      if obj.state_ ~= 4
        obj.raise_unexpected_state_(4);
      end

      obj.write_fixed_record_with_vlens_vector_(value);
      obj.state_ = 6;
    end

    % Ordinal 3
    function write_record_with_fixed_vectors(obj, value)
      if obj.state_ ~= 6
        obj.raise_unexpected_state_(6);
      end

      obj.write_record_with_fixed_vectors_(value);
      obj.state_ = 8;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"FixedVectors","sequence":[{"name":"fixedIntVector","type":{"vector":{"items":"int32","length":5}}},{"name":"fixedSimpleRecordVector","type":{"vector":{"items":"TestModel.SimpleRecord","length":3}}},{"name":"fixedRecordWithVlensVector","type":{"vector":{"items":"TestModel.RecordWithVlens","length":2}}},{"name":"recordWithFixedVectors","type":"TestModel.RecordWithFixedVectors"}]},"types":[{"name":"RecordWithFixedVectors","fields":[{"name":"fixedIntVector","type":{"vector":{"items":"int32","length":5}}},{"name":"fixedSimpleRecordVector","type":{"vector":{"items":"TestModel.SimpleRecord","length":3}}},{"name":"fixedRecordWithVlensVector","type":{"vector":{"items":"TestModel.RecordWithVlens","length":2}}}]},{"name":"RecordWithVlens","fields":[{"name":"a","type":{"vector":{"items":"TestModel.SimpleRecord"}}},{"name":"b","type":"int32"},{"name":"c","type":"int32"}]},{"name":"SimpleRecord","fields":[{"name":"x","type":"int32"},{"name":"y","type":"int32"},{"name":"z","type":"int32"}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_fixed_int_vector_(obj, value)
    write_fixed_simple_record_vector_(obj, value)
    write_fixed_record_with_vlens_vector_(obj, value)
    write_record_with_fixed_vectors_(obj, value)

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
        name = 'write_fixed_int_vector';
      elseif state == 2
        name = 'write_fixed_simple_record_vector';
      elseif state == 4
        name = 'write_fixed_record_with_vlens_vector';
      elseif state == 6
        name = 'write_record_with_fixed_vectors';
      else
        name = '<unknown>';
      end
    end
  end
end
