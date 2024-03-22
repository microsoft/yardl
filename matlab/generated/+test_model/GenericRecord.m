classdef GenericRecord < handle
  properties
    scalar_1
    scalar_2
    vector_1
    image_2
  end

  methods
    function obj = GenericRecord(scalar_1, scalar_2, vector_1, image_2)
      obj.scalar_1 = scalar_1;
      obj.scalar_2 = scalar_2;
      obj.vector_1 = vector_1;
      obj.image_2 = image_2;
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.GenericRecord') && ...
        isequal(obj.scalar_1, other.scalar_1) && ...
        isequal(obj.scalar_2, other.scalar_2) && ...
        isequal(obj.vector_1, other.vector_1) && ...
        isequal(obj.image_2, other.image_2);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

end