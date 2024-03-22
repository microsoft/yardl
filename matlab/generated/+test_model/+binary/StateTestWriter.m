% Binary writer for the StateTest protocol
classdef StateTestWriter < yardl.binary.BinaryProtocolWriter & test_model.StateTestWriterBase
  methods
    function obj = StateTestWriter(filename)
      obj@test_model.StateTestWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.StateTestWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_an_int_(obj, value)
      w = yardl.binary.Int32Serializer;
      w.write(obj.stream_, value);
    end

    function write_a_stream_(obj, value)
      w = yardl.binary.StreamSerializer(yardl.binary.Int32Serializer);
      w.write(obj.stream_, value);
    end

    function write_another_int_(obj, value)
      w = yardl.binary.Int32Serializer;
      w.write(obj.stream_, value);
    end
  end
end
