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
            w.write_optional_record(test_model.SimpleRecord(8, 9, 10));
            w.write_record_with_optional_fields(...
                test_model.RecordWithOptionalFields(int32(44), int32(44), yardl.Time.from_components(12, 34, 56, 0)));
            w.write_optional_record_with_optional_fields(...
                test_model.RecordWithOptionalFields(int32(12), int32(12), yardl.Time.from_components(11, 32, 26, 0)));
            w.close();
        end

        function testNestedRecords(testCase, format)
            w = create_validating_writer(testCase, format, 'NestedRecords');
            t = test_model.TupleWithRecords(...
                test_model.SimpleRecord(1, 2, 3), ...
                test_model.SimpleRecord(4, 5, 6));
            w.write_tuple_with_records(t);
            w.close();
        end

        function testVariableLengthVectors(testCase, format)
            w = create_validating_writer(testCase, format, 'Vlens');
            w.write_int_vector(int32([1, 2, 3]));
            w.write_complex_vector(single([1+2j, 3+4j]));
            rec = test_model.RecordWithVlens(...
                [test_model.SimpleRecord(1, 2, 3), test_model.SimpleRecord(4, 5, 6)], ...
                4, ...
                2 ...
            );
            w.write_record_with_vlens(rec);
            w.write_vlen_of_record_with_vlens([rec, rec]);
            w.close();
        end

        function testStrings(testCase, format)
            w = create_validating_writer(testCase, format, 'Strings');
            w.write_single_string("hello");
            w.write_rec_with_string(test_model.RecordWithStrings("Montréal", "臺北市"));
            w.close();
        end

        function testOptionalVectors(testCase, format)
            w = create_validating_writer(testCase, format, 'OptionalVectors');
            w.write_record_with_optional_vector(test_model.RecordWithOptionalVector());
            w.close();

            w = create_validating_writer(testCase, format, 'OptionalVectors');
            w.write_record_with_optional_vector(test_model.RecordWithOptionalVector(int32([1, 2, 3])));
            w.close();
        end

        function testFixedVectors(testCase, format)
            ints = int32([1, 2, 3, 4, 5]);
            SR = @test_model.SimpleRecord;
            simple_recs = [SR(1, 2, 3), SR(4, 5, 6), SR(7, 8, 9)];
            RV = @test_model.RecordWithVlens;
            vlens_recs = [RV([SR(1, 2, 3), SR(4, 5, 6)], 4, 2), RV([SR(7, 8, 9), SR(10, 11, 12)], 5, 3)];
            rec_with_fixed = test_model.RecordWithFixedVectors(ints, simple_recs, vlens_recs);

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
                    [SR(1, 2, 3), SR(4, 5, 6)]; ...
                    [SR(11, 12, 13), SR(14, 15, 16)]; ...
                    [SR(21, 22, 23), SR(24, 25, 26)] ...
                ];
            simple_recs = transpose(simple_recs);

            RV = @test_model.RecordWithVlens;
            vlens_recs = [...
                    [...
                        RV([SR(1, 2, 3), SR(7, 8, 9)], 13, 14), ...
                        RV([SR(7, 8, 9), SR(10, 11, 12)], 5, 3)]; ...
                    [...
                        RV([SR(31, 32, 33), SR(34, 35, 36), SR(37, 38, 39)], 213, 214), ...
                        RV([SR(41, 42, 43)], 313, 314)] ...
                ];
            vlens_recs = transpose(vlens_recs);

            rec_with_fixed = test_model.RecordWithFixedArrays(ints, simple_recs, vlens_recs);

            w = create_validating_writer(testCase, format, 'FixedArrays');
            w.write_ints(ints);
            w.write_fixed_simple_record_array(simple_recs);
            w.write_fixed_record_with_vlens_array(vlens_recs);
            w.write_record_with_fixed_arrays(rec_with_fixed);

            % TODO: Like in Python, named fixed arrays are kind of broken since it
            % doesn't seem to be possible to specify the shape of the array in the type
            named_fixed_array = transpose(int32([[1, 2, 3, 4]; [5, 6, 7, 8]]));
            w.write_named_array(named_fixed_array)
            w.close();
        end

        function testSubarrays(testCase, format)
            % TODO: Gotta figure out the (Python) subarray logic (in Matlab)
            testCase.assumeFail();

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

            ints = transpose(int32([...
                    [[1, 2, 3]; [4, 5, 6]]; ...
                    [[10, 20, 30]; [40, 50, 60]]; ...
                    [[100, 200, 300]; [400, 500, 600]] ...
                ]));
            w.write_nested_subarray(ints);

            ints = transpose(int32([[1, 2, 3]; [4, 5, 6]]));
            w.write_dynamic_with_fixed_vector_subarray(ints);

            images = transpose(int32([...
                    [[1, 2, 3]; [4, 5, 6]]; ...
                    [[10, 11, 12]; [13, 14, 15]]; ...
                ]));
            w.write_generic_subarray(images);

            w.close();
        end

        function testSubarraysInRecords(testCase, format)
            % TODO: Gotta figure out the (Python) subarray logic (in Matlab)
            testCase.assumeFail();
        end

        function testNDArrays(testCase, format)
            ints = transpose(int32([[1, 2, 3]; [4, 5, 6]]));

            SR = @test_model.SimpleRecord;
            simple_recs = [...
                    [SR(1, 2, 3), SR(4, 5, 6)]; ...
                    [SR(11, 12, 13), SR(14, 15, 16)]; ...
                    [SR(21, 22, 23), SR(24, 25, 26)] ...
                ];
            simple_recs = transpose(simple_recs);

            RV = @test_model.RecordWithVlens;
            vlens_recs = [...
                    [...
                        RV([SR(1, 2, 3), SR(7, 8, 9)], 13, 14), ...
                        RV([SR(7, 8, 9), SR(10, 11, 12)], 5, 3)]; ...
                    [...
                        RV([SR(31, 32, 33), SR(34, 35, 36), SR(37, 38, 39)], 213, 214), ...
                        RV([SR(41, 42, 43)], 313, 314)] ...
                ];
            vlens_recs = transpose(vlens_recs);

            rec_with_nds = test_model.RecordWithNDArrays(ints, simple_recs, vlens_recs);

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
                    [SR(1, 2, 3), SR(4, 5, 6), SR(7, 8, 9)]; ...
                    [SR(11, 12, 13), SR(14, 15, 16), SR(17, 18, 19)]; ...
                ]);
            RV = @test_model.RecordWithVlens;
            vlens_recs = transpose([...
                    [RV([SR(1, 2, 3)], -33, 44)]; ...
                    [RV([SR(8, 2, 9), SR(28, 3, 34)], 233, 347)] ...
                ]);
            rec_with_dynamic_nds = test_model.RecordWithDynamicNDArrays(ints, simple_recs, vlens_recs);

            w = create_validating_writer(testCase, format, 'DynamicNDArrays');
            w.write_ints(ints);
            w.write_simple_record_array(simple_recs);
            w.write_record_with_vlens_array(vlens_recs);
            w.write_record_with_dynamic_nd_arrays(rec_with_dynamic_nds);
            w.close();
        end

        function testMaps(testCase, format)
            % d = containers.Map('KeyType', 'char', 'ValueType', 'int32');
            d = dictionary();
            d("a") = int32(1);
            d("b") = 2;
            d("c") = 3;

            w = create_validating_writer(testCase, format, 'Maps');
            w.write_string_to_int(d);

            % TODO: Need to use R2022b 'dictionary' to properly build the map in MapSerializer.read
            %       i.e. an "empty" containers.Map() has a KeyType of 'char'
            % w.write_int_to_string(containers.Map(int32([1, 2, 3]), ["a", "b", "c"]));
            w.write_int_to_string(dictionary(int32([1, 2, 3]), ["a", "b", "c"]));

            % TODO: Need to use R2022b `dictionary` to store objects as values...
            % w.write_string_to_union(...
            %     containers.Map(...
            %         ["a", "b"], ...
            %         [test_model.StringOrInt32.Int32(1), test_model.StringOrInt32.String("2")] ...
            %     ) ...
            % );
            w.write_string_to_union(...
                dictionary(...
                    ["a", "b"], ...
                    [test_model.StringOrInt32.Int32(1), test_model.StringOrInt32.String("2")] ...
                ) ...
            );

            w.write_aliased_generic(d);
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
            w.write_int_or_simple_record(test_model.Int32OrSimpleRecord.SimpleRecord(test_model.SimpleRecord(1, 2, 3)));
            w.write_int_or_record_with_vlens(test_model.Int32OrRecordWithVlens.RecordWithVlens(test_model.RecordWithVlens([test_model.SimpleRecord(1, 2, 3)], 12, 13)));
            w.write_monosotate_or_int_or_simple_record(test_model.Int32OrSimpleRecord.Int32(6));
            w.write_record_with_unions(basic_types.RecordWithUnions(...
                basic_types.Int32OrString.Int32(7), ...
                basic_types.TimeOrDatetime.Datetime(yardl.DateTime.from_components(2025, 3, 4, 12, 34, 56, 0)), ...
                basic_types.GenericNullableUnion2.T1(basic_types.Fruits.APPLE) ...
            ));
            w.close();
        end

        function testEnums(testCase, format)
            w = create_validating_writer(testCase, format, 'Enums');
            w.write_single(test_model.Fruits.APPLE);
            w.write_vec([test_model.Fruits.APPLE, test_model.Fruits.BANANA, test_model.Fruits(233983)]);
            w.write_size(test_model.SizeBasedEnum.C);
            w.close();
        end

        function testFlags(testCase, format)
            w = create_validating_writer(testCase, format, 'Flags');

            mon_or_wed_or_fri = bitor(bitor(test_model.DaysOfWeek.MONDAY, test_model.DaysOfWeek.WEDNESDAY), test_model.DaysOfWeek.FRIDAY);
            w.write_days([...
                    test_model.DaysOfWeek.SUNDAY, ...
                    test_model.DaysOfWeek(mon_or_wed_or_fri), ...
... % bitor(bitor(test_model.DaysOfWeek.MONDAY, test_model.DaysOfWeek.WEDNESDAY), test_model.DaysOfWeek.FRIDAY), ...
                    test_model.DaysOfWeek(0), ...
                    test_model.DaysOfWeek(282839), ...
                    test_model.DaysOfWeek(234532) ...
                ]);

            w.write_formats([...
                    test_model.TextFormat.BOLD, ...
                    bitor(test_model.TextFormat.BOLD, test_model.TextFormat.ITALIC), ...
                    test_model.TextFormat.REGULAR, ...
                    test_model.TextFormat(232932) ...
                ]);

            w.close();
        end

        function testSimpleStreams(testCase, format)
            % TODO: Need to fix Stream read/write of multidimensional arrays (see write_fixed_vector)
            testCase.assumeFail();

            w = create_validating_writer(testCase, format, 'Streams');

            % TODO: MockWriter can't yet handle multiple stream write calls (the result of read is always the full concatenated stream)
            % w.write_int_data(int32(1:10));
            % w.write_int_data(int32(1:20));
            w.write_int_data(int32(1:30));

            w.write_optional_int_data([1, 2, yardl.None, 4, 5, yardl.None, 7, 8, 9, 10]);
            w.write_record_with_optional_vector_data([...
                test_model.RecordWithOptionalVector(), ...
                test_model.RecordWithOptionalVector(int32([1, 2, 3])), ...
                test_model.RecordWithOptionalVector(int32(1:10)) ...
            ]);
            w.write_fixed_vector(...
                transpose([...
                    [1, 2, 3]; [1, 2, 3]; [1, 2, 3]; [1, 2, 3];
                ])...
            );
            w.close();

            w = create_validating_writer(testCase, format, 'Streams');
            w.write_int_data([]);
            w.write_optional_int_data([]);
            w.write_record_with_optional_vector_data([]);
            w.write_fixed_vector([]);
            w.close();
        end

        function testStreamsOfUnions(testCase, format)
            w = create_validating_writer(testCase, format, 'StreamsOfUnions');
            w.write_int_or_simple_record([...
                test_model.Int32OrSimpleRecord.Int32(1), ...
                test_model.Int32OrSimpleRecord.SimpleRecord(test_model.SimpleRecord(1, 2, 3)), ...
                test_model.Int32OrSimpleRecord.Int32(2) ...
            ]);
            w.write_nullable_int_or_simple_record([...
                yardl.None, ...
                test_model.Int32OrSimpleRecord.Int32(1), ...
                test_model.Int32OrSimpleRecord.SimpleRecord(test_model.SimpleRecord(1, 2, 3)), ...
                yardl.None, ...
                test_model.Int32OrSimpleRecord.Int32(2), ...
                yardl.None ...
            ]);
            w.close();
        end

        function testStreamsOfAliasedUnions(testCase, format)
            w = create_validating_writer(testCase, format, 'StreamsOfAliasedUnions');
            w.write_int_or_simple_record([...
                test_model.AliasedIntOrSimpleRecord.Int32(1), ...
                test_model.AliasedIntOrSimpleRecord.SimpleRecord(test_model.SimpleRecord(1, 2, 3)), ...
                test_model.AliasedIntOrSimpleRecord.Int32(2) ...
            ]);
            w.write_nullable_int_or_simple_record([...
                yardl.None, ...
                test_model.AliasedNullableIntSimpleRecord.Int32(1), ...
                test_model.AliasedNullableIntSimpleRecord.SimpleRecord(test_model.SimpleRecord(1, 2, 3)), ...
                yardl.None, ...
                test_model.AliasedNullableIntSimpleRecord.Int32(2), ...
                yardl.None ...
            ]);
            w.close();
        end

        function testSimpleGenerics(testCase, format)
            w = create_validating_writer(testCase, format, 'SimpleGenerics');
            w.write_float_image(transpose(single([[1, 2]; [3, 4]])));
            w.write_int_image(transpose(int32([[1, 2]; [3, 4]])));
            w.write_int_image_alternate_syntax(transpose(int32([[1, 2]; [3, 4]])));
            w.write_string_image(transpose([["a", "b"]; ["c", "d"]]));
            w.write_int_float_tuple(test_model.MyTuple(int32(1), single(2)));
            w.write_float_float_tuple(test_model.MyTuple(single(1), single(2)));
            t = test_model.MyTuple(int32(1), single(2));
            w.write_int_float_tuple_alternate_syntax(t);
            w.write_int_string_tuple(test_model.MyTuple(int32(1), "2"));
            w.write_stream_of_type_variants([...
                test_model.ImageFloatOrImageDouble.ImageFloat(transpose(single([[1, 2]; [3, 4]]))), ...
                test_model.ImageFloatOrImageDouble.ImageDouble(transpose([[1, 2]; [3, 4]])) ...
            ]);
            w.close();
        end

        function testAdvancedGenerics(testCase, format)
            % TODO: Fix NDArraySerializers for nested NDArrays
            testCase.assumeFail();

            w = create_validating_writer(testCase, format, 'AdvancedGenerics');

            i1 = single([[3, 4, 5]; [6, 7, 8]]);
            i2 = single([[30, 40, 50]; [60, 70, 80]]);
            i3 = single([[300, 400, 500]; [600, 700, 800]]);
            i4 = single([[3000, 4000, 5000]; [6000, 7000, 8000]]);

            % TODO: How would I declare it with zeros first, then fill it with img_img_array[:] = ... ?  e.g. (in Python)
            img_img_array = i1;
            img_img_array(:, :, 2) = i2;
            img_img_array(:, :, 3) = i3;
            img_img_array(:, :, 4) = i4;

            w.write_float_image_image(img_img_array);

            w.write_generic_record_1(test_model.GenericRecord(...
                int32(1), ...
                "hello", ...
                int32([1, 2, 3]), ...
                [["abc", "def"]; ["a", "b"]] ...
            ));

            w.write_tuple_of_optionals(test_model.MyTuple(yardl.None, "hello"));
            w.write_tuple_of_optionals_alternate_syntax(test_model.MyTuple(int32(34), yardl.None));
            w.write_tuple_of_vectors(test_model.MyTuple(int32([1, 2, 3]), single([4, 5, 6])));
            w.close();
        end

        function testAliases(testCase, format)
            w = create_validating_writer(testCase, format, 'Aliases');
            w.write_aliased_string("hello");
            w.write_aliased_enum(test_model.Fruits.APPLE);
            w.write_aliased_open_generic(test_model.AliasedOpenGeneric("hello", test_model.Fruits.BANANA));
            w.write_aliased_closed_generic(test_model.AliasedClosedGeneric("hello", test_model.Fruits.PEAR));
            w.write_aliased_optional(int32(23));
            w.write_aliased_generic_optional(single(44));
            w.write_aliased_generic_union_2(test_model.AliasedGenericUnion2.T1("hello"));
            w.write_aliased_generic_vector(single([1.0, 33.0, 44.0]));
            w.write_aliased_generic_fixed_vector(single([1.0, 33.0, 44.0]));
            w.write_stream_of_aliased_generic_union_2([ ...
                test_model.AliasedGenericUnion2.T1("hello"), ...
                test_model.AliasedGenericUnion2.T2(test_model.Fruits.APPLE) ...
            ]);
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

    testfile = tempname;
    writer = create_writer(testfile);
    w = create_test_writer(testCase, writer, @() create_reader(testfile));
end
