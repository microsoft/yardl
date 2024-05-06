% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef OptionalSerializer < yardl.binary.TypeSerializer
    properties
        item_serializer_;
    end

    methods
        function self = OptionalSerializer(item_serializer)
            self.item_serializer_ = item_serializer;
        end

        function write(self, outstream, value)
            if isa(value, "yardl.Optional")
                if value.has_value()
                    outstream.write_byte(uint8(1));
                    self.item_serializer_.write(outstream, value.value());
                else
                    outstream.write_byte(uint8(0));
                    return
                end
            else
                outstream.write_byte(uint8(1));
                self.item_serializer_.write(outstream, value);
            end
        end

        function res = read(self, instream)
            % Returns either yardl.None or the inner optional value
            has_value = instream.read_byte();
            if has_value == 0
                res = yardl.None;
            else
                res = self.item_serializer_.read(instream);
            end
        end

        function c = get_class(~)
            c = "yardl.Optional";
        end
    end
end
