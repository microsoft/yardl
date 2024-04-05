% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef NDArraySerializer < yardl.binary.NDArraySerializerBase

    properties
        ndims_
    end

    methods
        function self = NDArraySerializer(item_serializer, ndims)
            self@yardl.binary.NDArraySerializerBase(item_serializer);
            self.ndims_ = ndims;
        end

        function write(self, outstream, values)
            if ndims(values) < self.ndims_
                throw(yardl.ValueError("Expected %d dimensions, got %d", self.ndims_, ndims(values)));
            end

            sz = size(values);

            flipped_shape = flip(sz);
            for dim = 1: self.ndims_
                len = flipped_shape(dim);
                outstream.write_unsigned_varint(len);
            end

            if ndims(values) == self.ndims_
                % This is an NDArray of scalars
                self.write_data_(outstream, values(:));
                return
            end

            % This is an NDArray of vectors/arrays
            inner_shape = sz(1:end-self.ndims_);
            outer_shape = sz(end-self.ndims_+1:end);
            values = reshape(values, [inner_shape prod(outer_shape)]);

            self.write_data_(outstream, values);
        end

        function value = read(self, instream)
            flipped_shape = zeros(1, self.ndims_);
            for dim = 1:self.ndims_
                flipped_shape(dim) = instream.read_unsigned_varint();
            end
            shape = flip(flipped_shape);

            value = self.read_data_(instream, shape);
        end
    end
end
