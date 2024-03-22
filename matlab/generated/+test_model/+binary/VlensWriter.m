% Binary writer for the Vlens protocol
classdef VlensWriter < yardl.binary.BinaryProtocolWriter & test_model.VlensWriterBase
  methods
    function obj = VlensWriter(filename)
      obj@test_model.VlensWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.VlensWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_int_vector_(obj, value)
      w = yardl.binary.VectorSerializer(yardl.binary.Int32Serializer);
      w.write(obj.stream_, value);
    end

    function write_complex_vector_(obj, value)
      w = yardl.binary.VectorSerializer(yardl.binary.Complexfloat32Serializer);
      w.write(obj.stream_, value);
    end

    function write_record_with_vlens_(obj, value)
      w = test_model.binary.RecordWithVlensSerializer();
      w.write(obj.stream_, value);
    end

    function write_vlen_of_record_with_vlens_(obj, value)
      w = yardl.binary.VectorSerializer(test_model.binary.RecordWithVlensSerializer());
      w.write(obj.stream_, value);
    end
  end
end
