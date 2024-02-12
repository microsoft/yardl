classdef StreamSerializer < yardl.binary.TypeSerializer
    properties
        element_serializer_;
    end

    methods
        function obj = StreamSerializer(element_serializer)
            obj.element_serializer_ = element_serializer;
        end

        function write(obj, outstream, value)
            outstream.write_unsigned_varint(length(value));
            % TODO: Optimize this?
            for i = 1:length(value)
                obj.element_serializer_.write(outstream, value(i));
            end
        end

        function res = read(obj, instream)
            count = instream.read_unsigned_varint();
            idx = 1;
            while count > 0
                for c = 1:count
                    % TODO: Optimize this "append" approach
                    res(idx) = obj.element_serializer_.read(instream);
                    idx = idx + 1;
                end
                count = instream.read_unsigned_varint();
            end
        end
    end
end
