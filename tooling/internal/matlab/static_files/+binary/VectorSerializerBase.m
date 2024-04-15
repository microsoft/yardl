% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef VectorSerializerBase < yardl.binary.TypeSerializer
    properties
        item_serializer_;
    end

    methods (Abstract, Access=protected)
        handle_write_count_(self, outstream, values)
        get_read_count_(self, instream)
    end

    methods (Access=protected)
        function self = VectorSerializerBase(item_serializer)
            self.item_serializer_ = item_serializer;
        end
    end

    methods
        function write(self, outstream, values)
            if iscolumn(values)
                values = transpose(values);
            end
            s = size(values);
            count = s(end);

            self.handle_write_count_(outstream, count);

            if iscell(values)
                % values is a cell array, so must be a vector of shape [1, COUNT]
                if ~isvector(s)
                    throw(yardl.ValueError("cell array must be a vector"));
                end
                for i = 1:count
                    self.item_serializer_.write(outstream, values{i});
                end
            else
                % values is an array, so must have shape [A, B, ..., COUNT]
                if self.item_serializer_.is_trivially_serializable()
                    self.item_serializer_.write_trivially(outstream, values);
                    return
                end

                if ndims(values) > 2
                    r = reshape(values, [], count);
                    for i = 1:count
                        val = reshape(r(:, i), s(1:end-1));
                        self.item_serializer_.write(outstream, val);
                    end
                else
                    for i = 1:count
                        self.item_serializer_.write(outstream, transpose(values(:, i)));
                    end
                end
            end
        end

        function res = read(self, instream)
            count = self.get_read_count_(instream);
            if count == 0
                res = yardl.allocate(self.get_class(), 0);
                return;
            end

            item_shape = self.item_serializer_.get_shape();
            if isempty(item_shape)
                res = cell(1, count);
                for i = 1:count
                    res{i} = self.item_serializer_.read(instream);
                end
                return
            end

            if self.item_serializer_.is_trivially_serializable()
                res = self.item_serializer_.read_trivially(instream, [prod(item_shape), count]);
            else
                res = yardl.allocate(self.get_class(), [prod(item_shape), count]);
                for i = 1:count
                    item = self.item_serializer_.read(instream);
                    res(:, i) = item(:);
                end
            end

            res = squeeze(reshape(res, [item_shape, count]));
            if iscolumn(res)
                res = transpose(res);
            end
        end

        function c = get_class(self)
            c = self.item_serializer_.get_class();
        end
    end
end
