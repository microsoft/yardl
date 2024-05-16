% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef BoolSerializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            arguments
                outstream (1,1) yardl.binary.CodedOutputStream
                value (1,1) {mustBeInRange(value, 0, 1)}
            end
            byte = cast(value, "uint8");
            outstream.write_bytes(byte);
        end

        function res = read(instream)
            byte = instream.read_byte();
            res = cast(byte, "logical");
        end

        function c = get_class()
            c = "logical";
        end

        function trivial = is_trivially_serializable()
            trivial = true;
        end
    end
end
