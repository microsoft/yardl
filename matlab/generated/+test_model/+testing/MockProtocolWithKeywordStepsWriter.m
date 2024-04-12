% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MockProtocolWithKeywordStepsWriter < matlab.mixin.Copyable & test_model.ProtocolWithKeywordStepsWriterBase
  properties
    testCase_
    expected_int
    expected_float
  end

  methods
    function obj = MockProtocolWithKeywordStepsWriter(testCase)
      obj.testCase_ = testCase;
      obj.expected_int = {};
      obj.expected_float = yardl.None;
    end

    function expect_write_int_(obj, value)
      if iscell(value)
        for n = 1:numel(value)
          obj.expected_int{end+1} = value{n};
        end
        return;
      end
      shape = size(value);
      lastDim = ndims(value);
      count = shape(lastDim);
      index = repelem({':'}, lastDim-1);
      for n = 1:count
        obj.expected_int{end+1} = value(index{:}, n);
      end
    end

    function expect_write_float_(obj, value)
      obj.expected_float = yardl.Optional(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.expected_int), "Expected call to write_int_ was not received");
      obj.testCase_.verifyEqual(obj.expected_float, yardl.None, "Expected call to write_float_ was not received");
    end
  end

  methods (Access=protected)
    function write_int_(obj, value)
      assert(iscell(value));
      assert(isscalar(value));
      obj.testCase_.verifyFalse(isempty(obj.expected_int), "Unexpected call to write_int_");
      obj.testCase_.verifyEqual(value{1}, obj.expected_int{1}, "Unexpected argument value for call to write_int_");
      obj.expected_int = obj.expected_int(2:end);
    end

    function write_float_(obj, value)
      obj.testCase_.verifyTrue(obj.expected_float.has_value(), "Unexpected call to write_float_");
      obj.testCase_.verifyEqual(value, obj.expected_float.value, "Unexpected argument value for call to write_float_");
      obj.expected_float = yardl.None;
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
