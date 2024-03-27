classdef Float32Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            % assert(class(value) == "single");
            assert(value <= realmax('single'));
            assert(value >= -realmax('single'));

            bytes = typecast(single(value), "uint8");
            outstream.write_bytes(bytes);
        end

        function res = read(instream)
            bytes = instream.read(4);
            res = typecast(bytes, "single");
        end

        function c = getClass()
            c = 'single';
        end
    end
end
