% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithVectorOfTimes < handle
  properties
    times
  end

  methods
    function self = RecordWithVectorOfTimes(times)
      if nargin > 0
        self.times = times;
      else
        self.times = yardl.Time.empty();
      end
    end

    function res = eq(self, other)
      res = ...
        isa(other, "test_model.RecordWithVectorOfTimes") && ...
        all([self.times] == [other.times]);
    end

    function res = ne(self, other)
      res = ~self.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithVectorOfTimes();
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
end
