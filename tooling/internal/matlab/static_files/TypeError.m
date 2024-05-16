% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

function err = TypeError(varargin)
    err = yardl.Exception("yardl:TypeError", varargin{:});
end
