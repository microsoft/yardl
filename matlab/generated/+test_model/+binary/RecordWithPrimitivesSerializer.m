% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithPrimitivesSerializer < yardl.binary.RecordSerializer
  methods
    function self = RecordWithPrimitivesSerializer()
      field_serializers{1} = yardl.binary.BoolSerializer;
      field_serializers{2} = yardl.binary.Int8Serializer;
      field_serializers{3} = yardl.binary.Uint8Serializer;
      field_serializers{4} = yardl.binary.Int16Serializer;
      field_serializers{5} = yardl.binary.Uint16Serializer;
      field_serializers{6} = yardl.binary.Int32Serializer;
      field_serializers{7} = yardl.binary.Uint32Serializer;
      field_serializers{8} = yardl.binary.Int64Serializer;
      field_serializers{9} = yardl.binary.Uint64Serializer;
      field_serializers{10} = yardl.binary.SizeSerializer;
      field_serializers{11} = yardl.binary.Float32Serializer;
      field_serializers{12} = yardl.binary.Float64Serializer;
      field_serializers{13} = yardl.binary.Complexfloat32Serializer;
      field_serializers{14} = yardl.binary.Complexfloat64Serializer;
      field_serializers{15} = yardl.binary.DateSerializer;
      field_serializers{16} = yardl.binary.TimeSerializer;
      field_serializers{17} = yardl.binary.DatetimeSerializer;
      self@yardl.binary.RecordSerializer('test_model.RecordWithPrimitives', field_serializers);
    end

    function write(self, outstream, value)
      arguments
        self
        outstream (1,1) yardl.binary.CodedOutputStream
        value (1,1) test_model.RecordWithPrimitives
      end
      self.write_(outstream, value.bool_field, value.int8_field, value.uint8_field, value.int16_field, value.uint16_field, value.int32_field, value.uint32_field, value.int64_field, value.uint64_field, value.size_field, value.float32_field, value.float64_field, value.complexfloat32_field, value.complexfloat64_field, value.date_field, value.time_field, value.datetime_field)
    end

    function value = read(self, instream)
      field_values = self.read_(instream);
      value = test_model.RecordWithPrimitives(field_values{:});
    end
  end
end
