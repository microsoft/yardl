classdef OptionalSerializer < yardl.binary.TypeSerializer
    properties
        element_serializer_;
    end

    methods
        function obj = OptionalSerializer(element_serializer)
            obj.element_serializer_ = element_serializer;
        end

        function write(obj, outstream, value)
            % if isa(value, 'yardl.None')
            %     outstream.write_byte_no_check(0);
            %     return
            % end
            % outstream.write_byte_no_check(1);
            % obj.element_serializer_.write(outstream, value);

            if isa(value, 'yardl.Optional')
                if value.has_value()
                    outstream.write_byte_no_check(1);
                    obj.element_serializer_.write(outstream, value.value());
                else
                    outstream.write_byte_no_check(0);
                    return
                end
            else
                outstream.write_byte_no_check(1);
                obj.element_serializer_.write(outstream, value);
            end
        end

        function res = read(obj, instream)
            has_value = instream.read_byte();
            if has_value == 0
                res = yardl.None;
            else
                res = obj.element_serializer_.read(instream);

                % value = obj.element_serializer_.read(instream);
                % res = yardl.Optional(value);
            end
        end

        function c = getClass(obj)
            % c = obj.element_serializer_.getClass();
            c = 'yardl.Optional';
        end
    end
end
