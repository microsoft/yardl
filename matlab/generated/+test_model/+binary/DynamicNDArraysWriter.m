% Binary writer for the DynamicNDArrays protocol
classdef DynamicNDArraysWriter < yardl.binary.BinaryProtocolWriter & test_model.DynamicNDArraysWriterBase
  methods
    function obj = DynamicNDArraysWriter(filename)
      obj@test_model.DynamicNDArraysWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.DynamicNDArraysWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_ints_(obj, value)
      w = yardl.binary.DynamicNDArraySerializer(yardl.binary.Int32Serializer);
      w.write(obj.stream_, value);
    end

    function write_simple_record_array_(obj, value)
      w = yardl.binary.DynamicNDArraySerializer(test_model.binary.SimpleRecordSerializer());
      w.write(obj.stream_, value);
    end

    function write_record_with_vlens_array_(obj, value)
      w = yardl.binary.DynamicNDArraySerializer(test_model.binary.RecordWithVlensSerializer());
      w.write(obj.stream_, value);
    end

    function write_record_with_dynamic_nd_arrays_(obj, value)
      w = test_model.binary.RecordWithDynamicNDArraysSerializer();
      w.write(obj.stream_, value);
    end
  end
end
