require 'atto'
include Atto::Test

require_relative '../lib/sitef' 

SITEF_NAME = '/tmp/sitef_test.sitef'
SITEF_NAME2 = '/tmp/sitef_test2.sitef'
TEST_DATA = []
TEST_DATA <<  { :nid => 10, :text => "I'm a 10 text man\nAnd I go on several\nlines." }
TEST_DATA <<  { :nid => 27, :text => "--\nA tricky text. 3 newlines:\n\n\nSlashes: \n\\\n:nid:27\n\n.text\n.." }

module Foo
class Try
  attr_reader :foo  
  attr_reader :bar
  
  def initialize (h = {})
    @foo = h['foo']
    @bar = h['bar']
  end
  
  def self.from_sitef(h)
    self.new(h)
  end
  
end
end

TEST_DATA2 = [ Foo::Try.new('foo' => 1, 'bar' => "Yeah"),  
Foo::Try.new('foo' => 1, 'bar' => "Oh")]
  

assert { Sitef } 
assert { Sitef.save_filename(SITEF_NAME, TEST_DATA) }

assert "Loading" do
  res, err = Sitef.load_filename(SITEF_NAME)
  p res
  res && err.empty? && res[0]["nid"] == "10" && res[1]["nid"] == "27"
end


assert "Text round trip" do
  res, err = Sitef.load_filename(SITEF_NAME)
  to = TEST_DATA[1][:text]
  tl = res[1]["text"]
  p to, tl
  to == tl
end

assert do
  Sitef.save_filename(SITEF_NAME2, TEST_DATA2, :foo, :bar) 
end


assert do
  res, err = Sitef.load_objects(SITEF_NAME2)
  p res
  res
end
