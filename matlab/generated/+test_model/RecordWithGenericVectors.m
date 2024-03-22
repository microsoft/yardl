classdef RecordWithGenericVectors < handle
  properties
    v
    av
  end

  methods
    function obj = RecordWithGenericVectors(v, av)
      obj.v = v;
      obj.av = av;
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithGenericVectors') && ...
        isequal(obj.v, other.v) && ...
        isequal(obj.av, other.av);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

end
