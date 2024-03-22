classdef FixedNDArraySerializer < yardl.binary.NDArraySerializerBase

    properties
        shape_
    end


    methods
        function self = FixedNDArraySerializer(element_serializer, shape)
            self@yardl.binary.NDArraySerializerBase(element_serializer);
            self.shape_ = shape;
        end

        function write(self, outstream, value)
            if size(value) ~= self.shape_
                expected = sprintf("%d ", self.shape_);
                actual = sprintf("%d ", size(value));
                throw(yardl.ValueError("Expected shape [%s], got [%s]", expected, actual));
            end

            self.write_data_(outstream, value);
        end

        function value = read(self, instream)
            value = self.read_data_(instream, self.shape_);
        end
    end
end
