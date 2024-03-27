classdef GeneratedTypesTest < matlab.unittest.TestCase
    methods (Test)

        % TODO: Add test for yardl.allocate for "each" kind of type

        function testDefaultRecordWithPrimitives(testCase)
            r = test_model.RecordWithPrimitives();

            testCase.verifyEqual(r.bool_field, false);
            testCase.verifyEqual(r.int32_field, int32(0));
            testCase.verifyEqual(r.date_field, yardl.Date());
            testCase.verifyEqual(r.time_field, yardl.Time());
            testCase.verifyEqual(r.datetime_field, yardl.DateTime());
        end

        function testDefaultRecordWithVectors(testCase)
            r = test_model.RecordWithVectors();

            testCase.verifyEqual(r.default_vector, int32([]));
            testCase.verifyEqual(r.default_vector_fixed_length, int32([0, 0, 0]));
            testCase.verifyEqual(r.default_vector, int32([]));
        end

        function testDefaultRecordWithArrays(testCase)
            import matlab.unittest.constraints.IsEmpty
            import matlab.unittest.constraints.IsInstanceOf

            r = test_model.RecordWithArrays();

            testCase.verifyEqual(r.default_array, int32([]));
            testCase.verifyThat(r.default_array_with_empty_dimension, IsEmpty);
            testCase.verifyThat(r.rank_1_array, IsEmpty);
            testCase.verifyThat(r.rank_2_array, IsEmpty);
            testCase.verifyEqual(size(r.rank_2_array), [0, 0]);
            testCase.verifyThat(r.rank_2_array_with_named_dimensions, IsEmpty);
            testCase.verifyEqual(size(r.rank_2_array_with_named_dimensions), [0, 0]);

            testCase.verifyEqual(r.rank_2_fixed_array, zeros(4, 3, 'int32'));
            testCase.verifyEqual(r.rank_2_fixed_array_with_named_dimensions, zeros(4, 3, 'int32'));

            testCase.verifyThat(r.dynamic_array, IsEmpty);
            testCase.verifyThat(r.dynamic_array, IsInstanceOf('int32'));

            testCase.verifyThat(r.array_of_vectors, IsInstanceOf('int32'));
            testCase.verifyEqual(size(r.array_of_vectors), [5, 4]);
        end

        function testDefaultRecordWithOptionalFields(testCase)
            r = test_model.RecordWithOptionalFields();

            testCase.verifyEqual(r.optional_int, yardl.None);
        end

        function testDefaultRecordWithUnionsOfContainers(testCase)
            r = test_model.RecordWithUnionsOfContainers();

            % testCase.verifyEqual(r.map_or_scalar, test_model.MapOrScalar.Map(containers.Map));
            testCase.verifyEqual(r.map_or_scalar, test_model.MapOrScalar.Map(dictionary));
            testCase.verifyEqual(r.vector_or_scalar, test_model.VectorOrScalar.Vector(int32([])));
            testCase.verifyEqual(r.array_or_scalar, test_model.ArrayOrScalar.Array(int32([])));
        end

        function testDefaultRecordWithAliasedGenerics(testCase)
            r = test_model.RecordWithAliasedGenerics();

            testCase.verifyEqual(r.my_strings.v1, "");
            testCase.verifyEqual(r.my_strings.v2, "");
            testCase.verifyEqual(r.aliased_strings.v1, "");
            testCase.verifyEqual(r.aliased_strings.v2, "");
        end

        function testDefaultRecordGenericEmpty(testCase)
            g1 = test_model.RecordWithOptionalGenericField();
            testCase.verifyEqual(g1.v, yardl.None);

            g1a = test_model.RecordWithAliasedOptionalGenericField();
            testCase.verifyEqual(g1a.v, g1.v);

            g2 = test_model.RecordWithOptionalGenericUnionField();
            testCase.verifyEqual(g2.v, yardl.None);

            g2a = test_model.RecordWithAliasedOptionalGenericUnionField();
            testCase.verifyEqual(g2a.v, g2.v)

            testCase.verifyError(@() test_model.MyTuple(), 'MATLAB:minrhs');
            rm = test_model.RecordWithGenericMaps();
            % testCase.verifyEqual(rm.m, containers.Map());
            testCase.verifyEqual(rm.m, dictionary);
            testCase.verifyEqual(rm.am, rm.m);
        end

        function testDefaultRecordWithGenericRequiredArguments(testCase)
            testCase.verifyError(@() test_model.RecordWithGenericArrays(), 'MATLAB:minrhs');
            testCase.verifyError(@() test_model.RecordWithGenericVectors(), 'MATLAB:minrhs');
            testCase.verifyError(@() test_model.RecordWithGenericFixedVectors(), 'MATLAB:minrhs');
        end

        function testDefaultRecordContainingNestedGenericRecords(testCase)
            r = test_model.RecordContainingNestedGenericRecords();

            g1 = test_model.RecordWithOptionalGenericField();
            g1a = test_model.RecordWithAliasedOptionalGenericField();
            g2 = test_model.RecordWithOptionalGenericUnionField();
            g2a = test_model.RecordWithAliasedOptionalGenericUnionField();
            g7 = test_model.RecordWithGenericMaps();

            testCase.verifyEqual(r.f1, g1);
            testCase.verifyEqual(r.f1a, g1a);
            testCase.verifyEqual(r.f2, g2);
            testCase.verifyEqual(r.f2a, g2a);

            testCase.verifyEqual(r.nested.g1, g1);
            testCase.verifyEqual(r.nested.g1a, g1a);
            testCase.verifyEqual(r.nested.g2, g2);
            testCase.verifyEqual(r.nested.g2a, g2a);

            testCase.verifyEqual(r.nested.g3.v1, "");
            testCase.verifyEqual(r.nested.g3.v2, int32(0));
            testCase.verifyEqual(r.nested.g3a.v1, r.nested.g3.v1);
            testCase.verifyEqual(r.nested.g3a.v2, r.nested.g3.v2);

            testCase.verifyEqual(r.nested.g4.v, int32([]));
            testCase.verifyEqual(r.nested.g4.av, r.nested.g4.v);

            testCase.verifyEqual(r.nested.g5.fv, int32([0, 0, 0]));
            testCase.verifyEqual(r.nested.g5.afv, r.nested.g5.fv);

            testCase.verifyEqual(r.nested.g6.nd, int32([]));
            testCase.verifyEqual(r.nested.g6.aliased_nd, r.nested.g6.nd);

            testCase.verifyEqual(r.nested.g6.fixed_nd, zeros(8, 16, 'int32'));
            testCase.verifyEqual(r.nested.g6.aliased_fixed_nd, r.nested.g6.fixed_nd);

            testCase.verifyEqual(r.nested.g6.dynamic_nd, int32([]));
            testCase.verifyEqual(r.nested.g6.aliased_dynamic_nd, r.nested.g6.dynamic_nd);

            testCase.verifyEqual(r.nested.g7, g7);
        end

    end
end
