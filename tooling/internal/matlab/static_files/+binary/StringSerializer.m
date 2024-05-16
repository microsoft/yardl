% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef StringSerializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            arguments
                outstream (1,1) yardl.binary.CodedOutputStream
                value (1,1) {mustBeTextScalar}
            end
            bytes = unicode2native(value, "utf-8");
            outstream.write_unsigned_varint(length(bytes));
            outstream.write_bytes(bytes);
        end

        function res = read(instream)
            len = instream.read_unsigned_varint();
            bytes = instream.read_bytes(len);
            res = convertCharsToStrings(native2unicode(bytes, "utf-8"));
        end

        function c = get_class()
            c = "string";
        end
    end
end
