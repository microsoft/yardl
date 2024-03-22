classdef RecordWithAliasedGenerics < handle
  properties
    my_strings
    aliased_strings
  end

  methods
    function obj = RecordWithAliasedGenerics(my_strings, aliased_strings)
      if nargin > 0
        obj.my_strings = my_strings;
        obj.aliased_strings = aliased_strings;
      else
        obj.my_strings = tuples.Tuple("", "");
        obj.aliased_strings = tuples.Tuple("", "");
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithAliasedGenerics') && ...
        isequal(obj.my_strings, other.my_strings) && ...
        isequal(obj.aliased_strings, other.aliased_strings);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = test_model.RecordWithAliasedGenerics();
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
