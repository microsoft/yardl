% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef CodedOutputStream < handle

    properties (Access=private)
        fid_
        owns_stream_
    end

    methods
        function self = CodedOutputStream(outfile)
            if isa(outfile, "string") || isa(outfile, "char")
                [fileId, errMsg] = fopen(outfile, "W");
                if fileId < 0
                    throw(yardl.binary.Exception(errMsg));
                end
                self.fid_ = fileId;
                self.owns_stream_ = true;
            else
                self.fid_ = outfile;
                self.owns_stream_ = false;
            end
        end

        function close(self)
            if self.owns_stream_ && self.fid_ > 2
                fclose(self.fid_);
                self.fid_ = -1;
            end
        end

        function write_bytes(self, bytes)
            arguments
                self
                bytes (1,:) {mustBeA(bytes, "uint8")}
            end
            fwrite(self.fid_, bytes, "uint8");
        end

        function write_byte(self, value)
            arguments
                self
                value (1,1) {mustBeA(value, "uint8")}
            end
            fwrite(self.fid_, value, "uint8");
        end

        function write_values_directly(self, values, precision)
            fwrite(self.fid_, values, precision);
        end

        function write_unsigned_varint(self, value)
            arguments
                self
                value (1,1) {mustBeInteger,mustBeNonnegative}
            end

            int_val = uint64(value);
            while true
                if int_val < 0x80
                    self.write_byte(uint8(int_val));
                    return
                end

                self.write_byte(uint8(bitor(bitand(int_val, uint64(0x7F)), uint64(0x80))));
                int_val = bitshift(int_val, -7);
            end
        end

        function res = zigzag_encode(~, value)
            int_val = int64(value);
            res = bitxor(bitshift(int_val, 1), bitshift(int_val, -63));
        end

        function write_signed_varint(self, value)
            arguments
                self
                value (1,1) {mustBeInteger}
            end
            self.write_unsigned_varint(self.zigzag_encode(value));
        end
    end
end
