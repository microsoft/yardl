% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef EnumWithKeywordSymbols < uint64
  methods (Static)
    function e = TRY
      e = test_model.EnumWithKeywordSymbols(2);
    end
    function e = CATCH
      e = test_model.EnumWithKeywordSymbols(1);
    end

    function z = zeros(varargin)
      elem = test_model.EnumWithKeywordSymbols(0);
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
