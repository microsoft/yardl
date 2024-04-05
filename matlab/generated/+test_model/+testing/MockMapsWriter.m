% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef MockMapsWriter < matlab.mixin.Copyable & test_model.MapsWriterBase
  properties
    testCase_
    write_string_to_int_written
    write_int_to_string_written
    write_string_to_union_written
    write_aliased_generic_written
  end

  methods
    function obj = MockMapsWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_string_to_int_written = yardl.None;
      obj.write_int_to_string_written = yardl.None;
      obj.write_string_to_union_written = yardl.None;
      obj.write_aliased_generic_written = yardl.None;
    end

    function expect_write_string_to_int_(obj, value)
      if obj.write_string_to_int_written.has_value()
        last_dim = ndims(value);
        obj.write_string_to_int_written = yardl.Optional(cat(last_dim, obj.write_string_to_int_written.value, value));
      else
        obj.write_string_to_int_written = yardl.Optional(value);
      end
    end

    function expect_write_int_to_string_(obj, value)
      if obj.write_int_to_string_written.has_value()
        last_dim = ndims(value);
        obj.write_int_to_string_written = yardl.Optional(cat(last_dim, obj.write_int_to_string_written.value, value));
      else
        obj.write_int_to_string_written = yardl.Optional(value);
      end
    end

    function expect_write_string_to_union_(obj, value)
      if obj.write_string_to_union_written.has_value()
        last_dim = ndims(value);
        obj.write_string_to_union_written = yardl.Optional(cat(last_dim, obj.write_string_to_union_written.value, value));
      else
        obj.write_string_to_union_written = yardl.Optional(value);
      end
    end

    function expect_write_aliased_generic_(obj, value)
      if obj.write_aliased_generic_written.has_value()
        last_dim = ndims(value);
        obj.write_aliased_generic_written = yardl.Optional(cat(last_dim, obj.write_aliased_generic_written.value, value));
      else
        obj.write_aliased_generic_written = yardl.Optional(value);
      end
    end

    function verify(obj)
      obj.testCase_.verifyEqual(obj.write_string_to_int_written, yardl.None, "Expected call to write_string_to_int_ was not received");
      obj.testCase_.verifyEqual(obj.write_int_to_string_written, yardl.None, "Expected call to write_int_to_string_ was not received");
      obj.testCase_.verifyEqual(obj.write_string_to_union_written, yardl.None, "Expected call to write_string_to_union_ was not received");
      obj.testCase_.verifyEqual(obj.write_aliased_generic_written, yardl.None, "Expected call to write_aliased_generic_ was not received");
    end
  end

  methods (Access=protected)
    function write_string_to_int_(obj, value)
      obj.testCase_.verifyTrue(obj.write_string_to_int_written.has_value(), "Unexpected call to write_string_to_int_");
      expected = obj.write_string_to_int_written.value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_string_to_int_");
      obj.write_string_to_int_written = yardl.None;
    end

    function write_int_to_string_(obj, value)
      obj.testCase_.verifyTrue(obj.write_int_to_string_written.has_value(), "Unexpected call to write_int_to_string_");
      expected = obj.write_int_to_string_written.value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_int_to_string_");
      obj.write_int_to_string_written = yardl.None;
    end

    function write_string_to_union_(obj, value)
      obj.testCase_.verifyTrue(obj.write_string_to_union_written.has_value(), "Unexpected call to write_string_to_union_");
      expected = obj.write_string_to_union_written.value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_string_to_union_");
      obj.write_string_to_union_written = yardl.None;
    end

    function write_aliased_generic_(obj, value)
      obj.testCase_.verifyTrue(obj.write_aliased_generic_written.has_value(), "Unexpected call to write_aliased_generic_");
      expected = obj.write_aliased_generic_written.value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_aliased_generic_");
      obj.write_aliased_generic_written = yardl.None;
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end
