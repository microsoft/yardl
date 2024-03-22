% Binary writer for the SubarraysInRecords protocol
classdef SubarraysInRecordsWriter < yardl.binary.BinaryProtocolWriter & test_model.SubarraysInRecordsWriterBase
  methods
    function obj = SubarraysInRecordsWriter(filename)
      obj@test_model.SubarraysInRecordsWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.SubarraysInRecordsWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_with_fixed_subarrays_(obj, value)
      w = yardl.binary.DynamicNDArraySerializer(test_model.binary.RecordWithFixedCollectionsSerializer());
      w.write(obj.stream_, value);
    end

    function write_with_vlen_subarrays_(obj, value)
      w = yardl.binary.DynamicNDArraySerializer(test_model.binary.RecordWithVlenCollectionsSerializer());
      w.write(obj.stream_, value);
    end
  end
end