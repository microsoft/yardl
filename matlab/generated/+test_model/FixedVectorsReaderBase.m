classdef FixedVectorsReaderBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = FixedVectorsReaderBase()
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
    function value = read_fixed_int_vector(obj)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      value = obj.read_fixed_int_vector_();
      obj.state_ = 2;
    end

    % Ordinal 1
    function value = read_fixed_simple_record_vector(obj)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      value = obj.read_fixed_simple_record_vector_();
      obj.state_ = 4;
    end

    % Ordinal 2
    function value = read_fixed_record_with_vlens_vector(obj)
      if obj.state_ ~= 4
        obj.raise_unexpected_state_(4);
      end

      value = obj.read_fixed_record_with_vlens_vector_();
      obj.state_ = 6;
    end

    % Ordinal 3
    function value = read_record_with_fixed_vectors(obj)
      if obj.state_ ~= 6
        obj.raise_unexpected_state_(6);
      end

      value = obj.read_record_with_fixed_vectors_();
      obj.state_ = 8;
    end

    function copy_to(obj, writer)
      writer.write_fixed_int_vector(obj.read_fixed_int_vector());
      writer.write_fixed_simple_record_vector(obj.read_fixed_simple_record_vector());
      writer.write_fixed_record_with_vlens_vector(obj.read_fixed_record_with_vlens_vector());
      writer.write_record_with_fixed_vectors(obj.read_record_with_fixed_vectors());
    end
  end

  methods (Static)
    function res = schema()
      res = test_model.FixedVectorsWriterBase.schema;
    end
  end

  methods (Abstract, Access=protected)
    read_fixed_int_vector_(obj, value)
    read_fixed_simple_record_vector_(obj, value)
    read_fixed_record_with_vlens_vector_(obj, value)
    read_record_with_fixed_vectors_(obj, value)

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
        name = 'read_fixed_int_vector';
      elseif state == 2
        name = 'read_fixed_simple_record_vector';
      elseif state == 4
        name = 'read_fixed_record_with_vlens_vector';
      elseif state == 6
        name = 'read_record_with_fixed_vectors';
      else
        name = '<unknown>';
      end
    end
  end
end
