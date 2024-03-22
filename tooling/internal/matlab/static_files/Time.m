% A basic time of day with nanosecond precision. It is not timezone-aware and is
% meant to represent a wall clock time.
classdef Time < handle

    properties (Access=private)
        nanoseconds_since_midnight_
    end

    methods
        function obj = Time(nanoseconds_since_midnight)
            if nargin > 0
                if nanoseconds_since_midnight < 0 || nanoseconds_since_midnight >= 24*60*60*1e9
                    throw(yardl.ValuError("Time must be between 00:00:00 and 23:59:59.999999999"));
                end
                obj.nanoseconds_since_midnight_ = nanoseconds_since_midnight;
            else
                obj.nanoseconds_since_midnight_ = 0;
            end
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

            eq = isa(other, 'yardl.Time') && ...
                all([obj.value] == [other.value]);
        end

        function ne = new(obj, other)
            ne = ~obj.eq(other);
        end
    end

    methods (Static)
        function t = from_datetime(value)
            nanoseconds_since_midnight = convertTo(value, 'epochtime', ...
                    'Epoch', datetime(value.Year, value.Month, value.Day, 'TimeZone', value.TimeZone), ...
                    'TicksPerSecond', 1e9);
            t = yardl.Time(nanoseconds_since_midnight);
        end

        function t = from_components(hour, minute, second, nanosecond)
            if ~(hour >= 0 && hour <= 23)
                throw(yardl.ValueError("hour must be between 0 and 23"));
            end
            if ~(minute >= 0 && minute <= 59)
                throw(yardl.ValueError("minute must be between 0 and 59"));
            end
            if ~(second >= 0 && second <= 59)
                throw(yardl.ValueError("second must be between 0 and 59"));
            end
            if ~(nanosecond >= 0 && nanosecond <= 999999999)
                throw(yardl.ValueError("nanosecond must be between 0 and 999999999"));
            end

            t = yardl.Time(hour * 60*60*1e9 + minute * 60*1e9 + second * 1e9 + nanosecond);
        end
    end

end
