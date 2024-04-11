% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef StreamsReaderBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = StreamsReaderBase()
      obj.state_ = 0;
    end

    function close(obj)
      obj.close_();
      if obj.state_ ~= 8
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
    function more = has_int_data(obj)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      more = obj.has_int_data_();
      if ~more
        obj.state_ = 2;
      end
    end

    function value = read_int_data(obj)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      value = obj.read_int_data_();
    end

    % Ordinal 1
    function more = has_optional_int_data(obj)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      more = obj.has_optional_int_data_();
      if ~more
        obj.state_ = 4;
      end
    end

    function value = read_optional_int_data(obj)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      value = obj.read_optional_int_data_();
    end

    % Ordinal 2
    function more = has_record_with_optional_vector_data(obj)
      if obj.state_ ~= 4
        obj.raise_unexpected_state_(4);
      end

      more = obj.has_record_with_optional_vector_data_();
      if ~more
        obj.state_ = 6;
      end
    end

    function value = read_record_with_optional_vector_data(obj)
      if obj.state_ ~= 4
        obj.raise_unexpected_state_(4);
      end

      value = obj.read_record_with_optional_vector_data_();
    end

    % Ordinal 3
    function more = has_fixed_vector(obj)
      if obj.state_ ~= 6
        obj.raise_unexpected_state_(6);
      end

      more = obj.has_fixed_vector_();
      if ~more
        obj.state_ = 8;
      end
    end

    function value = read_fixed_vector(obj)
      if obj.state_ ~= 6
        obj.raise_unexpected_state_(6);
      end

      value = obj.read_fixed_vector_();
    end

    function copy_to(obj, writer)
      while obj.has_int_data()
        item = obj.read_int_data();
        writer.write_int_data({item});
      end
      while obj.has_optional_int_data()
        item = obj.read_optional_int_data();
        writer.write_optional_int_data({item});
      end
      while obj.has_record_with_optional_vector_data()
        item = obj.read_record_with_optional_vector_data();
        writer.write_record_with_optional_vector_data({item});
      end
      while obj.has_fixed_vector()
        item = obj.read_fixed_vector();
        writer.write_fixed_vector({item});
      end
    end
  end

  methods (Static)
    function res = schema()
      res = test_model.StreamsWriterBase.schema;
    end
  end

  methods (Abstract, Access=protected)
    has_int_data_(obj)
    read_int_data_(obj)
    has_optional_int_data_(obj)
    read_optional_int_data_(obj)
    has_record_with_optional_vector_data_(obj)
    read_record_with_optional_vector_data_(obj)
    has_fixed_vector_(obj)
    read_fixed_vector_(obj)

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
        name = 'read_int_data';
      elseif state == 2
        name = 'read_optional_int_data';
      elseif state == 4
        name = 'read_record_with_optional_vector_data';
      elseif state == 6
        name = 'read_fixed_vector';
      else
        name = '<unknown>';
      end
    end
  end
end
