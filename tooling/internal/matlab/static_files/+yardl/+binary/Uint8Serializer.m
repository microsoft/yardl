classdef Uint8Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            assert(value <= intmax("uint8"));
            assert(value >= intmin("uint8"));
            % bytes = typecast(uint8(value), "uint8");
            outstream.write_bytes(uint8(value));
        end

        function res = read(instream)
            res = instream.read_byte();
        end
    end
end
