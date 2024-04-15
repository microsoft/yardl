% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef Float64Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            arguments
                outstream (1,1) yardl.binary.CodedOutputStream
                value (1,1) double
            end
            bytes = typecast(double(value), "uint8");
            outstream.write_bytes(bytes);
        end

        function res = read(instream)
            bytes = instream.read_bytes(8);
            res = typecast(bytes, "double");
        end

        function c = get_class()
            c = "double";
        end

        function trivial = is_trivially_serializable()
            trivial = true;
        end
    end
end
