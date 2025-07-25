% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef ComplexArraysReader < yardl.binary.BinaryProtocolReader & test_model.ComplexArraysReaderBase
  % Binary reader for the ComplexArrays protocol
  properties (Access=protected)
    floats_serializer
    doubles_serializer
  end

  methods
    function self = ComplexArraysReader(filename, options)
      arguments
        filename (1,1) string
        options.skip_completed_check (1,1) logical = false
      end
      self@test_model.ComplexArraysReaderBase(skip_completed_check=options.skip_completed_check);
      self@yardl.binary.BinaryProtocolReader(filename, test_model.ComplexArraysReaderBase.schema);
      self.floats_serializer = yardl.binary.DynamicNDArraySerializer(yardl.binary.Complexfloat32Serializer);
      self.doubles_serializer = yardl.binary.NDArraySerializer(yardl.binary.Complexfloat64Serializer, 2);
    end
  end

  methods (Access=protected)
    function value = read_floats_(self)
      value = self.floats_serializer.read(self.stream_);
    end

    function value = read_doubles_(self)
      value = self.doubles_serializer.read(self.stream_);
    end
  end
end
