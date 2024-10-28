% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

function err = RuntimeError(varargin)
    err = yardl.Exception("yardl:RuntimeError", varargin{:});
end
