% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

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

        function c = getClass()
            c = 'single';
        end

        function trivial = isTriviallySerializable()
            trivial = true;
        end
    end

    methods
        function writeTrivially(self, stream, values)
            rp = real(values);
            ip = imag(values);
            parts(1, :) = rp(:);
            parts(2, :) = ip(:);
            stream.write_values_directly(parts, self.getClass());
        end

        function res = readTrivially(obj, stream, shape)
            assert(ndims(shape) == 2);
            partshape = [2*shape(1) shape(2)];
            res = stream.read_values_directly(partshape, obj.getClass());
            res = reshape(res, [2 shape]);
            res = squeeze(complex(res(1, :, :), res(2, :, :)));
        end
    end
end
