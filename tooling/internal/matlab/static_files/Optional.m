% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef Optional < handle
    properties (SetAccess=protected)
        has_value
    end

    properties (Access=protected)
        value_
    end

    methods
        function self = Optional(varargin)
            if nargin > 0
                self.value_ = varargin{1};
                self.has_value = true;
            else
                self.has_value = false;
            end
        end

        function v = value(self)
            if ~self.has_value
                throw(yardl.ValueError("Optional type does not have a value"));
            end
            v = self.value_;
        end

        function eq = eq(a, b)
            if isa(a, 'yardl.Optional')
                if isa(b, 'yardl.Optional')
                    if all([a.has_value]) && all([b.has_value])
                        eq = isequal([a.value_], [b.value_]);
                    else
                        eq = [a.has_value] == [b.has_value];
                    end
                else
                    if all([a.has_value])
                        eq = isequal(b, [a.value_]);
                    else
                        eq = false;
                    end
                end
            else
                % b is the Optional
                if all([b.has_value])
                    eq = isequal(a, [b.value_]);
                else
                    eq = false;
                end
            end
        end

        function ne = ne(a, b)
            ne = ~eq(a, b);
        end

        function isequal = isequal(a, varargin)
            isequal = all(eq(a, [varargin{:}]));
        end

        function isequaln = isequaln(a, varargin)
            isequaln = isequal(a, [varargin{:}]);
        end
    end

    methods (Static)
        function z = zeros(varargin)
            elem = yardl.None;
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
