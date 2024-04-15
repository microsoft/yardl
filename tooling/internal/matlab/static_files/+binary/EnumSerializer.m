% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef EnumSerializer < yardl.binary.TypeSerializer
    properties
        classname_;
        constructor_;
        integer_serializer_;
    end

    methods
        function self = EnumSerializer(classname, classconstructor, integer_serializer)
            self.classname_ = classname;
            self.constructor_ = classconstructor;
            self.integer_serializer_ = integer_serializer;
        end

        function write(self, outstream, value)
            self.integer_serializer_.write(outstream, value);
        end

        function res = read(self, instream)
            int_value = self.integer_serializer_.read(instream);
            res = self.constructor_(int_value);
        end

        function c = get_class(self)
            c = self.classname_;
        end

        function trivial = is_trivially_serializable(self)
            trivial = self.integer_serializer_.is_trivially_serializable();
        end

        function write_trivially(self, outstream, values)
            self.integer_serializer_.write_trivially(outstream, values);
        end

        function res = read_trivially(self, instream, shape)
            res = self.integer_serializer_.read_trivially(instream, shape);
        end
    end
end
