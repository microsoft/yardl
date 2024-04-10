% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef ComputedFieldsTest < matlab.unittest.TestCase
    methods (Test)

        function testFieldAccess(testCase)
            r = test_model.RecordWithComputedFields();

            r.int_field = int32(42);
            testCase.verifyEqual(r.access_int_field(), int32(42));
            testCase.verifyEqual(r.access_other_computed_field(), r.access_int_field());

            r.string_field = "hello";
            testCase.verifyEqual(r.access_string_field(), "hello");

            r.tuple_field = test_model.MyTuple(1, 2);
            testCase.verifyEqual(r.access_tuple_field(), r.tuple_field);
            testCase.verifyEqual(r.access_nested_tuple_field(), r.tuple_field.v2);

            r.array_field = int32([[1; 2; 3], [4; 5; 6]]);
            testCase.verifyEqual(r.access_array_field(), r.array_field);
            testCase.verifyEqual(r.access_array_field_element(), r.array_field(2, 1));
            testCase.verifyEqual(r.access_array_field_element_by_name(), r.array_field(2, 1));

            r.vector_field = [1, 2, 3, 4];
            testCase.verifyEqual(r.access_vector_field, r.vector_field);
            testCase.verifyEqual(r.access_vector_field_element(), r.vector_field(2));
            r.vector_of_vectors_field = int32([[1; 2; 3], [4; 5; 6]]);
            testCase.verifyEqual(r.access_vector_of_vectors_field(), r.vector_of_vectors_field(3, 2));
            r.fixed_vector_of_vectors_field = int32([[1; 2; 3], [4; 5; 6]]);
            testCase.verifyEqual(r.access_fixed_vector_of_vectors_field(), r.fixed_vector_of_vectors_field(3, 2));

            r.map_field = dictionary(["hello", "world"], ["world", "bye"]);
            testCase.verifyEqual(r.access_map(), r.map_field);
            testCase.verifyEqual(r.access_map_entry(), "world");
            testCase.verifyEqual(r.access_map_entry_with_computed_field(), "world");
            testCase.verifyEqual(r.access_map_entry_with_computed_field_nested(), "bye");

            testCase.verifyError(@() r.access_missing_map_entry(), "MATLAB:dictionary:ScalarKeyNotFound");
        end

        function testLiterals(testCase)
            r = test_model.RecordWithComputedFields();
            testCase.verifyEqual(r.int_literal(), 42);
            testCase.verifyEqual(int64(r.large_negative_int64_literal()), -int64(0x4000000000000000));
            testCase.verifyEqual(uint64(r.large_u_int64_literal()), 0x8000000000000000);
            testCase.verifyEqual(r.string_literal(), "hello");
        end

        function testDimensionIndex(testCase)
            r = test_model.RecordWithComputedFields();

            testCase.verifyEqual(r.array_dimension_x_index(), 0);
            testCase.verifyEqual(r.array_dimension_y_index(), 1);

            r.string_field = "y";
            testCase.verifyEqual(r.array_dimension_index_from_string_field(), 1);

            r.string_field = "missing";
            testCase.verifyError(@() r.array_dimension_index_from_string_field(), "yardl:KeyError");
        end

        function testDimensionCount(testCase)
            r = test_model.RecordWithComputedFields();

            testCase.verifyEqual(r.array_dimension_count(), 2);

            r.dynamic_array_field = int32([[1; 2; 3], [4; 5; 6]]);
            testCase.verifyEqual(r.dynamic_array_dimension_count(), 2);
            r.dynamic_array_field = int32([1, 2, 3]);
            testCase.verifyEqual(r.dynamic_array_dimension_count(), 1);
        end

        function testVectorSize(testCase)
            r = test_model.RecordWithComputedFields();
            testCase.verifyEqual(r.vector_size(), 0);
            r.vector_field = [1, 2, 3, 4];
            testCase.verifyEqual(r.vector_size(), 4);

            testCase.verifyEqual(r.fixed_vector_size(), 3);
            testCase.verifyEqual(r.fixed_vector_of_vectors_size(), 2);
        end

        function testMapSize(testCase)
            r = test_model.RecordWithComputedFields();
            testCase.verifyEqual(r.map_size(), 0);
            r.map_field = dictionary(["hello", "bonjour"], ["world", "monde"]);
            testCase.verifyEqual(r.map_size(), 2);
        end

        function testArraySize(testCase)
            r = test_model.RecordWithComputedFields();
            r.array_field = int32([[1; 2; 3], [4; 5; 6]]);

            testCase.verifyEqual(r.array_size(), 6);
            testCase.verifyEqual(r.array_x_size(), 2);
            testCase.verifyEqual(r.array_y_size(), 3);
            testCase.verifyEqual(r.array_0_size(), 2);
            testCase.verifyEqual(r.array_1_size(), 3);

            testCase.verifyEqual(r.array_size_from_int_field(), 2);
            r.int_field = 1;
            testCase.verifyEqual(r.array_size_from_int_field(), 3);

            r.string_field = "x";
            testCase.verifyEqual(r.array_size_from_string_field(), 2);
            r.string_field = "y";
            testCase.verifyEqual(r.array_size_from_string_field(), 3);

            r.string_field = "missing";
            testCase.verifyError(@() r.array_size_from_string_field(), "yardl:KeyError");

            r.tuple_field.v1 = 1;
            testCase.verifyEqual(r.array_size_from_nested_int_field(), 3);
            testCase.verifyEqual(r.fixed_array_size(), 12);
            testCase.verifyEqual(r.fixed_array_x_size(), 3);
            testCase.verifyEqual(r.fixed_array_0_size(), 3);

            r.array_field_map_dimensions = int32([[1; 2; 3], [4; 5; 6]]);
            testCase.verifyEqual(r.array_field_map_dimensions_x_size(), 2);
        end

        function testSwitch(testCase)
            r = test_model.RecordWithComputedFields();
            r.optional_named_array = int32([[1; 2; 3], [4; 5; 6]]);
            testCase.verifyEqual(r.optional_named_array_length(), 6);
            testCase.verifyEqual(r.optional_named_array_length_with_discard(), 6);

            r.optional_named_array = yardl.None;
            testCase.verifyEqual(r.optional_named_array_length(), 0);
            testCase.verifyEqual(r.optional_named_array_length_with_discard(), 0);

            r.int_float_union = test_model.Int32OrFloat32.Int32(int32(42));
            testCase.verifyEqual(r.int_float_union_as_float(), single(42));
            r.int_float_union = test_model.Int32OrFloat32.Float32(single(42.9));
            testCase.verifyEqual(r.int_float_union_as_float(), single(42.9));

            r.nullable_int_float_union = yardl.None;
            testCase.verifyEqual(r.nullable_int_float_union_string(), "null");
            r.nullable_int_float_union = test_model.Int32OrFloat32.Int32(42);
            testCase.verifyEqual(r.nullable_int_float_union_string(), "int");
            r.nullable_int_float_union = test_model.Int32OrFloat32.Float32(42.9);
            testCase.verifyEqual(r.nullable_int_float_union_string(), "float");

            r.union_with_nested_generic_union = test_model.IntOrGenericRecordWithComputedFields.Int(42);
            testCase.verifyEqual(r.nested_switch(), -1);
            testCase.verifyEqual(r.use_nested_computed_field(), -1);

            g0 = basic_types.GenericRecordWithComputedFields(basic_types.T0OrT1.T0("hi"));
            r.union_with_nested_generic_union = ...
                test_model.IntOrGenericRecordWithComputedFields.GenericRecordWithComputedFields(g0);
            testCase.verifyEqual(r.nested_switch(), int32(10));
            testCase.verifyEqual(r.use_nested_computed_field(), int32(0));

            g1 = basic_types.GenericRecordWithComputedFields(basic_types.T0OrT1.T1(single(42.9)));
            r.union_with_nested_generic_union = ...
                test_model.IntOrGenericRecordWithComputedFields.GenericRecordWithComputedFields(g1);
            testCase.verifyEqual(r.nested_switch(), int32(20));
            testCase.verifyEqual(r.use_nested_computed_field(), int32(1));

            r.int_field = int32(10);
            testCase.verifyEqual(r.switch_over_single_value(), int32(10));


            gr = basic_types.GenericRecordWithComputedFields(basic_types.T0OrT1.T0(int32(42)));
            testCase.verifyEqual(gr.type_index(), 0);
            gr.f1 = basic_types.T0OrT1.T1(single(42.9));
            testCase.verifyEqual(gr.type_index(), 1);
        end

        function testArithmetic(testCase)
            r = test_model.RecordWithComputedFields();
            testCase.verifyEqual(r.arithmetic_1(), 3);
            testCase.verifyEqual(r.arithmetic_2(), 11);
            testCase.verifyEqual(r.arithmetic_3(), 13);

            r.array_field = int32([[1; 2; 3], [4; 5; 6]]);
            r.int_field = 1;
            testCase.verifyEqual(r.arithmetic_4(), 5);
            testCase.verifyEqual(r.arithmetic_5(), 3);

            testCase.verifyEqual(r.arithmetic_6(), 3.5);

            testCase.verifyEqual(r.arithmetic_7(), 49.0);

            r.complexfloat32_field = complex(2, 3);
            testCase.verifyEqual(r.arithmetic8(), single(complex(6, 9)));

            testCase.verifyEqual(r.arithmetic_9(), 2.2);
            testCase.verifyEqual(r.arithmetic_10(), 1e10 + 9e9);
        end

        function testCasting(testCase)
            r = test_model.RecordWithComputedFields();
            r.int_field = int32(42);
            testCase.verifyEqual(r.cast_int_to_float(), single(42));

            r.float32_field = 42.9;
            % NOTE: Matlab rounds floating point numbers when casting to integer
            testCase.verifyEqual(r.cast_float_to_int(), int32(43));

            testCase.verifyEqual(r.cast_power(), int32(49));

            r.complexfloat32_field = single(complex(2, 3));
            r.complexfloat64_field = complex(2, 3);
            testCase.verifyEqual(r.cast_complex32_to_complex64(), complex(2, 3));
            testCase.verifyEqual(r.cast_complex64_to_complex32(), single(complex(2, 3)));

            testCase.verifyEqual(r.cast_float_to_complex(), complex(single(66.6), 0));
        end

    end
end
