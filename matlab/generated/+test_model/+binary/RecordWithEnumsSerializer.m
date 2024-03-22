classdef RecordWithEnumsSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithEnumsSerializer()
      field_serializers{1} = yardl.binary.EnumSerializer('basic_types.Fruits', @basic_types.Fruits, yardl.binary.Int32Serializer);
      field_serializers{2} = yardl.binary.EnumSerializer('basic_types.DaysOfWeek', @basic_types.DaysOfWeek, yardl.binary.Int32Serializer);
      field_serializers{3} = yardl.binary.EnumSerializer('basic_types.TextFormat', @basic_types.TextFormat, yardl.binary.Uint64Serializer);
      obj@yardl.binary.RecordSerializer('test_model.RecordWithEnums', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithEnums'));
      obj.write_(outstream, value.enum, value.flags, value.flags_2)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithEnums(field_values{:});
    end
  end
end
