% Abstract writer for protocol NestedRecords
classdef (Abstract) NestedRecordsWriterBase < handle
  properties (Access=protected)
    state_
  end

  methods
    function obj = NestedRecordsWriterBase()
      obj.state_ = 0;
    end

    function close(obj)
      obj.close_();
      if obj.state_ ~= 2
        expected_method = obj.state_to_method_name_(bitand((int32(obj.state_) + 1), bitcmp(1, 'int8')));
        throw(yardl.ProtocolError("Protocol writer closed before all steps were called. Expected call to '%s'.", expected_method));
      end
    end

    % Ordinal 0
    function write_tuple_with_records(obj, value)
      if obj.state_ ~= 0
        obj.raise_unexpected_state_(0);
      end

      obj.write_tuple_with_records_(value);
      obj.state_ = 2;
    end
  end

  methods (Static)
    function res = schema()
      res = string('{"protocol":{"name":"NestedRecords","sequence":[{"name":"tupleWithRecords","type":"TestModel.TupleWithRecords"}]},"types":[{"name":"SimpleRecord","fields":[{"name":"x","type":"int32"},{"name":"y","type":"int32"},{"name":"z","type":"int32"}]},{"name":"TupleWithRecords","fields":[{"name":"a","type":"TestModel.SimpleRecord"},{"name":"b","type":"TestModel.SimpleRecord"}]}]}');
    end
  end

  methods (Abstract, Access=protected)
    write_tuple_with_records_(obj, value)

    end_stream_(obj)
    close_(obj)
  end

  methods (Access=private)
    function raise_unexpected_state_(obj, actual)
      expected_method = obj.state_to_method_name_(obj.state_);
      actual_method = obj.state_to_method_name_(actual);
      throw(yardl.ProtocolError("Expected call to '%s' but received call to '%s'", expected_method, actual_method));
    end

    function name = state_to_method_name_(obj, state)
      if state == 0
        name = 'write_tuple_with_records';
      else
        name = '<unknown>';
      end
    end
  end
end