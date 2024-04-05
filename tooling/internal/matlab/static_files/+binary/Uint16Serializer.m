% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef Uint16Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            assert(value <= intmax("uint16"));
            assert(value >= intmin("uint16"));
            value = uint16(value);
            outstream.write_unsigned_varint(value);
        end

        function res = read(instream)
            res = uint16(instream.read_unsigned_varint());
        end

        function c = getClass()
            c = 'uint16';
        end
    end
end
