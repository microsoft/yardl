classdef Union < handle
    properties (Access=protected)
        index_
        value_
    end

    methods
        function obj = Union(index, value)
            obj.index_ = index;
            obj.value_ = value;
        end

        function i = index(obj)
            i = obj.index_;
        end

        function v = value(obj)
            v = obj.value_;
        end
    end
end
