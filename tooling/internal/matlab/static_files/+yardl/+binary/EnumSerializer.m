classdef EnumSerializer < yardl.binary.TypeSerializer
    properties
        integer_serializer_;
        enum_class_;
    end

    methods
        function obj = EnumSerializer(integer_serializer, enum_class)
            obj.integer_serializer_ = integer_serializer;
            obj.enum_class_ = enum_class;
        end

        function write(obj, outstream, value)
            obj.integer_serializer_.write(outstream, value);
        end

        function res = read(obj, instream)
            int_value = obj.integer_serializer_.read(instream);
            res = obj.enum_class_(int_value);
        end
    end
end
