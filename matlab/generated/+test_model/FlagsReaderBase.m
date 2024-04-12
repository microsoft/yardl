% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef FlagsReaderBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = FlagsReaderBase()
      obj.state_ = 0;
    end

    function close(obj)
      obj.close_();
      if obj.state_ ~= 2
        expected_method = obj.state_to_method_name_(obj.state_);
        throw(yardl.ProtocolError("Protocol reader closed before all data was consumed. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function more = has_days(obj)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      more = obj.has_days_();
      if ~more
        obj.state_ = 1;
      end
    end

    function value = read_days(obj)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      value = obj.read_days_();
    end

    % Ordinal 1
    function more = has_formats(obj)
      if obj.state_ ~= 1
        obj.raise_unexpected_state_(1);
      end

      more = obj.has_formats_();
      if ~more
        obj.state_ = 2;
      end
    end

    function value = read_formats(obj)
      if obj.state_ ~= 1
        obj.raise_unexpected_state_(1);
      end

      value = obj.read_formats_();
    end

    function copy_to(obj, writer)
      while obj.has_days()
        item = obj.read_days();
        writer.write_days({item});
      end
      writer.end_days();
      while obj.has_formats()
        item = obj.read_formats();
        writer.write_formats({item});
      end
      writer.end_formats();
    end
  end

  methods (Static)
    function res = schema()
      res = test_model.FlagsWriterBase.schema;
    end
  end

  methods (Abstract, Access=protected)
    has_days_(obj)
    read_days_(obj)
    has_formats_(obj)
    read_formats_(obj)

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
        name = 'read_days';
      elseif state == 1
        name = 'read_formats';
      else
        name = '<unknown>';
      end
    end
  end
end
