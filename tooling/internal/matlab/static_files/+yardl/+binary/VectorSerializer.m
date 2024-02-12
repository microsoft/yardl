classdef VectorSerializer < yardl.binary.TypeSerializer
    properties
        element_serializer_;
    end

    methods
        function obj = VectorSerializer(element_serializer)
            obj.element_serializer_ = element_serializer;
        end

        function write(obj, outstream, value)
            outstream.write_unsigned_varint(length(value));
            for i = 1:length(value)
                obj.element_serializer_.write(outstream, value(i));
            end
        end

        function res = read(obj, instream)
            count = instream.read_unsigned_varint();
            for i = 1:count
                res(i) = obj.element_serializer_.read(instream);
            end
        end
    end
end
