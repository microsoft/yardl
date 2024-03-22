classdef TextFormat < uint64
  methods (Static)
    function e = REGULAR
      e = basic_types.TextFormat(0);
    end
    function e = BOLD
      e = basic_types.TextFormat(1);
    end
    function e = ITALIC
      e = basic_types.TextFormat(2);
    end
    function e = UNDERLINE
      e = basic_types.TextFormat(4);
    end
    function e = STRIKETHROUGH
      e = basic_types.TextFormat(8);
    end

    function z = zeros(varargin)
      elem = basic_types.TextFormat(0);
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
