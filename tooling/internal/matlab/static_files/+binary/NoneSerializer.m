% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef NoneSerializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
        end

        function res = read(instream)
            res = yardl.None;
        end

        function c = getClass()
            c = 'yardl.Optional';
        end
    end
end
