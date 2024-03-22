classdef RecordWithKeywordFields < handle
  properties
    int
    sizeof
    if
  end

  methods
    function obj = RecordWithKeywordFields(int, sizeof, if)
      obj.int = int;
      obj.sizeof = sizeof;
      obj.if = if;
    end

    function res = float(self)
      res = self.int;
      return
    end

    function res = double(self)
      res = self.float();
      return
    end

    function res = return(self)
      res = self.sizeof(1, 2);
      return
    end


    function res = eq(obj, other)
      res = ...
        isa(other, 'test_model.RecordWithKeywordFields') && ...
        all([obj.int] == [other.int]) && ...
        isequal(obj.sizeof, other.sizeof) && ...
        all([obj.if] == [other.if]);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

end