classdef NDArraySerializerBase < handle

    properties
        element_serializer_
    end

    methods (Abstract)
        write(self, outstream, value)
        read(self, instream)
    end

    methods
        function c = getClass(obj)
            c = obj.element_serializer_.getClass();
        end
    end

    methods (Access=protected)
        function self = NDArraySerializerBase(element_serializer)
            if isa(element_serializer, 'yardl.binary.FixedNDArraySerializer') || ...
                    isa(element_serializer, 'yardl.binary.FixedVectorSerializer')
                self.element_serializer_ = element_serializer.element_serializer_;
            else
                self.element_serializer_ = element_serializer;
            end
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

            if length(shape) > 1
                value = reshape(value, shape);
            end
        end
    end
end
