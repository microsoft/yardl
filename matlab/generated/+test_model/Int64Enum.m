classdef Int64Enum < int64
  methods (Static)
    function e = B
      e = test_model.Int64Enum(-4611686018427387904);
    end

    function z = zeros(varargin)
      elem = test_model.Int64Enum(0);
      if nargin == 0
        z = elem;
      elseif nargin == 1
        n = varargin{1};
        z = reshape(repelem(elem, n*n), [n, n]);
      else
        sz = [varargin{:}];
        z = reshape(repelem(elem, prod(sz)), sz);
      end
    end
  end
end
