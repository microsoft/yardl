% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef FixedVectorSerializer < yardl.binary.VectorSerializerBase
    properties
        length_;
    end

    methods (Access=protected)
        function handle_write_count_(self, ~, count)
            if count ~= self.length_
                throw(yardl.ValueError("Expected an array of length %d, got %d", self.length_, count));
            end
        end

        function count = get_read_count_(self, ~)
            count = self.length_;
        end
    end

    methods
        function self = FixedVectorSerializer(item_serializer, length)
            self@yardl.binary.VectorSerializerBase(item_serializer);
            self.length_ = length;
        end

        function s = get_shape(self)
            item_shape = self.item_serializer_.get_shape();
            if isempty(item_shape)
                s = [1, self.length_];
            elseif isscalar(item_shape)
                s = [item_shape self.length_];
            else
                s = [item_shape self.length_];
                s = s(s>1);
            end
        end

        function trivial = is_trivially_serializable(self)
            trivial = self.item_serializer_.is_trivially_serializable();
        end

        function write_trivially(self, outstream, values)
            self.item_serializer_.write_trivially(outstream, values);
        end

        function res = read_trivially(self, instream, shape)
            res = self.item_serializer_.read_trivially(instream, shape);
        end
    end
end
