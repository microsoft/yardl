classdef DateSerializer < yardl.binary.TypeSerializer
    methods (Static)
        function write(outstream, value)
            assert(isdatetime(value));
            dur = int32(days(value - EPOCH_ORDINAL_DAYS));
            outstream.write_signed_varint(dur);
        end

        function res = read(instream)
            days_since_epoch = instream.read_signed_varint();
            res = EPOCH_ORDINAL_DAYS + days(days_since_epoch);
        end
    end
end

function res = EPOCH_ORDINAL_DAYS
    res = datetime(1970, 1, 1);
end
