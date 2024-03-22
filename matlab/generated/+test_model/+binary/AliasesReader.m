% Binary reader for the Aliases protocol
classdef AliasesReader < yardl.binary.BinaryProtocolReader & test_model.AliasesReaderBase
  methods
    function obj = AliasesReader(filename)
      obj@test_model.AliasesReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.AliasesReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_aliased_string_(obj)
      r = yardl.binary.StringSerializer;
      value = r.read(obj.stream_);
    end

    function value = read_aliased_enum_(obj)
      r = yardl.binary.EnumSerializer('basic_types.Fruits', @basic_types.Fruits, yardl.binary.Int32Serializer);
      value = r.read(obj.stream_);
    end

    function value = read_aliased_open_generic_(obj)
      r = tuples.binary.TupleSerializer(yardl.binary.StringSerializer, yardl.binary.EnumSerializer('basic_types.Fruits', @basic_types.Fruits, yardl.binary.Int32Serializer));
      value = r.read(obj.stream_);
    end

    function value = read_aliased_closed_generic_(obj)
      r = tuples.binary.TupleSerializer(yardl.binary.StringSerializer, yardl.binary.EnumSerializer('basic_types.Fruits', @basic_types.Fruits, yardl.binary.Int32Serializer));
      value = r.read(obj.stream_);
    end

    function value = read_aliased_optional_(obj)
      r = yardl.binary.OptionalSerializer(yardl.binary.Int32Serializer);
      value = r.read(obj.stream_);
    end

    function value = read_aliased_generic_optional_(obj)
      r = yardl.binary.OptionalSerializer(yardl.binary.Float32Serializer);
      value = r.read(obj.stream_);
    end

    function value = read_aliased_generic_union_2_(obj)
      r = yardl.binary.UnionSerializer('basic_types.GenericUnion2', {yardl.binary.StringSerializer, yardl.binary.EnumSerializer('basic_types.Fruits', @basic_types.Fruits, yardl.binary.Int32Serializer)}, {@basic_types.GenericUnion2.T1, @basic_types.GenericUnion2.T2});
      value = r.read(obj.stream_);
    end

    function value = read_aliased_generic_vector_(obj)
      r = yardl.binary.VectorSerializer(yardl.binary.Float32Serializer);
      value = r.read(obj.stream_);
    end

    function value = read_aliased_generic_fixed_vector_(obj)
      r = yardl.binary.FixedVectorSerializer(yardl.binary.Float32Serializer, 3);
      value = r.read(obj.stream_);
    end

    function value = read_stream_of_aliased_generic_union_2_(obj)
      r = yardl.binary.StreamSerializer(yardl.binary.UnionSerializer('basic_types.GenericUnion2', {yardl.binary.StringSerializer, yardl.binary.EnumSerializer('basic_types.Fruits', @basic_types.Fruits, yardl.binary.Int32Serializer)}, {@basic_types.GenericUnion2.T1, @basic_types.GenericUnion2.T2}));
      value = r.read(obj.stream_);
    end
  end
end
