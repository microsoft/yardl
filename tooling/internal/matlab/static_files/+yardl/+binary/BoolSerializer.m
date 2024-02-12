classdef BoolSerializer < yardl.binary.TypeSerializer
    methods (Static)
        function write( outstream, value)
            assert(islogical(value));
            byte = cast(value, "uint8");
            outstream.write_bytes(byte);
        end

        function res = read(instream)
            byte = instream.read_byte();
            res = cast(byte, "logical");
        end
    end
end
