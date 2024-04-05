% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

function err = ValueError(varargin)
    err = yardl.Exception("yardl:ValueError", varargin{:});
end
