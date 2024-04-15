% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef Uint64Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            arguments
                outstream (1,1) yardl.binary.CodedOutputStream
                value (1,1) {mustBeInRange(value, 0, 18446744073709551615)}
            end
            outstream.write_unsigned_varint(value);
        end

        function res = read(instream)
            res = uint64(instream.read_unsigned_varint());
        end

        function c = get_class()
            c = "uint64";
        end
    end
end
