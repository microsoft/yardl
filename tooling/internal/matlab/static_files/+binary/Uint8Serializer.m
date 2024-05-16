% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef Uint8Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            arguments
                outstream (1,1) yardl.binary.CodedOutputStream
                value (1,1) {mustBeInRange(value, 0, 255)}
            end
            outstream.write_byte(uint8(value));
        end

        function res = read(instream)
            res = instream.read_byte();
        end

        function c = get_class()
            c = "uint8";
        end

        function trivial = is_trivially_serializable()
            trivial = true;
        end
    end
end
