classdef CodedOutputStream < handle

    properties
        fid
    end

    methods
        function obj = CodedOutputStream(fileId)
            obj.fid = fileId;
        end

        function close(obj)
            % flush...
            obj.fid = -1;
        end

        function write(obj, value)
            assert(isa(value, "uint8"));
            fwrite(obj.fid, value, "uint8");
        end

        function write_bytes(obj, bytes)
            fwrite(obj.fid, bytes, "uint8");
        end

        function write_byte_no_check(obj, value)
            assert(isscalar(value));
            assert(all(value >= 0));
            assert(all(value <= intmax("uint8")));

            fwrite(obj.fid, value, "uint8");
        end

        function write_unsigned_varint(obj, value)
            assert(isscalar(value));

            int_val = uint64(value);
            while true
                if int_val < 0x80
                    obj.write_byte_no_check(int_val);
                    return
                end

                obj.write_byte_no_check(bitor(bitand(int_val, uint64(0x7F)), uint64(0x80)));
                int_val = bitshift(int_val, -7);
            end
        end

        function res = zigzag_encode(~, value)
            int_val = int64(value);
            res = bitxor(bitshift(int_val, 1), bitshift(int_val, -63));
        end

        function write_signed_varint(obj, value)
            assert(isscalar(value));

            obj.write_unsigned_varint(obj.zigzag_encode(value));
        end

    end
end
