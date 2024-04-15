% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef Union < handle
    properties (Access=protected)
        index_
        value_
    end

    methods
        function self = Union(index, value)
            self.index_ = index;
            self.value_ = value;
        end

        function i = index(self)
            i = self.index_;
        end

        function v = value(self)
            v = self.value_;
        end
    end
end
