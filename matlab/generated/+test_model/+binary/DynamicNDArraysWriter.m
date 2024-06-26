% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef DynamicNDArraysWriter < yardl.binary.BinaryProtocolWriter & test_model.DynamicNDArraysWriterBase
  % Binary writer for the DynamicNDArrays protocol
  properties (Access=protected)
    ints_serializer
    simple_record_array_serializer
    record_with_vlens_array_serializer
    record_with_dynamic_nd_arrays_serializer
  end

  methods
    function self = DynamicNDArraysWriter(filename)
      self@test_model.DynamicNDArraysWriterBase();
      self@yardl.binary.BinaryProtocolWriter(filename, test_model.DynamicNDArraysWriterBase.schema);
      self.ints_serializer = yardl.binary.DynamicNDArraySerializer(yardl.binary.Int32Serializer);
      self.simple_record_array_serializer = yardl.binary.DynamicNDArraySerializer(test_model.binary.SimpleRecordSerializer());
      self.record_with_vlens_array_serializer = yardl.binary.DynamicNDArraySerializer(test_model.binary.RecordWithVlensSerializer());
      self.record_with_dynamic_nd_arrays_serializer = test_model.binary.RecordWithDynamicNDArraysSerializer();
    end
  end

  methods (Access=protected)
    function write_ints_(self, value)
      self.ints_serializer.write(self.stream_, value);
    end

    function write_simple_record_array_(self, value)
      self.simple_record_array_serializer.write(self.stream_, value);
    end

    function write_record_with_vlens_array_(self, value)
      self.record_with_vlens_array_serializer.write(self.stream_, value);
    end

    function write_record_with_dynamic_nd_arrays_(self, value)
      self.record_with_dynamic_nd_arrays_serializer.write(self.stream_, value);
    end
  end
end
