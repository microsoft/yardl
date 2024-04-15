% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef Uint16Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            arguments
                outstream (1,1) yardl.binary.CodedOutputStream
                value (1,1) {mustBeInRange(value, 0, 65535)}
            end
            outstream.write_unsigned_varint(value);
        end

        function res = read(instream)
            res = uint16(instream.read_unsigned_varint());
        end

        function c = get_class()
            c = "uint16";
        end
    end
end
