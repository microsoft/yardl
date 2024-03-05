classdef NDArraySerializer < yardl.binary.NDArraySerializerBase

    properties
        ndims_
    end


    methods
        function self = NDArraySerializer(element_serializer, ndims)
            self@yardl.binary.NDArraySerializerBase(element_serializer);
            self.ndims_ = ndims;
        end

        function write(self, outstream, value)
            if ndims(value) ~= self.ndims_
                throw(yardl.ValueError("Expected %s dimensions, got %s", self.ndims_, ndims(value)));
            end

            for dim = self.ndims_: -1: 1
                len = size(value, dim);
                outstream.write_unsigned_varint(len);
            end

            self.write_data_(outstream, value);
        end

        function value = read(self, instream)
            shape = zeros(1, self.ndims_);
            for dim = 1:self.ndims_
                shape(dim) = instream.read_unsigned_varint();
            end
            value = self.read_data_(instream, flip(shape));
        end
    end
end
