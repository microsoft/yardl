% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef StreamsReader < yardl.binary.BinaryProtocolReader & test_model.StreamsReaderBase
  % Binary reader for the Streams protocol
  properties (Access=protected)
    int_data_serializer
    optional_int_data_serializer
    record_with_optional_vector_data_serializer
    fixed_vector_serializer
  end

  methods
    function self = StreamsReader(filename, options)
      arguments
        filename (1,1) string
        options.skip_completed_check (1,1) logical = false
      end
      self@test_model.StreamsReaderBase(skip_completed_check=options.skip_completed_check);
      self@yardl.binary.BinaryProtocolReader(filename, test_model.StreamsReaderBase.schema);
      self.int_data_serializer = yardl.binary.StreamSerializer(yardl.binary.Int32Serializer);
      self.optional_int_data_serializer = yardl.binary.StreamSerializer(yardl.binary.OptionalSerializer(yardl.binary.Int32Serializer));
      self.record_with_optional_vector_data_serializer = yardl.binary.StreamSerializer(test_model.binary.RecordWithOptionalVectorSerializer());
      self.fixed_vector_serializer = yardl.binary.StreamSerializer(yardl.binary.FixedVectorSerializer(yardl.binary.Int32Serializer, 3));
    end
  end

  methods (Access=protected)
    function more = has_int_data_(self)
      more = self.int_data_serializer.hasnext(self.stream_);
    end

    function value = read_int_data_(self)
      value = self.int_data_serializer.read(self.stream_);
    end

    function more = has_optional_int_data_(self)
      more = self.optional_int_data_serializer.hasnext(self.stream_);
    end

    function value = read_optional_int_data_(self)
      value = self.optional_int_data_serializer.read(self.stream_);
    end

    function more = has_record_with_optional_vector_data_(self)
      more = self.record_with_optional_vector_data_serializer.hasnext(self.stream_);
    end

    function value = read_record_with_optional_vector_data_(self)
      value = self.record_with_optional_vector_data_serializer.read(self.stream_);
    end

    function more = has_fixed_vector_(self)
      more = self.fixed_vector_serializer.hasnext(self.stream_);
    end

    function value = read_fixed_vector_(self)
      value = self.fixed_vector_serializer.read(self.stream_);
    end
  end
end
