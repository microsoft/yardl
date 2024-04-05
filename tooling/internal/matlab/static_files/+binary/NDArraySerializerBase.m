% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef NDArraySerializerBase < yardl.binary.TypeSerializer

    properties
        item_serializer_
    end

    methods (Abstract)
        write(self, outstream, values)
        read(self, instream)
    end

    methods (Access=protected)
        function self = NDArraySerializerBase(item_serializer)
            self.item_serializer_ = item_serializer;
        end

        function write_data_(self, outstream, values)
            % values is an array of shape [A, B, ..., N], where
            % N is the "flattened" dimension of the NDArray, and
            % A, B, ... are the dimensions of the inner items.

            sz = size(values);

            if ndims(values) > 2
                count = sz(end);
                inner_shape = sz(1:end-1);
                r = reshape(values, [], count);
                for i = 1:count
                    val = reshape(r(:, i), inner_shape);
                    self.item_serializer_.write(outstream, val);
                end
            elseif isrow(values) || iscolumn(values)
                count = prod(sz);
                if iscell(values)
                    for i = 1:count
                        self.item_serializer_.write(outstream, values{i});
                    end
                else
                    for i = 1:count
                        self.item_serializer_.write(outstream, values(i));
                    end
                end
            else
                assert(ismatrix(values));
                count = sz(end);
                for i = 1:count
                    self.item_serializer_.write(outstream, values(:, i));
                end
            end
        end

        function res = read_data_(self, instream, shape)
            flat_length = prod(shape);

            item_shape = self.item_serializer_.getShape();

            if isempty(item_shape)
                res = cell(shape);
                for i = 1:flat_length
                    res{i} = self.item_serializer_.read(instream);
                end
            elseif isscalar(item_shape)
                res = yardl.allocate(self.getClass(), shape);
                for i = 1:flat_length
                    res(i) = self.item_serializer_.read(instream);
                end
                res = squeeze(res);
            else
                res = yardl.allocate(self.getClass(), [prod(item_shape), flat_length]);
                for i = 1:flat_length
                    item = self.item_serializer_.read(instream);
                    res(:, i) = item(:);
                end
                res = squeeze(reshape(res, [item_shape shape]));
            end
        end
    end

    methods
        function c = getClass(self)
            c = self.item_serializer_.getClass();
        end

        function s = getShape(~)
            s = [];
        end
    end
end
