% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MockFlagsWriter < matlab.mixin.Copyable & test_model.FlagsWriterBase
  properties
    testCase_
    expected_days
    expected_formats
  end

  methods
    function self = MockFlagsWriter(testCase)
      self.testCase_ = testCase;
      self.expected_days = {};
      self.expected_formats = {};
    end

    function expect_write_days_(self, value)
      if iscell(value)
        for n = 1:numel(value)
          self.expected_days{end+1} = value{n};
        end
        return;
      end
      shape = size(value);
      lastDim = ndims(value);
      count = shape(lastDim);
      index = repelem({':'}, lastDim-1);
      for n = 1:count
        self.expected_days{end+1} = value(index{:}, n);
      end
    end

    function expect_write_formats_(self, value)
      if iscell(value)
        for n = 1:numel(value)
          self.expected_formats{end+1} = value{n};
        end
        return;
      end
      shape = size(value);
      lastDim = ndims(value);
      count = shape(lastDim);
      index = repelem({':'}, lastDim-1);
      for n = 1:count
        self.expected_formats{end+1} = value(index{:}, n);
      end
    end

    function verify(self)
      self.testCase_.verifyTrue(isempty(self.expected_days), "Expected call to write_days_ was not received");
      self.testCase_.verifyTrue(isempty(self.expected_formats), "Expected call to write_formats_ was not received");
    end
  end

  methods (Access=protected)
    function write_days_(self, value)
      assert(iscell(value));
      assert(isscalar(value));
      self.testCase_.verifyFalse(isempty(self.expected_days), "Unexpected call to write_days_");
      self.testCase_.verifyEqual(value{1}, self.expected_days{1}, "Unexpected argument value for call to write_days_");
      self.expected_days = self.expected_days(2:end);
    end

    function write_formats_(self, value)
      assert(iscell(value));
      assert(isscalar(value));
      self.testCase_.verifyFalse(isempty(self.expected_formats), "Unexpected call to write_formats_");
      self.testCase_.verifyEqual(value{1}, self.expected_formats{1}, "Unexpected argument value for call to write_formats_");
      self.expected_formats = self.expected_formats(2:end);
    end

    function close_(self)
    end
    function end_stream_(self)
    end
  end
end
