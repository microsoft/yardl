classdef Int64Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            assert(value <= intmax("int64"));
            assert(value >= intmin("int64"));
            value = int64(value);
            outstream.write_signed_varint(value);
        end

        function res = read(instream)
            res = int64(instream.read_signed_varint());
        end
    end
end
