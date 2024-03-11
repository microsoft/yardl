classdef EqualityTest < matlab.unittest.TestCase
    methods (Test)

        function testSimpleEquality(testCase)
            a = SimpleRecord(1, 2, 3);
            b = SimpleRecord(1, 2, 3);
            testCase.verifyEqual(a, b);

            c = SimpleRecord(1, 2, 4);
            testCase.verifyNotEqual(a, c);
        end

        function testFlagsEquality(testCase)
            a = bitor(DaysOfWeek.MONDAY, DaysOfWeek.TUESDAY);
            b = bitor(DaysOfWeek.TUESDAY, DaysOfWeek.MONDAY);
            testCase.verifyEqual(a, b);

            c = DaysOfWeek(0);
            d = DaysOfWeek(0);
            testCase.verifyEqual(c, d);
            testCase.verifyNotEqual(a, c);

            e = DaysOfWeek(0xFFFF);
            f = DaysOfWeek(0xFFFF);
            testCase.verifyEqual(e, f);
        end

        function testEnumEquality(testCase)
            a = Fruits.APPLE;
            b = Fruits.APPLE;
            testCase.verifyEqual(a, b);

            c = Fruits(10000);
            d = Fruits(10000);
            testCase.verifyEqual(c, d);
        end

        function testRecordWithEnumEquality(testCase)
            a = RecordWithEnums(Fruits.APPLE, bitor(DaysOfWeek.SATURDAY, DaysOfWeek.SUNDAY), 0);
            b = RecordWithEnums(Fruits.APPLE, bitor(DaysOfWeek.SATURDAY, DaysOfWeek.SUNDAY), 0);
            testCase.verifyEqual(a, b);

            c = RecordWithEnums(Fruits.APPLE, DaysOfWeek.SATURDAY, 0);
            testCase.verifyNotEqual(a, c);
        end

        function testDateEquality(testCase)
            a = RecordWithPrimitives();
            a.date_field = yardl.Date.from_components(2020, 1, 1);
            b = RecordWithPrimitives();
            b.date_field = yardl.Date.from_components(2020, 1, 1);

            testCase.verifyEqual(a, b);

            c = RecordWithPrimitives();
            c.date_field = yardl.Date.from_components(2020, 1, 2);
            testCase.verifyNotEqual(a, c);
        end

        function testTimeEquality(testCase)
            a = RecordWithPrimitives();
            a.time_field = yardl.Time.from_components(12, 22, 44, 0);
            b = RecordWithPrimitives();
            b.time_field = yardl.Time.from_components(12, 22, 44, 0);

            testCase.verifyEqual(a, b);

            c = RecordWithPrimitives();
            c.time_field = yardl.Time.from_components(12, 22, 45, 0);
            testCase.verifyNotEqual(a, c);
        end

        function testDateTimeEquality(testCase)
            a = RecordWithPrimitives();
            a.datetime_field = yardl.DateTime.from_components(2020, 1, 1, 12, 22, 44, 0);
            b = RecordWithPrimitives();
            b.datetime_field = yardl.DateTime.from_components(2020, 1, 1, 12, 22, 44, 0);

            testCase.verifyEqual(a, b);

            c = RecordWithPrimitives();
            c.datetime_field = yardl.DateTime.from_components(2020, 1, 1, 12, 22, 45, 0);
            testCase.verifyNotEqual(a, c);
        end

        function testStringEquality(testCase)
            a = RecordWithStrings("a", "b");
            b = RecordWithStrings("a", "b");
            testCase.verifyEqual(a, b);

            c = RecordWithStrings("a", "c");
            testCase.verifyNotEqual(a, c);
        end


        function testRecordWithPrimitiveVectorsEquality(testCase)
            a = RecordWithVectors( ...
                [1, 2], ...
                [1, 2, 3], ...
                [[1, 2], [3, 4]] ...
            );

            b = RecordWithVectors( ...
                [1, 2], ...
                [1, 2, 3], ...
                [[1, 2], [3, 4]] ...
            );

            testCase.verifyEqual(a, b);
        end

        function testOptionalIntEquality(testCase)
            a = RecordWithOptionalFields();
            a.optional_time = yardl.Time.from_components(1, 1, 1, 1);
            b = RecordWithOptionalFields();
            b.optional_time = yardl.Time.from_components(1, 1, 1, 1);
            testCase.verifyEqual(a, b);

            c = RecordWithOptionalFields();
            c.optional_time = yardl.Time.from_components(1, 1, 1, 2);
            testCase.verifyNotEqual(a, c);
            testCase.verifyNotEqual(b, c);

            d = RecordWithOptionalFields();
            e = RecordWithOptionalFields();
            testCase.verifyEqual(d, e);
            testCase.verifyNotEqual(a, d);
        end

        function testTimeVectorEquality(testCase)
            a = RecordWithVectorOfTimes(...
                [yardl.Time.from_components(1, 1, 1, 1), yardl.Time.from_components(1, 1, 1, 1)] ...
            );
            b = RecordWithVectorOfTimes(...
                [yardl.Time.from_components(1, 1, 1, 1), yardl.Time.from_components(1, 1, 1, 1)] ...
            );
            testCase.verifyEqual(a, b);

            c = RecordWithVectorOfTimes(...
                [yardl.Time.from_components(1, 1, 1, 1), yardl.Time.from_components(1, 1, 1, 2)] ...
            );
            testCase.verifyNotEqual(a, c);
        end

        function testSimpleArrayEquality(testCase)
            a = RecordWithArrays();
            a.default_array = int32([1, 2, 3]);
            b = RecordWithArrays();
            b.default_array = int32([1, 2, 3]);
            testCase.verifyEqual(a, b);

            c = RecordWithArrays();
            c.default_array = int32([1, 2, 4]);
            testCase.verifyNotEqual(a, c);
        end

        % function testSimpleUnionEquality(testCase)
        %     a = basic_types.RecordWithUnions();
        %     a.null_or_int_or_string = yardl.None;
        %     b = basic_types.RecordWithUnions();
        %     b.null_or_int_or_string = yardl.None;
        %     testCase.verifyEqual(a, b);

        %     c = basic_types.RecordWithUnions();
        %     c.null_or_int_or_string = basic_types.Int32OrString.Int32(1);
        %     d = basic_types.RecordWithUnions();
        %     d.null_or_int_or_string = basic_types.Int32OrString.Int32(1);
        %     testCase.verifyEqual(c, d);
        %     testCase.verifyNotEqual(a, c);

        %     e = basic_types.RecordWithUnions();
        %     e.null_or_int_or_string = basic_types.Int32OrString.String("hello");
        %     d = basic_types.RecordWithUnions();
        %     d.null_or_int_or_string = basic_types.Int32OrString.String("hello");
        %     testCase.verifyEqual(e, d);
        %     testCase.verifyNotEqual(a, e);
        %     testCase.verifyNotEqual(c, e);
        % end

    end
end
