classdef NDArraySerializerBase < handle

    properties
        element_serializer_
    end

    methods (Abstract)
        write(self, outstream, value)
        read(self, instream)
    end

    methods (Access=protected)
        function self = NDArraySerializerBase(element_serializer)
            self.element_serializer_ = element_serializer;
        end

        function write_data_(self, outstream, value)
            flat_value = value(:);
            for i = 1:length(flat_value)
                self.element_serializer_.write(outstream, flat_value(i));
            end
        end

        function value = read_data_(self, instream, shape)
            flat_length = prod(shape);
            % if self.element_serializer_.is_trivially_serializable()...

            for i = 1:flat_length
                value(i) = self.element_serializer_.read(instream);
            end

            value = reshape(value, shape);
        end
    end
end
