% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray < yardl.Union
  methods (Static)
    function res = Int32(value)
      res = test_model.Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray(1, value);
    end

    function res = Float32(value)
      res = test_model.Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray(2, value);
    end

    function res = String(value)
      res = test_model.Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray(3, value);
    end

    function res = SimpleRecord(value)
      res = test_model.Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray(4, value);
    end

    function res = NamedFixedNDArray(value)
      res = test_model.Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray(5, value);
    end

    function z = zeros(varargin)
      elem = test_model.Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray(0, yardl.None);
      if nargin == 0
        z = elem;
        return;
      end
      sz = [varargin{:}];
      if isscalar(sz)
        sz = [sz, sz];
      end
      z = reshape(repelem(elem, prod(sz)), sz);
    end
  end

  methods
    function res = isInt32(self)
      res = self.index == 1;
    end

    function res = isFloat32(self)
      res = self.index == 2;
    end

    function res = isString(self)
      res = self.index == 3;
    end

    function res = isSimpleRecord(self)
      res = self.index == 4;
    end

    function res = isNamedFixedNDArray(self)
      res = self.index == 5;
    end

    function eq = eq(self, other)
      eq = isa(other, "test_model.Int32OrFloat32OrStringOrSimpleRecordOrNamedFixedNDArray") && all([self.index_] == [other.index_], 'all') && all([self.value] == [other.value], 'all');
    end

    function ne = ne(self, other)
      ne = ~self.eq(other);
    end

    function t = tag(self)
      tags_ = ["Int32", "Float32", "String", "SimpleRecord", "NamedFixedNDArray"];
      t = tags_(self.index_);
    end
  end
end
