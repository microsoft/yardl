% A basic datetime with nanosecond precision, always in UTC.
classdef DateTime < handle

    properties (Access=private)
        nanoseconds_since_epoch_
    end

    methods
        function obj = DateTime(nanoseconds_since_epoch)
            obj.nanoseconds_since_epoch_ = nanoseconds_since_epoch;
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

            if isa(other, 'yardl.DateTime')
                eq = all([obj.value] == [other.value]);
            else
                eq = false;
            end
        end
    end

    methods (Static)
        function dt = from_datetime(value)
            value.TimeZone = 'UTC';
            nanoseconds_since_epoch = convertTo(value, 'epochtime', 'TicksPerSecond', 1e9);
            dt = yardl.DateTime(nanoseconds_since_epoch);
        end
    end

end
