classdef NoneSerializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
        end

        function res = read(instream)
            res = yardl.None;
        end
    end
end
