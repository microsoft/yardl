% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef VlensWriter < yardl.binary.BinaryProtocolWriter & test_model.VlensWriterBase
  % Binary writer for the Vlens protocol
  properties (Access=protected)
    int_vector_serializer
    complex_vector_serializer
    record_with_vlens_serializer
    vlen_of_record_with_vlens_serializer
  end

  methods
    function self = VlensWriter(filename)
      self@test_model.VlensWriterBase();
      self@yardl.binary.BinaryProtocolWriter(filename, test_model.VlensWriterBase.schema);
      self.int_vector_serializer = yardl.binary.VectorSerializer(yardl.binary.Int32Serializer);
      self.complex_vector_serializer = yardl.binary.VectorSerializer(yardl.binary.Complexfloat32Serializer);
      self.record_with_vlens_serializer = test_model.binary.RecordWithVlensSerializer();
      self.vlen_of_record_with_vlens_serializer = yardl.binary.VectorSerializer(test_model.binary.RecordWithVlensSerializer());
    end
  end

  methods (Access=protected)
    function write_int_vector_(self, value)
      self.int_vector_serializer.write(self.stream_, value);
    end

    function write_complex_vector_(self, value)
      self.complex_vector_serializer.write(self.stream_, value);
    end

    function write_record_with_vlens_(self, value)
      self.record_with_vlens_serializer.write(self.stream_, value);
    end

    function write_vlen_of_record_with_vlens_(self, value)
      self.vlen_of_record_with_vlens_serializer.write(self.stream_, value);
    end
  end
end
