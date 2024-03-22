% Abstract writer for protocol NDArrays
classdef (Abstract) NDArraysWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = NDArraysWriterBase()
      obj.state_ = 0;
    end

    function close(obj)
      obj.close_();
      if obj.state_ ~= 10
        expected_method = obj.state_to_method_name_(bitand((int32(obj.state_) + 1), bitcmp(1, 'int8')));
        throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function write_ints(obj, value)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      obj.write_ints_(value);
      obj.state_ = 2;
    end

    % Ordinal 1
    function write_simple_record_array(obj, value)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      obj.write_simple_record_array_(value);
      obj.state_ = 4;
    end

    % Ordinal 2
    function write_record_with_vlens_array(obj, value)
      if obj.state_ ~= 4
        obj.raise_unexpected_state_(4);
      end

      obj.write_record_with_vlens_array_(value);
      obj.state_ = 6;
    end

    % Ordinal 3
    function write_record_with_nd_arrays(obj, value)
      if obj.state_ ~= 6
        obj.raise_unexpected_state_(6);
      end

      obj.write_record_with_nd_arrays_(value);
      obj.state_ = 8;
    end

    % Ordinal 4
    function write_named_array(obj, value)
      if obj.state_ ~= 8
        obj.raise_unexpected_state_(8);
      end

      obj.write_named_array_(value);
      obj.state_ = 10;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"NDArrays","sequence":[{"name":"ints","type":{"array":{"items":"int32","dimensions":2}}},{"name":"simpleRecordArray","type":{"array":{"items":"TestModel.SimpleRecord","dimensions":2}}},{"name":"recordWithVlensArray","type":{"array":{"items":"TestModel.RecordWithVlens","dimensions":2}}},{"name":"recordWithNDArrays","type":"TestModel.RecordWithNDArrays"},{"name":"namedArray","type":"TestModel.NamedNDArray"}]},"types":[{"name":"NamedNDArray","type":{"array":{"items":"int32","dimensions":[{"name":"dimA"},{"name":"dimB"}]}}},{"name":"RecordWithNDArrays","fields":[{"name":"ints","type":{"array":{"items":"int32","dimensions":2}}},{"name":"fixedSimpleRecordArray","type":{"array":{"items":"TestModel.SimpleRecord","dimensions":2}}},{"name":"fixedRecordWithVlensArray","type":{"array":{"items":"TestModel.RecordWithVlens","dimensions":2}}}]},{"name":"RecordWithVlens","fields":[{"name":"a","type":{"vector":{"items":"TestModel.SimpleRecord"}}},{"name":"b","type":"int32"},{"name":"c","type":"int32"}]},{"name":"SimpleRecord","fields":[{"name":"x","type":"int32"},{"name":"y","type":"int32"},{"name":"z","type":"int32"}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_ints_(obj, value)
    write_simple_record_array_(obj, value)
    write_record_with_vlens_array_(obj, value)
    write_record_with_nd_arrays_(obj, value)
    write_named_array_(obj, value)

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
        name = 'write_ints';
      elseif state == 2
        name = 'write_simple_record_array';
      elseif state == 4
        name = 'write_record_with_vlens_array';
      elseif state == 6
        name = 'write_record_with_nd_arrays';
      elseif state == 8
        name = 'write_named_array';
      else
        name = '<unknown>';
      end
    end
  end
end
