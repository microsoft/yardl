% A basic datetime with nanosecond precision, always in UTC.
classdef DateTime < handle

    properties (Access=private)
        nanoseconds_since_epoch_
    end

    methods
        function obj = DateTime(nanoseconds_since_epoch)
            if nargin > 0
                obj.nanoseconds_since_epoch_ = nanoseconds_since_epoch;
            else
                obj.nanoseconds_since_epoch_ = 0;
            end
        end

        function value = value(obj)
            value = obj.nanoseconds_since_epoch_;
        end

        function dt = to_datetime(obj)
            dt = datetime(obj.nanoseconds_since_epoch_, 'ConvertFrom', 'epochtime', 'TicksPerSecond', 1e9);
        end

        function eq = eq(obj, other)
            if isa(other, 'datetime')
                other = yardl.DateTime.from_datetime(other);
            end

            eq = isa(other, 'yardl.DateTime') && ...
                all([obj.value] == [other.value]);
        end

        function ne = new(obj, other)
            ne = ~obj.eq(other);
        end
    end

    methods (Static)
        function dt = from_datetime(value)
            value.TimeZone = 'UTC';
            nanoseconds_since_epoch = convertTo(value, 'epochtime', 'TicksPerSecond', 1e9);
            dt = yardl.DateTime(nanoseconds_since_epoch);
        end

        function dt = from_components(year, month, day, hour, minute, second, nanosecond)
            if ~(nanosecond >= 0 && nanosecond < 999999999)
                throw(yardl.ValueError("nanosecond must be in 0..1e9"));
            end
            mdt = datetime(year, month, day, hour, minute, second, 'TimeZone', 'UTC');
            seconds_since_epoch = convertTo(mdt, 'epochtime');
            dt = yardl.DateTime(seconds_since_epoch * 1e9 + nanosecond);
        end
    end

end
