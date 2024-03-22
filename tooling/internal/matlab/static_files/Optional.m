classdef Optional < handle
    properties (SetAccess=protected)
        has_value
        value
    end

    methods
        function obj = Optional(varargin)
            if nargin > 0 && ~isa(varargin{1}, 'yardl.Optional')
                obj.value = varargin{1};
                obj.has_value = true;
            else
                obj.value = NaN;
                obj.has_value = false;
            end
        end

        function v = get.value(obj)
            if ~obj.has_value
                throw(yardl.ValueError("Optional type does not have a value"));
            end
            v = obj.value;
        end

        function eq = eq(a, b)
            % if isa(a, 'yardl.Optional')
            %     if isa(b, 'yardl.Optional')
            %         if a.has_value && b.has_value
            %             eq = a.value == b.value;
            %         else
            %             eq = a.has_value == b.has_value;
            %         end
            %     else
            %         eq = a.has_value && b == a.value;
            %     end
            % else
            %     % b is the Optional
            %     eq = b.has_value && a == b.value;
            % end

            if isa(a, 'yardl.Optional')
                if isa(b, 'yardl.Optional')
                    if all([a.has_value]) && all([b.has_value])
                        % eq = [a.value] == [b.value];
                        eq = isequal([a.value], [b.value]);
                    else
                        eq = [a.has_value] == [b.has_value];
                    end
                else
                    if all([a.has_value])
                        % eq = b == [a.value];
                        eq = isequal(b, [a.value]);
                    else
                        eq = false;
                    end
                end
            else
                % b is the Optional
                if all([b.has_value])
                    % eq = a == [b.value];
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
            elseif nargin == 1
                n = varargin{1};
                z = reshape(repelem(elem, n*n), [n, n]);
            else
                sz = [varargin{:}];
                z = reshape(repelem(elem, prod(sz)), sz);
            end
        end
    end
end
