classdef RecordWithAliasedGenericsSerializer < yardl.binary.RecordSerializer
  methods
    function obj = RecordWithAliasedGenericsSerializer()
      field_serializers{1} = tuples.binary.TupleSerializer(yardl.binary.StringSerializer, yardl.binary.StringSerializer);
      field_serializers{2} = tuples.binary.TupleSerializer(yardl.binary.StringSerializer, yardl.binary.StringSerializer);
      obj@yardl.binary.RecordSerializer('test_model.RecordWithAliasedGenerics', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.RecordWithAliasedGenerics'));
      obj.write_(outstream, value.my_strings, value.aliased_strings)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.RecordWithAliasedGenerics(field_values{:});
    end
  end
end