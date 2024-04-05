classdef YardlTypesTest < matlab.unittest.TestCase

    methods (Test)

        function testDateTimeFromValidDatetime(testCase)
            dt = yardl.DateTime.from_datetime(datetime(2020, 2, 29, 12, 22, 44, .111222));
            testCase.verifyEqual(dt.value(), int64(1582978964000111222));
        end

        function testDateTimeFromValidComponents(testCase)
            dt = yardl.DateTime.from_components(2020, 2, 29, 12, 22, 44, 111222);
            testCase.verifyEqual(dt.value(), int64(1582978964000111222));
        end

        function testDateTimeFromInvalidComponents(testCase)
            % nanosecond out of range
            testCase.verifyError(@() yardl.DateTime.from_components(2021, 2, 15, 12, 22, 44, 9999999999999999), 'yardl:ValueError');
        end

        function testDateTimeFromInt(testCase)
            dt = yardl.DateTime(1577967764111222333);
            mdt = dt.to_datetime();
            testCase.verifyEqual(mdt.Year, 2020);
            testCase.verifyEqual(mdt.Month, 1);
            testCase.verifyEqual(mdt.Day, 2);
            testCase.verifyEqual(mdt.Hour, 12);
            testCase.verifyEqual(mdt.Minute, 22);
            testCase.verifyEqual(mdt.Second * 1e9, 44111222333, 'RelTol', 1e-8);
        end

        function testDateFromValidDatetime(testCase)
            d = yardl.Date.from_datetime(datetime(2020, 2, 29));
            testCase.verifyEqual(d.value(), 18321);
        end

        function testDateFromValidComponents(testCase)
            d = yardl.Date.from_components(2020, 2, 29);
            md = d.to_datetime();
            testCase.verifyEqual(md.Year, 2020);
            testCase.verifyEqual(md.Month, 2);
            testCase.verifyEqual(md.Day, 29);
        end

        function testTimeFromValidDatetime(testCase)
            t = yardl.Time.from_datetime(datetime(2024, 3, 11, 12, 22, 44, 111.222));
            mdt = t.to_datetime();
            testCase.verifyEqual(mdt.Hour, 12);
            testCase.verifyEqual(mdt.Minute, 22);
            testCase.verifyEqual(mdt.Second * 1e9, 44111222000, 'RelTol', 1e-8);
        end

        function testTimeFromValidComponents(testCase)
            t = yardl.Time.from_components(12, 22, 44, 111222333);
            mdt = t.to_datetime();
            testCase.verifyEqual(mdt.Hour, 12);
            testCase.verifyEqual(mdt.Minute, 22);
            testCase.verifyEqual(mdt.Second * 1e9, 44111222333, 'RelTol', 1e-8);
        end

        function testTimeFromInvalidComponents(testCase)
            % hour out of range
            testCase.verifyError(@() yardl.Time.from_components(24, 0, 0, 0), 'yardl:ValueError');

            % minute out of range
            testCase.verifyError(@() yardl.Time.from_components(12, -3, 0, 0), 'yardl:ValueError');

            % second out of range
            testCase.verifyError(@() yardl.Time.from_components(12, 22, 74, 0), 'yardl:ValueError');

            % nanosecond out of range
            testCase.verifyError(@() yardl.Time.from_components(12, 22, 44, 9999999999999999), 'yardl:ValueError');
        end

        function testOptionals(testCase)
            % None.has_value == false
            testCase.verifyEqual(yardl.None().has_value, false);
            % Optional(X).has_value == true
            testCase.verifyEqual(yardl.Optional(42).has_value, true);

            % ERROR: None.value
            testCase.verifyError(@() yardl.None().value, 'yardl:ValueError');

            % ERROR: None.has_value = true
            function setHasValueTrue()
                o = yardl.None;
                o.has_value = true;
            end
            testCase.verifyError(@setHasValueTrue, 'MATLAB:class:SetProhibited');

            % ERROR: Optional(x).has_value = false
            function setHasValueFalse()
                o = yardl.Optional(42);
                o.has_value = false;
            end
            testCase.verifyError(@setHasValueFalse, 'MATLAB:class:SetProhibited');

            % None == None
            testCase.verifyEqual(yardl.None, yardl.None);

            % Optional(X) == Optional(X)
            testCase.verifyEqual(yardl.Optional(int16(42)), yardl.Optional(int16(42)));
            testCase.verifyEqual(yardl.Optional("hello"), yardl.Optional("hello"));
            testCase.verifyEqual(yardl.Optional([]), yardl.Optional([]));
            testCase.verifyEqual(yardl.Optional([[1, 2]; [3, 4]]), yardl.Optional([[1, 2]; [3, 4]]));

            % Optional(X) == X
            testCase.verifyTrue(yardl.Optional(int16(42)) == int16(42));
            testCase.verifyTrue(yardl.Optional("hello") == "hello");
            testCase.verifyTrue(isequal(yardl.Optional([]), []));
            testCase.verifyTrue(isequal(yardl.Optional([[1, 2]; [3, 4]]), [[1, 2]; [3, 4]]));

            % X == Optional(X)
            testCase.verifyEqual(int16(42), yardl.Optional(int16(42)));
            testCase.verifyEqual("hello", yardl.Optional("hello"));
            testCase.verifyEqual([], yardl.Optional([]));
            testCase.verifyEqual([[1, 2]; [3, 4]], yardl.Optional([[1, 2]; [3, 4]]));

            % None ~= Optional(X)
            testCase.verifyNotEqual(yardl.None, yardl.Optional(42));

            % Optional(X) ~= None
            testCase.verifyNotEqual(yardl.Optional(42), yardl.None);

            % None ~= X
            testCase.verifyNotEqual(yardl.None, 42);
            % X ~= None
            testCase.verifyNotEqual(42, yardl.None);

            % [Optionals] == [Optionals]
            os = arrayfun(@yardl.Optional, 1:5);
            testCase.verifyEqual(os, os);

            % [Nones] == [Nones]
            nones = repelem(yardl.None, 5);
            testCase.verifyEqual(nones, nones);

            % [Optional, None, ...] == [Optional, None, ...]
            mixed1 = [1, 2, yardl.None, 4, 5, yardl.None];
            testCase.verifyEqual(mixed1, mixed1);
            mixed2 = [yardl.Optional(1), yardl.Optional(2), yardl.None, yardl.Optional(4), yardl.Optional(5), yardl.None];
            testCase.verifyEqual(mixed1, mixed2);
            testCase.verifyEqual(mixed2, mixed1);

            % [Optional, None, ...] ~= [None, Optional, ...]
            mixed3 = [1, 2, yardl.None, 4, 5, 6];
            testCase.verifyNotEqual(mixed1, mixed3);
            testCase.verifyNotEqual(mixed3, mixed1);
        end
    end
end
