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

    end
end
