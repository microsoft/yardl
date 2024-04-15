% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef Optional < handle
    properties (SetAccess=protected)
        has_value
        value
    end

    methods
        function self = Optional(varargin)
            if nargin > 0
                self.value = varargin{1};
                self.has_value = true;
            else
                self.value = NaN;
                self.has_value = false;
            end
        end

        function v = get.value(self)
            if ~self.has_value
                throw(yardl.ValueError("Optional type does not have a value"));
            end
            v = self.value;
        end

        function eq = eq(a, b)
            if isa(a, 'yardl.Optional')
                if isa(b, 'yardl.Optional')
                    if all([a.has_value]) && all([b.has_value])
                        eq = isequal([a.value], [b.value]);
                    else
                        eq = [a.has_value] == [b.has_value];
                    end
                else
                    if all([a.has_value])
                        eq = isequal(b, [a.value]);
                    else
                        eq = false;
                    end
                end
            else
                % b is the Optional
                if all([b.has_value])
                    eq = isequal(a, [b.value]);
                else
                    eq = false;
                end
            end
        end

        function ne = ne(a, b)
            ne = ~eq(a, b);
        end

        function isequal = isequal(a, b)
            isequal = all(eq(a, b));
        end

        function isequaln = isequaln(a, b)
            isequaln = all(eq(a, b));
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
