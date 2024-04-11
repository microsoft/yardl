% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef FixedArraysReader < yardl.binary.BinaryProtocolReader & test_model.FixedArraysReaderBase
  % Binary reader for the FixedArrays protocol
  properties (Access=protected)
    ints_serializer
    fixed_simple_record_array_serializer
    fixed_record_with_vlens_array_serializer
    record_with_fixed_arrays_serializer
    named_array_serializer
  end

  methods
    function obj = FixedArraysReader(filename)
      obj@test_model.FixedArraysReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.FixedArraysReaderBase.schema);
      obj.ints_serializer = yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [3, 2]);
      obj.fixed_simple_record_array_serializer = yardl.binary.FixedNDArraySerializer(test_model.binary.SimpleRecordSerializer(), [2, 3]);
      obj.fixed_record_with_vlens_array_serializer = yardl.binary.FixedNDArraySerializer(test_model.binary.RecordWithVlensSerializer(), [2, 2]);
      obj.record_with_fixed_arrays_serializer = test_model.binary.RecordWithFixedArraysSerializer();
      obj.named_array_serializer = yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [4, 2]);
    end
  end

  methods (Access=protected)
    function value = read_ints_(obj)
      value = obj.ints_serializer.read(obj.stream_);
    end

    function value = read_fixed_simple_record_array_(obj)
      value = obj.fixed_simple_record_array_serializer.read(obj.stream_);
    end

    function value = read_fixed_record_with_vlens_array_(obj)
      value = obj.fixed_record_with_vlens_array_serializer.read(obj.stream_);
    end

    function value = read_record_with_fixed_arrays_(obj)
      value = obj.record_with_fixed_arrays_serializer.read(obj.stream_);
    end

    function value = read_named_array_(obj)
      value = obj.named_array_serializer.read(obj.stream_);
    end
  end
end
