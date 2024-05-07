% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef GenericRecordWithComputedFields < handle
  properties
    f1
  end

  methods
    function self = GenericRecordWithComputedFields(kwargs)
      arguments
        kwargs.f1;
      end
      if ~isfield(kwargs, "f1")
        throw(yardl.TypeError("Missing required keyword argument 'f1'"))
      end
      self.f1 = kwargs.f1;
    end

    function res = type_index(self)
      var1 = self.f1;
      if isa(var1, "basic_types.T0OrT1") && var1.index == 1
        res = 0;
        return
      end
      if isa(var1, "basic_types.T0OrT1") && var1.index == 2
        res = 1;
        return
      end
      throw(yardl.RuntimeError("Unexpected union case"))
    end


    function res = eq(self, other)
      res = ...
        isa(other, "basic_types.GenericRecordWithComputedFields") && ...
        isequal(self.f1, other.f1);
    end

    function res = ne(self, other)
      res = ~self.eq(other);
    end
  end

end
