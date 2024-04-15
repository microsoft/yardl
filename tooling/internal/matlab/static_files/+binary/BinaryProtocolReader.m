% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef BinaryProtocolReader < handle

    properties (Access=protected)
        stream_
    end

    methods
        function self = BinaryProtocolReader(infile, expected_schema)
            self.stream_ = yardl.binary.CodedInputStream(infile);
            magic_bytes = self.stream_.read_bytes(length(yardl.binary.MAGIC_BYTES));
            if magic_bytes ~= yardl.binary.MAGIC_BYTES
                throw(yardl.binary.Exception("Invalid magic bytes"));
            end

            version = read_fixed_int32(self.stream_);
            if version ~= yardl.binary.CURRENT_BINARY_FORMAT_VERSION
                throw(yardl.binary.Exception("Invalid binary format version"));
            end

            s = yardl.binary.StringSerializer();
            schema = s.read(self.stream_);
            if ~isempty(expected_schema) & schema ~= expected_schema
                throw(yardl.binary.Exception("Invalid schema"));
            end
        end
    end

    methods (Access=protected)
        function close_(self)
            self.stream_.close();
        end
    end
end

function res = read_fixed_int32(stream)
    res = typecast(stream.read_bytes(4), "int32");
end
