classdef GeneratedTypesTest < matlab.unittest.TestCase
    methods (Test)

        function testDefaultRecordWithPrimitives(testCase)
            r = RecordWithPrimitives();

            assert(r.bool_field == false);
            assert(r.int32_field == int32(0));
            assert(r.date_field == datetime(1970, 1, 1));
            assert(r.time_field == yardl.Time(0));
            assert(r.datetime_field == yardl.DateTime(0));
        end

        function testDefaultRecordWithVectors(testCase)
            r = RecordWithVectors();

            assert(isequal(r.default_vector, int32([])));
            assert(isequal(r.default_vector_fixed_length, int32([0, 0, 0])));
            assert(isequal(r.default_vector, int32([])));
        end

        function testDefaultRecordWithArrays(testCase)
            r = RecordWithArrays();

            assert(isequal(r.default_array, int32([])));
            assert(isempty(r.default_array_with_empty_dimension));
            assert(isempty(r.rank_1_array));
            assert(isempty(r.rank_2_array));
            assert(isequal(size(r.rank_2_array), [0, 0]));
            assert(isempty(r.rank_2_array_with_named_dimensions));
            assert(isequal(size(r.rank_2_array_with_named_dimensions), [0, 0]));

            assert(isequal(r.rank_2_fixed_array, zeros(4, 3, 'int32')));
            assert(isequal(r.rank_2_fixed_array_with_named_dimensions, zeros(4, 3, 'int32')));

            assert(isempty(r.dynamic_array));
            assert(isa(r.dynamic_array, 'int32'));

            assert(isa(r.array_of_vectors, 'int32'));
            assert(isequal(size(r.array_of_vectors), [5, 4]));
        end

        function testDefaultRecordWithOptionalFields(testCase)
            r = RecordWithOptionalFields();

            assert(r.optional_int == yardl.None);
        end

        function testDefaultRecordWithUnionsOfContainers(testCase)
            r = RecordWithUnionsOfContainers();

            assert(r.map_or_scalar == MapOrScalar.Map(containers.Map));
            assert(r.vector_or_scalar == VectorOrScalar.Vector(int32([])));
            assert(r.array_or_scalar == ArrayOrScalar.Array(int32([])));
        end

        function testDefaultRecordWithAliasedGenerics(testCase)
            r = RecordWithAliasedGenerics();

            assert(r.my_strings.v1 == "");
            assert(r.my_strings.v2 == "");
            assert(r.aliased_strings.v1 == "");
            assert(r.aliased_strings.v2 == "");
        end

        function testDefaultRecordGenericEmpty(testCase)
            g1 = RecordWithOptionalGenericField();
            testCase.verifyEqual(g1.v, yardl.None);

            g1a = RecordWithAliasedOptionalGenericField();
            testCase.verifyEqual(g1a.v, g1.v);

            g2 = RecordWithOptionalGenericUnionField();
            testCase.verifyEqual(g2.v, yardl.None);

            g2a = RecordWithAliasedOptionalGenericUnionField();
            testCase.verifyEqual(g2a.v, g2.v)

            testCase.verifyError(@() MyTuple(), 'MATLAB:minrhs');
            rm = RecordWithGenericMaps();
            testCase.verifyEqual(rm.m, containers.Map());
            testCase.verifyEqual(rm.am, rm.m);
        end

        function testDefaultRecordWithGenericRequiredArguments(testCase)
            testCase.verifyError(@() RecordWithGenericArrays(), 'MATLAB:minrhs');
            testCase.verifyError(@() RecordWithGenericVectors(), 'MATLAB:minrhs');
            testCase.verifyError(@() RecordWithGenericFixedVectors(), 'MATLAB:minrhs');
        end

        function testDefaultRecordContainingNestedGenericRecords(testCase)
            r = RecordContainingNestedGenericRecords();

            g1 = RecordWithOptionalGenericField();
            g1a = RecordWithAliasedOptionalGenericField();
            g2 = RecordWithOptionalGenericUnionField();
            g2a = RecordWithAliasedOptionalGenericUnionField();
            g7 = RecordWithGenericMaps();

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
