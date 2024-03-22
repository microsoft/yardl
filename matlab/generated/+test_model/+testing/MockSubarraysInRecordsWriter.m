classdef MockSubarraysInRecordsWriter < test_model.SubarraysInRecordsWriterBase
  properties
    testCase_
    write_with_fixed_subarrays_written
    write_with_vlen_subarrays_written
  end

  methods
    function obj = MockSubarraysInRecordsWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_with_fixed_subarrays_written = Node.empty();
      obj.write_with_vlen_subarrays_written = Node.empty();
    end

    function expect_write_with_fixed_subarrays_(obj, value)
      obj.write_with_fixed_subarrays_written(end+1) = Node(value);
    end

    function expect_write_with_vlen_subarrays_(obj, value)
      obj.write_with_vlen_subarrays_written(end+1) = Node(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_with_fixed_subarrays_written), "Expected call to write_with_fixed_subarrays_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_with_vlen_subarrays_written), "Expected call to write_with_vlen_subarrays_ was not received");
    end
  end

  methods (Access=protected)
    function write_with_fixed_subarrays_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_with_fixed_subarrays_written), "Unexpected call to write_with_fixed_subarrays_");
      expected = obj.write_with_fixed_subarrays_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_with_fixed_subarrays_");
      obj.write_with_fixed_subarrays_written = obj.write_with_fixed_subarrays_written(2:end);
    end

    function write_with_vlen_subarrays_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_with_vlen_subarrays_written), "Unexpected call to write_with_vlen_subarrays_");
      expected = obj.write_with_vlen_subarrays_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_with_vlen_subarrays_");
      obj.write_with_vlen_subarrays_written = obj.write_with_vlen_subarrays_written(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
