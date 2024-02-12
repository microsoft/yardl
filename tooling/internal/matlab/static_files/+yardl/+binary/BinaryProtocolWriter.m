classdef BinaryProtocolWriter < handle

    properties (Access=protected)
        fid_
        stream_
    end

    methods
        function obj = BinaryProtocolWriter(filename, schema)
            % TODO: What if user wants to append? Should take fid as arg
            [fileId, errMsg] = fopen(filename, "w");
            if fileId < 0
                throw(yardl.Exception(errMsg));
            end

            obj.fid_ = fileId;
            obj.stream_ = yardl.binary.CodedOutputStream(fileId);

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
            if obj.fid_ > 2
                obj.stream_.close();
                fclose(obj.fid_);
                obj.fid_ = -1;
            end
        end
    end
end

function write_fixed_int32(stream, value)
    assert(value >= intmin("int32"));
    assert(value <= intmax("int32"));
    value = int32(value);
    stream.write(typecast(value, "uint8"));
end
