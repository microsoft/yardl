% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithEnums < handle
  properties
    enum
    flags
    flags_2
    rec
  end

  methods
    function self = RecordWithEnums(kwargs)
      arguments
        kwargs.enum;
        kwargs.flags = basic_types.DaysOfWeek(0);
        kwargs.flags_2 = basic_types.TextFormat.REGULAR;
        kwargs.rec;
      end
      if ~isfield(kwargs, "enum")
        throw(yardl.TypeError("Missing required keyword argument 'enum'"))
      end
      self.enum = kwargs.enum;
      self.flags = kwargs.flags;
      self.flags_2 = kwargs.flags_2;
      if ~isfield(kwargs, "rec")
        throw(yardl.TypeError("Missing required keyword argument 'rec'"))
      end
      self.rec = kwargs.rec;
    end

    function res = eq(self, other)
      res = ...
        isa(other, "test_model.RecordWithEnums") && ...
        isequal({self.enum}, {other.enum}) && ...
        isequal({self.flags}, {other.flags}) && ...
        isequal({self.flags_2}, {other.flags_2}) && ...
        isequal({self.rec}, {other.rec});
    end

    function res = ne(self, other)
      res = ~self.eq(other);
    end

    function res = isequal(self, other)
      res = all(eq(self, other));
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithEnums(enum=yardl.None, rec=yardl.None);
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
