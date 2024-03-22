classdef SmallBenchmarkRecordSerializer < yardl.binary.RecordSerializer
  methods
    function obj = SmallBenchmarkRecordSerializer()
      field_serializers{1} = yardl.binary.Float64Serializer;
      field_serializers{2} = yardl.binary.Float32Serializer;
      field_serializers{3} = yardl.binary.Float32Serializer;
      obj@yardl.binary.RecordSerializer('test_model.SmallBenchmarkRecord', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.SmallBenchmarkRecord'));
      obj.write_(outstream, value.a, value.b, value.c)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.SmallBenchmarkRecord(field_values{:});
    end
  end
end
