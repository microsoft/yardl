% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithEnums < handle
  properties
    enum
    flags
    flags_2
  end

  methods
    function self = RecordWithEnums(enum, flags, flags_2)
      if nargin > 0
        self.enum = enum;
        self.flags = flags;
        self.flags_2 = flags_2;
      else
        self.enum = basic_types.Fruits.APPLE;
        self.flags = basic_types.DaysOfWeek(0);
        self.flags_2 = basic_types.TextFormat.REGULAR;
      end
    end

    function res = eq(self, other)
      res = ...
        isa(other, "test_model.RecordWithEnums") && ...
        all([self.enum] == [other.enum]) && ...
        all([self.flags] == [other.flags]) && ...
        all([self.flags_2] == [other.flags_2]);
    end

    function res = ne(self, other)
      res = ~self.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithEnums();
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
