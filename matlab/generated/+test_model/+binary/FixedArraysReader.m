% Binary reader for the FixedArrays protocol
classdef FixedArraysReader < yardl.binary.BinaryProtocolReader & test_model.FixedArraysReaderBase
  methods
    function obj = FixedArraysReader(filename)
      obj@test_model.FixedArraysReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.FixedArraysReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_ints_(obj)
      r = yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [3, 2]);
      value = r.read(obj.stream_);
    end

    function value = read_fixed_simple_record_array_(obj)
      r = yardl.binary.FixedNDArraySerializer(test_model.binary.SimpleRecordSerializer(), [2, 3]);
      value = r.read(obj.stream_);
    end

    function value = read_fixed_record_with_vlens_array_(obj)
      r = yardl.binary.FixedNDArraySerializer(test_model.binary.RecordWithVlensSerializer(), [2, 2]);
      value = r.read(obj.stream_);
    end

    function value = read_record_with_fixed_arrays_(obj)
      r = test_model.binary.RecordWithFixedArraysSerializer();
      value = r.read(obj.stream_);
    end

    function value = read_named_array_(obj)
      r = yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [4, 2]);
      value = r.read(obj.stream_);
    end
  end
end