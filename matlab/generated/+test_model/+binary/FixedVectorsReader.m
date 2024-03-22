% Binary reader for the FixedVectors protocol
classdef FixedVectorsReader < yardl.binary.BinaryProtocolReader & test_model.FixedVectorsReaderBase
  methods
    function obj = FixedVectorsReader(filename)
      obj@test_model.FixedVectorsReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.FixedVectorsReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_fixed_int_vector_(obj)
      r = yardl.binary.FixedVectorSerializer(yardl.binary.Int32Serializer, 5);
      value = r.read(obj.stream_);
    end

    function value = read_fixed_simple_record_vector_(obj)
      r = yardl.binary.FixedVectorSerializer(test_model.binary.SimpleRecordSerializer(), 3);
      value = r.read(obj.stream_);
    end

    function value = read_fixed_record_with_vlens_vector_(obj)
      r = yardl.binary.FixedVectorSerializer(test_model.binary.RecordWithVlensSerializer(), 2);
      value = r.read(obj.stream_);
    end

    function value = read_record_with_fixed_vectors_(obj)
      r = test_model.binary.RecordWithFixedVectorsSerializer();
      value = r.read(obj.stream_);
    end
  end
end
