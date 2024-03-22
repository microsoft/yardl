classdef RecordWithUnionsSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithUnionsSerializer()
      field_serializers{1} = yardl.binary.UnionSerializer('basic_types.Int32OrString', {yardl.binary.NoneSerializer, yardl.binary.Int32Serializer, yardl.binary.StringSerializer}, {yardl.None, @basic_types.Int32OrString.Int32, @basic_types.Int32OrString.String});
      field_serializers{2} = yardl.binary.UnionSerializer('basic_types.TimeOrDatetime', {yardl.binary.TimeSerializer, yardl.binary.DatetimeSerializer}, {@basic_types.TimeOrDatetime.Time, @basic_types.TimeOrDatetime.Datetime});
      field_serializers{3} = yardl.binary.UnionSerializer('basic_types.GenericNullableUnion2', {yardl.binary.NoneSerializer, yardl.binary.EnumSerializer('basic_types.Fruits', @basic_types.Fruits, yardl.binary.Int32Serializer), yardl.binary.EnumSerializer('basic_types.DaysOfWeek', @basic_types.DaysOfWeek, yardl.binary.Int32Serializer)}, {yardl.None, @basic_types.GenericNullableUnion2.T1, @basic_types.GenericNullableUnion2.T2});
      obj@yardl.binary.RecordSerializer('basic_types.RecordWithUnions', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'basic_types.RecordWithUnions'));
      obj.write_(outstream, value.null_or_int_or_string, value.date_or_datetime, value.null_or_fruits_or_days_of_week)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = basic_types.RecordWithUnions(field_values{:});
    end
  end
end
