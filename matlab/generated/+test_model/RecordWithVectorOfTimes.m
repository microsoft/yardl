classdef RecordWithVectorOfTimes < handle
  properties
    times
  end

  methods
    function obj = RecordWithVectorOfTimes(times)
      if nargin > 0
        obj.times = times;
      else
        obj.times = yardl.Time.empty();
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithVectorOfTimes') && ...
        all([obj.times] == [other.times]);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithVectorOfTimes();
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
