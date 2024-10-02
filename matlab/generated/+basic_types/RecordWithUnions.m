% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithUnions < handle
  properties
    null_or_int_or_string
    date_or_datetime
    null_or_fruits_or_days_of_week
    record_or_int
  end

  methods
    function self = RecordWithUnions(kwargs)
      arguments
        kwargs.null_or_int_or_string = yardl.None;
        kwargs.date_or_datetime = basic_types.TimeOrDatetime.Time(yardl.Time());
        kwargs.null_or_fruits_or_days_of_week = yardl.None;
        kwargs.record_or_int = basic_types.RecordWithStringOrInt32.RecordWithString(basic_types.RecordWithString());
      end
      self.null_or_int_or_string = kwargs.null_or_int_or_string;
      self.date_or_datetime = kwargs.date_or_datetime;
      self.null_or_fruits_or_days_of_week = kwargs.null_or_fruits_or_days_of_week;
      self.record_or_int = kwargs.record_or_int;
    end

    function res = eq(self, other)
      res = ...
        isa(other, "basic_types.RecordWithUnions") && ...
        isequal(self.null_or_int_or_string, other.null_or_int_or_string) && ...
        isequal(self.date_or_datetime, other.date_or_datetime) && ...
        isequal(self.null_or_fruits_or_days_of_week, other.null_or_fruits_or_days_of_week) && ...
        isequal(self.record_or_int, other.record_or_int);
    end

    function res = ne(self, other)
      res = ~self.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = basic_types.RecordWithUnions();
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
end
