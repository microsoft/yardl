% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef TypeSerializer < handle
    methods (Static, Abstract)
        write(self, stream, value)
        res = read(self, stream)
        c = get_class(self)
    end

    methods (Static)
        function s = get_shape()
            s = 1;
        end

        function trivial = is_trivially_serializable()
            trivial = false;
        end
    end

    methods
        function write_trivially(self, stream, values)
            if ~self.is_trivially_serializable()
                throw(yardl.TypeError("Not implemented for non-trivially-serializable types"));
            end
            stream.write_values_directly(values, self.get_class());
        end

        function res = read_trivially(self, stream, shape)
            if ~self.is_trivially_serializable()
                throw(yardl.TypeError("Not implemented for non-trivially-serializable types"));
            end
            res = stream.read_values_directly(shape, self.get_class());
        end
    end
end
