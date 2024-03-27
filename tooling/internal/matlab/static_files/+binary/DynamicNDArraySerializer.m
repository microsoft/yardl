classdef DynamicNDArraySerializer < yardl.binary.NDArraySerializerBase

    methods
        function self = DynamicNDArraySerializer(item_serializer)
            self@yardl.binary.NDArraySerializerBase(item_serializer);
        end

        function write(self, outstream, value)
            outstream.write_unsigned_varint(ndims(value));
            for dim = ndims(value): -1: 1
                len = size(value, dim);
                outstream.write_unsigned_varint(len);
            end

            self.write_data_(outstream, value);
        end

        function value = read(self, instream)
            ndims = instream.read_unsigned_varint();
            shape = zeros(1, ndims);
            for dim = 1:ndims
                shape(dim) = instream.read_unsigned_varint();
            end
            value = self.read_data_(instream, flip(shape));
        end
    end
end
