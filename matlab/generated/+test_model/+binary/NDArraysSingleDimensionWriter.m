% Binary writer for the NDArraysSingleDimension protocol
classdef NDArraysSingleDimensionWriter < yardl.binary.BinaryProtocolWriter & test_model.NDArraysSingleDimensionWriterBase
  methods
    function obj = NDArraysSingleDimensionWriter(filename)
      obj@test_model.NDArraysSingleDimensionWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.NDArraysSingleDimensionWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_ints_(obj, value)
      w = yardl.binary.NDArraySerializer(yardl.binary.Int32Serializer, 1);
      w.write(obj.stream_, value);
    end

    function write_simple_record_array_(obj, value)
      w = yardl.binary.NDArraySerializer(test_model.binary.SimpleRecordSerializer(), 1);
      w.write(obj.stream_, value);
    end

    function write_record_with_vlens_array_(obj, value)
      w = yardl.binary.NDArraySerializer(test_model.binary.RecordWithVlensSerializer(), 1);
      w.write(obj.stream_, value);
    end

    function write_record_with_nd_arrays_(obj, value)
      w = test_model.binary.RecordWithNDArraysSingleDimensionSerializer();
      w.write(obj.stream_, value);
    end
  end
end
