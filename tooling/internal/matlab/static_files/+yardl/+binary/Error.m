% function err = Exception(msg, A)
%     identifier = "yardl:binary:error";
%     if nargin < 2
%         err = MException(identifier, msg);
%     else
%         err = MException(identifier, msg, A);
%     end
% end

function err = Error(varargin)
    err = yardl.Exception("yardl:binary:Error", varargin{:});
end
