classdef EnumSerializer < yardl.binary.TypeSerializer
    properties
        classname_;
        constructor_;
        integer_serializer_;
    end

    methods
        function obj = EnumSerializer(classname, classconstructor, integer_serializer)
            obj.classname_ = classname;
            obj.constructor_ = classconstructor;
            obj.integer_serializer_ = integer_serializer;
        end

        function write(obj, outstream, value)
            obj.integer_serializer_.write(outstream, value);
        end

        function res = read(obj, instream)
            int_value = obj.integer_serializer_.read(instream);
            res = obj.constructor_(int_value);
        end

        function c = getClass(obj)
            c = obj.classname_;
        end
    end
end
