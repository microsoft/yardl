% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef ProtocolStateTestReader < test_model.StateTestReaderBase

    methods (Access=protected)

        function res = read_an_int_(self)
            res = -2;
        end

        function more = has_a_stream_(self)
            persistent n
            if isempty(n)
                n = 0;
                more = true;
            else
                more = false;
            end
            n = n + 1;
        end

        function res = read_a_stream_(self)
            res = [-1, -2, -3];
        end

        function res = read_another_int_(self)
            res = -4;
        end

        function close_(self)
        end

    end
end
