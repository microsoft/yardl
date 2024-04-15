% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef DynamicNDArraySerializer < yardl.binary.NDArraySerializerBase

    methods
        function self = DynamicNDArraySerializer(item_serializer)
            self@yardl.binary.NDArraySerializerBase(item_serializer);
        end

        function write(self, outstream, values)
            item_shape = self.item_serializer_.get_shape();
            shape = size(values);
            if isempty(item_shape)
                % values is an array of variable-length vectors or arrays
                values = values(:);
            elseif isscalar(item_shape)
                % values is an array of scalars
                values = values(:);
            else
                % values is an array of fixed-length vectors or arrays
                item_shape = item_shape(item_shape > 1);
                outer_shape = shape(length(item_shape) + 1:end);
                values = reshape(values, [item_shape prod(outer_shape)]);
                shape = outer_shape;
            end

            outstream.write_unsigned_varint(length(shape));
            flipped_shape = flip(shape);
            for dim = 1:length(flipped_shape)
                outstream.write_unsigned_varint(flipped_shape(dim));
            end

            self.write_data_(outstream, values);
        end

        function value = read(self, instream)
            ndims = instream.read_unsigned_varint();
            flipped_shape = zeros(1, ndims);
            for dim = 1:ndims
                flipped_shape(dim) = instream.read_unsigned_varint();
            end
            shape = flip(flipped_shape);
            value = self.read_data_(instream, shape);
        end
    end
end
