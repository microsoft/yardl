classdef EnumsReaderBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = EnumsReaderBase()
      obj.state_ = 0;
    end

    function close(obj)
      obj.close_();
      if obj.state_ ~= 6
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
    function value = read_single(obj)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      value = obj.read_single_();
      obj.state_ = 2;
    end

    % Ordinal 1
    function value = read_vec(obj)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      value = obj.read_vec_();
      obj.state_ = 4;
    end

    % Ordinal 2
    function value = read_size(obj)
      if obj.state_ ~= 4
        obj.raise_unexpected_state_(4);
      end

      value = obj.read_size_();
      obj.state_ = 6;
    end

    function copy_to(obj, writer)
      writer.write_single(obj.read_single());
      writer.write_vec(obj.read_vec());
      writer.write_size(obj.read_size());
    end
  end

  methods (Static)
    function res = schema()
      res = test_model.EnumsWriterBase.schema;
    end
  end

  methods (Abstract, Access=protected)
    read_single_(obj, value)
    read_vec_(obj, value)
    read_size_(obj, value)

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
        name = 'read_single';
      elseif state == 2
        name = 'read_vec';
      elseif state == 4
        name = 'read_size';
      else
        name = '<unknown>';
      end
    end
  end
end