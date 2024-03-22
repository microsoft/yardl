% Binary writer for the Subarrays protocol
classdef SubarraysWriter < yardl.binary.BinaryProtocolWriter & test_model.SubarraysWriterBase
  methods
    function obj = SubarraysWriter(filename)
      obj@test_model.SubarraysWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.SubarraysWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_dynamic_with_fixed_int_subarray_(obj, value)
      w = yardl.binary.DynamicNDArraySerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [3]));
      w.write(obj.stream_, value);
    end

    function write_dynamic_with_fixed_float_subarray_(obj, value)
      w = yardl.binary.DynamicNDArraySerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Float32Serializer, [3]));
      w.write(obj.stream_, value);
    end

    function write_known_dim_count_with_fixed_int_subarray_(obj, value)
      w = yardl.binary.NDArraySerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [3]), 1);
      w.write(obj.stream_, value);
    end

    function write_known_dim_count_with_fixed_float_subarray_(obj, value)
      w = yardl.binary.NDArraySerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Float32Serializer, [3]), 1);
      w.write(obj.stream_, value);
    end

    function write_fixed_with_fixed_int_subarray_(obj, value)
      w = yardl.binary.FixedNDArraySerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [3]), [2]);
      w.write(obj.stream_, value);
    end

    function write_fixed_with_fixed_float_subarray_(obj, value)
      w = yardl.binary.FixedNDArraySerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Float32Serializer, [3]), [2]);
      w.write(obj.stream_, value);
    end

    function write_nested_subarray_(obj, value)
      w = yardl.binary.DynamicNDArraySerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [3]), [2]));
      w.write(obj.stream_, value);
    end

    function write_dynamic_with_fixed_vector_subarray_(obj, value)
      w = yardl.binary.DynamicNDArraySerializer(yardl.binary.FixedVectorSerializer(yardl.binary.Int32Serializer, 3));
      w.write(obj.stream_, value);
    end

    function write_generic_subarray_(obj, value)
      w = yardl.binary.NDArraySerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [3]), 2);
      w.write(obj.stream_, value);
    end
  end
end
