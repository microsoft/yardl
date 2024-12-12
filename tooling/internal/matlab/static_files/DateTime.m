% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef DateTime < handle
    % A basic datetime with nanosecond precision, always in UTC.

    properties (SetAccess=private)
        nanoseconds_since_epoch
    end

    methods
        function self = DateTime(nanoseconds_since_epoch)
            if nargin > 0
                self.nanoseconds_since_epoch = nanoseconds_since_epoch;
            else
                self.nanoseconds_since_epoch = 0;
            end
        end

        function value = value(self)
            value = self.nanoseconds_since_epoch;
        end

        function dt = to_datetime(self)
            dt = datetime(self.nanoseconds_since_epoch, 'ConvertFrom', 'epochtime', 'TicksPerSecond', 1e9);
        end

        function eq = eq(self, other)
            if isa(other, 'datetime')
                other = yardl.DateTime.from_datetime(other);
            end

            eq = isa(other, 'yardl.DateTime') && ...
                all([self.nanoseconds_since_epoch] == [other.nanoseconds_since_epoch]);
        end

        function ne = ne(self, other)
            ne = ~self.eq(other);
        end

        function isequal = isequal(self, other)
            isequal = all(eq(self, other));
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

        function dt = now(~)
            dt = yardl.DateTime.from_datetime(datetime('now'));
        end

        function z = zeros(varargin)
            elem = yardl.DateTime(0);
            if nargin == 0
                z = elem;
                return
            end

            sz = [varargin{:}];
            if isscalar(sz)
                sz = [sz, sz];
            end
            z = reshape(repelem(elem, prod(sz)), sz);
        end
    end

end
