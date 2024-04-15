% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef Int32Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            arguments
                outstream (1,1) yardl.binary.CodedOutputStream
                value (1,1) {mustBeInRange(value, -2147483648, 2147483647)}
            end
            outstream.write_signed_varint(value);
        end

        function res = read(instream)
            res = int32(instream.read_signed_varint());
        end

        function c = get_class()
            c = "int32";
        end
    end
end
