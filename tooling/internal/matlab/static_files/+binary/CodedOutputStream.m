% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef CodedOutputStream < handle

    properties (Access=private)
        fid_
        owns_stream_
    end

    methods
        function self = CodedOutputStream(output)
            if isa(output, "string") || isa(output, "char")
                [fileId, errMsg] = fopen(output, "W");
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

        function write(self, value)
            assert(isa(value, "uint8"));
            fwrite(self.fid_, value, "uint8");
        end

        function write_bytes(self, bytes)
            fwrite(self.fid_, bytes, "uint8");
        end

        function write_byte_no_check(self, value)
            assert(isscalar(value));
            assert(all(value >= 0));
            assert(all(value <= intmax("uint8")));

            fwrite(self.fid_, value, "uint8");
        end

        function write_values_directly(self, values, precision)
            fwrite(self.fid_, values, precision);
        end

        function write_unsigned_varint(self, value)
            assert(isscalar(value));

            int_val = uint64(value);
            while true
                if int_val < 0x80
                    self.write_byte_no_check(int_val);
                    return
                end

                self.write_byte_no_check(bitor(bitand(int_val, uint64(0x7F)), uint64(0x80)));
                int_val = bitshift(int_val, -7);
            end
        end

        function res = zigzag_encode(~, value)
            int_val = int64(value);
            res = bitxor(bitshift(int_val, 1), bitshift(int_val, -63));
        end

        function write_signed_varint(self, value)
            assert(isscalar(value));

            self.write_unsigned_varint(self.zigzag_encode(value));
        end

    end
end
