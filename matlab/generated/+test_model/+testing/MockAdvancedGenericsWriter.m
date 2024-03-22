classdef MockAdvancedGenericsWriter < test_model.AdvancedGenericsWriterBase
  properties
    testCase_
    write_float_image_image_written
    write_generic_record_1_written
    write_tuple_of_optionals_written
    write_tuple_of_optionals_alternate_syntax_written
    write_tuple_of_vectors_written
  end

  methods
    function obj = MockAdvancedGenericsWriter(testCase)
      obj.testCase_ = testCase;
      obj.write_float_image_image_written = Node.empty();
      obj.write_generic_record_1_written = Node.empty();
      obj.write_tuple_of_optionals_written = Node.empty();
      obj.write_tuple_of_optionals_alternate_syntax_written = Node.empty();
      obj.write_tuple_of_vectors_written = Node.empty();
    end

    function expect_write_float_image_image_(obj, value)
      obj.write_float_image_image_written(end+1) = Node(value);
    end

    function expect_write_generic_record_1_(obj, value)
      obj.write_generic_record_1_written(end+1) = Node(value);
    end

    function expect_write_tuple_of_optionals_(obj, value)
      obj.write_tuple_of_optionals_written(end+1) = Node(value);
    end

    function expect_write_tuple_of_optionals_alternate_syntax_(obj, value)
      obj.write_tuple_of_optionals_alternate_syntax_written(end+1) = Node(value);
    end

    function expect_write_tuple_of_vectors_(obj, value)
      obj.write_tuple_of_vectors_written(end+1) = Node(value);
    end

    function verify(obj)
      obj.testCase_.verifyTrue(isempty(obj.write_float_image_image_written), "Expected call to write_float_image_image_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_generic_record_1_written), "Expected call to write_generic_record_1_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_tuple_of_optionals_written), "Expected call to write_tuple_of_optionals_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_tuple_of_optionals_alternate_syntax_written), "Expected call to write_tuple_of_optionals_alternate_syntax_ was not received");
      obj.testCase_.verifyTrue(isempty(obj.write_tuple_of_vectors_written), "Expected call to write_tuple_of_vectors_ was not received");
    end
  end

  methods (Access=protected)
    function write_float_image_image_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_float_image_image_written), "Unexpected call to write_float_image_image_");
      expected = obj.write_float_image_image_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_float_image_image_");
      obj.write_float_image_image_written = obj.write_float_image_image_written(2:end);
    end

    function write_generic_record_1_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_generic_record_1_written), "Unexpected call to write_generic_record_1_");
      expected = obj.write_generic_record_1_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_generic_record_1_");
      obj.write_generic_record_1_written = obj.write_generic_record_1_written(2:end);
    end

    function write_tuple_of_optionals_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_tuple_of_optionals_written), "Unexpected call to write_tuple_of_optionals_");
      expected = obj.write_tuple_of_optionals_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_tuple_of_optionals_");
      obj.write_tuple_of_optionals_written = obj.write_tuple_of_optionals_written(2:end);
    end

    function write_tuple_of_optionals_alternate_syntax_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_tuple_of_optionals_alternate_syntax_written), "Unexpected call to write_tuple_of_optionals_alternate_syntax_");
      expected = obj.write_tuple_of_optionals_alternate_syntax_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_tuple_of_optionals_alternate_syntax_");
      obj.write_tuple_of_optionals_alternate_syntax_written = obj.write_tuple_of_optionals_alternate_syntax_written(2:end);
    end

    function write_tuple_of_vectors_(obj, value)
      obj.testCase_.verifyTrue(~isempty(obj.write_tuple_of_vectors_written), "Unexpected call to write_tuple_of_vectors_");
      expected = obj.write_tuple_of_vectors_written(1).value;
      obj.testCase_.verifyEqual(value, expected, "Unexpected argument value for call to write_tuple_of_vectors_");
      obj.write_tuple_of_vectors_written = obj.write_tuple_of_vectors_written(2:end);
    end

    function close_(obj)
    end
    function end_stream_(obj)
    end
  end
end