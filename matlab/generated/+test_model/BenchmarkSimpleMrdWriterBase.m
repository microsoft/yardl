% Abstract writer for protocol BenchmarkSimpleMrd
classdef (Abstract) BenchmarkSimpleMrdWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = BenchmarkSimpleMrdWriterBase()
      obj.state_ = 0;
    end

    function close(obj)
      if obj.state_ == 1
        obj.end_stream_();
        obj.close_();
        return
      end
      obj.close_();
      if obj.state_ ~= 2
        expected_method = obj.state_to_method_name_(bitand((int32(obj.state_) + 1), bitcmp(1, 'int8')));
        throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function write_data(obj, value)
      if bitand(int32(obj.state_), bitcmp(1, 'int8')) ~= 0
        obj.raise_unexpected_state_(0);
      end

      obj.write_data_(value);
      obj.state_ = 1;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"BenchmarkSimpleMrd","sequence":[{"name":"data","type":{"stream":{"items":[{"tag":"acquisition","explicitTag":true,"type":"TestModel.SimpleAcquisition"},{"tag":"image","explicitTag":true,"type":{"name":"Image.Image","typeArguments":["float32"]}}]}}}]},"types":[{"name":"Image","typeParameters":["T"],"type":{"array":{"items":"T","dimensions":[{"name":"x"},{"name":"y"}]}}},{"name":"SimpleAcquisition","fields":[{"name":"flags","type":"uint64"},{"name":"idx","type":"TestModel.SimpleEncodingCounters"},{"name":"data","type":{"array":{"items":"complexfloat32","dimensions":2}}},{"name":"trajectory","type":{"array":{"items":"float32","dimensions":2}}}]},{"name":"SimpleEncodingCounters","fields":[{"name":"e1","type":[null,"uint32"]},{"name":"e2","type":[null,"uint32"]},{"name":"slice","type":[null,"uint32"]},{"name":"repetition","type":[null,"uint32"]}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_data_(obj, value)

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
        name = 'write_data';
      else
        name = '<unknown>';
      end
    end
  end
end
