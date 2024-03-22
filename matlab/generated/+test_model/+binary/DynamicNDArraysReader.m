% Binary reader for the DynamicNDArrays protocol
classdef DynamicNDArraysReader < yardl.binary.BinaryProtocolReader & test_model.DynamicNDArraysReaderBase
  methods
    function obj = DynamicNDArraysReader(filename)
      obj@test_model.DynamicNDArraysReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.DynamicNDArraysReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_ints_(obj)
      r = yardl.binary.DynamicNDArraySerializer(yardl.binary.Int32Serializer);
      value = r.read(obj.stream_);
    end

    function value = read_simple_record_array_(obj)
      r = yardl.binary.DynamicNDArraySerializer(test_model.binary.SimpleRecordSerializer());
      value = r.read(obj.stream_);
    end

    function value = read_record_with_vlens_array_(obj)
      r = yardl.binary.DynamicNDArraySerializer(test_model.binary.RecordWithVlensSerializer());
      value = r.read(obj.stream_);
    end

    function value = read_record_with_dynamic_nd_arrays_(obj)
      r = test_model.binary.RecordWithDynamicNDArraysSerializer();
      value = r.read(obj.stream_);
    end
  end
end