classdef TupleSerializer < yardl.binary.RecordSerializer
  methods
    function obj = TupleSerializer(t1_serializer, t2_serializer)
      field_serializers{1} = t1_serializer;
      field_serializers{2} = t2_serializer;
      obj@yardl.binary.RecordSerializer('tuples.Tuple', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'tuples.Tuple'));
      obj.write_(outstream, value.v1, value.v2)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = tuples.Tuple(field_values{:});
    end
  end
end
