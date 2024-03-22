% classdef None < handle
%     methods
%         function eq = eq(~, other)
%             eq = isa(other, 'yardl.None');
%         end

%         function ne = ne(obj, other)
%             ne = ~eq(obj, other);
%         end
%     end
% end

% classdef None < yardl.Optional
%     methods
%         function obj = None()
%             obj.value_ = NaN;
%             obj.has_value_ = false;
%         end

%         % function eq = eq(~, other)
%         %     eq = isa(other, 'yardl.None');
%         % end

%         % function ne = ne(obj, other)
%         %     ne = ~eq(obj, other);
%         % end
%     end
% end

function n = None()
    n = yardl.Optional();
end
