% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef Int16Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            arguments
                outstream (1,1) yardl.binary.CodedOutputStream
                value (1,1) {mustBeInRange(value, -32768, 32767)}
            end
            outstream.write_signed_varint(value);
        end

        function res = read(instream)
            res = int16(instream.read_signed_varint());
        end

        function c = get_class()
            c = "int16";
        end
    end
end
