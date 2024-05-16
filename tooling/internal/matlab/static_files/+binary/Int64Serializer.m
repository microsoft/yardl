% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef Int64Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            arguments
                outstream (1,1) yardl.binary.CodedOutputStream
                value (1,1) {mustBeInRange(value, -9223372036854775808, 9223372036854775807)}
            end
            outstream.write_signed_varint(value);
        end

        function res = read(instream)
            res = int64(instream.read_signed_varint());
        end

        function c = get_class()
            c = "int64";
        end
    end
end
