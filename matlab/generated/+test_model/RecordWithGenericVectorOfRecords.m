% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithGenericVectorOfRecords < handle
  properties
    v
  end

  methods
    function self = RecordWithGenericVectorOfRecords(v)
      self.v = v;
    end

    function res = eq(self, other)
      res = ...
        isa(other, "test_model.RecordWithGenericVectorOfRecords") && ...
        isequal(self.v, other.v);
    end

    function res = ne(self, other)
      res = ~self.eq(other);
    end
  end

end
