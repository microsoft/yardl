classdef EqualityTest < matlab.unittest.TestCase
    methods (Test)

        % TODO: Add tests for equality of *arrays* of each type,
        %   since Matlab eq method applies to both scalar and non-scalar values

        function testSimpleEquality(testCase)
            a = test_model.SimpleRecord(1, 2, 3);
            b = test_model.SimpleRecord(1, 2, 3);
            testCase.verifyEqual(a, b);

            c = test_model.SimpleRecord(1, 2, 4);
            testCase.verifyNotEqual(a, c);
        end

        function testFlagsEquality(testCase)
            a = bitor(test_model.DaysOfWeek.MONDAY, test_model.DaysOfWeek.TUESDAY);
            b = bitor(test_model.DaysOfWeek.TUESDAY, test_model.DaysOfWeek.MONDAY);
            testCase.verifyEqual(a, b);

            c = test_model.DaysOfWeek(0);
            d = test_model.DaysOfWeek(0);
            testCase.verifyEqual(c, d);
            testCase.verifyNotEqual(a, c);

            e = test_model.DaysOfWeek(0xFFFF);
            f = test_model.DaysOfWeek(0xFFFF);
            testCase.verifyEqual(e, f);
        end

        function testEnumEquality(testCase)
            a = test_model.Fruits.APPLE;
            b = test_model.Fruits.APPLE;
            testCase.verifyEqual(a, b);

            c = test_model.Fruits(10000);
            d = test_model.Fruits(10000);
            testCase.verifyEqual(c, d);
        end

        function testRecordWithEnumEquality(testCase)
            a = test_model.RecordWithEnums(test_model.Fruits.APPLE, bitor(test_model.DaysOfWeek.SATURDAY, test_model.DaysOfWeek.SUNDAY), 0);
            b = test_model.RecordWithEnums(test_model.Fruits.APPLE, bitor(test_model.DaysOfWeek.SATURDAY, test_model.DaysOfWeek.SUNDAY), 0);
            testCase.verifyEqual(a, b);

            c = test_model.RecordWithEnums(test_model.Fruits.APPLE, test_model.DaysOfWeek.SATURDAY, 0);
            testCase.verifyNotEqual(a, c);
        end

        function testDateEquality(testCase)
            a = test_model.RecordWithPrimitives();
            a.date_field = yardl.Date.from_components(2020, 1, 1);
            b = test_model.RecordWithPrimitives();
            b.date_field = yardl.Date.from_components(2020, 1, 1);

            testCase.verifyEqual(a, b);

            c = test_model.RecordWithPrimitives();
            c.date_field = yardl.Date.from_components(2020, 1, 2);
            testCase.verifyNotEqual(a, c);
        end

        function testTimeEquality(testCase)
            a = test_model.RecordWithPrimitives();
            a.time_field = yardl.Time.from_components(12, 22, 44, 0);
            b = test_model.RecordWithPrimitives();
            b.time_field = yardl.Time.from_components(12, 22, 44, 0);

            testCase.verifyEqual(a, b);

            c = test_model.RecordWithPrimitives();
            c.time_field = yardl.Time.from_components(12, 22, 45, 0);
            testCase.verifyNotEqual(a, c);
        end

        function testDateTimeEquality(testCase)
            a = test_model.RecordWithPrimitives();
            a.datetime_field = yardl.DateTime.from_components(2020, 1, 1, 12, 22, 44, 0);
            b = test_model.RecordWithPrimitives();
            b.datetime_field = yardl.DateTime.from_components(2020, 1, 1, 12, 22, 44, 0);

            testCase.verifyEqual(a, b);

            c = test_model.RecordWithPrimitives();
            c.datetime_field = yardl.DateTime.from_components(2020, 1, 1, 12, 22, 45, 0);
            testCase.verifyNotEqual(a, c);
        end

        function testStringEquality(testCase)
            a = test_model.RecordWithStrings("a", "b");
            b = test_model.RecordWithStrings("a", "b");
            testCase.verifyEqual(a, b);

            c = test_model.RecordWithStrings("a", "c");
            testCase.verifyNotEqual(a, c);
        end


        function testRecordWithPrimitiveVectorsEquality(testCase)
            a = test_model.RecordWithVectors( ...
                [1, 2], ...
                [1, 2, 3], ...
                [[1, 2], [3, 4]] ...
            );

            b = test_model.RecordWithVectors( ...
                [1, 2], ...
                [1, 2, 3], ...
                [[1, 2], [3, 4]] ...
            );

            testCase.verifyEqual(a, b);
        end

        function testOptionalIntEquality(testCase)
            a = test_model.RecordWithOptionalFields();
            a.optional_time = yardl.Time.from_components(1, 1, 1, 1);
            b = test_model.RecordWithOptionalFields();
            b.optional_time = yardl.Time.from_components(1, 1, 1, 1);
            testCase.verifyEqual(a, b);

            c = test_model.RecordWithOptionalFields();
            c.optional_time = yardl.Time.from_components(1, 1, 1, 2);
            testCase.verifyNotEqual(a, c);
            testCase.verifyNotEqual(b, c);

            d = test_model.RecordWithOptionalFields();
            e = test_model.RecordWithOptionalFields();
            testCase.verifyEqual(d, e);
            testCase.verifyNotEqual(a, d);
        end

        function testTimeVectorEquality(testCase)
            a = test_model.RecordWithVectorOfTimes(...
                [yardl.Time.from_components(1, 1, 1, 1), yardl.Time.from_components(1, 1, 1, 1)] ...
            );
            b = test_model.RecordWithVectorOfTimes(...
                [yardl.Time.from_components(1, 1, 1, 1), yardl.Time.from_components(1, 1, 1, 1)] ...
            );
            testCase.verifyEqual(a, b);

            c = test_model.RecordWithVectorOfTimes(...
                [yardl.Time.from_components(1, 1, 1, 1), yardl.Time.from_components(1, 1, 1, 2)] ...
            );
            testCase.verifyNotEqual(a, c);
        end

        function testSimpleArrayEquality(testCase)
            a = test_model.RecordWithArrays();
            a.default_array = int32([1, 2, 3]);
            b = test_model.RecordWithArrays();
            b.default_array = int32([1, 2, 3]);
            testCase.verifyEqual(a, b);

            c = test_model.RecordWithArrays();
            c.default_array = int32([1, 2, 4]);
            testCase.verifyNotEqual(a, c);
        end

        function testSimpleUnionEquality(testCase)
            a = basic_types.RecordWithUnions();
            a.null_or_int_or_string = yardl.None;
            b = basic_types.RecordWithUnions();
            b.null_or_int_or_string = yardl.None;
            testCase.verifyEqual(a, b);

            c = basic_types.RecordWithUnions();
            c.null_or_int_or_string = basic_types.Int32OrString.Int32(1);
            d = basic_types.RecordWithUnions();
            d.null_or_int_or_string = basic_types.Int32OrString.Int32(1);
            testCase.verifyEqual(c, d);
            testCase.verifyNotEqual(a, c);

            e = basic_types.RecordWithUnions();
            e.null_or_int_or_string = basic_types.Int32OrString.String("hello");
            d = basic_types.RecordWithUnions();
            d.null_or_int_or_string = basic_types.Int32OrString.String("hello");
            testCase.verifyEqual(e, d);
            testCase.verifyNotEqual(a, e);
            testCase.verifyNotEqual(c, e);
        end

        function testTimeUnionEquality(testCase)
            a = basic_types.RecordWithUnions();
            a.date_or_datetime = basic_types.TimeOrDatetime.Time(yardl.Time.from_components(1, 1, 1, 1));
            b = basic_types.RecordWithUnions();
            b.date_or_datetime = basic_types.TimeOrDatetime.Time(yardl.Time.from_components(1, 1, 1, 1));
            testCase.verifyEqual(a, b);
        end

        function testGenericEquality(testCase)
            a = test_model.GenericRecord(1, 2.0, [1, 2, 3], single([[1.1, 2.2], [3.3, 4.4]]));
            b = test_model.GenericRecord(1, 2.0, [1, 2, 3], single([[1.1, 2.2], [3.3, 4.4]]));
            testCase.verifyEqual(a, b);

            c = test_model.MyTuple(42.0, "hello, world");
            d = tuples.Tuple(42.0, "hello, world");
            testCase.verifyTrue(c == d)

            e = test_model.AliasedTuple(42.0, "hello, world");
            testCase.verifyTrue(c == e);
        end

    end
end
