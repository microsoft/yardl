% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MockUnionsWriter < test_model.UnionsWriterBase
  properties
    testCase_
    write_int_or_simple_record_written
    write_int_or_record_with_vlens_written
    write_monosotate_or_int_or_simple_record_written
    write_record_with_unions_written
  end

  methods
    function obj = MockUnionsWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_int_or_simple_record_written = Node.empty();
      obj.write_int_or_record_with_vlens_written = Node.empty();
      obj.write_monosotate_or_int_or_simple_record_written = Node.empty();
      obj.write_record_with_unions_written = Node.empty();
    end

    function expect_write_int_or_simple_record_(obj, value)
      obj.write_int_or_simple_record_written(end+1) = Node(value);
    end

    function expect_write_int_or_record_with_vlens_(obj, value)
      obj.write_int_or_record_with_vlens_written(end+1) = Node(value);
    end

    function expect_write_monosotate_or_int_or_simple_record_(obj, value)
      obj.write_monosotate_or_int_or_simple_record_written(end+1) = Node(value);
    end

    function expect_write_record_with_unions_(obj, value)
      obj.write_record_with_unions_written(end+1) = Node(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_int_or_simple_record_written), "Expected call to write_int_or_simple_record_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_int_or_record_with_vlens_written), "Expected call to write_int_or_record_with_vlens_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_monosotate_or_int_or_simple_record_written), "Expected call to write_monosotate_or_int_or_simple_record_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_record_with_unions_written), "Expected call to write_record_with_unions_ was not received");
    end
  end

  methods (Access=protected)
    function write_int_or_simple_record_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_int_or_simple_record_written), "Unexpected call to write_int_or_simple_record_");
      expected = obj.write_int_or_simple_record_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_int_or_simple_record_");
      obj.write_int_or_simple_record_written = obj.write_int_or_simple_record_written(2:end);
    end

    function write_int_or_record_with_vlens_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_int_or_record_with_vlens_written), "Unexpected call to write_int_or_record_with_vlens_");
      expected = obj.write_int_or_record_with_vlens_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_int_or_record_with_vlens_");
      obj.write_int_or_record_with_vlens_written = obj.write_int_or_record_with_vlens_written(2:end);
    end

    function write_monosotate_or_int_or_simple_record_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_monosotate_or_int_or_simple_record_written), "Unexpected call to write_monosotate_or_int_or_simple_record_");
      expected = obj.write_monosotate_or_int_or_simple_record_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_monosotate_or_int_or_simple_record_");
      obj.write_monosotate_or_int_or_simple_record_written = obj.write_monosotate_or_int_or_simple_record_written(2:end);
    end

    function write_record_with_unions_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_record_with_unions_written), "Unexpected call to write_record_with_unions_");
      expected = obj.write_record_with_unions_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_record_with_unions_");
      obj.write_record_with_unions_written = obj.write_record_with_unions_written(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
