% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef HelloWorldReader < yardl.binary.BinaryProtocolReader & sandbox.HelloWorldReaderBase
  % Binary reader for the HelloWorld protocol
  properties (Access=protected)
    data_serializer
  end

  methods
    function self = HelloWorldReader(filename)
      self@sandbox.HelloWorldReaderBase();
      self@yardl.binary.BinaryProtocolReader(filename, sandbox.HelloWorldReaderBase.schema);
      self.data_serializer = yardl.binary.StreamSerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Complexfloat64Serializer, [2]));
    end
  end

  methods (Access=protected)
    function more = has_data_(self)
      more = self.data_serializer.hasnext(self.stream_);
    end

    function value = read_data_(self)
      value = self.data_serializer.read(self.stream_);
    end
  end
end
