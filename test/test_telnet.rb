require 'atto'
include Atto::Test

require_relative '../lib/telnet'
require_relative '../lib/rfc1143' 

include Telnet::Codes

assert { Telnet } 

class TestClient
  attr_reader :buffer
  attr_reader :iac
  attr_reader :out

  def initialize
    @buffer = ''
    @iac    = nil
    @out    = ''
  end



  # Telnet event handlers
  def telnet_event(type, *args)
    puts("Telnet event received by client: #{type}, #{args}")
    meth = "telnet_#{type}".to_sym
    self.send(meth, *args)
  end
  
  def telnet_send_data(zbuf)    
    @out << zbuf
  end
  
  def telnet_environment(fb, vars)
    p fb,vars
  end
  
  def telnet_environment(fb, vars)  
    p fb,vars
  end
  
  
  def telnet_mssp(vars)
    @mssp_vars = vars
  end
  
  def telnet_ttype_is(term)
    @term = term
    p "term #{@term}"
  end
    
  def telnet_ttype_send(term)
    p "term #{term} sent"
  end
  
  
  def telnet_compress(compress)  
    p "compress #{compress} set"
  end  
    
    
  def telnet_subnegotiate(sb_telopt, buffer)
    p "received subnegotiate #{sb_telopt} #{buffer}"
  end  
  
  def telnet_data(data)
    @buffer << data
    p "received data #{data}"
  end  
  
  def telnet_iac(byte)  
    p "received iac #{byte}"
  end 

  def telnet_will(opt)  
    p "received will #{opt}"
  end 
  
  def telnet_do(opt)  
    p "received do #{opt}"
  end 
  
  def telnet_wont(opt)  
    p "received wont #{opt}"
  end 
  
  def telnet_dont(opt)  
    p "received dont #{opt}"
  end 

end




assert do
  cl = TestClient.new()
  tn = Telnet.new(cl)
  tn
end

assert do
  cl = TestClient.new()
  tn = Telnet.new(cl)
  tn.telnet_receive("Hello")
  tn.telnet_receive(" World")
  cl.buffer == "Hello World"
end


assert do
  cl = TestClient.new()
  tn = Telnet.new(cl)
  tn.telnet_receive([TELNET_IAC, TELNET_TELOPT_ECHO].pack('c*'))  
end


assert do
  cl = TestClient.new()
  tn = Telnet.new(cl)
  tn.telnet_send_negotiate(TELNET_DO, TELNET_TELOPT_TTYPE)
  p cl.out
end

assert do
  cl = TestClient.new()
  tn = Telnet.new(cl)
  tn.telnet_receive([TELNET_IAC, TELNET_WILL, TELNET_TELOPT_NAWS].pack('c*'))
  p cl.out
end



