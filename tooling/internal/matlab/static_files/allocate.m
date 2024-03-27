function res = allocate(classname, varargin)
    if classname == "string"
        res = strings(varargin{:});
    else
        res = zeros(varargin{:}, classname);
    end
end
