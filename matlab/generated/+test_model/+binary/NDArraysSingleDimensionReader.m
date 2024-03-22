% Binary reader for the NDArraysSingleDimension protocol
classdef NDArraysSingleDimensionReader < yardl.binary.BinaryProtocolReader & test_model.NDArraysSingleDimensionReaderBase
  methods
    function obj = NDArraysSingleDimensionReader(filename)
      obj@test_model.NDArraysSingleDimensionReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.NDArraysSingleDimensionReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_ints_(obj)
      r = yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 1);
      value = r.read(obj.stream_);
    end

    function value = read_simple_record_array_(obj)
      r = yardl.binary.NDArraySerializer(test_model.binary.SimpleRecordSerializer(), 1);
      value = r.read(obj.stream_);
    end

    function value = read_record_with_vlens_array_(obj)
      r = yardl.binary.NDArraySerializer(test_model.binary.RecordWithVlensSerializer(), 1);
      value = r.read(obj.stream_);
    end

    function value = read_record_with_nd_arrays_(obj)
      r = test_model.binary.RecordWithNDArraysSingleDimensionSerializer();
      value = r.read(obj.stream_);
    end
  end
end
