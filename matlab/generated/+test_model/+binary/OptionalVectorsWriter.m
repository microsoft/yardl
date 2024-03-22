% Binary writer for the OptionalVectors protocol
classdef OptionalVectorsWriter < yardl.binary.BinaryProtocolWriter & test_model.OptionalVectorsWriterBase
  methods
    function obj = OptionalVectorsWriter(filename)
      obj@test_model.OptionalVectorsWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.OptionalVectorsWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_record_with_optional_vector_(obj, value)
      w = test_model.binary.RecordWithOptionalVectorSerializer();
      w.write(obj.stream_, value);
    end
  end
end
