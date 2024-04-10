% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef CodedStreamTest < matlab.unittest.TestCase

    methods (Test)

        function testScalarByte(testCase)
            % w = yardl.binary.CodedOutputStream(tempname)
            % for i = 0:15
            %     w.
            % end
        end


        function testVarShort(testCase)
            entries = int16([0, 1, 5, 33, 0x7E, 0x7F, 0x80, 0x81, 255, 256, 257, 838, 0x3FFF, 0x4000, 0x4001, 0x7FFF]);

            filename = tempname;
            w = yardl.binary.CodedOutputStream(filename);
            for i = 1:length(entries)
                yardl.binary.Int16Serializer.write(w, entries(i));
                yardl.binary.Int16Serializer.write(w, -entries(i));
            end
            w.close();

            r = yardl.binary.CodedInputStream(filename);
            for i = 1:length(entries)
                testCase.verifyEqual(yardl.binary.Int16Serializer.read(r), entries(i));
                testCase.verifyEqual(yardl.binary.Int16Serializer.read(r), -entries(i));
            end
            r.close();
        end

        function testVarUShort(testCase)
            entries = uint16([0, 1, 5, 33, 0x7E, 0x7F, 0x80, 0x81, 255, 256, 257, 838, 0x3FFF, 0x4000, 0x4001, 0x7FFF, 0x8000, 0x8001, 0xFFFF]);

            filename = tempname;
            w = yardl.binary.CodedOutputStream(filename);
            for i = 1:length(entries)
                yardl.binary.Uint16Serializer.write(w, entries(i));
            end
            w.close();

            r = yardl.binary.CodedInputStream(filename);
            for i = 1:length(entries)
                testCase.verifyEqual(yardl.binary.Uint16Serializer.read(r), entries(i));
            end
            r.close();
        end

        function testVarIntegers(testCase)
            entries = uint32([ 0, 1, 5, 33, 0x7E, 0x7F, 0x80, 0x81, 255, 256, 257, 838, 283928, 2847772, 3443, 0x7FFFFFFF, 0xFFFFFFFF]);

            filename = tempname;
            w = yardl.binary.CodedOutputStream(filename);
            for i = 1:length(entries)
                yardl.binary.Uint32Serializer.write(w, entries(i));

                yardl.binary.Int32Serializer.write(w, int32(entries(i)));
                yardl.binary.Int32Serializer.write(w, -int32(entries(i)));

                yardl.binary.Uint64Serializer.write(w, uint64(entries(i)));
                yardl.binary.Uint64Serializer.write(w, bitor(uint64(entries(i)), uint64(0x800000000)));

                yardl.binary.Int64Serializer.write(w, int64(entries(i)));
                yardl.binary.Int64Serializer.write(w, -int64(entries(i)));
                yardl.binary.Int64Serializer.write(w, bitor(int64(entries(i)), int64(0x400000000)));
                yardl.binary.Int64Serializer.write(w, -bitor(int64(entries(i)), int64(0x400000000)));
            end
            w.close();

            r = yardl.binary.CodedInputStream(filename);
            for i = 1:length(entries)
                testCase.verifyEqual(yardl.binary.Uint32Serializer.read(r), entries(i));

                testCase.verifyEqual(yardl.binary.Int32Serializer.read(r), int32(entries(i)));
                testCase.verifyEqual(yardl.binary.Int32Serializer.read(r), -int32(entries(i)));

                testCase.verifyEqual(yardl.binary.Uint64Serializer.read(r), uint64(entries(i)));
                testCase.verifyEqual(yardl.binary.Uint64Serializer.read(r), bitor(uint64(entries(i)), uint64(0x800000000)));

                testCase.verifyEqual(yardl.binary.Int64Serializer.read(r), int64(entries(i)));
                testCase.verifyEqual(yardl.binary.Int64Serializer.read(r), -int64(entries(i)));
                testCase.verifyEqual(yardl.binary.Int64Serializer.read(r), bitor(int64(entries(i)), int64(0x400000000)));
                testCase.verifyEqual(yardl.binary.Int64Serializer.read(r), -bitor(int64(entries(i)), int64(0x400000000)));
            end
            r.close();
        end

    end
end
