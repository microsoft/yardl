classdef RecordContainingGenericRecords < handle
  properties
    g1
    g1a
    g2
    g2a
    g3
    g3a
    g4
    g5
    g6
    g7
  end

  methods
    function obj = RecordContainingGenericRecords(g1, g1a, g2, g2a, g3, g3a, g4, g5, g6, g7)
      obj.g1 = g1;
      obj.g1a = g1a;
      obj.g2 = g2;
      obj.g2a = g2a;
      obj.g3 = g3;
      obj.g3a = g3a;
      obj.g4 = g4;
      obj.g5 = g5;
      obj.g6 = g6;
      obj.g7 = g7;
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordContainingGenericRecords') && ...
        isequal(obj.g1, other.g1) && ...
        isequal(obj.g1a, other.g1a) && ...
        isequal(obj.g2, other.g2) && ...
        isequal(obj.g2a, other.g2a) && ...
        isequal(obj.g3, other.g3) && ...
        isequal(obj.g3a, other.g3a) && ...
        isequal(obj.g4, other.g4) && ...
        isequal(obj.g5, other.g5) && ...
        isequal(obj.g6, other.g6) && ...
        isequal(obj.g7, other.g7);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

end
