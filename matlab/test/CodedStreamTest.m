% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef CodedStreamTest < matlab.unittest.TestCase

    methods (Test)

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

            delete(filename);
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

            delete(filename);
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

            delete(filename);
        end

        function testBytes(testCase)
            filename = tempname;

            val = 42;
            unsigned = uint8([0, 1, 5, 33, 0x7E, 0x7F, 0x80, 0x81, 255]);
            signed = int8([0, 1, 5, 33, 0x7E, 0x7F, -0x80, -0x7F, -1]);

            unsignedSerializer = yardl.binary.Uint8Serializer;
            signedSerializer = yardl.binary.Int8Serializer;

            w = yardl.binary.CodedOutputStream(filename);
            w.write_byte(uint8(val));
            testCase.verifyError(@() w.write_byte(-1), 'MATLAB:validators:mustBeA');
            testCase.verifyError(@() w.write_byte(256), 'MATLAB:validators:mustBeA');

            w.write_bytes(unsigned);
            testCase.verifyError(@() w.write_bytes(signed), 'MATLAB:validators:mustBeA');

            unsignedSerializer.write(w, val);
            unsignedSerializer.write_trivially(w, unsigned);

            signedSerializer.write(w, -val);
            signedSerializer.write_trivially(w, signed);
            w.close();

            r = yardl.binary.CodedInputStream(filename);
            testCase.verifyEqual(r.read_byte(), uint8(val));
            testCase.verifyEqual(r.read_bytes(length(unsigned)), unsigned);
            testCase.verifyEqual(unsignedSerializer.read(r), uint8(val));
            testCase.verifyEqual(unsignedSerializer.read_trivially(r, length(unsigned)), unsigned);
            testCase.verifyEqual(signedSerializer.read(r), int8(-val));
            testCase.verifyEqual(signedSerializer.read_trivially(r, length(signed)), signed);
            r.close();

            delete(filename);
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
            float32Serializer.write_trivially(w, fs);
            float64Serializer.write_trivially(w, ds);
            complexFloat32Serializer.write(w, cf);
            complexFloat64Serializer.write(w, df);
            complexFloat32Serializer.write_trivially(w, cfs);
            complexFloat64Serializer.write_trivially(w, dfs);
            w.close();

            r = yardl.binary.CodedInputStream(filename);
            testCase.verifyEqual(float32Serializer.read(r), f);
            testCase.verifyEqual(float64Serializer.read(r), d);
            testCase.verifyEqual(float32Serializer.read_trivially(r, size(fs)), fs);
            testCase.verifyEqual(float64Serializer.read_trivially(r, size(ds)), ds);
            testCase.verifyEqual(complexFloat32Serializer.read(r), cf);
            testCase.verifyEqual(complexFloat64Serializer.read(r), df);
            testCase.verifyEqual(complexFloat32Serializer.read_trivially(r, size(cfs)), cfs);
            testCase.verifyEqual(complexFloat64Serializer.read_trivially(r, size(dfs)), dfs);
            r.close();

            delete(filename);
        end

        function testBools(testCase)
            filename = tempname;

            bools = [true, false, true, false];
            boolSerializer = yardl.binary.BoolSerializer;

            w = yardl.binary.CodedOutputStream(filename);
            boolSerializer.write(w, true);
            boolSerializer.write_trivially(w, bools);
            w.close();

            r = yardl.binary.CodedInputStream(filename);
            testCase.verifyEqual(boolSerializer.read(r), true);
            testCase.verifyEqual(boolSerializer.read_trivially(r, size(bools)), bools);
            r.close();

            delete(filename);
        end

        function testNonTrivialErrors(testCase)
            filename = tempname;

            int16Serializer = yardl.binary.Int16Serializer;
            int32Serializer = yardl.binary.Int32Serializer;
            int64Serializer = yardl.binary.Int64Serializer;
            uint16Serializer = yardl.binary.Uint16Serializer;
            uint32Serializer = yardl.binary.Uint32Serializer;
            uint64Serializer = yardl.binary.Uint64Serializer;
            stringSerializer = yardl.binary.StringSerializer;

            w = yardl.binary.CodedOutputStream(filename);
            testCase.verifyError(@() int16Serializer.write_trivially(w, []), 'yardl:TypeError');
            testCase.verifyError(@() int32Serializer.write_trivially(w, []), 'yardl:TypeError');
            testCase.verifyError(@() int64Serializer.write_trivially(w, []), 'yardl:TypeError');
            testCase.verifyError(@() uint16Serializer.write_trivially(w, []), 'yardl:TypeError');
            testCase.verifyError(@() uint32Serializer.write_trivially(w, []), 'yardl:TypeError');
            testCase.verifyError(@() uint64Serializer.write_trivially(w, []), 'yardl:TypeError');
            testCase.verifyError(@() stringSerializer.write_trivially(w, []), 'yardl:TypeError');
            w.close();

            delete(filename);
        end
    end
end
