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
    function obj = VlensWriter(filename)
      obj@test_model.VlensWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.VlensWriterBase.schema);
      obj.int_vector_serializer = yardl.binary.VectorSerializer(yardl.binary.Int32Serializer);
      obj.complex_vector_serializer = yardl.binary.VectorSerializer(yardl.binary.Complexfloat32Serializer);
      obj.record_with_vlens_serializer = test_model.binary.RecordWithVlensSerializer();
      obj.vlen_of_record_with_vlens_serializer = yardl.binary.VectorSerializer(test_model.binary.RecordWithVlensSerializer());
    end
  end

  methods (Access=protected)
    function write_int_vector_(obj, value)
      obj.int_vector_serializer.write(obj.stream_, value);
    end

    function write_complex_vector_(obj, value)
      obj.complex_vector_serializer.write(obj.stream_, value);
    end

    function write_record_with_vlens_(obj, value)
      obj.record_with_vlens_serializer.write(obj.stream_, value);
    end

    function write_vlen_of_record_with_vlens_(obj, value)
      obj.vlen_of_record_with_vlens_serializer.write(obj.stream_, value);
    end
  end
end
