classdef TypeSerializer < handle
    methods (Static, Abstract)
        write(obj, stream, value)
        res = read(obj, stream)
    end
end
