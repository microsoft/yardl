% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

function err = KeyError(varargin)
    err = yardl.Exception("yardl:KeyError", varargin{:});
end
