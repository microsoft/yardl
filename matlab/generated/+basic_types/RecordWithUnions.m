classdef RecordWithUnions < handle
  properties
    null_or_int_or_string
    date_or_datetime
    null_or_fruits_or_days_of_week
  end

  methods
    function obj = RecordWithUnions(null_or_int_or_string, date_or_datetime, null_or_fruits_or_days_of_week)
      if nargin > 0
        obj.null_or_int_or_string = null_or_int_or_string;
        obj.date_or_datetime = date_or_datetime;
        obj.null_or_fruits_or_days_of_week = null_or_fruits_or_days_of_week;
      else
        obj.null_or_int_or_string = yardl.None;
        obj.date_or_datetime = basic_types.TimeOrDatetime.Time(yardl.Time());
        obj.null_or_fruits_or_days_of_week = yardl.None;
      end
    end

    function res = eq(obj, other)
      res = ...
        isa(other, 'basic_types.RecordWithUnions') && ...
        all([obj.null_or_int_or_string] == [other.null_or_int_or_string]) && ...
        all([obj.date_or_datetime] == [other.date_or_datetime]) && ...
        all([obj.null_or_fruits_or_days_of_week] == [other.null_or_fruits_or_days_of_week]);
    end

    function res = ne(obj, other)
      res = ~obj.eq(other);
    end
  end

  methods (Static)
    function z = zeros(varargin)
      elem = basic_types.RecordWithUnions();
      if nargin == 0
        z = elem;
      elseif nargin == 1
        n = varargin{1};
        z = reshape(repelem(elem, n*n), [n, n]);
      else
        sz = [varargin{:}];
        z = reshape(repelem(elem, prod(sz)), sz);
      end
    end
  end
end
