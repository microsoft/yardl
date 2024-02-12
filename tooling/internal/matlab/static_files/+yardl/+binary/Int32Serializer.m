classdef Int32Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            assert(value <= intmax("int32"));
            assert(value >= intmin("int32"));
            value = int32(value);
            outstream.write_signed_varint(value);
        end

        function res = read(instream)
            res = int32(instream.read_signed_varint());
        end
    end
end
