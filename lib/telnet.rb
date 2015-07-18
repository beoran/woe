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

  attr_reader :telopts
  
  def initialize(client)
    @client     = client
    @telopts    = {}    # Telopt support.
    @rfc1143    = {}    # RFC1143 support.
    @buffer     = ""    # Subrequest buffer
    @state      = :data # state of telnet protocol parser.
    @sb_telopt  = nil;  # current subnegotiation
    @compress   = false # compression state
    @zdeflate   = Zlib::Deflate.new() # Deflate stream for compression2 support.
    @zinflate   = Zlib::Inflate.new() # Inflate stream for compression2 support.
  end
  
  # Closes the telnet connection, send last compressed data if needed.
  def close
    if @compress 
      zbuf = @zdeflate.flush(Zlib::FINISH)
      @client.telnet_send_data(zbuf)
    end
    @zdeflate.close
    @zinflate.close    
  end
  
  # Send an event to the client to notify it of a state change or of data
  def send_event(type, *data)
    @client.telnet_event(type, *data)
  end
  
  # Sends unescaped data to client, possibly compressing it if needed
  def send_raw(buf)
    if @compress
      @zdeflate << buf
      # for short messages the "compressed" stream wil actually be 
      # bigger than the uncompressed one, but that's unavoidable
      # due to the streaming nature of network connections.
      zbuf = @zdeflate.flush(Zlib::SYNC_FLUSH)
    else
      zbuf = buf
    end
    # Don't use send_event here, since that's only for events received
    @client.telnet_send_data(zbuf)
  end
  
  # Send data to client (escapes IAC bytes) 
  def send_escaped(buf)
    iac = TELNET_IAC.chr
    self.send_raw(buf.gsub("#{iac}", "#{iac}#{iac}"))
  end
  
  # Send negotiation bytes
  
  # negotiation bytes 
  def send_negotiate(cmd, telopt)
    bytes = ""
    bytes << TELNET_IAC
    bytes << cmd
    bytes << telopt
    send_raw(bytes)
  end  
  
  # 
  
  # Check if we support a particular telsopt using the RFC1143 state
  def us_support(telopt)
    have = @rfc1143[telopt]
    return false unless have
    return (have.telopt == telopt) && have.us == :yes 
  end
  
  # Check if the remote supports a telopt (and it is enabled)
  def him_support(telopt)
    have = @rfc1143[telopt]
    return false unless have
    return (have.telopt == telopt) && have.him == :yes 
  end
  
  # Set that we support an option (using the RFC1143 state)
  def set_support(telopt, support=true, us = :no, him = :no)
    rfc1143_set(telopt, support=true, us = :no, him = :no)
  end
   
  # retrieve RFC1143 option state
  def rfc1143_get(telopt)
    @rfc1143[telopt]
  end
    
  # save RFC1143 option state
  def rfc1143_set(telopt, support=true, us = :no, him = :no)
    agree = support
    @rfc1143[telopt] = RFC1143.new(telopt, us, him, agree)
    return @rfc1143[telopt]
  end
  
  
  # RFC1143 telnet option negotiation helper
  def rfc1143_negotiate(telopt)
    q = rfc1143_get(telopt)
    return nil, nil unless q
    
    case @state
    when :will
      return q.handle_will 
    when :wont
      return q.handle_wont
    when :do
      return q.handle_do
    when :dont
      return q.handle_dont
    end  
  end
  
  # Performs a telnet negotiation
  def do_negotiate(telopt)
    res, arg = rfc1143_negotiate(telopt)
    send_event(@state, telopt, res, arg)
  end
  
  
  # Process a subnegotiation buffer for a naws event
  def subnegotiate_naws(buffer)
    # Some clients, like Gnome-Mud can't even get this right. Grrr!
    if buffer.nil? || buffer.empty? || buffer.size != 4
      log_info("Bad NAWS negotiation: #{buffer}")
      return nil
    end
    arr   = buffer.bytes.to_a
    w     = (arr[0] << 8) + arr[1]
    h     = (arr[2] << 8) + arr[3]
    send_event(:naws, w, h)
  end
  

  # Storage for environment values
  class Environment 
    attr_accessor :type
    attr_accessor :value
    
    def initialize(type, value)
      @type   = type
      @value  = value
    end
  end


  # process an ENVIRON/NEW-ENVIRON subnegotiation buffer
  def subnegotiate_environ(buffer)
    vars  = []
    cmd   = ""
    arr   = buffer.bytes.to_a
    fb    = arr.first  
    # first byte must be a valid command 
    if fb != TELNET_ENVIRON_SEND && fb != TELNET_ENVIRON_IS && fb != TELNET_ENVIRON_INFO
      log_error("telopt environment subneg command not valid")
      return 0
    end
    
    cmd << fb    
    
    if (buffer.size == 1) 
      send_event(:environment, fb, vars)
      return false
    end
        
    # Second byte must be VAR or USERVAR, if present
    sb = arr[1]
    if sb != TELNET_ENVIRON_VAR && fb != TELNET_ENVIRON_USEVAR
      log_error("telopt environment subneg missing variable type")
      return false
    end
    
    # ensure last byte is not an escape byte (makes parsing later easier) 
    lb = arr.last
    if lb == TELNET_ENVIRON_ESC
      log_error("telopt environment subneg ends with ESC")
      return false
    end

    var    = nil
    index  = 1
    escape = false
    
    arr.shift
    
    arr.each do | c | 
      case c
      when TELNET_ENVIRON_VAR
      when TELNET_ENVIRON_VALUE
      when TELNET_ENVIRON_USERVAR
        if escape
          escape = false
          var.value << c
        elsif var
          vars << var
          var = Environment.new(c, "")
        else
          var = Environment.new(c, "")        
        end
      when TELNET_ENVIRON_ESC
        escape = true
      else
        var.value << c  
      end # case
    end # each
    
    send_event(:environment, fb, vars)    
    return false
  end



