classdef FixedVectorSerializer < yardl.binary.TypeSerializer
    properties
        item_serializer_;
        length_;
    end

    methods
        function obj = FixedVectorSerializer(item_serializer, length)
            obj.item_serializer_ = item_serializer;
            obj.length_ = length;
        end

        function write(obj, outstream, values)
            s = size(values);
            count = s(end);

            if count ~= obj.length_
                throw(yardl.ValueError("Expected an array of length %d, got %d", obj.length_, count));
            end

            if iscell(values)
                assert(s(1) == 1);
                for i = 1:count
                    obj.item_serializer_.write(outstream, values{i});
                end
            else
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
            item_shape = obj.item_serializer_.getShape();
            if isempty(item_shape)
                res = cell(1, obj.length_);
                for i = 1:obj.length_
                    res{i} = obj.item_serializer_.read(instream);
                end
            elseif isscalar(item_shape)
                % res = zeros(obj.getShape(), obj.getClass());
                res = yardl.allocate(obj.getClass(), obj.getShape());
                for i = 1:obj.length_
                    res(i) = obj.item_serializer_.read(instream);
                end
            else
                % res = zeros([prod(item_shape), obj.length_], obj.getClass());
                res = yardl.allocate(obj.getClass(), [prod(item_shape), obj.length_]);
                for i = 1:obj.length_
                    item = obj.item_serializer_.read(instream);
                    res(:, i) = item(:);
                end
                res = squeeze(reshape(res, [item_shape, obj.length_]));
            end
        end

        function c = getClass(obj)
            c = obj.item_serializer_.getClass();
        end

        function s = getShape(obj)
            item_shape = obj.item_serializer_.getShape();
            if isempty(item_shape)
                s = [1, obj.length_];
            elseif isscalar(item_shape)
                s = [item_shape obj.length_];
            else
                s = [item_shape obj.length_];
                s = s(s>1);
            end
        end
    end
end
