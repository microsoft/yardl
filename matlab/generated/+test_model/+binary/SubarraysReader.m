% Binary reader for the Subarrays protocol
classdef SubarraysReader < yardl.binary.BinaryProtocolReader & test_model.SubarraysReaderBase
  methods
    function obj = SubarraysReader(filename)
      obj@test_model.SubarraysReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.SubarraysReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_dynamic_with_fixed_int_subarray_(obj)
      r = yardl.binary.DynamicNDArraySerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [3]));
      value = r.read(obj.stream_);
    end

    function value = read_dynamic_with_fixed_float_subarray_(obj)
      r = yardl.binary.DynamicNDArraySerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Float32Serializer, [3]));
      value = r.read(obj.stream_);
    end

    function value = read_known_dim_count_with_fixed_int_subarray_(obj)
      r = yardl.binary.NDArraySerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [3]), 1);
      value = r.read(obj.stream_);
    end

    function value = read_known_dim_count_with_fixed_float_subarray_(obj)
      r = yardl.binary.NDArraySerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Float32Serializer, [3]), 1);
      value = r.read(obj.stream_);
    end

    function value = read_fixed_with_fixed_int_subarray_(obj)
      r = yardl.binary.FixedNDArraySerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [3]), [2]);
      value = r.read(obj.stream_);
    end

    function value = read_fixed_with_fixed_float_subarray_(obj)
      r = yardl.binary.FixedNDArraySerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Float32Serializer, [3]), [2]);
      value = r.read(obj.stream_);
    end

    function value = read_nested_subarray_(obj)
      r = yardl.binary.DynamicNDArraySerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [3]), [2]));
      value = r.read(obj.stream_);
    end

    function value = read_dynamic_with_fixed_vector_subarray_(obj)
      r = yardl.binary.DynamicNDArraySerializer(yardl.binary.FixedVectorSerializer(yardl.binary.Int32Serializer, 3));
      value = r.read(obj.stream_);
    end

    function value = read_generic_subarray_(obj)
      r = yardl.binary.NDArraySerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [3]), 2);
      value = r.read(obj.stream_);
    end
  end
end
