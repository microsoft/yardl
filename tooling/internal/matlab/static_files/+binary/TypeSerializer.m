% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef TypeSerializer < handle
    methods (Static, Abstract)
        write(obj, stream, value)
        res = read(obj, stream)
        c = getClass(obj)
    end

    methods (Static)
        function s = getShape()
            s = 1;
        end
    end
end
