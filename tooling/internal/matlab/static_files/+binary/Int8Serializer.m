% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef Int8Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            assert(value <= intmax("int8"));
            assert(value >= intmin("int8"));
            bytes = typecast(int8(value), "uint8");
            outstream.write_bytes(bytes);
        end

        function res = read(instream)
            res = typecast(instream.read_byte(), "int8");
        end

        function c = getClass()
            c = 'int8';
        end

        function trivial = isTriviallySerializable()
            trivial = true;
        end
    end
end
