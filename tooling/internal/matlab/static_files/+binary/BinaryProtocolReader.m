% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef BinaryProtocolReader < handle

    properties (Access=protected)
        stream_
    end

    methods
        function obj = BinaryProtocolReader(input, expected_schema)
            obj.stream_ = yardl.binary.CodedInputStream(input);
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
                throw(yardl.binary.Exception("Invalid schema"));
            end
        end
    end

    methods (Access=protected)
        function close_(obj)
            obj.stream_.close();
        end
    end
end

function res = read_fixed_int32(stream)
    res = typecast(stream.read(4), "int32");
end
