% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MockStreamsOfUnionsWriter < matlab.mixin.Copyable & test_model.StreamsOfUnionsWriterBase
  properties
    testCase_
    expected_int_or_simple_record
    expected_nullable_int_or_simple_record
  end

  methods
    function obj = MockStreamsOfUnionsWriter(testCase)
      obj.testCase_ = testCase;
      obj.expected_int_or_simple_record = {};
      obj.expected_nullable_int_or_simple_record = {};
    end

    function expect_write_int_or_simple_record_(obj, value)
      if iscell(value)
        for n = 1:numel(value)
          obj.expected_int_or_simple_record{end+1} = value{n};
        end
        return;
      end
      shape = size(value);
      lastDim = ndims(value);
      count = shape(lastDim);
      index = repelem({':'}, lastDim-1);
      for n = 1:count
        obj.expected_int_or_simple_record{end+1} = value(index{:}, n);
      end
    end

    function expect_write_nullable_int_or_simple_record_(obj, value)
      if iscell(value)
        for n = 1:numel(value)
          obj.expected_nullable_int_or_simple_record{end+1} = value{n};
        end
        return;
      end
      shape = size(value);
      lastDim = ndims(value);
      count = shape(lastDim);
      index = repelem({':'}, lastDim-1);
      for n = 1:count
        obj.expected_nullable_int_or_simple_record{end+1} = value(index{:}, n);
      end
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.expected_int_or_simple_record), "Expected call to write_int_or_simple_record_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.expected_nullable_int_or_simple_record), "Expected call to write_nullable_int_or_simple_record_ was not received");
    end
  end

  methods (Access=protected)
    function write_int_or_simple_record_(obj, value)
      assert(iscell(value));
      assert(isscalar(value));
      obj.testCase_.verifyFalse(isempty(obj.expected_int_or_simple_record), "Unexpected call to write_int_or_simple_record_");
      obj.testCase_.verifyEqual(value{1}, obj.expected_int_or_simple_record{1}, "Unexpected argument value for call to write_int_or_simple_record_");
      obj.expected_int_or_simple_record = obj.expected_int_or_simple_record(2:end);
    end

    function write_nullable_int_or_simple_record_(obj, value)
      assert(iscell(value));
      assert(isscalar(value));
      obj.testCase_.verifyFalse(isempty(obj.expected_nullable_int_or_simple_record), "Unexpected call to write_nullable_int_or_simple_record_");
      obj.testCase_.verifyEqual(value{1}, obj.expected_nullable_int_or_simple_record{1}, "Unexpected argument value for call to write_nullable_int_or_simple_record_");
      obj.expected_nullable_int_or_simple_record = obj.expected_nullable_int_or_simple_record(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
