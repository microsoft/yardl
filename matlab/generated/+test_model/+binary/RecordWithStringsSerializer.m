classdef RecordWithStringsSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithStringsSerializer()
      field_serializers{1} = yardl.binary.StringSerializer;
      field_serializers{2} = yardl.binary.StringSerializer;
      obj@yardl.binary.RecordSerializer('test_model.RecordWithStrings', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithStrings'));
      obj.write_(outstream, value.a, value.b)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithStrings(field_values{:});
    end
  end
end
