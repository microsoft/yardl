% Binary writer for the Aliases protocol
classdef AliasesWriter < yardl.binary.BinaryProtocolWriter & test_model.AliasesWriterBase
  methods
    function obj = AliasesWriter(filename)
      obj@test_model.AliasesWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.AliasesWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_aliased_string_(obj, value)
      w = yardl.binary.StringSerializer;
      w.write(obj.stream_, value);
    end

    function write_aliased_enum_(obj, value)
      w = yardl.binary.EnumSerializer('basic_types.Fruits', @basic_types.Fruits, yardl.binary.Int32Serializer);
      w.write(obj.stream_, value);
    end

    function write_aliased_open_generic_(obj, value)
      w = tuples.binary.TupleSerializer(yardl.binary.StringSerializer, yardl.binary.EnumSerializer('basic_types.Fruits', @basic_types.Fruits, yardl.binary.Int32Serializer));
      w.write(obj.stream_, value);
    end

    function write_aliased_closed_generic_(obj, value)
      w = tuples.binary.TupleSerializer(yardl.binary.StringSerializer, yardl.binary.EnumSerializer('basic_types.Fruits', @basic_types.Fruits, yardl.binary.Int32Serializer));
      w.write(obj.stream_, value);
    end

    function write_aliased_optional_(obj, value)
      w = yardl.binary.OptionalSerializer(yardl.binary.Int32Serializer);
      w.write(obj.stream_, value);
    end

    function write_aliased_generic_optional_(obj, value)
      w = yardl.binary.OptionalSerializer(yardl.binary.Float32Serializer);
      w.write(obj.stream_, value);
    end

    function write_aliased_generic_union_2_(obj, value)
      w = yardl.binary.UnionSerializer('basic_types.GenericUnion2', {yardl.binary.StringSerializer, yardl.binary.EnumSerializer('basic_types.Fruits', @basic_types.Fruits, yardl.binary.Int32Serializer)}, {@basic_types.GenericUnion2.T1, @basic_types.GenericUnion2.T2});
      w.write(obj.stream_, value);
    end

    function write_aliased_generic_vector_(obj, value)
      w = yardl.binary.VectorSerializer(yardl.binary.Float32Serializer);
      w.write(obj.stream_, value);
    end

    function write_aliased_generic_fixed_vector_(obj, value)
      w = yardl.binary.FixedVectorSerializer(yardl.binary.Float32Serializer, 3);
      w.write(obj.stream_, value);
    end

    function write_stream_of_aliased_generic_union_2_(obj, value)
      w = yardl.binary.StreamSerializer(yardl.binary.UnionSerializer('basic_types.GenericUnion2', {yardl.binary.StringSerializer, yardl.binary.EnumSerializer('basic_types.Fruits', @basic_types.Fruits, yardl.binary.Int32Serializer)}, {@basic_types.GenericUnion2.T1, @basic_types.GenericUnion2.T2}));
      w.write(obj.stream_, value);
    end
  end
end
