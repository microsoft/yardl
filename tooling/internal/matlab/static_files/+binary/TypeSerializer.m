% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef TypeSerializer < handle
    methods (Static, Abstract)
        write(obj, stream, value)
        res = read(obj, stream)
        c = getClass(obj)
    end

    methods (Static)
        function s = getShape()
            s = 1;
        end

        function trivial = isTriviallySerializable()
            trivial = false;
        end
    end

    methods
        function writeTrivially(obj, stream, values)
            if ~obj.isTriviallySerializable()
                error("Not implemented for non-trivially-serializable types")
            end
            stream.write_values_directly(values, obj.getClass());
        end

        function res = readTrivially(obj, stream, shape)
            if ~obj.isTriviallySerializable()
                error("Not implemented for non-trivially-serializable types")
            end
            res = stream.read_values_directly(shape, obj.getClass());
        end
    end
end
