classdef Date < handle

    properties (Access=private)
        days_since_epoch_
    end

    methods
        function obj = Date(days_since_epoch)
            if nargin > 0
                obj.days_since_epoch_ = days_since_epoch;
            else
                obj.days_since_epoch_ = 0;
            end
        end

        function value = value(obj)
            value = obj.days_since_epoch_;
        end

        function dt = to_datetime(obj)
            dt = datetime(0, 'ConvertFrom', 'epochtime') + days(obj.days_since_epoch_);
        end

        function eq = eq(obj, other)
            if isa(other, 'datetime')
                other = yardl.Date.from_datetime(other);
            end

            eq = isa(other, 'yardl.Date') && ...
                all([obj.value] == [other.value]);
        end
    end

    methods (Static)
        function d = from_datetime(value)
            dur = value - datetime(0, 'ConvertFrom', 'epochtime');
            d = yardl.Date(days(dur));
        end

        function d = from_components(y, m, d)
            d = yardl.Date.from_datetime(datetime(y, m, d));
        end
    end

end
