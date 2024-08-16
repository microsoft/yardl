% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithNoDefaultEnum < handle
  properties
    enum
  end

  methods
    function self = RecordWithNoDefaultEnum(kwargs)
      arguments
        kwargs.enum;
      end
      if ~isfield(kwargs, "enum")
        throw(yardl.TypeError("Missing required keyword argument 'enum'"))
      end
      self.enum = kwargs.enum;
    end

    function res = eq(self, other)
      res = ...
        isa(other, "test_model.RecordWithNoDefaultEnum") && ...
        isequal(self.enum, other.enum);
    end

    function res = ne(self, other)
      res = ~self.eq(other);
    end
  end

end