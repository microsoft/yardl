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
                throw(yardl.ValueError("Expected shape %s, got %s", self.shape_, size(value)));
            end

            self.write_data_(outstream, value);
        end

        function value = read(self, instream)
            value = self.read_data_(instream, self.shape_);
        end
    end
end
