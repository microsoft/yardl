% Alternative to Matlab's `ndims` function
% Collapses dimensions of size 1 (making it behave like numpy.array.ndim)
function c = dimension_count(arr)
    s = size(arr);
    c = length(s(s > 1));
end
