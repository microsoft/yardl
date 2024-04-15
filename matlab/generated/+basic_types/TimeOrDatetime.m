% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef TimeOrDatetime < yardl.Union
  methods (Static)
    function res = Time(value)
      res = basic_types.TimeOrDatetime(1, value);
    end

    function res = Datetime(value)
      res = basic_types.TimeOrDatetime(2, value);
    end

    function z = zeros(varargin)
      elem = basic_types.TimeOrDatetime(0, yardl.None);
      if nargin == 0
        z = elem;
        return;
      end
      sz = [varargin{:}];
      if isscalar(sz)
        sz = [sz, sz];
      end
      z = reshape(repelem(elem, prod(sz)), sz);
    end
  end

  methods
    function res = isTime(self)
      res = self.index == 1;
    end

    function res = isDatetime(self)
      res = self.index == 2;
    end

    function eq = eq(self, other)
      eq = isa(other, 'basic_types.TimeOrDatetime') && other.index == self.index && other.value == self.value;
    end

    function ne = ne(self, other)
      ne = ~self.eq(other);
    end
  end
end
