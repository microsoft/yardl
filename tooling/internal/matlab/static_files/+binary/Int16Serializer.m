% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef Int16Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            assert(value <= intmax("int16"));
            assert(value >= intmin("int16"));
            value = int16(value);
            outstream.write_signed_varint(value);
        end

        function res = read(instream)
            res = int16(instream.read_signed_varint());
        end

        function c = getClass()
            c = 'int16';
        end
    end
end
