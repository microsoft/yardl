classdef VectorSerializer < yardl.binary.TypeSerializer
    properties
        item_serializer_;
    end

    methods
        function obj = VectorSerializer(item_serializer)
            obj.item_serializer_ = item_serializer;
        end

        function write(obj, outstream, values)
            s = size(values);
            count = s(end);
            outstream.write_unsigned_varint(count);

            if iscell(values)
                % values is a cell array, so must be a vector of shape [1, COUNT]
                assert(s(1) == 1);
                for i = 1:count
                    obj.item_serializer_.write(obj, values{i});
                end
            else
                % values is an array, so must have shape [A, B, ..., COUNT]
                if ndims(values) > 2
                    r = reshape(values, [], count);
                    for i = 1:count
                        val = reshape(r(:, i), s(1:end-1));
                        obj.item_serializer_.write(outstream, val);
                    end
                else
                    for i = 1:count
                        % obj.item_serializer_.write(outstream, values(:, i));
                        obj.item_serializer_.write(outstream, transpose(values(:, i)));
                    end
                end
            end
        end

        function res = read(obj, instream)
            count = instream.read_unsigned_varint();
            if count == 0
                % res = zeros(0, obj.getClass());
                res = yardl.allocate(obj.getClass(), 0);
                return;
            end

            item_shape = obj.item_serializer_.getShape();
            if isempty(item_shape)
                res = cell(1, count);
                for i = 1:count
                    res{i} = obj.item_serializer_.read(instream);
                end
            elseif isscalar(item_shape)
                % res = zeros([1 count], obj.getClass());
                res = yardl.allocate(obj.getClass(), [1 count]);
                for i = 1:count
                    res(i) = obj.item_serializer_.read(instream);
                end
                res = squeeze(res);
            else
                % res = zeros([prod(item_shape), count], obj.getClass());
                res = yardl.allocate(obj.getClass(), [prod(item_shape), count]);
                for i = 1:count
                    item = obj.item_serializer_.read(instream);
                    res(:, i) = item(:);
                end
                res = squeeze(reshape(res, [item_shape, total_count]));
            end
        end

        function c = getClass(obj)
            c = obj.item_serializer_.getClass;
        end

        function s = getShape(~)
            s = [];
        end
    end
end
