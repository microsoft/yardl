classdef Complexfloat64Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            assert(class(value) == "double");
            real_bytes = typecast(double(real(value)), "uint8");
            imag_bytes = typecast(double(imag(value)), "uint8");
            outstream.write_bytes(real_bytes);
            outstream.write_bytes(imag_bytes);
        end

        function res = read(instream)
            real_bytes = instream.read(8);
            imag_bytes = instream.read(8);
            res = complex(typecast(real_bytes, "double"), typecast(imag_bytes, "double"));
        end

        function c = getClass()
            c = 'double';
        end
    end
end
