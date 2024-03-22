% Abstract writer for protocol Aliases
classdef (Abstract) AliasesWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = AliasesWriterBase()
      obj.state_ = 0;
    end

    function close(obj)
      if obj.state_ == 19
        obj.end_stream_();
        obj.close_();
        return
      end
      obj.close_();
      if obj.state_ ~= 20
        expected_method = obj.state_to_method_name_(bitand((int32(obj.state_) + 1), bitcmp(1, 'int8')));
        throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function write_aliased_string(obj, value)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      obj.write_aliased_string_(value);
      obj.state_ = 2;
    end

    % Ordinal 1
    function write_aliased_enum(obj, value)
      if obj.state_ ~= 2
        obj.raise_unexpected_state_(2);
      end

      obj.write_aliased_enum_(value);
      obj.state_ = 4;
    end

    % Ordinal 2
    function write_aliased_open_generic(obj, value)
      if obj.state_ ~= 4
        obj.raise_unexpected_state_(4);
      end

      obj.write_aliased_open_generic_(value);
      obj.state_ = 6;
    end

    % Ordinal 3
    function write_aliased_closed_generic(obj, value)
      if obj.state_ ~= 6
        obj.raise_unexpected_state_(6);
      end

      obj.write_aliased_closed_generic_(value);
      obj.state_ = 8;
    end

    % Ordinal 4
    function write_aliased_optional(obj, value)
      if obj.state_ ~= 8
        obj.raise_unexpected_state_(8);
      end

      obj.write_aliased_optional_(value);
      obj.state_ = 10;
    end

    % Ordinal 5
    function write_aliased_generic_optional(obj, value)
      if obj.state_ ~= 10
        obj.raise_unexpected_state_(10);
      end

      obj.write_aliased_generic_optional_(value);
      obj.state_ = 12;
    end

    % Ordinal 6
    function write_aliased_generic_union_2(obj, value)
      if obj.state_ ~= 12
        obj.raise_unexpected_state_(12);
      end

      obj.write_aliased_generic_union_2_(value);
      obj.state_ = 14;
    end

    % Ordinal 7
    function write_aliased_generic_vector(obj, value)
      if obj.state_ ~= 14
        obj.raise_unexpected_state_(14);
      end

      obj.write_aliased_generic_vector_(value);
      obj.state_ = 16;
    end

    % Ordinal 8
    function write_aliased_generic_fixed_vector(obj, value)
      if obj.state_ ~= 16
        obj.raise_unexpected_state_(16);
      end

      obj.write_aliased_generic_fixed_vector_(value);
      obj.state_ = 18;
    end

    % Ordinal 9
    function write_stream_of_aliased_generic_union_2(obj, value)
      if bitand(int32(obj.state_), bitcmp(1, 'int8')) ~= 18
        obj.raise_unexpected_state_(18);
      end

      obj.write_stream_of_aliased_generic_union_2_(value);
      obj.state_ = 19;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"Aliases","sequence":[{"name":"aliasedString","type":"TestModel.AliasedString"},{"name":"aliasedEnum","type":"TestModel.AliasedEnum"},{"name":"aliasedOpenGeneric","type":{"name":"TestModel.AliasedOpenGeneric","typeArguments":["TestModel.AliasedString","TestModel.AliasedEnum"]}},{"name":"aliasedClosedGeneric","type":"TestModel.AliasedClosedGeneric"},{"name":"aliasedOptional","type":"TestModel.AliasedOptional"},{"name":"aliasedGenericOptional","type":{"name":"TestModel.AliasedGenericOptional","typeArguments":["float32"]}},{"name":"aliasedGenericUnion2","type":{"name":"TestModel.AliasedGenericUnion2","typeArguments":["TestModel.AliasedString","TestModel.AliasedEnum"]}},{"name":"aliasedGenericVector","type":{"name":"TestModel.AliasedGenericVector","typeArguments":["float32"]}},{"name":"aliasedGenericFixedVector","type":{"name":"TestModel.AliasedGenericFixedVector","typeArguments":["float32"]}},{"name":"streamOfAliasedGenericUnion2","type":{"stream":{"items":{"name":"TestModel.AliasedGenericUnion2","typeArguments":["TestModel.AliasedString","TestModel.AliasedEnum"]}}}}]},"types":[{"name":"Fruits","values":[{"symbol":"apple","value":0},{"symbol":"banana","value":1},{"symbol":"pear","value":2}]},{"name":"GenericUnion2","typeParameters":["T1","T2"],"type":[{"tag":"T1","type":"T1"},{"tag":"T2","type":"T2"}]},{"name":"GenericVector","typeParameters":["T"],"type":{"vector":{"items":"T"}}},{"name":"MyTuple","typeParameters":["T1","T2"],"type":{"name":"Tuples.Tuple","typeArguments":["T1","T2"]}},{"name":"AliasedClosedGeneric","type":{"name":"TestModel.AliasedTuple","typeArguments":["TestModel.AliasedString","TestModel.AliasedEnum"]}},{"name":"AliasedEnum","type":"TestModel.Fruits"},{"name":"AliasedGenericFixedVector","typeParameters":["T"],"type":{"vector":{"items":"T","length":3}}},{"name":"AliasedGenericOptional","typeParameters":["T"],"type":[null,"T"]},{"name":"AliasedGenericUnion2","typeParameters":["T1","T2"],"type":{"name":"BasicTypes.GenericUnion2","typeArguments":["T1","T2"]}},{"name":"AliasedGenericVector","typeParameters":["T"],"type":{"name":"BasicTypes.GenericVector","typeArguments":["T"]}},{"name":"AliasedOpenGeneric","typeParameters":["T1","T2"],"type":{"name":"TestModel.AliasedTuple","typeArguments":["T1","T2"]}},{"name":"AliasedOptional","type":[null,"int32"]},{"name":"AliasedString","type":"string"},{"name":"AliasedTuple","typeParameters":["T1","T2"],"type":{"name":"TestModel.MyTuple","typeArguments":["T1","T2"]}},{"name":"Fruits","type":"BasicTypes.Fruits"},{"name":"MyTuple","typeParameters":["T1","T2"],"type":{"name":"BasicTypes.MyTuple","typeArguments":["T1","T2"]}},{"name":"Tuple","typeParameters":["T1","T2"],"fields":[{"name":"v1","type":"T1"},{"name":"v2","type":"T2"}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_aliased_string_(obj, value)
    write_aliased_enum_(obj, value)
    write_aliased_open_generic_(obj, value)
    write_aliased_closed_generic_(obj, value)
    write_aliased_optional_(obj, value)
    write_aliased_generic_optional_(obj, value)
    write_aliased_generic_union_2_(obj, value)
    write_aliased_generic_vector_(obj, value)
    write_aliased_generic_fixed_vector_(obj, value)
    write_stream_of_aliased_generic_union_2_(obj, value)

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
        name = 'write_aliased_string';
      elseif state == 2
        name = 'write_aliased_enum';
      elseif state == 4
        name = 'write_aliased_open_generic';
      elseif state == 6
        name = 'write_aliased_closed_generic';
      elseif state == 8
        name = 'write_aliased_optional';
      elseif state == 10
        name = 'write_aliased_generic_optional';
      elseif state == 12
        name = 'write_aliased_generic_union_2';
      elseif state == 14
        name = 'write_aliased_generic_vector';
      elseif state == 16
        name = 'write_aliased_generic_fixed_vector';
      elseif state == 18
        name = 'write_stream_of_aliased_generic_union_2';
      else
        name = '<unknown>';
      end
    end
  end
end