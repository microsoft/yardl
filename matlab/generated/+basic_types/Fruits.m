% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef Fruits < uint64
  methods (Static)
    function v = APPLE
      v = basic_types.Fruits(0);
    end
    function v = BANANA
      v = basic_types.Fruits(1);
    end
    function v = PEAR
      v = basic_types.Fruits(2);
    end

    function z = zeros(varargin)
      elem = basic_types.Fruits(0);
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
