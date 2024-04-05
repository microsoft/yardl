% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

function err = Exception(id, varargin)
    err = MException(id, varargin{:});
end
