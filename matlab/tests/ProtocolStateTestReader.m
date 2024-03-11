classdef ProtocolStateTestReader < StateTestReaderBase

    methods (Access=protected)

        function res = read_an_int_(self)
            res = -2;
        end

        function res =read_a_stream_(self)
            res = [-1, -2, -3];
        end

        function res =read_another_int_(self)
            res = -4;
        end

        function close_(self)
        end

    end
end
