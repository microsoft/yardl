classdef Uint16Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            assert(value <= intmax("uint16"));
            assert(value >= intmin("uint16"));
            value = uint16(value);
            outstream.write_unsigned_varint(value);
        end

        function res = read(instream)
            res = instream.read_unsigned_varint();
        end
    end
end
