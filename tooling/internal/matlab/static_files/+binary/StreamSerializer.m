% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef StreamSerializer < yardl.binary.TypeSerializer
    properties
        item_serializer_;
        items_remaining_;
    end

    methods
        function self = StreamSerializer(item_serializer)
            self.item_serializer_ = item_serializer;
            self.items_remaining_ = 0;
        end

        function write(self, outstream, values)
            if isempty(values)
                return;
            end

            if iscolumn(values)
                values = transpose(values);
            end
            s = size(values);
            count = s(end);
            outstream.write_unsigned_varint(count);

            if iscell(values)
                if ~isvector(s)
                    throw(yardl.ValueError("cell array must be a vector"));
                end
                for i = 1:count
                    self.item_serializer_.write(outstream, values{i});
                end
            else
                if ndims(values) > 2
                    r = reshape(values, [], count);
                    inner_shape = s(1:end-1);
                    for i = 1:count
                        val = reshape(r(:, i), inner_shape);
                        self.item_serializer_.write(outstream, val);
                    end
                else
                    for i = 1:count
                        self.item_serializer_.write(outstream, transpose(values(:, i)));
                    end
                end
            end
        end

        function res = hasnext(self, instream)
            if self.items_remaining_ <= 0
                self.items_remaining_ = instream.read_unsigned_varint();
                if self.items_remaining_ <= 0
                    res = false;
                    return;
                end
            end
            res = true;
        end

        function res = read(self, instream)
            if self.items_remaining_ <= 0
                throw(yardl.RuntimeError("Stream has been exhausted"));
            end

            res = self.item_serializer_.read(instream);
            self.items_remaining_ = self.items_remaining_ - 1;
        end

        function c = get_class(self)
            c = self.item_serializer_.get_class();
        end

        function s = get_shape(~)
            s = [];
        end
    end
end
