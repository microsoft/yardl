% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef Uint8Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            assert(value <= intmax("uint8"));
            assert(value >= intmin("uint8"));
            outstream.write_bytes(uint8(value));
        end

        function res = read(instream)
            res = instream.read_byte();
        end

        function c = getClass()
            c = 'uint8';
        end

        function trivial = isTriviallySerializable()
            trivial = true;
        end
    end
end
