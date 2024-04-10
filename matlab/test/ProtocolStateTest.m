% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef ProtocolStateTest < matlab.unittest.TestCase
    methods (Test)

        function testProperSequenceWrite(testCase)
            w = ProtocolStateTestWriter();
            w.write_an_int(1);
            w.write_a_stream([1, 2, 3]);
            w.write_another_int(3);
            w.close();
        end

        function testProperSequenceWriteEmptyStream(testCase)
            w = ProtocolStateTestWriter();
            w.write_an_int(1);
            w.write_a_stream([]);
            w.write_another_int(3);
            w.close();
        end

        function testProperSequenceWriteMultipleStreams(testCase)
            w = ProtocolStateTestWriter();
            w.write_an_int(1);
            w.write_a_stream([1, 2, 3]);
            w.write_a_stream([4, 5, 6]);
            w.write_another_int(3);
            w.close();
        end

        function testSequenceWriteMissingFirstStep(testCase)
            w = ProtocolStateTestWriter();
            testCase.verifyError(@() w.write_a_stream([1, 2, 3]), 'yardl:ProtocolError');
        end

        function testSequenceWritePrematureClose(testCase)
            w = ProtocolStateTestWriter();
            w.write_an_int(1);
            w.write_a_stream([1, 2, 3]);
            testCase.verifyError(@() w.close(), 'yardl:ProtocolError');
        end


        function testProperSequenceRead(testCase)
            r = ProtocolStateTestReader();
            r.read_an_int();
            r.read_a_stream();
            r.read_another_int();
            r.close();
        end

        function testReadSkipFirstStep(testCase)
            r = ProtocolStateTestReader();
            testCase.verifyError(@() r.read_a_stream(), 'yardl:ProtocolError');
        end

        function testReadCloseEarly(testCase)
            r = ProtocolStateTestReader();
            r.read_an_int();
            r.read_a_stream();
            testCase.verifyError(@() r.close(), 'yardl:ProtocolError');
        end

    end
end
