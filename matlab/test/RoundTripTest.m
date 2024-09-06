% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef RoundTripTest < matlab.unittest.TestCase
    properties (TestParameter)
        % format = {"binary", "ndjson"};
        format = {"binary"};
    end

    methods (Test)

        function testScalarPrimitives(testCase, format)
            rec = test_model.RecordWithPrimitives();
            rec.bool_field = true;
            rec.int8_field = int8(-88);
            rec.uint8_field = uint8(88);
            rec.int16_field = int16(-1616);
            rec.uint16_field = uint16(1616);
            rec.int32_field = int32(-3232);
            rec.uint32_field = uint32(3232);
            rec.int64_field = int64(-64646464);
            rec.uint64_field = uint64(64646464);
            rec.size_field = uint64(64646464);
            rec.float32_field = single(32.0);
            rec.float64_field = double(64.64);
            rec.complexfloat32_field = single(complex(32.0, 64.0));
            rec.complexfloat64_field = 64.64 + 32.32j;
            rec.date_field = yardl.Date.from_components(2024, 4, 2);
            rec.time_field = yardl.Time.from_components(12, 34, 56, 0);
            rec.datetime_field = yardl.DateTime.from_components(2024, 4, 2, 12, 34, 56, 111222333);

            w = create_validating_writer(testCase, format, 'Scalars');
            w.write_int32(int32(42))
            w.write_record(rec)
            w.close();
        end

        function testScalarOptionals(testCase, format)
            w = create_validating_writer(testCase, format, 'ScalarOptionals');
            w.write_optional_int(yardl.None);
            w.write_optional_record(yardl.None);
            w.write_record_with_optional_fields(test_model.RecordWithOptionalFields());
            w.write_optional_record_with_optional_fields(yardl.None);
            w.close();

            w = create_validating_writer(testCase, format, 'ScalarOptionals');
            w.write_optional_int(int32(55));
            w.write_optional_record(test_model.SimpleRecord(x=8, y=9, z=10));
            w.write_record_with_optional_fields(...
                test_model.RecordWithOptionalFields(optional_int=int32(44), ...
                    optional_int_alternate_syntax=int32(44), ...
                    optional_time=yardl.Time.from_components(12, 34, 56, 0)));
            w.write_optional_record_with_optional_fields(...
                test_model.RecordWithOptionalFields(optional_int=int32(12), ...
                    optional_int_alternate_syntax=int32(12), ...
                    optional_time=yardl.Time.from_components(11, 32, 26, 0)));
            w.close();
        end

        function testNestedRecords(testCase, format)
            w = create_validating_writer(testCase, format, 'NestedRecords');
            t = test_model.TupleWithRecords(...
                a=test_model.SimpleRecord(x=1, y=2, z=3), ...
                b=test_model.SimpleRecord(x=4, y=5, z=6));
            w.write_tuple_with_records(t);
            w.close();
        end

        function testVariableLengthVectors(testCase, format)
            w = create_validating_writer(testCase, format, 'Vlens');
            w.write_int_vector(int32([1, 2, 3]));
            w.write_complex_vector(single([1+2j, 3+4j]));
            rec = test_model.RecordWithVlens(...
                a=[test_model.SimpleRecord(x=1, y=2, z=3), test_model.SimpleRecord(x=4, y=5, z=6)], ...
                b=4, ...
                c=2 ...
            );
            w.write_record_with_vlens(rec);
            w.write_vlen_of_record_with_vlens([rec, rec]);
            w.close();
        end

        function testStrings(testCase, format)
            w = create_validating_writer(testCase, format, 'Strings');
            w.write_single_string("hello");
            w.write_rec_with_string(test_model.RecordWithStrings(a="Montréal", b="臺北市"));
            w.close();
        end

        function testOptionalVectors(testCase, format)
            w = create_validating_writer(testCase, format, 'OptionalVectors');
            w.write_record_with_optional_vector(test_model.RecordWithOptionalVector());
            w.close();

            w = create_validating_writer(testCase, format, 'OptionalVectors');
            w.write_record_with_optional_vector(test_model.RecordWithOptionalVector(optional_vector=int32([1, 2, 3])));
            w.close();
        end

        function testFixedVectors(testCase, format)
            ints = int32([1, 2, 3, 4, 5]);
            SR = @test_model.SimpleRecord;
            simple_recs = [SR(x=1, y=2, z=3), SR(x=4, y=5, z=6), SR(x=7, y=8, z=9)];
            RV = @test_model.RecordWithVlens;
            vlens_recs = [...
                RV(a=[SR(x=1, y=2, z=3), SR(x=4, y=5, z=6)], b=4, c=2), ...
                RV(a=[SR(x=7, y=8, z=9), SR(x=10, y=11, z=12)], b=5, c=3)];
            rec_with_fixed = test_model.RecordWithFixedVectors(...
                fixed_int_vector=ints, ...
                fixed_simple_record_vector=simple_recs, ...
                fixed_record_with_vlens_vector=vlens_recs);

            w = create_validating_writer(testCase, format, 'FixedVectors');
            w.write_fixed_int_vector(ints);
            w.write_fixed_simple_record_vector(simple_recs);
            w.write_fixed_record_with_vlens_vector(vlens_recs);
            w.write_record_with_fixed_vectors(rec_with_fixed);
            w.close();
        end

        function testFixedArrays(testCase, format)
            ints = int32([[1, 2, 3]; [4, 5, 6]]);
            ints = transpose(ints);

            SR = @test_model.SimpleRecord;
            simple_recs = [...
                    [SR(x=1, y=2, z=3), SR(x=4, y=5, z=6)]; ...
                    [SR(x=11, y=12, z=13), SR(x=14, y=15, z=16)]; ...
                    [SR(x=21, y=22, z=23), SR(x=24, y=25, z=26)] ...
                ];
            simple_recs = transpose(simple_recs);

            RV = @test_model.RecordWithVlens;
            vlens_recs = [...
                    [...
                        RV(a=[SR(x=1, y=2, z=3), SR(x=7, y=8, z=9)], b=13, c=14), ...
                        RV(a=[SR(x=7, y=8, z=9), SR(x=10, y=11, z=12)], b=5, c=3)]; ...
                    [...
                        RV(a=[SR(x=31, y=32, z=33), SR(x=34, y=35, z=36), SR(x=37, y=38, z=39)], b=213, c=214), ...
                        RV(a=[SR(x=41, y=42, z=43)], b=313, c=314)] ...
                ];
            vlens_recs = transpose(vlens_recs);

            rec_with_fixed = test_model.RecordWithFixedArrays(...
                ints=ints, ...
                fixed_simple_record_array=simple_recs, ...
                fixed_record_with_vlens_array=vlens_recs);

            w = create_validating_writer(testCase, format, 'FixedArrays');
            w.write_ints(ints);
            w.write_fixed_simple_record_array(simple_recs);
            w.write_fixed_record_with_vlens_array(vlens_recs);
            w.write_record_with_fixed_arrays(rec_with_fixed);

            named_fixed_array = transpose(int32([[1, 2, 3, 4]; [5, 6, 7, 8]]));
            w.write_named_array(named_fixed_array)
            w.close();
        end

        function testSubarrays(testCase, format)
            w = create_validating_writer(testCase, format, 'Subarrays');

            ints = transpose(int32([[1, 2, 3]; [4, 5, 6]]));
            w.write_dynamic_with_fixed_int_subarray(ints);

            floats = transpose(single([[1, 2, 3]; [4, 5, 6]]));
            w.write_dynamic_with_fixed_float_subarray(floats);

            ints = transpose(int32([[1, 2, 3]; [4, 5, 6]]));
            w.write_known_dim_count_with_fixed_int_subarray(ints);

            floats = transpose(single([[1, 2, 3]; [4, 5, 6]]));
            w.write_known_dim_count_with_fixed_float_subarray(floats);

            ints = transpose(int32([[1, 2, 3]; [4, 5, 6]]));
            w.write_fixed_with_fixed_int_subarray(ints);

            floats = transpose(single([[1, 2, 3]; [4, 5, 6]]));
            w.write_fixed_with_fixed_float_subarray(floats);

            nested(:, :, 1) = int32([[1; 2; 3], [4; 5; 6]]);
            nested(:, :, 2) = [[10; 20; 30], [40; 50; 60]];
            nested(:, :, 3) = [[100; 200; 300], [400; 500; 600]];
            w.write_nested_subarray(nested);

            ints = transpose(int32([[1, 2, 3]; [4, 5, 6]]));
            w.write_dynamic_with_fixed_vector_subarray(ints);

            image(:, 1, 1) = int32([1; 2; 3]);
            image(:, 1, 2) = [4; 5; 6];
            image(:, 2, 1) = [10; 11; 12];
            image(:, 2, 2) = [13; 14; 15];
            w.write_generic_subarray(image);

            w.close();
        end

        function testSubarraysInRecords(testCase, format)
            RF = @test_model.RecordWithFixedCollections;
            records_with_fixed = [...
                RF(fixed_vector=int32([1, 2, 3]), fixed_array=int32([[11; 12; 13], [14; 15; 16]])), ...
                RF(fixed_vector=int32([101, 102, 103]), fixed_array=int32([[1011; 1012; 1013], [1014; 1015; 1016]])), ...
            ];

            RV = @test_model.RecordWithVlenCollections;
            records_with_vlens = [...
                RV(vector=int32([1, 2, 3]), array=int32([[11, 12, 13]; [14, 15, 16]])), ...])
                RV(vector=int32([101, 102, 103]), array=int32([[1011, 1012, 1013]; [1014, 1015, 1016]])), ...
            ];

            w = create_validating_writer(testCase, format, 'SubarraysInRecords');
            w.write_with_fixed_subarrays(records_with_fixed);
            w.write_with_vlen_subarrays(records_with_vlens);
            w.close();
        end

        function testNDArrays(testCase, format)
            ints = transpose(int32([[1, 2, 3]; [4, 5, 6]]));

            SR = @test_model.SimpleRecord;
            simple_recs = [...
                    [SR(x=1, y=2, z=3), SR(x=4, y=5, z=6)]; ...
                    [SR(x=11, y=12, z=13), SR(x=14, y=15, z=16)]; ...
                    [SR(x=21, y=22, z=23), SR(x=24, y=25, z=26)] ...
                ];
            simple_recs = transpose(simple_recs);

            RV = @test_model.RecordWithVlens;
            vlens_recs = [...
                    [...
                        RV(a=[SR(x=1, y=2, z=3), SR(x=7, y=8, z=9)], b=13, c=14), ...
                        RV(a=[SR(x=7, y=8, z=9), SR(x=10, y=11, z=12)], b=5, c=3)]; ...
                    [...
                        RV(a=[SR(x=31, y=32, z=33), SR(x=34, y=35, z=36), SR(x=37, y=38, z=39)], b=213, c=214), ...
                        RV(a=[SR(x=41, y=42, z=43)], b=313, c=314)] ...
                ];
            vlens_recs = transpose(vlens_recs);

            rec_with_nds = test_model.RecordWithNDArrays(...
                ints=ints, ...
                fixed_simple_record_array=simple_recs, ...
                fixed_record_with_vlens_array=vlens_recs);

            w = create_validating_writer(testCase, format, 'NDArrays');
            w.write_ints(ints);
            w.write_simple_record_array(simple_recs);
            w.write_record_with_vlens_array(vlens_recs);
            w.write_record_with_nd_arrays(rec_with_nds);
            w.write_named_array(ints);
            w.close();
        end

        function testDynamicNDArrays(testCase, format)
            ints = transpose(int32([[1, 2, 3]; [4, 5, 6]; [7, 8, 9]; [10, 11, 12]]));
            SR = @test_model.SimpleRecord;
            simple_recs = transpose([...
                    [SR(x=1, y=2, z=3), SR(x=4, y=5, z=6), SR(x=7, y=8, z=9)]; ...
                    [SR(x=11, y=12, z=13), SR(x=14, y=15, z=16), SR(x=17, y=18, z=19)]; ...
                ]);
            RV = @test_model.RecordWithVlens;
            vlens_recs = transpose([...
                    [RV(a=[SR(x=1, y=2, z=3)], b=-33, c=44)]; ...
                    [RV(a=[SR(x=8, y=2, z=9), SR(x=28, y=3, z=34)], b=233, c=347)] ...
                ]);
            rec_with_dynamic_nds = test_model.RecordWithDynamicNDArrays(...
                ints=ints, ...
                simple_record_array=simple_recs, ...
                record_with_vlens_array=vlens_recs);

            w = create_validating_writer(testCase, format, 'DynamicNDArrays');
            w.write_ints(ints);
            w.write_simple_record_array(simple_recs);
            w.write_record_with_vlens_array(vlens_recs);
            w.write_record_with_dynamic_nd_arrays(rec_with_dynamic_nds);
            w.close();
        end

        function testMultiDArrays(testCase, format)
            % ch=8, z=2, y=64, x=32
            img = zeros(32, 64, 2, 8, 'single');
            w = create_validating_writer(testCase, format, 'MultiDArrays');

            w.write_images({...
                img(:,:,:,1), ...
                img(:,:,1,1), ...
                img(:,1,1,:), ...
                img(1,1,1,1) ...
            });
            w.end_images();

            w.write_frames({...
                img(:,:,1,1), ...
                img(:,:,2,1) ...
            });
            w.end_frames();

            w.close();
        end

        function testMaps(testCase, format)
            d = dictionary();
            d("a") = int32(1);
            d("b") = 2;
            d("c") = 3;

            w = create_validating_writer(testCase, format, 'Maps');
            w.write_string_to_int(d);
            w.write_int_to_string(dictionary(int32([1, 2, 3]), ["a", "b", "c"]));
            w.write_string_to_union(...
                dictionary(...
                    ["a", "b"], ...
                    [test_model.StringOrInt32.Int32(1), test_model.StringOrInt32.String("2")] ...
                ) ...
            );
            w.write_aliased_generic(d);

            r1 = test_model.RecordWithMaps(set_1=dictionary(uint32(1), uint32(1), uint32(2), uint32(2)), set_2=dictionary(int32(-1), true, int32(3), false));
            r2 = test_model.RecordWithMaps(set_1=dictionary(uint32(1), uint32(2), uint32(2), uint32(1)), set_2=dictionary(int32(-1), false, int32(3), true));
            w.write_records([r1, r2]);
            w.close();

            % Now again for "empty" maps
            w = create_validating_writer(testCase, format, 'Maps');
            w.write_string_to_int(dictionary());
            w.write_int_to_string(dictionary());
            w.write_string_to_union(dictionary());
            w.write_aliased_generic(dictionary());
            w.write_records([test_model.RecordWithMaps(set_1=dictionary(), set_2=dictionary())]);

            w.close();
        end

        function testUnions(testCase, format)
            w = create_validating_writer(testCase, format, 'Unions');
            w.write_int_or_simple_record(test_model.Int32OrSimpleRecord.Int32(1));
            w.write_int_or_record_with_vlens(test_model.Int32OrRecordWithVlens.Int32(2));
            w.write_monosotate_or_int_or_simple_record(yardl.None);
            w.write_record_with_unions(basic_types.RecordWithUnions());
            w.close();

            w = create_validating_writer(testCase, format, 'Unions');
            w.write_int_or_simple_record(test_model.Int32OrSimpleRecord.SimpleRecord(test_model.SimpleRecord(x=1, y=2, z=3)));
            w.write_int_or_record_with_vlens(test_model.Int32OrRecordWithVlens.RecordWithVlens(test_model.RecordWithVlens(a=[test_model.SimpleRecord(x=1, y=2, z=3)], b=12, c=13)));
            w.write_monosotate_or_int_or_simple_record(test_model.Int32OrSimpleRecord.Int32(6));
            w.write_record_with_unions(basic_types.RecordWithUnions(...
                null_or_int_or_string= basic_types.Int32OrString.Int32(7), ...
                date_or_datetime= basic_types.TimeOrDatetime.Datetime(yardl.DateTime.from_components(2025, 3, 4, 12, 34, 56, 0)), ...
                null_or_fruits_or_days_of_week= basic_types.GenericNullableUnion2.T1(basic_types.Fruits.APPLE) ...
            ));
            w.close();
        end

        function testEnums(testCase, format)
            w = create_validating_writer(testCase, format, 'Enums');
            w.write_single(test_model.Fruits.APPLE);
            w.write_vec([test_model.Fruits.APPLE, test_model.Fruits.BANANA, test_model.Fruits(233983)]);
            w.write_size(test_model.SizeBasedEnum.C);
            w.write_rec(test_model.RecordWithEnums(...
                    enum=test_model.Fruits.PEAR, ...
                    rec=test_model.RecordWithNoDefaultEnum(enum=test_model.Fruits.BANANA)));
            w.close();
        end

        function testFlags(testCase, format)
            w = create_validating_writer(testCase, format, 'Flags');

            import test_model.DaysOfWeek;
            mon_or_wed_or_fri = DaysOfWeek(DaysOfWeek.MONDAY, DaysOfWeek.WEDNESDAY, DaysOfWeek.FRIDAY);
            w.write_days([...
                    test_model.DaysOfWeek.SUNDAY, ...
                    test_model.DaysOfWeek(mon_or_wed_or_fri), ...
                    test_model.DaysOfWeek(0), ...
                    test_model.DaysOfWeek(282839), ...
                    test_model.DaysOfWeek(234532) ...
                ]);
            w.end_days();

            w.write_formats([...
                    test_model.TextFormat.BOLD, ...
                    test_model.TextFormat.BOLD.with_flags(test_model.TextFormat.ITALIC), ...
                    test_model.TextFormat.REGULAR, ...
                    test_model.TextFormat(232932) ...
                ]);
            w.end_formats();

            w.close();
        end

        function testSimpleStreams(testCase, format)
            % Non-empty streams
            w = create_validating_writer(testCase, format, 'Streams');
            w.write_int_data(int32(1:10));
            w.write_int_data(int32(42));
            w.write_int_data(int32(1:20));
            w.end_int_data();
            w.write_optional_int_data([1, 2, yardl.None, 4, 5, yardl.None, 7, 8, 9, 10]);
            w.end_optional_int_data();
            w.write_record_with_optional_vector_data([...
                test_model.RecordWithOptionalVector(), ...
                test_model.RecordWithOptionalVector(optional_vector=int32([1, 2, 3])), ...
                test_model.RecordWithOptionalVector(optional_vector=int32(1:10)) ...
            ]);
            w.end_record_with_optional_vector_data();
            w.write_fixed_vector(...
                repelem({int32([1, 2, 3])}, 4)...
            );
            w.end_fixed_vector();
            w.close();

            % Mixed empty/non-empty streams
            w = create_validating_writer(testCase, format, 'Streams');
            w.write_int_data(int32(1:10))
            w.write_int_data([]);
            w.end_int_data();
            w.write_optional_int_data([]);
            w.write_optional_int_data([1, 2, yardl.None, 4, 5, yardl.None, 7, 8, 9, 10]);
            w.end_optional_int_data();
            w.write_record_with_optional_vector_data([]);
            w.end_record_with_optional_vector_data();
            % w.write_fixed_vector(int32(repmat([1;2;3], 1,4)));
            w.write_fixed_vector(repelem({int32([1, 2, 3])}, 4));
            w.end_fixed_vector();
            w.close();

            % % All empty streams
            w = create_validating_writer(testCase, format, 'Streams');
            w.write_int_data([]);
            w.end_int_data();
            w.write_optional_int_data([]);
            w.end_optional_int_data();
            w.write_record_with_optional_vector_data([]);
            w.end_record_with_optional_vector_data();
            w.write_fixed_vector([]);
            w.end_fixed_vector();
            w.close();

            w = create_validating_writer(testCase, format, 'Streams');
            w.end_int_data();
            w.end_optional_int_data();
            w.end_record_with_optional_vector_data();
            w.end_fixed_vector();
            w.close();
        end

        function testStreamsOfUnions(testCase, format)
            w = create_validating_writer(testCase, format, 'StreamsOfUnions');
            w.write_int_or_simple_record([...
                test_model.Int32OrSimpleRecord.Int32(1), ...
                test_model.Int32OrSimpleRecord.SimpleRecord(test_model.SimpleRecord(x=1, y=2, z=3)), ...
                test_model.Int32OrSimpleRecord.Int32(2) ...
            ]);
            w.end_int_or_simple_record();
            w.write_nullable_int_or_simple_record([...
                yardl.None, ...
                test_model.Int32OrSimpleRecord.Int32(1), ...
                test_model.Int32OrSimpleRecord.SimpleRecord(test_model.SimpleRecord(x=1, y=2, z=3)), ...
                yardl.None, ...
                test_model.Int32OrSimpleRecord.Int32(2), ...
                yardl.None ...
            ]);
            w.end_nullable_int_or_simple_record();
            w.write_many_cases([...
                test_model.Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray.Int32(3), ...
                test_model.Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray.Float32(7.0), ...
                test_model.Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray.String("hello"), ...
                test_model.Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray.SimpleRecord(test_model.SimpleRecord(x=1, y=2, z=3)), ...
                test_model.Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray.NamedFixedNDArray(transpose(int32([[1, 2, 3, 4]; [5, 6, 7, 8]]))), ...
            ]);
            w.end_many_cases();
            w.close();
        end

        function testStreamsOfAliasedUnions(testCase, format)
            w = create_validating_writer(testCase, format, 'StreamsOfAliasedUnions');
            w.write_int_or_simple_record([...
                test_model.AliasedIntOrSimpleRecord.Int32(1), ...
                test_model.AliasedIntOrSimpleRecord.SimpleRecord(test_model.SimpleRecord(x=1, y=2, z=3)), ...
                test_model.AliasedIntOrSimpleRecord.Int32(2) ...
            ]);
            w.end_int_or_simple_record();
            w.write_nullable_int_or_simple_record([...
                yardl.None, ...
                test_model.AliasedNullableIntSimpleRecord.Int32(1), ...
                test_model.AliasedNullableIntSimpleRecord.SimpleRecord(test_model.SimpleRecord(x=1, y=2, z=3)), ...
                yardl.None, ...
                test_model.AliasedNullableIntSimpleRecord.Int32(2), ...
                yardl.None ...
            ]);
            w.end_nullable_int_or_simple_record();
            w.close();
        end

        function testSimpleGenerics(testCase, format)
            w = create_validating_writer(testCase, format, 'SimpleGenerics');
            w.write_float_image(transpose(single([[1, 2]; [3, 4]])));
            w.write_int_image(transpose(int32([[1, 2]; [3, 4]])));
            w.write_int_image_alternate_syntax(transpose(int32([[1, 2]; [3, 4]])));
            w.write_string_image(transpose([["a", "b"]; ["c", "d"]]));
            w.write_int_float_tuple(test_model.MyTuple(v1=int32(1), v2=single(2)));
            w.write_float_float_tuple(test_model.MyTuple(v1=single(1), v2=single(2)));
            t = test_model.MyTuple(v1=int32(1), v2=single(2));
            w.write_int_float_tuple_alternate_syntax(t);
            w.write_int_string_tuple(test_model.MyTuple(v1=int32(1), v2="2"));
            w.write_stream_of_type_variants([...
                test_model.ImageFloatOrImageDouble.ImageFloat(transpose(single([[1, 2]; [3, 4]]))), ...
                test_model.ImageFloatOrImageDouble.ImageDouble(transpose([[1, 2]; [3, 4]])) ...
            ]);
            w.end_stream_of_type_variants();
            w.close();
        end

        function testAdvancedGenerics(testCase, format)
            w = create_validating_writer(testCase, format, 'AdvancedGenerics');

            i1 = single([[3, 4, 5]; [6, 7, 8]]);
            i2 = single([[30, 40, 50]; [60, 70, 80]]);
            i3 = single([[300, 400, 500]; [600, 700, 800]]);
            i4 = single([[3000, 4000, 5000]; [6000, 7000, 8000]]);

            img_img_array{1, 1} = i1;
            img_img_array{1, 2} = i2;
            img_img_array{2, 1} = i3;
            img_img_array{2, 2} = i4;

            w.write_float_image_image(img_img_array);

            w.write_generic_record_1(test_model.GenericRecord(...
                scalar_1=int32(1), ...
                scalar_2="hello", ...
                vector_1=int32([1, 2, 3]), ...
                image_2=[["abc", "def"]; ["a", "b"]] ...
            ));

            w.write_tuple_of_optionals(test_model.MyTuple(v1=yardl.None, v2="hello"));
            w.write_tuple_of_optionals_alternate_syntax(test_model.MyTuple(v1=int32(34), v2=yardl.None));
            w.write_tuple_of_vectors(test_model.MyTuple(v1=int32([1, 2, 3]), v2=single([4, 5, 6])));
            w.close();
        end

        function testAliases(testCase, format)
            w = create_validating_writer(testCase, format, 'Aliases');
            w.write_aliased_string("hello");
            w.write_aliased_enum(test_model.Fruits.APPLE);
            w.write_aliased_open_generic(test_model.AliasedOpenGeneric(v1="hello", v2=test_model.Fruits.BANANA));
            w.write_aliased_closed_generic(test_model.AliasedClosedGeneric(v1="hello", v2=test_model.Fruits.PEAR));
            w.write_aliased_optional(int32(23));
            w.write_aliased_generic_optional(single(44));
            w.write_aliased_generic_union_2(test_model.AliasedGenericUnion2.T1("hello"));
            w.write_aliased_generic_vector(single([1.0, 33.0, 44.0]));
            w.write_aliased_generic_fixed_vector(single([1.0, 33.0, 44.0]));
            w.write_stream_of_aliased_generic_union_2([ ...
                test_model.AliasedGenericUnion2.T1("hello"), ...
                test_model.AliasedGenericUnion2.T2(test_model.Fruits.APPLE) ...
            ]);
            w.end_stream_of_aliased_generic_union_2();
            w.close();
        end
    end
end


function w = create_validating_writer(testCase, format, protocol)
    writer_name = "test_model." + format + "." + protocol + "Writer";
    reader_name = "test_model." + format + "." + protocol + "Reader";
    test_writer_name = "test_model.testing.Test" + protocol + "Writer";

    create_writer = str2func(writer_name);
    create_reader = str2func(reader_name);
    create_test_writer = str2func(test_writer_name);

    w = create_test_writer(testCase, format, @(f) create_writer(f), @(f) create_reader(f));
end
