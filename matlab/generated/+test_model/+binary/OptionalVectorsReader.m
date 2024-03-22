% Binary reader for the OptionalVectors protocol
classdef OptionalVectorsReader < yardl.binary.BinaryProtocolReader & test_model.OptionalVectorsReaderBase
  methods
    function obj = OptionalVectorsReader(filename)
      obj@test_model.OptionalVectorsReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.OptionalVectorsReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_record_with_optional_vector_(obj)
      r = test_model.binary.RecordWithOptionalVectorSerializer();
      value = r.read(obj.stream_);
    end
  end
end
