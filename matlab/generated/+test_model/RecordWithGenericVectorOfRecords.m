% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithGenericVectorOfRecords < handle
  properties
    v
  end

  methods
    function obj = RecordWithGenericVectorOfRecords(v)
      obj.v = v;
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithGenericVectorOfRecords') && ...
        isequal(obj.v, other.v);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

end
