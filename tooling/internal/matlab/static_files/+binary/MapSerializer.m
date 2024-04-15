% Copyright (c) Microsoft Corporation.
% Licensed under the MIT License.

classdef MapSerializer < yardl.binary.TypeSerializer
    properties
        key_serializer_;
        value_serializer_;
    end

    methods
        function self = MapSerializer(key_serializer, value_serializer)
            self.key_serializer_ = key_serializer;
            self.value_serializer_ = value_serializer;
        end

        function write(self, outstream, value)
            arguments
                self (1,1)
                outstream (1,1) yardl.binary.CodedOutputStream
                value (1,1) dictionary
            end
            count = numEntries(value);
            outstream.write_unsigned_varint(count);
            ks = keys(value);
            vs = values(value);
            for i = 1:count
                self.key_serializer_.write(outstream, ks(i));
                self.value_serializer_.write(outstream, vs(i));
            end
        end

        function res = read(self, instream)
            count = instream.read_unsigned_varint();
            res = dictionary();
            for i = 1:count
                k = self.key_serializer_.read(instream);
                v = self.value_serializer_.read(instream);
                res(k) = v;
            end
        end

        function c = get_class(~)
            c = "dictionary";
        end
    end
end
