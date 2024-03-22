classdef RecordWithOptionalFields < handle
  properties
    optional_int
    optional_int_alternate_syntax
    optional_time
  end

  methods
    function obj = RecordWithOptionalFields(optional_int, optional_int_alternate_syntax, optional_time)
      if nargin > 0
        obj.optional_int = optional_int;
        obj.optional_int_alternate_syntax = optional_int_alternate_syntax;
        obj.optional_time = optional_time;
      else
        obj.optional_int = yardl.None;
        obj.optional_int_alternate_syntax = yardl.None;
        obj.optional_time = yardl.None;
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithOptionalFields') && ...
        all([obj.optional_int] == [other.optional_int]) && ...
        all([obj.optional_int_alternate_syntax] == [other.optional_int_alternate_syntax]) && ...
        all([obj.optional_time] == [other.optional_time]);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithOptionalFields();
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
