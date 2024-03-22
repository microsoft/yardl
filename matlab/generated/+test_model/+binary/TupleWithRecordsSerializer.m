classdef TupleWithRecordsSerializer < yardl.binary.RecordSerializer
  methods
    function obj = TupleWithRecordsSerializer()
      field_serializers{1} = test_model.binary.SimpleRecordSerializer();
      field_serializers{2} = test_model.binary.SimpleRecordSerializer();
      obj@yardl.binary.RecordSerializer('test_model.TupleWithRecords', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.TupleWithRecords'));
      obj.write_(outstream, value.a, value.b)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.TupleWithRecords(field_values{:});
    end
  end
end