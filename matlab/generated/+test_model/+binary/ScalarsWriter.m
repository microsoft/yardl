% Binary writer for the Scalars protocol
classdef ScalarsWriter < yardl.binary.BinaryProtocolWriter & test_model.ScalarsWriterBase
  methods
    function obj = ScalarsWriter(filename)
      obj@test_model.ScalarsWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.ScalarsWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_int32_(obj, value)
      w = yardl.binary.Int32Serializer;
      w.write(obj.stream_, value);
    end

    function write_record_(obj, value)
      w = test_model.binary.RecordWithPrimitivesSerializer();
      w.write(obj.stream_, value);
    end
  end
end