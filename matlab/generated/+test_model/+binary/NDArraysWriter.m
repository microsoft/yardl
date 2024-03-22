% Binary writer for the NDArrays protocol
classdef NDArraysWriter < yardl.binary.BinaryProtocolWriter & test_model.NDArraysWriterBase
  methods
    function obj = NDArraysWriter(filename)
      obj@test_model.NDArraysWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.NDArraysWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_ints_(obj, value)
      w = yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 2);
      w.write(obj.stream_, value);
    end

    function write_simple_record_array_(obj, value)
      w = yardl.binary.NDArraySerializer(test_model.binary.SimpleRecordSerializer(), 2);
      w.write(obj.stream_, value);
    end

    function write_record_with_vlens_array_(obj, value)
      w = yardl.binary.NDArraySerializer(test_model.binary.RecordWithVlensSerializer(), 2);
      w.write(obj.stream_, value);
    end

    function write_record_with_nd_arrays_(obj, value)
      w = test_model.binary.RecordWithNDArraysSerializer();
      w.write(obj.stream_, value);
    end

    function write_named_array_(obj, value)
      w = yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 2);
      w.write(obj.stream_, value);
    end
  end
end
