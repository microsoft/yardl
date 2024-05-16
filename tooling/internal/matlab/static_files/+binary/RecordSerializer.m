% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef RecordSerializer < yardl.binary.TypeSerializer

    properties
        field_serializers
        classname
    end

    methods
        function self = RecordSerializer(classname, field_serializers)
            self.classname = classname;
            self.field_serializers = field_serializers;
        end

        function c = get_class(self)
            c = self.classname;
        end
    end

    methods (Access=protected)
        function write_(self, outstream, varargin)
            for i = 1:nargin-2
                fs = self.field_serializers{i};
                field_value = varargin{i};
                fs.write(outstream, field_value);
            end
        end

        function res = read_(self, instream)
            res = cell(size(self.field_serializers));
            for i = 1:length(self.field_serializers)
                fs = self.field_serializers{i};
                res{i} = fs.read(instream);
            end
        end
    end
end
