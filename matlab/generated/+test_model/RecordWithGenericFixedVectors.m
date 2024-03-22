classdef RecordWithGenericFixedVectors < handle
  properties
    fv
    afv
  end

  methods
    function obj = RecordWithGenericFixedVectors(fv, afv)
      obj.fv = fv;
      obj.afv = afv;
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithGenericFixedVectors') && ...
        isequal(obj.fv, other.fv) && ...
        isequal(obj.afv, other.afv);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

end