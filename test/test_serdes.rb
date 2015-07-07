
require 'atto'
include Atto::Test

require_relative '../lib/sitef' 
require_relative '../lib/serdes' 

class Try
  include Serdes
  
  serdes_reader   :id  
  serdes_accessor :foo
  serdes_reader   :bar
  serdes_reader   :hp, :Integer
  
end

@t = nil

assert do 
  @t = Try.new(:id => 'try79', :foo => 'fooo', :bar => "bar\nbar\bar", :hp => 45) 
  @t
end

assert do
  @t.save_one
end

assert do 
  @tl = Try.load_one('try79')  
  @tl && @tl.id == 'try79' && @tl.foo == @t.foo && @tl.bar == @t.bar && @tl.hp == @t.hp
end


