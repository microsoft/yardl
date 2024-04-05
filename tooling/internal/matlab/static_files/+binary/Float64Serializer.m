% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef Float64Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            assert(class(value) == "double");
            bytes = typecast(double(value), "uint8");
            outstream.write_bytes(bytes);
        end

        function res = read(instream)
            bytes = instream.read(8);
            res = typecast(bytes, "double");
        end

        function c = getClass()
            c = 'double';
        end
    end
end
