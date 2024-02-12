% A basic time of day with nanosecond precision. It is not timezone-aware and is
% meant to represent a wall clock time.
classdef Time < handle

    properties (Access=private)
        nanoseconds_since_midnight_
    end

    methods
        function obj = Time(nanoseconds_since_midnight)
            obj.nanoseconds_since_midnight_ = nanoseconds_since_midnight;
        end

        function value = value(obj)
            value = obj.nanoseconds_since_midnight_;
        end

        function dt = to_datetime(obj)
            dt = datetime(obj.nanoseconds_since_midnight_, 'ConvertFrom', 'epochtime', 'Epoch', datetime('today'), 'TicksPerSecond', 1e9);
        end

        function eq = eq(obj, other)
            if isa(other, 'datetime')
                other = yardl.Time.from_datetime(other);
            end

            if isa(other, 'yardl.Time')
                eq = all([obj.value] == [other.value]);
            else
                eq = false;
            end
        end
    end

    methods (Static)
        function t = from_datetime(value)
            nanoseconds_since_midnight = convertTo(value, 'epochtime', 'Epoch', datetime('today', 'TimeZone', value.TimeZone), 'TicksPerSecond', 1e9);
            t = yardl.Time(nanoseconds_since_midnight);
        end
    end

end
