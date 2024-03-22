classdef CodedInputStream < handle

    properties (Access=private)
        fid_
        owns_stream_
    end

    methods
        function self = CodedInputStream(input)
            if isa(input, "string") || isa(input, "char")
                [fileId, errMsg] = fopen(input, "r");
                if fileId < 0
                    throw(yardl.binary.Exception(errMsg));
                end
                self.fid_ = fileId;
                self.owns_stream_ = true;
            else
                self.fid_ = input;
                self.owns_stream_ = false;
            end
        end

        function close(self)
            if self.owns_stream_ && self.fid_ > 2
                fclose(self.fid_);
                self.fid_ = -1;
            end
        end

        % In Python, this uses struct packing for any selfect...
        function res = read(self, count)
            res = fread(self.fid_, count, "*uint8");
        end

        function res = read_byte(self)
            res = fread(self.fid_, 1, "*uint8");
        end

        function res = read_unsigned_varint(self)
            res = uint64(0);
            shift = uint8(0);

            while true
                byte = self.read_byte();
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

        function res = read_signed_varint(self)
            res = self.zigzag_decode(self.read_unsigned_varint());
        end

    end
end
