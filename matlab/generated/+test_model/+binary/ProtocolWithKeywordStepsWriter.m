% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef ProtocolWithKeywordStepsWriter < yardl.binary.BinaryProtocolWriter & test_model.ProtocolWithKeywordStepsWriterBase
  % Binary writer for the ProtocolWithKeywordSteps protocol
  properties (Access=protected)
    int_serializer
    float_serializer
  end

  methods
    function obj = ProtocolWithKeywordStepsWriter(filename)
      obj@test_model.ProtocolWithKeywordStepsWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.ProtocolWithKeywordStepsWriterBase.schema);
      obj.int_serializer = yardl.binary.StreamSerializer(test_model.binary.RecordWithKeywordFieldsSerializer());
      obj.float_serializer = yardl.binary.EnumSerializer('test_model.EnumWithKeywordSymbols', @test_model.EnumWithKeywordSymbols, yardl.binary.Int32Serializer);
    end
  end

  methods (Access=protected)
    function write_int_(obj, value)
      obj.int_serializer.write(obj.stream_, value);
    end

    function write_float_(obj, value)
      obj.float_serializer.write(obj.stream_, value);
    end
  end
end
