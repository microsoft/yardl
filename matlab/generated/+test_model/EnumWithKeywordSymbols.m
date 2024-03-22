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
