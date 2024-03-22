classdef RecordWithEnums < handle
  properties
    enum
    flags
    flags_2
  end

  methods
    function obj = RecordWithEnums(enum, flags, flags_2)
      if nargin > 0
        obj.enum = enum;
        obj.flags = flags;
        obj.flags_2 = flags_2;
      else
        obj.enum = basic_types.Fruits.APPLE;
        obj.flags = basic_types.DaysOfWeek(0);
        obj.flags_2 = basic_types.TextFormat.REGULAR;
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithEnums') && ...
        all([obj.enum] == [other.enum]) && ...
        all([obj.flags] == [other.flags]) && ...
        all([obj.flags_2] == [other.flags_2]);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithEnums();
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