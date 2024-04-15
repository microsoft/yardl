% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef Complexfloat64Serializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            arguments
                outstream (1,1) yardl.binary.CodedOutputStream
                value (1,1) double
            end
            real_bytes = typecast(double(real(value)), "uint8");
            imag_bytes = typecast(double(imag(value)), "uint8");
            outstream.write_bytes(real_bytes);
            outstream.write_bytes(imag_bytes);
        end

        function res = read(instream)
            real_bytes = instream.read_bytes(8);
            imag_bytes = instream.read_bytes(8);
            res = complex(typecast(real_bytes, "double"), typecast(imag_bytes, "double"));
        end

        function c = get_class()
            c = "double";
        end

        function trivial = is_trivially_serializable()
            trivial = true;
        end
    end

    methods
        function write_trivially(self, stream, values)
            rp = real(values);
            ip = imag(values);
            parts(1, :) = rp(:);
            parts(2, :) = ip(:);
            stream.write_values_directly(parts, self.get_class());
        end

        function res = read_trivially(self, stream, shape)
            assert(ndims(shape) == 2);
            partshape = [2*shape(1) shape(2)];
            res = stream.read_values_directly(partshape, self.get_class());
            res = reshape(res, [2 shape]);
            res = squeeze(complex(res(1, :, :), res(2, :, :)));
        end
    end
end
