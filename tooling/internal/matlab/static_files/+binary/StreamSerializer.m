% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef StreamSerializer < yardl.binary.TypeSerializer
    properties
        item_serializer_;
        items_remaining_;
    end

    methods
        function obj = StreamSerializer(item_serializer)
            obj.item_serializer_ = item_serializer;
            obj.items_remaining_ = 0;
        end

        function write(obj, outstream, values)
            if isempty(values)
                return;
            end

            if iscolumn(values)
                values = transpose(values);
            end
            s = size(values);
            count = s(end);
            outstream.write_unsigned_varint(count);

            if iscell(values)
                assert(s(1) == 1);
                for i = 1:count
                    obj.item_serializer_.write(outstream, values{i});
                end
            else
                if ndims(values) > 2
                    r = reshape(values, [], count);
                    inner_shape = s(1:end-1);
                    for i = 1:count
                        val = reshape(r(:, i), inner_shape);
                        obj.item_serializer_.write(outstream, val);
                    end
                else
                    for i = 1:count
                        obj.item_serializer_.write(outstream, transpose(values(:, i)));
                    end
                end
            end
        end

        function res = hasnext(obj, instream)
            if obj.items_remaining_ <= 0
                obj.items_remaining_ = instream.read_unsigned_varint();
                if obj.items_remaining_ <= 0
                    res = false;
                    return;
                end
            end
            res = true;
        end

        function res = read(obj, instream)
            if obj.items_remaining_ <= 0
                throw(yardl.RuntimeError("Stream has been exhausted"));
            end

            res = obj.item_serializer_.read(instream);
            obj.items_remaining_ = obj.items_remaining_ - 1;
        end

        function c = getClass(obj)
            c = obj.item_serializer_.getClass();
        end

        function s = getShape(~)
            s = [];
        end
    end
end
