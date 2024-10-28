% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef Int8Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            arguments
                outstream (1,1) yardl.binary.CodedOutputStream
                value (1,1) {mustBeInRange(value, -128, 127)}
            end
            byte = typecast(int8(value), "uint8");
            outstream.write_byte(byte);
        end

        function res = read(instream)
            res = typecast(instream.read_byte(), "int8");
        end

        function c = get_class()
            c = "int8";
        end

        function trivial = is_trivially_serializable()
            trivial = true;
        end
    end
end
