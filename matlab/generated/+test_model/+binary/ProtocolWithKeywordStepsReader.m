% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef ProtocolWithKeywordStepsReader < yardl.binary.BinaryProtocolReader & test_model.ProtocolWithKeywordStepsReaderBase
  % Binary reader for the ProtocolWithKeywordSteps protocol
  properties (Access=protected)
    int_serializer
    float_serializer
  end

  methods
    function self = ProtocolWithKeywordStepsReader(filename, options)
      arguments
        filename (1,1) string
        options.skip_completed_check (1,1) logical = false
      end
      self@test_model.ProtocolWithKeywordStepsReaderBase(skip_completed_check=options.skip_completed_check);
      self@yardl.binary.BinaryProtocolReader(filename, test_model.ProtocolWithKeywordStepsReaderBase.schema);
      self.int_serializer = yardl.binary.StreamSerializer(test_model.binary.RecordWithKeywordFieldsSerializer());
      self.float_serializer = yardl.binary.EnumSerializer('test_model.EnumWithKeywordSymbols', @test_model.EnumWithKeywordSymbols, yardl.binary.Int32Serializer);
    end
  end

  methods (Access=protected)
    function more = has_int_(self)
      more = self.int_serializer.hasnext(self.stream_);
    end

    function value = read_int_(self)
      value = self.int_serializer.read(self.stream_);
    end

    function value = read_float_(self)
      value = self.float_serializer.read(self.stream_);
    end
  end
end
