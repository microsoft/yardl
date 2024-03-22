classdef BenchmarkSimpleMrdReaderBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = BenchmarkSimpleMrdReaderBase()
      obj.state_ = 0;
    end

    function close(obj)
      obj.close_();
      if obj.state_ ~= 2
        if mod(obj.state_, 2) == 1
          previous_method = obj.state_to_method_name_(obj.state_ - 1);
          throw(yardl.ProtocolError("Protocol reader closed before all data was consumed. The iterable returned by '%s' was not fully consumed.", previous_method));
        else
          expected_method = obj.state_to_method_name_(obj.state_);
          throw(yardl.ProtocolError("Protocol reader closed before all data was consumed. Expected call to '%s'.", expected_method));
        end
      end
    end

    % Ordinal 0
    function value = read_data(obj)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      value = obj.read_data_();
      obj.state_ = 2;
    end

    function copy_to(obj, writer)
      writer.write_data(obj.read_data());
    end
  end

  methods (Static)
    function res = schema()
      res = test_model.BenchmarkSimpleMrdWriterBase.schema;
    end
  end

  methods (Abstract, Access=protected)
    read_data_(obj, value)

    close_(obj)
  end

  methods (Access=private)
    function raise_unexpected_state_(obj, actual)
      actual_method = obj.state_to_method_name_(actual);
      expected_method = obj.state_to_method_name_(obj.state_);
      throw(yardl.ProtocolError("Expected call to '%s' but received call to '%s'.", expected_method, actual_method));
    end

    function name = state_to_method_name_(obj, state)
      if state == 0
        name = 'read_data';
      else
        name = '<unknown>';
      end
    end
  end
end
