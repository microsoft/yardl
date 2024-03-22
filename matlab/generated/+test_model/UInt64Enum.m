classdef UInt64Enum < uint64
  methods (Static)
    function e = A
      e = test_model.UInt64Enum(9223372036854775808);
    end

    function z = zeros(varargin)
      elem = test_model.UInt64Enum(0);
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
