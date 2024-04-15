% This file was generated by the "yardl" tool. DO NOT EDIT.

classdef RecordWithUnions < handle
  properties
    null_or_int_or_string
    date_or_datetime
    null_or_fruits_or_days_of_week
  end

  methods
    function self = RecordWithUnions(null_or_int_or_string, date_or_datetime, null_or_fruits_or_days_of_week)
      if nargin > 0
        self.null_or_int_or_string = null_or_int_or_string;
        self.date_or_datetime = date_or_datetime;
        self.null_or_fruits_or_days_of_week = null_or_fruits_or_days_of_week;
      else
        self.null_or_int_or_string = yardl.None;
        self.date_or_datetime = basic_types.TimeOrDatetime.Time(yardl.Time());
        self.null_or_fruits_or_days_of_week = yardl.None;
      end
    end

    function res = eq(self, other)
      res = ...
        isa(other, 'basic_types.RecordWithUnions') && ...
        all([self.null_or_int_or_string] == [other.null_or_int_or_string]) && ...
        all([self.date_or_datetime] == [other.date_or_datetime]) && ...
        all([self.null_or_fruits_or_days_of_week] == [other.null_or_fruits_or_days_of_week]);
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
