classdef Fruits < uint64
  methods (Static)
    function e = APPLE
      e = basic_types.Fruits(0);
    end
    function e = BANANA
      e = basic_types.Fruits(1);
    end
    function e = PEAR
      e = basic_types.Fruits(2);
    end

    function z = zeros(varargin)
      elem = basic_types.Fruits(0);
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
