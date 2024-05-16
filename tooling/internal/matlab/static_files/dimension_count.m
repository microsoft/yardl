% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

function c = dimension_count(arr)
    % Alternative to Matlab's `ndims` function
    % Collapses dimensions of size 1 (making it behave like numpy.array.ndim)
    s = size(arr);
    c = length(s(s > 1));
end
