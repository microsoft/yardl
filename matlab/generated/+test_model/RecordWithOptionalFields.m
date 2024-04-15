% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithOptionalFields < handle
  properties
    optional_int
    optional_int_alternate_syntax
    optional_time
  end

  methods
    function self = RecordWithOptionalFields(optional_int, optional_int_alternate_syntax, optional_time)
      if nargin > 0
        self.optional_int = optional_int;
        self.optional_int_alternate_syntax = optional_int_alternate_syntax;
        self.optional_time = optional_time;
      else
        self.optional_int = yardl.None;
        self.optional_int_alternate_syntax = yardl.None;
        self.optional_time = yardl.None;
      end
    end

    function res = eq(self, other)
      res = ...
        isa(other, 'test_model.RecordWithOptionalFields') && ...
        all([self.optional_int] == [other.optional_int]) && ...
        all([self.optional_int_alternate_syntax] == [other.optional_int_alternate_syntax]) && ...
        all([self.optional_time] == [other.optional_time]);
    end

    function res = ne(self, other)
      res = ~self.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithOptionalFields();
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
