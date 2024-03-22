classdef RecordWithGenericArrays < handle
  properties
    nd
    fixed_nd
    dynamic_nd
    aliased_nd
    aliased_fixed_nd
    aliased_dynamic_nd
  end

  methods
    function obj = RecordWithGenericArrays(nd, fixed_nd, dynamic_nd, aliased_nd, aliased_fixed_nd, aliased_dynamic_nd)
      obj.nd = nd;
      obj.fixed_nd = fixed_nd;
      obj.dynamic_nd = dynamic_nd;
      obj.aliased_nd = aliased_nd;
      obj.aliased_fixed_nd = aliased_fixed_nd;
      obj.aliased_dynamic_nd = aliased_dynamic_nd;
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithGenericArrays') && ...
        isequal(obj.nd, other.nd) && ...
        isequal(obj.fixed_nd, other.fixed_nd) && ...
        isequal(obj.dynamic_nd, other.dynamic_nd) && ...
        isequal(obj.aliased_nd, other.aliased_nd) && ...
        isequal(obj.aliased_fixed_nd, other.aliased_fixed_nd) && ...
        isequal(obj.aliased_dynamic_nd, other.aliased_dynamic_nd);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

end
