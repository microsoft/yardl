classdef FixedNDArraySerializer < yardl.binary.NDArraySerializerBase

    properties
        shape_
    end

    methods
        function self = FixedNDArraySerializer(item_serializer, shape)
            self@yardl.binary.NDArraySerializerBase(item_serializer);
            self.shape_ = shape;
        end

        function write(self, outstream, values)
            sz = size(values);

            if numel(values) == prod(self.shape_)
                % This is an NDArray of scalars
                self.write_data_(outstream, values(:));
                return;
            end

            if length(sz) < length(self.shape_)
                expected = sprintf("%d ", self.shape_);
                actual = sprintf("%d ", sz);
                throw(yardl.ValueError("Expected shape [%s], got [%s]", expected, actual));
            end

            fixedSize = sz(end-length(self.shape_)+1:end);
            if fixedSize ~= self.shape_
                expected = sprintf("%d ", self.shape_);
                actual = sprintf("%d ", fixedSize);
                throw(yardl.ValueError("Expected shape [%s], got [%s]", expected, actual));
            end

            inner_shape = sz(1:end-length(self.shape_));
            values = reshape(values, [inner_shape prod(self.shape_)]);

            self.write_data_(outstream, values);
        end

        function value = read(self, instream)
            % TODO: Better way of doing this? Possibly needed elsewhere to?
            %       self.shape_ should really be non-scalar...
            if isscalar(self.shape_)
                value = self.read_data_(instream, [1 self.shape_]);
            else
                value = self.read_data_(instream, self.shape_);
            end
        end

        function s = getShape(obj)
            item_shape = obj.item_serializer_.getShape();
            if isempty(item_shape)
                s = obj.shape_;
            elseif isscalar(item_shape)
                s = [item_shape obj.shape_ ];
            else
                s = [item_shape obj.shape_];
                s = s(s>1);
            end
        end
    end
end
