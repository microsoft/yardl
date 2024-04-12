% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef AdvancedGenericsReaderBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = AdvancedGenericsReaderBase()
      obj.state_ = 0;
    end

    function close(obj)
      obj.close_();
      if obj.state_ ~= 5
        expected_method = obj.state_to_method_name_(obj.state_);
        throw(yardl.ProtocolError("Protocol reader closed before all data was consumed. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function value = read_float_image_image(obj)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      value = obj.read_float_image_image_();
      obj.state_ = 1;
    end

    % Ordinal 1
    function value = read_generic_record_1(obj)
      if obj.state_ ~= 1
        obj.raise_unexpected_state_(1);
      end

      value = obj.read_generic_record_1_();
      obj.state_ = 2;
    end

    % Ordinal 2
    function value = read_tuple_of_optionals(obj)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      value = obj.read_tuple_of_optionals_();
      obj.state_ = 3;
    end

    % Ordinal 3
    function value = read_tuple_of_optionals_alternate_syntax(obj)
      if obj.state_ ~= 3
        obj.raise_unexpected_state_(3);
      end

      value = obj.read_tuple_of_optionals_alternate_syntax_();
      obj.state_ = 4;
    end

    % Ordinal 4
    function value = read_tuple_of_vectors(obj)
      if obj.state_ ~= 4
        obj.raise_unexpected_state_(4);
      end

      value = obj.read_tuple_of_vectors_();
      obj.state_ = 5;
    end

    function copy_to(obj, writer)
      writer.write_float_image_image(obj.read_float_image_image());
      writer.write_generic_record_1(obj.read_generic_record_1());
      writer.write_tuple_of_optionals(obj.read_tuple_of_optionals());
      writer.write_tuple_of_optionals_alternate_syntax(obj.read_tuple_of_optionals_alternate_syntax());
      writer.write_tuple_of_vectors(obj.read_tuple_of_vectors());
    end
  end

  methods (Static)
    function res = schema()
      res = test_model.AdvancedGenericsWriterBase.schema;
    end
  end

  methods (Abstract, Access=protected)
    read_float_image_image_(obj)
    read_generic_record_1_(obj)
    read_tuple_of_optionals_(obj)
    read_tuple_of_optionals_alternate_syntax_(obj)
    read_tuple_of_vectors_(obj)

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
        name = 'read_float_image_image';
      elseif state == 1
        name = 'read_generic_record_1';
      elseif state == 2
        name = 'read_tuple_of_optionals';
      elseif state == 3
        name = 'read_tuple_of_optionals_alternate_syntax';
      elseif state == 4
        name = 'read_tuple_of_vectors';
      else
        name = '<unknown>';
      end
    end
  end
end
