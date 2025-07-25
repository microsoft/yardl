% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef BenchmarkInt256x256Reader < yardl.binary.BinaryProtocolReader & test_model.BenchmarkInt256x256ReaderBase
  % Binary reader for the BenchmarkInt256x256 protocol
  properties (Access=protected)
    int256x256_serializer
  end

  methods
    function self = BenchmarkInt256x256Reader(filename, options)
      arguments
        filename (1,1) string
        options.skip_completed_check (1,1) logical = false
      end
      self@test_model.BenchmarkInt256x256ReaderBase(skip_completed_check=options.skip_completed_check);
      self@yardl.binary.BinaryProtocolReader(filename, test_model.BenchmarkInt256x256ReaderBase.schema);
      self.int256x256_serializer = yardl.binary.StreamSerializer(yardl.binary.FixedNDArraySerializer(yardl.binary.Int32Serializer, [256, 256]));
    end
  end

  methods (Access=protected)
    function more = has_int256x256_(self)
      more = self.int256x256_serializer.hasnext(self.stream_);
    end

    function value = read_int256x256_(self)
      value = self.int256x256_serializer.read(self.stream_);
    end
  end
end
