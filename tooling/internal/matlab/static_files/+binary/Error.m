% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

function err = Error(varargin)
    err = yardl.Exception("yardl:binary:Error", varargin{:});
end
