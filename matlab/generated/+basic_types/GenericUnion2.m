classdef GenericUnion2 < yardl.Union
  methods (Static)
    function res = T1(value)
      res = basic_types.GenericUnion2(1, value);
    end

    function res = T2(value)
      res = basic_types.GenericUnion2(2, value);
    end

    function z = zeros(varargin)
      elem = basic_types.GenericUnion2(0, yardl.None);
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
      eq = isa(other, 'basic_types.GenericUnion2') && other.index == self.index && other.value == self.value;
    end

    function ne = ne(self, other)
      ne = ~self.eq(other);
    end
  end
end
