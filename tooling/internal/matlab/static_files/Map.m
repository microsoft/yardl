% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef Map < handle
    properties (SetAccess=protected)
        dict
    end

    methods

        function self = Map(varargin)
            self.dict = dictionary(varargin{:});
        end

        function k = keys(self)
            k = keys(self.dict);
        end

        function v = values(self)
            v = values(self.dict);
        end

        function v = lookup(self, key)
            v = self.dict(key);
        end

        function insert(self, key, value)
            self.dict(key) = value;
        end

        function remove(self, key)
            self.dict(key) = [];
        end

        function n = numEntries(self)
            n = numEntries(self.dict);
        end

        function res = eq(a, b)
            if isa(b, 'yardl.Map')
                res = isequal({a.dict}, {b.dict});
            elseif isa(b, 'dictionary')
                res = isequal({a.dict}, {b});
            else
                res = false;
            end
        end

        function ne = ne(a, b)
            ne = ~eq(a, b);
        end

        function isequal = isequal(self, other)
            isequal = all(eq(self, other));
        end
    end

    methods (Static)
        function z = zeros(varargin)
            elem = yardl.Map();
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
