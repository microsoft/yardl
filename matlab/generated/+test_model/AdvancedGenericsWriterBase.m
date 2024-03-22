% Abstract writer for protocol AdvancedGenerics
classdef (Abstract) AdvancedGenericsWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = AdvancedGenericsWriterBase()
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
    function write_float_image_image(obj, value)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      obj.write_float_image_image_(value);
      obj.state_ = 2;
    end

    % Ordinal 1
    function write_generic_record_1(obj, value)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      obj.write_generic_record_1_(value);
      obj.state_ = 4;
    end

    % Ordinal 2
    function write_tuple_of_optionals(obj, value)
      if obj.state_ ~= 4
        obj.raise_unexpected_state_(4);
      end

      obj.write_tuple_of_optionals_(value);
      obj.state_ = 6;
    end

    % Ordinal 3
    function write_tuple_of_optionals_alternate_syntax(obj, value)
      if obj.state_ ~= 6
        obj.raise_unexpected_state_(6);
      end

      obj.write_tuple_of_optionals_alternate_syntax_(value);
      obj.state_ = 8;
    end

    % Ordinal 4
    function write_tuple_of_vectors(obj, value)
      if obj.state_ ~= 8
        obj.raise_unexpected_state_(8);
      end

      obj.write_tuple_of_vectors_(value);
      obj.state_ = 10;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"AdvancedGenerics","sequence":[{"name":"floatImageImage","type":{"name":"TestModel.Image","typeArguments":[{"name":"TestModel.Image","typeArguments":["float32"]}]}},{"name":"genericRecord1","type":{"name":"TestModel.GenericRecord","typeArguments":["int32","string"]}},{"name":"tupleOfOptionals","type":{"name":"TestModel.MyTuple","typeArguments":[[null,"int32"],[null,"string"]]}},{"name":"tupleOfOptionalsAlternateSyntax","type":{"name":"TestModel.MyTuple","typeArguments":[[null,"int32"],[null,"string"]]}},{"name":"tupleOfVectors","type":{"name":"TestModel.MyTuple","typeArguments":[{"vector":{"items":"int32"}},{"vector":{"items":"float32"}}]}}]},"types":[{"name":"MyTuple","typeParameters":["T1","T2"],"type":{"name":"Tuples.Tuple","typeArguments":["T1","T2"]}},{"name":"Image","typeParameters":["T"],"type":{"array":{"items":"T","dimensions":[{"name":"x"},{"name":"y"}]}}},{"name":"GenericRecord","typeParameters":["T1","T2"],"fields":[{"name":"scalar1","type":"T1"},{"name":"scalar2","type":"T2"},{"name":"vector1","type":{"vector":{"items":"T1"}}},{"name":"image2","type":{"name":"TestModel.Image","typeArguments":["T2"]}}]},{"name":"Image","typeParameters":["T"],"type":{"name":"Image.Image","typeArguments":["T"]}},{"name":"MyTuple","typeParameters":["T1","T2"],"type":{"name":"BasicTypes.MyTuple","typeArguments":["T1","T2"]}},{"name":"Tuple","typeParameters":["T1","T2"],"fields":[{"name":"v1","type":"T1"},{"name":"v2","type":"T2"}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_float_image_image_(obj, value)
    write_generic_record_1_(obj, value)
    write_tuple_of_optionals_(obj, value)
    write_tuple_of_optionals_alternate_syntax_(obj, value)
    write_tuple_of_vectors_(obj, value)

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
        name = 'write_float_image_image';
      elseif state == 2
        name = 'write_generic_record_1';
      elseif state == 4
        name = 'write_tuple_of_optionals';
      elseif state == 6
        name = 'write_tuple_of_optionals_alternate_syntax';
      elseif state == 8
        name = 'write_tuple_of_vectors';
      else
        name = '<unknown>';
      end
    end
  end
end
