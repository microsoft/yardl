classdef Complexfloat32Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            % assert(class(value) == "single");
            assert(real(value) <= realmax('single'));
            assert(imag(value) <= realmax('single'));
            assert(real(value) >= -realmax('single'));
            assert(imag(value) >= -realmax('single'));

            real_bytes = typecast(single(real(value)), "uint8");
            imag_bytes = typecast(single(imag(value)), "uint8");
            outstream.write_bytes(real_bytes);
            outstream.write_bytes(imag_bytes);
        end

        function res = read(instream)
            real_bytes = instream.read(4);
            imag_bytes = instream.read(4);
            res = complex(typecast(real_bytes, "single"), typecast(imag_bytes, "single"));
        end
    end
end
