% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef NoneSerializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(~, ~)
        end

        function res = read(~)
            res = yardl.None;
        end

        function c = get_class()
            c = "yardl.Optional";
        end
    end
end
