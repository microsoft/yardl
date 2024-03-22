% Binary reader for the Vlens protocol
classdef VlensReader < yardl.binary.BinaryProtocolReader & test_model.VlensReaderBase
  methods
    function obj = VlensReader(filename)
      obj@test_model.VlensReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.VlensReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_int_vector_(obj)
      r = yardl.binary.VectorSerializer(yardl.binary.Int32Serializer);
      value = r.read(obj.stream_);
    end

    function value = read_complex_vector_(obj)
      r = yardl.binary.VectorSerializer(yardl.binary.Complexfloat32Serializer);
      value = r.read(obj.stream_);
    end

    function value = read_record_with_vlens_(obj)
      r = test_model.binary.RecordWithVlensSerializer();
      value = r.read(obj.stream_);
    end

    function value = read_vlen_of_record_with_vlens_(obj)
      r = yardl.binary.VectorSerializer(test_model.binary.RecordWithVlensSerializer());
      value = r.read(obj.stream_);
    end
  end
end
