% Binary writer for the FixedArrays protocol
classdef FixedArraysWriter < yardl.binary.BinaryProtocolWriter & test_model.FixedArraysWriterBase
  methods
    function obj = FixedArraysWriter(filename)
      obj@test_model.FixedArraysWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.FixedArraysWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_ints_(obj, value)
      w = yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [3, 2]);
      w.write(obj.stream_, value);
    end

    function write_fixed_simple_record_array_(obj, value)
      w = yardl.binary.FixedNDArraySerializer(test_model.binary.SimpleRecordSerializer(), [2, 3]);
      w.write(obj.stream_, value);
    end

    function write_fixed_record_with_vlens_array_(obj, value)
      w = yardl.binary.FixedNDArraySerializer(test_model.binary.RecordWithVlensSerializer(), [2, 2]);
      w.write(obj.stream_, value);
    end

    function write_record_with_fixed_arrays_(obj, value)
      w = test_model.binary.RecordWithFixedArraysSerializer();
      w.write(obj.stream_, value);
    end

    function write_named_array_(obj, value)
      w = yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [4, 2]);
      w.write(obj.stream_, value);
    end
  end
end
