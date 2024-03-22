% Binary writer for the ProtocolWithKeywordSteps protocol
classdef ProtocolWithKeywordStepsWriter < yardl.binary.BinaryProtocolWriter & test_model.ProtocolWithKeywordStepsWriterBase
  methods
    function obj = ProtocolWithKeywordStepsWriter(filename)
      obj@test_model.ProtocolWithKeywordStepsWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.ProtocolWithKeywordStepsWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_int_(obj, value)
      w = yardl.binary.StreamSerializer(test_model.binary.RecordWithKeywordFieldsSerializer());
      w.write(obj.stream_, value);
    end

    function write_float_(obj, value)
      w = yardl.binary.EnumSerializer('test_model.EnumWithKeywordSymbols', @test_model.EnumWithKeywordSymbols, yardl.binary.Int32Serializer);
      w.write(obj.stream_, value);
    end
  end
end
