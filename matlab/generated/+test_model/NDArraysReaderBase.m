classdef NDArraysReaderBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = NDArraysReaderBase()
      obj.state_ = 0;
    end

    function close(obj)
      obj.close_();
      if obj.state_ ~= 10
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
    function value = read_ints(obj)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      value = obj.read_ints_();
      obj.state_ = 2;
    end

    % Ordinal 1
    function value = read_simple_record_array(obj)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      value = obj.read_simple_record_array_();
      obj.state_ = 4;
    end

    % Ordinal 2
    function value = read_record_with_vlens_array(obj)
      if obj.state_ ~= 4
        obj.raise_unexpected_state_(4);
      end

      value = obj.read_record_with_vlens_array_();
      obj.state_ = 6;
    end

    % Ordinal 3
    function value = read_record_with_nd_arrays(obj)
      if obj.state_ ~= 6
        obj.raise_unexpected_state_(6);
      end

      value = obj.read_record_with_nd_arrays_();
      obj.state_ = 8;
    end

    % Ordinal 4
    function value = read_named_array(obj)
      if obj.state_ ~= 8
        obj.raise_unexpected_state_(8);
      end

      value = obj.read_named_array_();
      obj.state_ = 10;
    end

    function copy_to(obj, writer)
      writer.write_ints(obj.read_ints());
      writer.write_simple_record_array(obj.read_simple_record_array());
      writer.write_record_with_vlens_array(obj.read_record_with_vlens_array());
      writer.write_record_with_nd_arrays(obj.read_record_with_nd_arrays());
      writer.write_named_array(obj.read_named_array());
    end
  end

  methods (Static)
    function res = schema()
      res = test_model.NDArraysWriterBase.schema;
    end
  end

  methods (Abstract, Access=protected)
    read_ints_(obj, value)
    read_simple_record_array_(obj, value)
    read_record_with_vlens_array_(obj, value)
    read_record_with_nd_arrays_(obj, value)
    read_named_array_(obj, value)

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
        name = 'read_ints';
      elseif state == 2
        name = 'read_simple_record_array';
      elseif state == 4
        name = 'read_record_with_vlens_array';
      elseif state == 6
        name = 'read_record_with_nd_arrays';
      elseif state == 8
        name = 'read_named_array';
      else
        name = '<unknown>';
      end
    end
  end
end
