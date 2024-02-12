classdef StringSerializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            % if ischar(value)
            %     value = convertCharsToStrings(value);
            % end
            bytes = unicode2native(value, "utf-8");
            outstream.write_unsigned_varint(length(bytes));
            outstream.write_bytes(bytes);
        end

        function res = read(instream)
            len = instream.read_unsigned_varint();
            bytes = instream.read(len);
            res = convertCharsToStrings(native2unicode(bytes, "utf-8"));
        end
    end
end
