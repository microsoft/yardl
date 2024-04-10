% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef StreamSerializer < yardl.binary.TypeSerializer
    properties
        item_serializer_;
    end

    methods
        function obj = StreamSerializer(item_serializer)
            obj.item_serializer_ = item_serializer;
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

        function res = read(obj, instream)
            count = instream.read_unsigned_varint();
            if count == 0
                res = [];
                return;
            end

            item_shape = obj.item_serializer_.getShape();
            if isempty(item_shape)
                res = cell(1, count);
                idx = 1;
                while count > 0
                    for c = 1:count
                        res{idx} = obj.item_serializer_.read(instream);
                        idx = idx + 1;
                    end
                    count = instream.read_unsigned_varint();
                end
                return
            end

            res = yardl.allocate(obj.getClass(), [prod(item_shape), count]);
            total_count = 0;
            while count > 0
                for c = 1:count
                    idx = total_count + c;
                    item = obj.item_serializer_.read(instream);
                    res(:, idx) = item(:);
                end

                total_count = total_count + count;
                count = instream.read_unsigned_varint();
            end

            res = squeeze(reshape(res, [item_shape, total_count]));
            if iscolumn(res)
                res = transpose(res);
            end
        end

        function c = getClass(obj)
            c = obj.item_serializer_.getClass();
        end

        function s = getShape(~)
            s = [];
        end
    end
end
