classdef RecordSerializer < yardl.binary.TypeSerializer

    properties
        field_serializers
        classname
    end

    methods
        function obj = RecordSerializer(classname, field_serializers)
            obj.classname = classname;
            obj.field_serializers = field_serializers;
        end

        function c = getClass(obj)
            c = obj.classname;
        end
    end

    methods (Access=protected)
        function write_(obj, outstream, varargin)
            for i = 1:nargin-2
                fs = obj.field_serializers{i};
                field_value = varargin{i};
                fs.write(outstream, field_value);
            end
        end

        function res = read_(obj, instream)
            res = cell(size(obj.field_serializers));
            for i = 1:length(obj.field_serializers)
                fs = obj.field_serializers{i};
                res{i} = fs.read(instream);
            end
        end
    end
end
