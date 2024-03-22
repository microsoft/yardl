% Binary writer for the Strings protocol
classdef StringsWriter < yardl.binary.BinaryProtocolWriter & test_model.StringsWriterBase
  methods
    function obj = StringsWriter(filename)
      obj@test_model.StringsWriterBase();
      obj@yardl.binary.BinaryProtocolWriter(filename, test_model.StringsWriterBase.schema);
    end
  end

  methods (Access=protected)
    function write_single_string_(obj, value)
      w = yardl.binary.StringSerializer;
      w.write(obj.stream_, value);
    end

    function write_rec_with_string_(obj, value)
      w = test_model.binary.RecordWithStringsSerializer();
      w.write(obj.stream_, value);
    end
  end
end
