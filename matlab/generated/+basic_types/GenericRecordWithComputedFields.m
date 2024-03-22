classdef GenericRecordWithComputedFields < handle
  properties
    f1
  end

  methods
    function obj = GenericRecordWithComputedFields(f1)
      obj.f1 = f1;
    end

    function res = type_index(self)
      var1 = self.f1;
      if var1.index == 1
        res = 0;
        return
      end
      if var1.index == 2
        res = 1;
        return
      end
      throw(yardl.RuntimeError("Unexpected union case"))
    end


    function res = eq(obj, other)
      res = ...
        isa(other, 'basic_types.GenericRecordWithComputedFields') && ...
        isequal(obj.f1, other.f1);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

end
