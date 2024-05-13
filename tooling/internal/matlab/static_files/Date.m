% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef Date < handle

    properties (SetAccess=private)
        days_since_epoch
    end

    methods
        function self = Date(days_since_epoch)
            if nargin > 0
                self.days_since_epoch = days_since_epoch;
            else
                self.days_since_epoch = 0;
            end
        end

        function value = value(self)
            value = self.days_since_epoch;
        end

        function dt = to_datetime(self)
            dt = datetime(0, 'ConvertFrom', 'epochtime') + days(self.days_since_epoch);
        end

        function eq = eq(self, other)
            if isa(other, 'datetime')
                other = yardl.Date.from_datetime(other);
            end

            eq = isa(other, 'yardl.Date') && ...
                all([self.value] == [other.value]);
        end

        function ne = new(self, other)
            ne = ~self.eq(other);
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

        function z = zeros(varargin)
            elem = yardl.Date(0);
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
