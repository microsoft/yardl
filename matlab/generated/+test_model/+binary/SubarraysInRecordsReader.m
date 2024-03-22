% Binary reader for the SubarraysInRecords protocol
classdef SubarraysInRecordsReader < yardl.binary.BinaryProtocolReader & test_model.SubarraysInRecordsReaderBase
  methods
    function obj = SubarraysInRecordsReader(filename)
      obj@test_model.SubarraysInRecordsReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.SubarraysInRecordsReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_with_fixed_subarrays_(obj)
      r = yardl.binary.DynamicNDArraySerializer(test_model.binary.RecordWithFixedCollectionsSerializer());
      value = r.read(obj.stream_);
    end

    function value = read_with_vlen_subarrays_(obj)
      r = yardl.binary.DynamicNDArraySerializer(test_model.binary.RecordWithVlenCollectionsSerializer());
      value = r.read(obj.stream_);
    end
  end
end
