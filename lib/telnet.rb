require 'zlib'
require_relative 'telnet/codes'
require_relative 'rfc1143'
require_relative 'monolog'



# This Telnet class implements a subset of the Telnet protocol.
#
class Telnet
  include Monolog
  include Telnet::Codes

  # Allowed telnet state codes
  STATES = [:data, :iac, :will, :wont, :do, :dont, :sb, :sb_data, :sb_data_iac]
  
  # Helper structs
  Telopt = Struct.new(:telopt, :us, :him)

  
  
  def initialize(client)
    @client     = client
    @telopts    = {}    # Telopt support.
    @rfc1143    = {}    # RFC1143 support.
    @buffer     = ""    # Subrequest buffer
    @state      = :data # state of telnet protocol parser.
    @sb_telopt  = nil;  # current subnegotiation
    @compress   = false # compression state
  end
  
  # Wait for input from the server 
  def wait_for_input
    return Fiber.yield
  end
  
  # Called when client data should be filtered before being passed to the server
  def client_to_server(data)
    result = ""
    data.each_byte do | b |
    iac    = TELNET_IAC.chr  
      
      case @buffer
        when /\A#{iac}#{iac}\Z/
        
        # ongoing negotiation
        when /\A#{iac}\Z/
          return nil
        else
          
        
      end
    end
  end
  
  # Sends unescaped data to client, possibly compressing it if needed
  def send_raw(buf)
    if @compress
      zbuf = Zlib.deflate(buf)
    else
      zbuf = buf
    end
    @client.send_data(zbuf)
  end
  
  # Send data to client (escapes IAC bytes) 
  def send_escaped(buf)
    iac = TELNET_IAC.chr
    self.send_raw(buf.gsub("#{iac}", "#{iac}#{iac}")
  end
  
  # Send negotiation bytes
  
/* send negotiation bytes */
  def send_negotiate(cmd, telopt)
    bytes = ""
    bytes << TELNET_IAC
    bytes << cmd
    bytes << telopt
    send_raw(bytes)
  end
  
  # Check if we support a particular telopt;
  def us_support(telopt)
    have = @telopts[telopt] 
    return false unless have
    return (have.telopt == telopt) && have.us 
  end
  
  # Check if the remote supports a telopt
  def him_support(telopt)
    have = @telopts[telopt] 
    return false unless have
    return (have.telopt == telopt) && have.him 
  end
 
  
  # retrieve RFC1143 option state
  def rfc1143_get(telopt)
    @rfc1143[telopt] 
  end
  
  
  # save RFC1143 option state
  def rfc1143_set(telopt, us, him)
    agree = we_support(telopt)
    @rfc1143[telopt] = RFC1143.new(telopt, us, him, agree)
    return @rfc1143[telopt]
  end
  
  
  # RFC1143 telnet option negotiation helper
  def rfc1143_negotiate_will(rfc1143)
    return rfc1143.handle_will
  end
  
  # RFC1143 telnet option negotiation helper
  def rfc1143_negotiate_wont(rfc1143)
     return rfc1143.handle_wont
  end
  
  # RFC1143 telnet option negotiation helper
  def rfc1143_negotiate_do(rfc1143)
    return rfc1143.handle_do
  end
  
  # RFC1143 telnet option negotiation helper
  def rfc1143_negotiate_dont(rfc1143)
    return rfc1143.handle_dont
  end
    
  # RFC1143 telnet option negotiation helper
  def rfc1143_negotiate(telopt)
    q = rfc1143_get(telopt)
    return nil, nil unless q
    
    case @state
    when :will
      return rfc1143_negotiate_will(q)    
    when :wont
      return rfc1143_negotiate_wont(q)    
    when :do
      return rfc1143_negotiate_do(q)    
    when :dont
      return rfc1143_negotiate_dont(q)    
    end  
  end
  
  def do_negotiate(telopt)
    res, arg = rfc1143_negotiate(telopt)
    return unless res
    if res == :error
      log_error(arg)
    else
      send_negotiate(res, arg)
    end
  end
  
  
  
  # Called when server data should be filtered before being passed to the client
  def server_to_client(data)
  
  end
end
