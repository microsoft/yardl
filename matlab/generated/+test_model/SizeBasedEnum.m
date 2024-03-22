classdef SizeBasedEnum < uint64
  methods (Static)
    function e = A
      e = test_model.SizeBasedEnum(0);
    end
    function e = B
      e = test_model.SizeBasedEnum(1);
    end
    function e = C
      e = test_model.SizeBasedEnum(2);
    end

    function z = zeros(varargin)
      elem = test_model.SizeBasedEnum(0);
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
