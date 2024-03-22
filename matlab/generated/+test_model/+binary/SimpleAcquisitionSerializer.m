classdef SimpleAcquisitionSerializer < yardl.binary.RecordSerializer
  methods
    function obj = SimpleAcquisitionSerializer()
      field_serializers{1} = yardl.binary.Uint64Serializer;
      field_serializers{2} = test_model.binary.SimpleEncodingCountersSerializer();
      field_serializers{3} = yardl.binary.NDArraySerializer(yardl.binary.Complexfloat32Serializer, 2);
      field_serializers{4} = yardl.binary.NDArraySerializer(yardl.binary.Float32Serializer, 2);
      obj@yardl.binary.RecordSerializer('test_model.SimpleAcquisition', field_serializers);
    end

    function write(obj, outstream, value)
      assert(isa(value, 'test_model.SimpleAcquisition'));
      obj.write_(outstream, value.flags, value.idx, value.data, value.trajectory)
    end

    function value = read(obj, instream)
      field_values = obj.read_(instream);
      value = test_model.SimpleAcquisition(field_values{:});
    end
  end
end
