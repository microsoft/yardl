% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef Uint32Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            arguments
                outstream (1,1) yardl.binary.CodedOutputStream
                value (1,1) {mustBeInRange(value, 0, 4294967295)}
            end
            outstream.write_unsigned_varint(value);
        end

        function res = read(instream)
            res = uint32(instream.read_unsigned_varint());
        end

        function c = get_class()
            c = "uint32";
        end
    end
end