# process an MSSP subnegotiation buffer
def subnegotiate_mssp(buffer)
  telnet_event_t ev;
  struct telnet_environ_t *values;
  char *var = 0;
  char *c, *last, *out;
  size_t i, count;
  unsigned char next_type;
  
  if buffer.size < 1
    return 0
  end
  
  arr   = buffer.bytes.to_a
  fb    = arr.first  
  # first byte must be a valid command
  if fb != TELNET_MSSSP_VAR
    log_error("telopt MSSP subneg data not valid")
    return false
  end
  
  vars    = {}
  var     = ""
  val     = ""
  mstate  = :var
  while index <  arr.size
    c     = arr[index]
    case c
    when TELNET_MSSP_VAR
      mstate = :var
      if mstate == :val
        vars[var] = val
        var = ""
        val = ""
      end      
    when TELNET_MSSP_VAL
      mstate = :val
    else
      if mstate == :var
        var << c  
      elsif mstate == :val
        val << c  
      end      
    end # case
    index += 1
  end # while
  
  send_event(:mssp, vars)
  return false
end


# parse ZMP command subnegotiation buffers 
def subnegotiate_zmp(buffer)
  args = []
  arg  = ""
  
  buffer.each_byte do |b|  
    if b == 0
      args << arg
      arg = ""
    else
      arg << byte
    end
  end
  send_event(:zmp, vars)
  return false
end

# parse TERMINAL-TYPE command subnegotiation buffers
def subnegotiate_ttype(buffer)
  # make sure request is not empty
  if buffer.size == 0
    log_error("Incomplete TERMINAL-TYPE request");
    return 0
  end
  
  arr   = buffer.bytes
  fb    = arr.first
  term  = nil 
  
  if fb == TELNET_TTYPE_IS
    term = buffer[1, buffer.size]
    send_event(:ttype_is, term)
  elsif fb == TELNET_TTYPE_SEND
    term = buffer[1, buffer.size]
    send_event(:ttype_send, term)
  else
    log_error("TERMINAL-TYPE request has invalid type")
    return false
  end
  return false
end


# process a subnegotiation buffer; returns true if the current buffer
# must be aborted and reprocessed due to COMPRESS2 being activated

