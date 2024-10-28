% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef VectorSerializer < yardl.binary.VectorSerializerBase

    methods (Access=protected)
        function handle_write_count_(~, outstream, count)
            outstream.write_unsigned_varint(count);
        end

        function count = get_read_count_(~, instream)
            count = instream.read_unsigned_varint();
        end
    end

    methods
        function self = VectorSerializer(item_serializer)
            self@yardl.binary.VectorSerializerBase(item_serializer);
        end

        function s = get_shape(~)
            s = [];
        end
    end
end
