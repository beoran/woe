require 'atto'
include Atto::Test

require_relative '../lib/rfc1143' 

assert { RFC1143 }

sm = RFC1143.new(:echo, :no, :no, true)

assert { sm }
assert { sm.telopt  == :echo }
assert { sm.us      == :no }
assert { sm.him     == :no }
assert { sm.agree   == true }

assert do 
  sm = RFC1143.new(:echo, :no, :no, true)
  res, arg = sm.handle_will
  res == :send_do
  arg == :echo
end

assert do 
  sm = RFC1143.new(:echo, :no, :no, false)
  res, arg = sm.handle_will
  res == :send_dont
  arg == :echo
end