def do_subnegotiate(buffer)
  case @sb_telopt
  when TELNET_TELOPT_COMPRESS2
    # received COMPRESS2 begin marker, setup our zlib box and
    # start handling the compressed stream if it's not already.
    @compress = true
    send_event(:compress, @compress)
    return true
  # specially handled subnegotiation telopt types
  when TELNET_TELOPT_ZMP
    return subnegotiate_zmp(buffer)
  when TELNET_TELOPT_TTYPE
    return subnegotiate_ttype(buffer)
  when TELNET_TELOPT_ENVIRON  
    return subnegotiate_environ(buffer)
  when TELNET_TELOPT_NEW_ENVIRON
    return subnegotiate_environ(buffer)
  when TELNET_TELOPT_MSSP
    return subnegotiate_mssp(buffer)
  when TELNET_TELOPT_NAWS
    return subnegotiate_naws(buffer)
  else
    send_event(:subnegotiate, @sb_telopt, buffer)
    return false
  end
end


  
  def process_byte(byte) 
    # p "process_byte, #{@state} #{byte}"
    case @state
    # regular data
    when :data
      if byte == TELNET_IAC
        # receive buffered bytes as data and go to IAC state if it's notempty
        send_event(:data, @buffer) unless @buffer.empty?
        @buffer = ""
        @state = :iac
      else
        @buffer << byte
      end
    # IAC received before
    when :iac
      case byte
      # subnegotiation
      when TELNET_SB
        @state = :sb
      # negotiation commands
      when TELNET_WILL
        @state = :will
      when TELNET_WONT
        @state = :wont
      when TELNET_DO
        @state = :do
      when TELNET_DONT
        @state = :dont
      # IAC escaping 
      when TELNET_IAC
        @buffer << TELNET_IAC.chr
        send_event(:data, @buffer) unless @buffer.empty?
        @buffer = ""
        @state = :data
      # some other command
      else
        send_event(:iac, byte)
        @state = :data
      end

    # negotiation received before
    when :will, :wont, :do, :dont
      do_negotiate(byte)
      @state = :data
    # subnegotiation started, determine option to subnegotiate
    when :sb
      @sb_telopt = byte
      @state     = :sb_data
    # subnegotiation data, buffer bytes until the end request 
    when :sb_data
      # IAC command in subnegotiation -- either IAC SE or IAC IAC
      if (byte == TELNET_IAC)
        @state = :sb_data_iac
      elsif (@sb_telopt == TELNET_TELOPT_COMPRESS && byte == TELNET_WILL)
        # MCCPv1 defined an invalid subnegotiation sequence (IAC SB 85 WILL SE) 
        # to start compression. Catch and discard this case, only support 
        # MMCPv2.
        @state = data
      else 
        @buffer << byte
      end

    # IAC received inside a subnegotiation
    when :sb_data_iac
      case byte
        # end subnegotiation
        when TELNET_SE
          @state = :data
          # process subnegotiation
          compress = do_subnegotiate(@buffer)
          # if compression was negotiated, the rest of the stream is compressed
          # and processing it requires decompressing it. Return true to signal 
          # this.
          @buffer = ""
          return true if compress
        # escaped IAC byte
        when TELNET_IAC
        # push IAC into buffer */
          @buffer << byte
          @state = :sb_data
        # something else -- protocol error.  attempt to process
        # content in subnegotiation buffer, then evaluate the
        # given command as an IAC code.
        else
          log_error("Unexpected byte after IAC inside SB: %d", byte)
          @state = :iac
          # subnegotiate with the buffer anyway, even though it's an error
          compress = do_subnegotiate(@buffer)
          # if compression was negotiated, the rest of the stream is compressed
          # and processing it requires decompressing it. Return true to signal 
          # this.
          @buffer = ""
          return true if compress
        end
    when :data  
      # buffer any other bytes
      @buffer << byte
    else 
      # programing error, shouldn't happen
      raise "Error in telet state machine!"
    end
    # return false to signal compression needn't start
    return false
  end
  
  def process_bytes(bytes)
    # I have a feeling this way of handling strings isn't very efficient.. :p
    arr = bytes.bytes.to_a
    byte = arr.shift
    while byte
      compress = process_byte(byte)
      if compress
        # paper over this for a while... 
        new_bytes = Zlib.inflate(arr.pack('c*')) rescue nil
        if new_bytes
          arr = new_bytes.bytes.to_a
        end
      end
      byte = arr.shift    
    end
    send_event(:data, @buffer) unless @buffer.empty?
    @buffer = ""
  end
  
  # Call this when the server receives data from the client
  def telnet_receive(data)
    # the COMPRESS2 protocol seems to be half-duplex in that only 
    # the server's data stream is compressed (unless maybe if the client
    # is asked to also compress with a DO command ?)
    process_bytes(data)
  end
  
  # Send a bytes array (raw) to the client
  def telnet_send_bytes(*bytes)
    s     = bytes.pack('C*')
    send_raw(s)
  end
  
  # send an iac command 
  def telnet_send_iac(cmd)
    telnet_send_bytes(TELNET_IAC, cmd)
  end

  # send negotiation
  def telnet_send_negotiate(cmd, telopt)
    # get current option states
    q = rfc1143_get(telopt)
    unless q
      rfc1143_set(telopt)
      q = rfc1143_get(telopt)
    end
    
    act, arg = nil, nil
    case cmd
      when TELNET_WILL
        act, arg = q.send_will
      when TELNET_WONT
        act, arg = q.send_wont
      when TELNET_DO
        act, arg = q.send_do
      when TELNET_DONT
        act, arg = q.send_dont    
    end
        
    return false unless act    
    telnet_send_bytes(TELNET_IAC, act, telopt)
  end
        

  # send non-command data (escapes IAC bytes)
  def telnet_send(buffer)
    send_escaped(buffer)
  end
  
  # send subnegotiation header
  def telnet_begin_sb(telopt)
    telnet_send_bytes(TELNET_IAC, TELNET_SB, telopt)
  end

  # send subnegotiation ending
  def telnet_end_sb()
    telnet_send_bytes(TELNET_IAC, TELNET_SE)
  end


  # send complete subnegotiation
  def telnet_subnegotiation(telopt, buffer = nil)
    telnet_send_bytes(TELNET_IAC, TELNET_SB, telopt)
    telnet_send(buffer) if buffer;
    telnet_send_bytes(TELNET_IAC, TELNET_SE)
  end
  
  # start compress2 compression
  def telnet_begin_compress2() 
    telnet_send_bytes(TELNET_IAC, TELNET_SB, TELNET_TELOPT_COMPRESS2, TELNET_IAC, TELNET_SE);
    @compress = true
  end
  
  # send formatted data
  def telnet_raw_printf(fmt, *args)
    buf   = sprintf(fmt, *args)
    telnet_send(buf)
  end


  # send formatted data with \r and \n translation in addition to IAC IAC 
  def telnet_printf(fmt, *args)
    crlf  = "\r\n"
    clnul = "\r\0"
    buf   = sprintf(fmt, *args)
    buf.gsub!("\r", crnul)
    buf.gsub!("\n", crlf)
    telnet_send(buf)
  end

  # begin NEW-ENVIRON subnegotation
  def telnet_begin_newenviron(cmd)
    telnet_begin_sb(TELNET_TELOPT_NEW_ENVIRON)
    telnet_send_bytes(cmd)
  end
  
  # send a NEW-ENVIRON value
  def telnet_newenviron_value(type, value)
    telnet_send_bytes(type)
    telnet_send(string)
  end
  
  # send TERMINAL-TYPE SEND command
  def telnet_ttype_send() 
    telnet_send_bytes(TELNET_IAC, TELNET_SB, TELNET_TELOPT_TTYPE, TELNET_TTYPE_SEND, TELNET_IAC, TELNET_SE)
  end  
  
  # send TERMINAL-TYPE IS command 
  def telnet_ttype_is(ttype)
    telnet_send_bytes(TELNET_IAC, TELNET_SB, TELNET_TELOPT_TTYPE, TELNET_TTYPE_IS)
    telnet_send(ttype)
  end
  
  # send MSSP data
  def telnet_send_mssp(mssp)
    buf = ""
    mssp.each do | key, val| 
      buf << TELNET_MSSP_VAR.chr
      buf << key
      buf << TELNET_MSSP_VAL.chr
      buf << val      
    end
    telnet_subnegotiation(TELNET_TELOPT_MSSP, buf)
  end

end
