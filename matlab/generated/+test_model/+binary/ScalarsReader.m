% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef ScalarsReader < yardl.binary.BinaryProtocolReader & test_model.ScalarsReaderBase
  % Binary reader for the Scalars protocol
  properties (Access=protected)
    int32_serializer
    record_serializer
  end

  methods
    function self = ScalarsReader(filename, options)
      arguments
        filename (1,1) string
        options.skip_completed_check (1,1) logical = false
      end
      self@test_model.ScalarsReaderBase(skip_completed_check=options.skip_completed_check);
      self@yardl.binary.BinaryProtocolReader(filename, test_model.ScalarsReaderBase.schema);
      self.int32_serializer = yardl.binary.Int32Serializer;
      self.record_serializer = test_model.binary.RecordWithPrimitivesSerializer();
    end
  end

  methods (Access=protected)
    function value = read_int32_(self)
      value = self.int32_serializer.read(self.stream_);
    end

    function value = read_record_(self)
      value = self.record_serializer.read(self.stream_);
    end
  end
end
