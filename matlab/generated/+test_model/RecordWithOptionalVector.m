classdef RecordWithOptionalVector < handle
  properties
    optional_vector
  end

  methods
    function obj = RecordWithOptionalVector(optional_vector)
      if nargin > 0
        obj.optional_vector = optional_vector;
      else
        obj.optional_vector = yardl.None;
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithOptionalVector') && ...
        all([obj.optional_vector] == [other.optional_vector]);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithOptionalVector();
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
