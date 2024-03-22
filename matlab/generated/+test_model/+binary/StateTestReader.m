% Binary reader for the StateTest protocol
classdef StateTestReader < yardl.binary.BinaryProtocolReader & test_model.StateTestReaderBase
  methods
    function obj = StateTestReader(filename)
      obj@test_model.StateTestReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.StateTestReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_an_int_(obj)
      r = yardl.binary.Int32Serializer;
      value = r.read(obj.stream_);
    end

    function value = read_a_stream_(obj)
      r = yardl.binary.StreamSerializer(yardl.binary.Int32Serializer);
      value = r.read(obj.stream_);
    end

    function value = read_another_int_(obj)
      r = yardl.binary.Int32Serializer;
      value = r.read(obj.stream_);
    end
  end
end
