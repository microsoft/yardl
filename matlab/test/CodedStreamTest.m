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

        function testFloats(testCase)
            filename = tempname;

            f = single(42.75);
            d = double(-42.75);
            fs = rand(123, 89, 'single');
            ds = rand(123, 89);

            cf = single(complex(7, -8));
            df = complex(7, -8);
            cfs = complex(rand(123, 89, 'single'), rand(123, 89, 'single'));
            dfs = complex(rand(123, 89), rand(123, 89));

            float32Serializer = yardl.binary.Float32Serializer;
            float64Serializer = yardl.binary.Float64Serializer;
            complexFloat32Serializer = yardl.binary.Complexfloat32Serializer;
            complexFloat64Serializer = yardl.binary.Complexfloat64Serializer;

            w = yardl.binary.CodedOutputStream(filename);
            float32Serializer.write(w, f);
            float64Serializer.write(w, d);
            float32Serializer.writeTrivially(w, fs);
            float64Serializer.writeTrivially(w, ds);
            complexFloat32Serializer.write(w, cf);
            complexFloat64Serializer.write(w, df);
            complexFloat32Serializer.writeTrivially(w, cfs);
            complexFloat64Serializer.writeTrivially(w, dfs);
            w.close();

            r = yardl.binary.CodedInputStream(filename);
            testCase.verifyEqual(float32Serializer.read(r), f);
            testCase.verifyEqual(float64Serializer.read(r), d);
            testCase.verifyEqual(float32Serializer.readTrivially(r, size(fs)), fs);
            testCase.verifyEqual(float64Serializer.readTrivially(r, size(ds)), ds);
            testCase.verifyEqual(complexFloat32Serializer.read(r), cf);
            testCase.verifyEqual(complexFloat64Serializer.read(r), df);
            testCase.verifyEqual(complexFloat32Serializer.readTrivially(r, size(cfs)), cfs);
            testCase.verifyEqual(complexFloat64Serializer.readTrivially(r, size(dfs)), dfs);
            r.close();
        end
    end
end
