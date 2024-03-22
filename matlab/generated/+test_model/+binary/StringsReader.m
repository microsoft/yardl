% Binary reader for the Strings protocol
classdef StringsReader < yardl.binary.BinaryProtocolReader & test_model.StringsReaderBase
  methods
    function obj = StringsReader(filename)
      obj@test_model.StringsReaderBase();
      obj@yardl.binary.BinaryProtocolReader(filename, test_model.StringsReaderBase.schema);
    end
  end

  methods (Access=protected)
    function value = read_single_string_(obj)
      r = yardl.binary.StringSerializer;
      value = r.read(obj.stream_);
    end

    function value = read_rec_with_string_(obj)
      r = test_model.binary.RecordWithStringsSerializer();
      value = r.read(obj.stream_);
    end
  end
end