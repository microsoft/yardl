% Abstract writer for protocol SimpleGenerics
classdef (Abstract) SimpleGenericsWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = SimpleGenericsWriterBase()
      obj.state_ = 0;
    end

    function close(obj)
      if obj.state_ == 17
        obj.end_stream_();
        obj.close_();
        return
      end
      obj.close_();
      if obj.state_ ~= 18
        expected_method = obj.state_to_method_name_(bitand((int32(obj.state_) + 1), bitcmp(1, 'int8')));
        throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function write_float_image(obj, value)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      obj.write_float_image_(value);
      obj.state_ = 2;
    end

    % Ordinal 1
    function write_int_image(obj, value)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      obj.write_int_image_(value);
      obj.state_ = 4;
    end

    % Ordinal 2
    function write_int_image_alternate_syntax(obj, value)
      if obj.state_ ~= 4
        obj.raise_unexpected_state_(4);
      end

      obj.write_int_image_alternate_syntax_(value);
      obj.state_ = 6;
    end

    % Ordinal 3
    function write_string_image(obj, value)
      if obj.state_ ~= 6
        obj.raise_unexpected_state_(6);
      end

      obj.write_string_image_(value);
      obj.state_ = 8;
    end

    % Ordinal 4
    function write_int_float_tuple(obj, value)
      if obj.state_ ~= 8
        obj.raise_unexpected_state_(8);
      end

      obj.write_int_float_tuple_(value);
      obj.state_ = 10;
    end

    % Ordinal 5
    function write_float_float_tuple(obj, value)
      if obj.state_ ~= 10
        obj.raise_unexpected_state_(10);
      end

      obj.write_float_float_tuple_(value);
      obj.state_ = 12;
    end

    % Ordinal 6
    function write_int_float_tuple_alternate_syntax(obj, value)
      if obj.state_ ~= 12
        obj.raise_unexpected_state_(12);
      end

      obj.write_int_float_tuple_alternate_syntax_(value);
      obj.state_ = 14;
    end

    % Ordinal 7
    function write_int_string_tuple(obj, value)
      if obj.state_ ~= 14
        obj.raise_unexpected_state_(14);
      end

      obj.write_int_string_tuple_(value);
      obj.state_ = 16;
    end

    % Ordinal 8
    function write_stream_of_type_variants(obj, value)
      if bitand(int32(obj.state_), bitcmp(1, 'int8')) ~= 16
        obj.raise_unexpected_state_(16);
      end

      obj.write_stream_of_type_variants_(value);
      obj.state_ = 17;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"SimpleGenerics","sequence":[{"name":"floatImage","type":"Image.FloatImage"},{"name":"intImage","type":"Image.IntImage"},{"name":"intImageAlternateSyntax","type":{"name":"TestModel.Image","typeArguments":["int32"]}},{"name":"stringImage","type":{"name":"TestModel.Image","typeArguments":["string"]}},{"name":"intFloatTuple","type":{"name":"Tuples.Tuple","typeArguments":["int32","float32"]}},{"name":"floatFloatTuple","type":{"name":"Tuples.Tuple","typeArguments":["float32","float32"]}},{"name":"intFloatTupleAlternateSyntax","type":{"name":"Tuples.Tuple","typeArguments":["int32","float32"]}},{"name":"intStringTuple","type":{"name":"Tuples.Tuple","typeArguments":["int32","string"]}},{"name":"streamOfTypeVariants","type":{"stream":{"items":[{"tag":"imageFloat","explicitTag":true,"type":"Image.FloatImage"},{"tag":"imageDouble","explicitTag":true,"type":{"name":"TestModel.Image","typeArguments":["float64"]}}]}}}]},"types":[{"name":"FloatImage","type":{"name":"Image.Image","typeArguments":["float32"]}},{"name":"Image","typeParameters":["T"],"type":{"array":{"items":"T","dimensions":[{"name":"x"},{"name":"y"}]}}},{"name":"IntImage","type":{"name":"Image.Image","typeArguments":["int32"]}},{"name":"Image","typeParameters":["T"],"type":{"name":"Image.Image","typeArguments":["T"]}},{"name":"Tuple","typeParameters":["T1","T2"],"fields":[{"name":"v1","type":"T1"},{"name":"v2","type":"T2"}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_float_image_(obj, value)
    write_int_image_(obj, value)
    write_int_image_alternate_syntax_(obj, value)
    write_string_image_(obj, value)
    write_int_float_tuple_(obj, value)
    write_float_float_tuple_(obj, value)
    write_int_float_tuple_alternate_syntax_(obj, value)
    write_int_string_tuple_(obj, value)
    write_stream_of_type_variants_(obj, value)

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
        name = 'write_float_image';
      elseif state == 2
        name = 'write_int_image';
      elseif state == 4
        name = 'write_int_image_alternate_syntax';
      elseif state == 6
        name = 'write_string_image';
      elseif state == 8
        name = 'write_int_float_tuple';
      elseif state == 10
        name = 'write_float_float_tuple';
      elseif state == 12
        name = 'write_int_float_tuple_alternate_syntax';
      elseif state == 14
        name = 'write_int_string_tuple';
      elseif state == 16
        name = 'write_stream_of_type_variants';
      else
        name = '<unknown>';
      end
    end
  end
end
