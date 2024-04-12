% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MockBenchmarkSimpleMrdWriter < matlab.mixin.Copyable & test_model.BenchmarkSimpleMrdWriterBase
  properties
    testCase_
    expected_data
  end

  methods
    function obj = MockBenchmarkSimpleMrdWriter(testCase)
      obj.testCase_ = testCase;
      obj.expected_data = {};
    end

    function expect_write_data_(obj, value)
      if iscell(value)
        for n = 1:numel(value)
          obj.expected_data{end+1} = value{n};
        end
        return;
      end
      shape = size(value);
      lastDim = ndims(value);
      count = shape(lastDim);
      index = repelem({':'}, lastDim-1);
      for n = 1:count
        obj.expected_data{end+1} = value(index{:}, n);
      end
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.expected_data), "Expected call to write_data_ was not received");
    end
  end

  methods (Access=protected)
    function write_data_(obj, value)
      assert(iscell(value));
      assert(isscalar(value));
      obj.testCase_.verifyFalse(isempty(obj.expected_data), "Unexpected call to write_data_");
      obj.testCase_.verifyEqual(value{1}, obj.expected_data{1}, "Unexpected argument value for call to write_data_");
      obj.expected_data = obj.expected_data(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
