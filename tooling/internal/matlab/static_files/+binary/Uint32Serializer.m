classdef Uint32Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            assert(value <= intmax("uint32"));
            assert(value >= intmin("uint32"));
            value = uint32(value);
            outstream.write_unsigned_varint(value);
        end

        function res = read(instream)
            res = uint32(instream.read_unsigned_varint());
        end

        function c = getClass()
            c = 'uint32';
        end
    end
end
