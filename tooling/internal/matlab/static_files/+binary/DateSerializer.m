classdef DateSerializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            if isa(value, 'datetime')
                value = yardl.Date.from_datetime(value).value;
            elseif isa(value, 'yardl.Date')
                value = value.value;
            else
                throw(yardl.TypeError("Expected datetime or yardl.Date, got %s", class(value)));
            end
            outstream.write_signed_varint(value);
        end

        function res = read(instream)
            value = instream.read_signed_varint();
            res = yardl.Date(value);
        end

        function c = getClass()
            c = 'yardl.Date';
        end
    end
end
