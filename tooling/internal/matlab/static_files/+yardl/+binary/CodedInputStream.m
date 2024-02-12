classdef CodedInputStream < handle

    properties
        fileId
    end

    methods
        function obj = CodedInputStream(fileId)
            obj.fileId = fileId;
        end

        function close(obj)
            % flush...
            obj.fileId = -1;
        end

        % In Python, this uses struct packing for any object...
        function res = read(obj, count)
            res = fread(obj.fileId, count, "*uint8");
        end

        function res = read_byte(obj)
            res = fread(obj.fileId, 1, "*uint8");
        end

        function res = read_unsigned_varint(obj)
            res = uint64(0);
            shift = uint8(0);

            while true
                byte = obj.read_byte();
                res = bitor(res, bitshift(uint64(bitand(byte, 0x7F)), shift));
                if byte < 0x80
                    return
                end
                shift = shift + 7;
            end
        end

        function res = zigzag_decode(~, value)
            value = uint64(value);
            % res = int64(bitxor(bitshift(value, -1), -bitand(value, 1)));
            res = bitxor(int64(bitshift(value, -1)), -int64(bitand(value, 1)));
        end

        function res = read_signed_varint(obj)
            res = obj.zigzag_decode(obj.read_unsigned_varint());
        end

    end
end
