classdef ProtocolStateTestWriter < test_model.StateTestWriterBase

    methods (Access=protected)

        function write_an_int_(self, value)
        end

        function write_a_stream_(self, value)
        end

        function write_another_int_(self, value)
        end

        function end_stream_(self)
        end

        function close_(self)
        end

    end
end
