% function err = Exception(id, msg, A)
%     if nargin < 3
%         err = MException(id, msg);
%     else
%         err = MException(id, msg, A);
%     end
% end

function err = Exception(id, varargin)
    err = MException(id, varargin{:});
end
