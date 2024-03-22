classdef Uint64Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            assert(value <= intmax("uint64"));
            assert(value >= intmin("uint64"));
            value = uint64(value);
            outstream.write_unsigned_varint(value);
        end

        function res = read(instream)
            res = uint64(instream.read_unsigned_varint());
        end

        function c = getClass()
            c = 'uint64';
        end
    end
end
