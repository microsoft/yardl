classdef DatetimeSerializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            if isa(value, 'datetime')
                value = yardl.DateTime.from_datetime(value).value;
            elseif isa(value, 'yardl.DateTime')
                value = value.value;
            else
                throw(yardl.TypeError("Expected datetime or yardl.DateTime, got %s", class(value)));
            end
            outstream.write_signed_varint(value);
        end

        function res = read(instream)
            value = instream.read_signed_varint();
            res = yardl.DateTime(value);
        end

        function c = getClass()
            c = 'yardl.DateTime';
        end
    end
end
