% Binary reader for the NDArrays protocol
classdef NDArraysReader < yardl.binary.BinaryProtocolReader & test_model.NDArraysReaderBase
  methods
    function obj = NDArraysReader(filename)
      obj@test_model.NDArraysReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.NDArraysReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_ints_(obj)
      r = yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 2);
      value = r.read(obj.stream_);
    end

    function value = read_simple_record_array_(obj)
      r = yardl.binary.NDArraySerializer(test_model.binary.SimpleRecordSerializer(), 2);
      value = r.read(obj.stream_);
    end

    function value = read_record_with_vlens_array_(obj)
      r = yardl.binary.NDArraySerializer(test_model.binary.RecordWithVlensSerializer(), 2);
      value = r.read(obj.stream_);
    end

    function value = read_record_with_nd_arrays_(obj)
      r = test_model.binary.RecordWithNDArraysSerializer();
      value = r.read(obj.stream_);
    end

    function value = read_named_array_(obj)
      r = yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 2);
      value = r.read(obj.stream_);
    end
  end
end