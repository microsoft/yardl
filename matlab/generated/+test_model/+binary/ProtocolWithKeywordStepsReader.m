% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef ProtocolWithKeywordStepsReader < yardl.binary.BinaryProtocolReader & test_model.ProtocolWithKeywordStepsReaderBase
  % Binary reader for the ProtocolWithKeywordSteps protocol
  methods
    function obj = ProtocolWithKeywordStepsReader(filename)
      obj@test_model.ProtocolWithKeywordStepsReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.ProtocolWithKeywordStepsReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_int_(obj)
      r = yardl.binary.StreamSerializer(test_model.binary.RecordWithKeywordFieldsSerializer());
      value = r.read(obj.stream_);
    end

    function value = read_float_(obj)
      r = yardl.binary.EnumSerializer('test_model.EnumWithKeywordSymbols', @test_model.EnumWithKeywordSymbols, yardl.binary.Int32Serializer);
      value = r.read(obj.stream_);
    end
  end
end
