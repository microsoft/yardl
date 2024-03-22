classdef Tuple < handle
  properties
    v1
    v2
  end

  methods
    function obj = Tuple(v1, v2)
      obj.v1 = v1;
      obj.v2 = v2;
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'tuples.Tuple') && ...
        isequal(obj.v1, other.v1) && ...
        isequal(obj.v2, other.v2);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

end