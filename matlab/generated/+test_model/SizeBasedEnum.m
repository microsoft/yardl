% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef SizeBasedEnum < uint64
  methods (Static)
    function v = A
      v = test_model.SizeBasedEnum(0);
    end
    function v = B
      v = test_model.SizeBasedEnum(1);
    end
    function v = C
      v = test_model.SizeBasedEnum(2);
    end

    function z = zeros(varargin)
      elem = test_model.SizeBasedEnum(0);
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
