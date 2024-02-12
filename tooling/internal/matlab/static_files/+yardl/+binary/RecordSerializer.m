classdef RecordSerializer < handle

    properties
        field_serializers
    end

    methods

        function obj = RecordSerializer(field_serializers)
            obj.field_serializers = field_serializers;
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
