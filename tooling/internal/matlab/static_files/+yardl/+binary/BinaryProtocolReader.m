classdef BinaryProtocolReader < handle

    properties (Access=protected)
        fid_
        stream_
    end

    methods
        function obj = BinaryProtocolReader(filename, expected_schema)
            [fileId, errMsg] = fopen(filename, "r");
            if fileId < 0
                throw(yardl.binary.Exception(errMsg));
            end

            obj.fid_ = fileId;
            obj.stream_ = yardl.binary.CodedInputStream(fileId);

            magic_bytes = obj.stream_.read(length(yardl.binary.MAGIC_BYTES));
            if magic_bytes ~= yardl.binary.MAGIC_BYTES
                throw(yardl.binary.Exception("Invalid magic bytes"));
            end

            version = read_fixed_int32(obj.stream_);
            if version ~= yardl.binary.CURRENT_BINARY_FORMAT_VERSION
                throw(yardl.binary.Exception("Invalid binary format version"));
            end

            s = yardl.binary.StringSerializer();
            schema = s.read(obj.stream_);
            if ~isempty(expected_schema) & schema ~= expected_schema
                fprintf("Expected schema: %s\n", expected_schema);
                fprintf("Actual schema:   %s\n", schema);
                throw(yardl.binary.Exception("Invalid schema"));
            end
        end
    end

    methods (Access=protected)
        function close_(obj)
            if obj.fid_ > 2
                obj.stream_.close();
                fclose(obj.fid_);
                obj.fid_ = -1;
            end
        end
    end
end

function res = read_fixed_int32(stream)
    res = typecast(stream.read(4), "int32");
end
