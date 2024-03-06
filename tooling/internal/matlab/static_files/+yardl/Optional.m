classdef Optional < handle
    properties (Access=protected)
        has_value_
        value_
    end

    methods
        function obj = Optional(varargin)
            if nargin > 0
                obj.value_ = varargin{1};
                obj.has_value_ = true;
            else
                obj.has_value_ = false;
            end
        end

        function has_value = has_value(obj)
            has_value = obj.has_value_;
        end

        function value = value(obj)
            if ~obj.has_value()
                throw(yardl.TypeError("Optional type does not have a value"));
            end
            value = obj.value_;
        end
    end
end
