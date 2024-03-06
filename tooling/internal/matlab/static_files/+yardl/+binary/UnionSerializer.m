classdef UnionSerializer < handle

    properties (Access=protected)
        union_class_
        case_serializers_
        case_factories_
        offset_
    end

    methods

        function obj = UnionSerializer(union_class, case_serializers, case_factories)
            obj.union_class_ = union_class;
            obj.case_serializers_ = case_serializers;
            obj.case_factories_ = case_factories;

            if isa(case_serializers{1}, 'yardl.None')
                obj.offset_ = 1;
            else
                obj.offset_ = 0;
            end
        end

        function write(obj, outstream, value)
            if isa(value, 'yardl.None')
                if isa(obj.case_serializers_{1}, 'yardl.None')
                    outstream.write_byte_no_check(0);
                    return;
                else
                    throw(yardl.TypeError("None is not valid for this union type"))
                end
            end

            if ~isa(value, obj.union_class_)
                throw(yardl.TypeError("Expected union value of type %s, got %s", obj.union_class_, class(value)))
            end

            tag_index = uint8(value.index + obj.offset_);
            outstream.write_byte_no_check(tag_index-1);

            serializer = obj.case_serializers_{tag_index};
            serializer.write(outstream, value.value);
        end

        function res = read(obj, instream)
            case_index = instream.read_byte() + 1;

            if case_index == 1 && obj.offset_ == 1
                res = yardl.None;
                return
            end

            serializer = obj.case_serializers_{case_index};
            value = serializer.read(instream);

            factory = obj.case_factories_{case_index};
            res = factory(value);
        end
    end
end
