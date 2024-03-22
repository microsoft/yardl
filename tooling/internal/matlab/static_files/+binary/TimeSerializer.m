classdef TimeSerializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            if isa(value, 'datetime')
                value = yardl.Time.from_datetime(value).value;
            elseif isa(value, 'yardl.Time')
                value = value.value;
            else
                throw(yardl.TypeError("Expected datetime or yardl.Time, got %s", class(value)));
            end

            outstream.write_signed_varint(value);
        end

        function res = read(instream)
            value = instream.read_signed_varint();
            res = yardl.Time(value);
        end

        function c = getClass()
            c = 'yardl.Time';
        end
    end
end
