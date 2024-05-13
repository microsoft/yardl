% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef Union < handle & matlab.mixin.CustomDisplay
    properties (Access=protected)
        index_
    end

    properties (SetAccess=protected)
        value
    end

    methods
        function self = Union(index, value)
            self.index_ = index;
            self.value = value;
        end

        function i = index(self)
            i = self.index_;
        end
    end

    methods (Abstract)
        t = tag(self)
    end

    methods (Access=protected)
        function displayScalarObject(obj)
            className = matlab.mixin.CustomDisplay.getClassNameForHeader(obj);
            header = sprintf('%s.%s\n', className, obj.tag());
            disp(header)
            propgroup = getPropertyGroups(obj);
            matlab.mixin.CustomDisplay.displayPropertyGroups(obj,propgroup)
        end
    end
end
