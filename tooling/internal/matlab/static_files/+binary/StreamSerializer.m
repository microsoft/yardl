classdef StreamSerializer < yardl.binary.TypeSerializer
    properties
        element_serializer_;
    end

    methods
        function obj = StreamSerializer(element_serializer)
            obj.element_serializer_ = element_serializer;
        end

        function write(obj, outstream, value)
            len = size(value, 2);
            outstream.write_unsigned_varint(len);
            for i = 1:len
                obj.element_serializer_.write(outstream, value(:, i));
            end
        end

        function res = read(obj, instream)
            count = instream.read_unsigned_varint();
            if count == 0
                res = [];
                return;
            end

            % Preallocate the result vector (THIS DOESN'T WORK if the element is a vector or array)
            res = zeros(1, count, obj.getClass);
            idx = 1;
            while count > 0
                for c = 1:count
                    res(:, idx) = obj.element_serializer_.read(instream);
                    idx = idx + 1;
                end
                count = instream.read_unsigned_varint();
            end
        end

        function c = getClass(obj)
            c = obj.element_serializer_.getClass();
        end
    end
end
