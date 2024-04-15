% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef UnionSerializer < handle

    properties (Access=protected)
        classname_
        case_serializers_
        case_factories_
        offset_
    end

    methods
        function self = UnionSerializer(union_class, case_serializers, case_factories)
            self.classname_ = union_class;
            self.case_serializers_ = case_serializers;
            self.case_factories_ = case_factories;

            if isa(case_serializers{1}, "yardl.binary.NoneSerializer")
                self.offset_ = 1;
            else
                self.offset_ = 0;
            end
        end

        function write(self, outstream, value)
            if isa(value, "yardl.Optional")
                if ~isa(self.case_serializers_{1}, "yardl.binary.NoneSerializer")
                    throw(yardl.TypeError("Optional is not valid for this union type"))
                end

                if value.has_value()
                    value = value.value;
                else
                    outstream.write_byte(uint8(0));
                    return;
                end
            end

            if ~isa(value, self.classname_)
                throw(yardl.TypeError("Expected union value of type %s, got %s", self.classname_, class(value)))
            end

            tag_index = uint8(value.index + self.offset_);
            outstream.write_byte(tag_index-1);

            serializer = self.case_serializers_{tag_index};
            serializer.write(outstream, value.value);
        end

        function res = read(self, instream)
            case_index = instream.read_byte() + 1;

            if case_index == 1 && self.offset_ == 1
                res = yardl.None;
                return
            end

            serializer = self.case_serializers_{case_index};
            value = serializer.read(instream);

            factory = self.case_factories_{case_index};
            res = factory(value);
        end

        function c = get_class(self)
            if isa(self.case_serializers_{1}, "yardl.binary.NoneSerializer")
                c = "yardl.Optional";
            else
                c = self.classname_;
            end
        end

    end

    methods (Static)
        function s = get_shape()
            s = 1;
        end
    end
end
