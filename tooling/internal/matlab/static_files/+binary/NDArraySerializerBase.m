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

            if ~iscell(values) && self.item_serializer_.is_trivially_serializable()
                self.item_serializer_.write_trivially(outstream, values);
                return;
            end

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
                count = sz(end);
                for i = 1:count
                    self.item_serializer_.write(outstream, values(:, i));
                end
            end
        end

        function res = read_data_(self, instream, shape)
            flat_length = prod(shape);

            item_shape = self.item_serializer_.get_shape();

            if isempty(item_shape)
                res = cell(shape);
                for i = 1:flat_length
                    res{i} = self.item_serializer_.read(instream);
                end
                return
            end

            if self.item_serializer_.is_trivially_serializable()
                res = self.item_serializer_.read_trivially(instream, [prod(item_shape), flat_length]);
            else
                res = yardl.allocate(self.get_class(), [prod(item_shape), flat_length]);
                for i = 1:flat_length
                    item = self.item_serializer_.read(instream);
                    res(:, i) = item(:);
                end
            end

            % Tricky reshaping to remove unnecessary singleton dimensions in
            %   subarrays, arrays of fixed vectors, etc.
            item_shape = item_shape(item_shape > 1);
            if isempty(item_shape) && isscalar(shape)
                item_shape = 1;
            end
            newshape = [item_shape shape];
            res = reshape(res, newshape);

            if iscolumn(res)
                res = transpose(res);
            end
        end
    end

    methods
        function c = get_class(self)
            c = self.item_serializer_.get_class();
        end

        function s = get_shape(~)
            s = [];
        end
    end
end
