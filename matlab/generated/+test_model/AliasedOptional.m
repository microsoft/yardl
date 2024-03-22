function o = AliasedOptional(value) 
  assert(isa(value, 'int32'));
  o = yardl.Optional(value);
end
