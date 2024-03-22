classdef DaysOfWeek < uint64
  methods (Static)
    function e = MONDAY
      e = basic_types.DaysOfWeek(1);
    end
    function e = TUESDAY
      e = basic_types.DaysOfWeek(2);
    end
    function e = WEDNESDAY
      e = basic_types.DaysOfWeek(4);
    end
    function e = THURSDAY
      e = basic_types.DaysOfWeek(8);
    end
    function e = FRIDAY
      e = basic_types.DaysOfWeek(16);
    end
    function e = SATURDAY
      e = basic_types.DaysOfWeek(32);
    end
    function e = SUNDAY
      e = basic_types.DaysOfWeek(64);
    end

    function z = zeros(varargin)
      elem = basic_types.DaysOfWeek(0);
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
