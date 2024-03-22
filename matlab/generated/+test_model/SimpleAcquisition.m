classdef SimpleAcquisition < handle
  properties
    flags
    idx
    data
    trajectory
  end

  methods
    function obj = SimpleAcquisition(flags, idx, data, trajectory)
      if nargin > 0
        obj.flags = flags;
        obj.idx = idx;
        obj.data = data;
        obj.trajectory = trajectory;
      else
        obj.flags = uint64(0);
        obj.idx = test_model.SimpleEncodingCounters();
        obj.data = single.empty(0, 0);
        obj.trajectory = single.empty(0, 0);
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.SimpleAcquisition') && ...
        all([obj.flags] == [other.flags]) && ...
        all([obj.idx] == [other.idx]) && ...
        isequal(obj.data, other.data) && ...
        isequal(obj.trajectory, other.trajectory);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.SimpleAcquisition();
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