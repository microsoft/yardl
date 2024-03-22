% Binary writer for the FixedVectors protocol
classdef FixedVectorsWriter < yardl.binary.BinaryProtocolWriter & test_model.FixedVectorsWriterBase
  methods
    function obj = FixedVectorsWriter(filename)
      obj@test_model.FixedVectorsWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.FixedVectorsWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_fixed_int_vector_(obj, value)
      w = yardl.binary.FixedVectorSerializer(yardl.binary.Int32Serializer, 5);
      w.write(obj.stream_, value);
    end

    function write_fixed_simple_record_vector_(obj, value)
      w = yardl.binary.FixedVectorSerializer(test_model.binary.SimpleRecordSerializer(), 3);
      w.write(obj.stream_, value);
    end

    function write_fixed_record_with_vlens_vector_(obj, value)
      w = yardl.binary.FixedVectorSerializer(test_model.binary.RecordWithVlensSerializer(), 2);
      w.write(obj.stream_, value);
    end

    function write_record_with_fixed_vectors_(obj, value)
      w = test_model.binary.RecordWithFixedVectorsSerializer();
      w.write(obj.stream_, value);
    end
  end
end
