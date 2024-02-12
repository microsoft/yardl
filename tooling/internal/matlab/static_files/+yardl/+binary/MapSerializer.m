classdef MapSerializer < yardl.binary.TypeSerializer
    properties
        key_serializer_;
        value_serializer_;
    end

    methods
        function obj = MapSerializer(key_serializer, value_serializer)
            obj.key_serializer_ = key_serializer;
            obj.value_serializer_ = value_serializer;
        end

        function write(obj, outstream, value)
            % assert(isa(value, 'containers.Map'))
            % OR, starting in R2022, Mathworks recommends using `dictionary`
            % assert(isa(value, 'dictionary'))

            outstream.write_unsigned_varint(length(value));
            ks = keys(value);
            vs = values(value);
            for i = 1:length(value)
                obj.key_serializer_.write(outstream, ks{i});
                obj.value_serializer_.write(outstream, vs{i});
            end
        end

        function res = read(obj, instream)
            count = instream.read_unsigned_varint();
            % TODO: If we can require R2022, should use `dictionary`
            res = containers.Map;
            for i = 1:count
                k = obj.key_serializer_.read(instream);
                v = obj.value_serializer_.read(instream);
                res(k) = v;
            end
        end
    end
end
