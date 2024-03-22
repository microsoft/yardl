% Binary reader for the Scalars protocol
classdef ScalarsReader < yardl.binary.BinaryProtocolReader & test_model.ScalarsReaderBase
  methods
    function obj = ScalarsReader(filename)
      obj@test_model.ScalarsReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.ScalarsReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_int32_(obj)
      r = yardl.binary.Int32Serializer;
      value = r.read(obj.stream_);
    end

    function value = read_record_(obj)
      r = test_model.binary.RecordWithPrimitivesSerializer();
      value = r.read(obj.stream_);
    end
  end
end
