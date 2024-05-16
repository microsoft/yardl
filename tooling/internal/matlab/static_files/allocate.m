% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

function res = allocate(classname, varargin)
    % Wrapper around zeros, used to preallocate arrays of yardl objects
    if classname == "string"
        res = strings(varargin{:});
    else
        res = zeros(varargin{:}, classname);
    end
end
