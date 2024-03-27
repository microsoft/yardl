classdef BinaryProtocolWriter < handle

    properties (Access=protected)
        stream_
    end

    methods
        function obj = BinaryProtocolWriter(output, schema)
            obj.stream_ = yardl.binary.CodedOutputStream(output);
            obj.stream_.write_bytes(yardl.binary.MAGIC_BYTES);
            write_fixed_int32(obj.stream_, yardl.binary.CURRENT_BINARY_FORMAT_VERSION);
            s = yardl.binary.StringSerializer();
            s.write(obj.stream_, schema);
        end
    end

    methods (Access=protected)
        function end_stream_(obj)
            obj.stream_.write_byte_no_check(0);
        end

        function close_(obj)
            obj.stream_.close();
        end
    end
end

function write_fixed_int32(stream, value)
    assert(value >= intmin("int32"));
    assert(value <= intmax("int32"));
    value = int32(value);
    stream.write(typecast(value, "uint8"));
end
