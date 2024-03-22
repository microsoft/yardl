classdef GenericUnion3Alternate < yardl.Union
  methods (Static)
    function res = U(value)
      res = test_model.GenericUnion3Alternate(1, value);
    end

    function res = V(value)
      res = test_model.GenericUnion3Alternate(2, value);
    end

    function res = W(value)
      res = test_model.GenericUnion3Alternate(3, value);
    end

    function z = zeros(varargin)
      elem = test_model.GenericUnion3Alternate(0, yardl.None);
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
      eq = isa(other, 'test_model.GenericUnion3Alternate') && other.index == self.index && other.value == self.value;
    end

    function ne = ne(self, other)
      ne = ~self.eq(other);
    end
  end
end
