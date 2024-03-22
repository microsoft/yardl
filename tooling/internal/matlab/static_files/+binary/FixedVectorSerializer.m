classdef FixedVectorSerializer < yardl.binary.TypeSerializer
    properties
        element_serializer_;
        length_;
    end

    methods
        function obj = FixedVectorSerializer(element_serializer, length)
            obj.element_serializer_ = element_serializer;
            obj.length_ = length;
        end

        function write(obj, outstream, value)
            if length(value) ~= obj.length_
                throw(yardl.ValueError("Expected an array of length %d, got %d", obj.length_, length(value)));
            end

            for i = 1:obj.length_
                obj.element_serializer_.write(outstream, value(i));
            end
        end

        function res = read(obj, instream)
            for i = 1:obj.length_
                res(i) = obj.element_serializer_.read(instream);
            end
        end

        function c = getClass(obj)
            c = obj.element_serializer_.getClass();
        end
    end
end
