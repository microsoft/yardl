% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef Float32Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            arguments
                outstream (1,1) yardl.binary.CodedOutputStream
                value (1,1) single
            end

            bytes = typecast(single(value), "uint8");
            outstream.write_bytes(bytes);
        end

        function res = read(instream)
            bytes = instream.read_bytes(4);
            res = typecast(bytes, "single");
        end

        function c = get_class()
            c = "single";
        end

        function trivial = is_trivially_serializable()
            trivial = true;
        end
    end
end
