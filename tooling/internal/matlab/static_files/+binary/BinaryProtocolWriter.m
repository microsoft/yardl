% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef BinaryProtocolWriter < handle

    properties (Access=protected)
        stream_
    end

    methods
        function self = BinaryProtocolWriter(outfile, schema)
            self.stream_ = yardl.binary.CodedOutputStream(outfile);
            self.stream_.write_bytes(yardl.binary.MAGIC_BYTES);
            write_fixed_int32(self.stream_, yardl.binary.CURRENT_BINARY_FORMAT_VERSION);
            s = yardl.binary.StringSerializer();
            s.write(self.stream_, schema);
        end
    end

    methods (Access=protected)
        function end_stream_(self)
            self.stream_.write_byte(uint8(0));
        end

        function close_(self)
            self.stream_.close();
        end
    end
end

function write_fixed_int32(stream, value)
    arguments
        stream (1,1) yardl.binary.CodedOutputStream
        value (1,1) {mustBeA(value, "int32")}
    end
    stream.write_bytes(typecast(value, "uint8"));
end
