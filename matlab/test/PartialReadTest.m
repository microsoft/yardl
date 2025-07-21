% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef PartialReadTest < matlab.unittest.TestCase

    methods (Static)
        function filename = generateStream()
            filename = tempname();
            w = test_model.binary.StateTestWriter(filename);
            w.write_an_int(42);
            w.write_a_stream([1, 2, 3, 4, 5]);
            w.end_a_stream();
            w.write_another_int(153);
            w.close();
        end
    end

    methods (Test)

        function testSkipSteps(testCase)
            filename = PartialReadTest.generateStream();
            r = test_model.binary.StateTestReader(filename, skip_completed_check=true);
            r.read_an_int();
            r.close();

            r = test_model.binary.StateTestReader(filename, skip_completed_check=true);
            r.read_an_int();
            while r.has_a_stream()
                item = r.read_a_stream();
            end
            r.close();
        end

        function testSkipStreamItems(testCase)
            filename = PartialReadTest.generateStream();
            r = test_model.binary.StateTestReader(filename, skip_completed_check=true);
            r.read_an_int();
            count = 0;
            while r.has_a_stream()
                item = r.read_a_stream();
                % Skip remaining stream items
                if count == 2
                    break;
                end
                count = count + 1;
            end
            r.close();
        end

    end
end
