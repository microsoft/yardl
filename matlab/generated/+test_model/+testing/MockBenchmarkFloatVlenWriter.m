% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MockBenchmarkFloatVlenWriter < matlab.mixin.Copyable & test_model.BenchmarkFloatVlenWriterBase
  properties
    testCase_
    expected_float_array
  end

  methods
    function self = MockBenchmarkFloatVlenWriter(testCase)
      self.testCase_ = testCase;
      self.expected_float_array = {};
    end

    function expect_write_float_array_(self, value)
      if iscell(value)
        for n = 1:numel(value)
          self.expected_float_array{end+1} = value{n};
        end
        return;
      end
      shape = size(value);
      lastDim = ndims(value);
      count = shape(lastDim);
      index = repelem({':'}, lastDim-1);
      for n = 1:count
        self.expected_float_array{end+1} = value(index{:}, n);
      end
    end

    function verify(self)
      self.testCase_.verifyTrue(isempty(self.expected_float_array), "Expected call to write_float_array_ was not received");
    end
  end

  methods (Access=protected)
    function write_float_array_(self, value)
      assert(iscell(value));
      assert(isscalar(value));
      self.testCase_.verifyFalse(isempty(self.expected_float_array), "Unexpected call to write_float_array_");
      self.testCase_.verifyEqual(value{1}, self.expected_float_array{1}, "Unexpected argument value for call to write_float_array_");
      self.expected_float_array = self.expected_float_array(2:end);
    end

    function close_(self)
    end
    function end_stream_(self)
    end
  end
end
