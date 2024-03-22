classdef ArrayOrScalar < yardl.Union
  methods (Static)
    function res = Array(value)
      res = test_model.ArrayOrScalar(1, value);
    end

    function res = Scalar(value)
      res = test_model.ArrayOrScalar(2, value);
    end

    function z = zeros(varargin)
      elem = test_model.ArrayOrScalar(0, yardl.None);
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

  methods
    function eq = eq(self, other)
      eq = isa(other, 'test_model.ArrayOrScalar') && other.index == self.index && other.value == self.value;
    end

    function ne = ne(self, other)
      ne = ~self.eq(other);
    end
  end
end
