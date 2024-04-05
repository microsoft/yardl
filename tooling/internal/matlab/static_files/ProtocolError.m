% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

function err = ProtocolError(varargin)
    err = yardl.Exception("yardl:ProtocolError", varargin{:});
end
