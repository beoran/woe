package telnet

import "bytes"
import "compression/zlib"
import "github.com/beoran/woe/monolog"

// This Telnet struct implements a subset of the Telnet protocol.

// Telnet states
type TelnetState int

const (
    data_state TelnetState  = iota,
    iac_state               = iota,
    will_state              = iota,
    wont_state              = iota,
    do_state                = iota,
    dont_state              = iota
    sb_state                = iota,
    sb_data_state           = iota,
    sb_data_iac_state       = iota
)


// Telnet event types

type Event interface {
    isEvent()
}

type DataEvent struct {
    Data [] byte
}

func (me DataEvent) isEvent() {}

type NAWSEvent struct {
    W   int
    H   int
}

func (me NAWSEvent) isEvent() {}


type TTypeEvent struct {
    Telopt  byte
    Name    string
}

func (me TTypeEvent) isEvent() {}


type SubnegotioateEvent struct {
    Telopt    byte
    Buffer [] byte
}

func (me SubnegotiateEvent) isEvent() {}


type IACEvent struct {
    Telopt    byte
}

func (me IACEvent) isEvent() {}


type CompressionEvent struct {
    Compress  bool
}

func (me CompressionEvent) isEvent() {}


type EnvironmentEvent struct {
    Telopt    byte
    Vars      map[string] string
}

func (me EnvironmentEvent) isEvent() {}


type MSSPEvent struct {
    Telopt    byte
    Vars      map[string] string
}

func (me MSSPEvent) isEvent() {}


type ZMPEvent struct {
    Telopt    byte
    Vars      map[string] string
}

func (me ZMPEvent) isEvent() {}


type EventChannel chan[Event]


type Telopt struct {
    telopt byte
    us     byte
    him    byte
}
    
type Telnet struct { 
  Events            EventChannel  
  telopts map[byte] Telopt 
  state             TelnetState 
  compress          bool
  zwriter           Writer
  zreader           Reader
  buffer          []byte
  sb_telopt         byte
}

func New() telnet * Telnet {
    events     := make(EventChannel, 64)
    telopts    := make (map[byte] Telopt)
    state      := data_state
    compress   := false
    zwriter    := nil
    zreader    := nil
    buffer     := make([]byte, 1024, 1024)
    sb_telopt  := 0
    telnet      = &Telnet { events, telopts, state, compress, 
        zwriter, zreader, buffer, sb_telopt
    }
    return telnet
}

// Starts compresssion
func (me * Telnet) StartCompression() {
    // var zwbuf  bytes.Buffer
    // me.zwriter = zlib.NewWriter(&zwbuf);
}
  
// Closes the telnet connection, send last compressed data if needed.
func (me * Telnet) Close() { 
    if me.compress  {
        me.zwriter.Close()
        me.zreader.Close()    
    }
}

// Filters raw text, only compressing it if needed. 
func (me * Telnet) FilterRaw(in []byte, out chan []byte) {
    // XXX Handle compression here later
    out <- in
} 

// Filters text, escaping IAC bytes. 
func (me * Telnet) FilterRaw(in []byte, out chan []byte) {
    buffer := make([]byte, len(in) * 2, len(in) * 2) 
    outdex := 0
    /* Double IAC characters to escape them. */
    for index := 0; index < len(in) ; index++ {
        now := in[index]
        if now == TELNET_IAC {
            buffer[outdex] = TELNET_IAC; 
            outdex++;    
        }
        buffer[outdex] = TELNET_IAC;
        outdex++;
    }
    out <- buffer
} 

// Send negotiation bytes
func (me * Telnet) SendNegotiate(cmd byte, telopt byte, out chan []byte) {
    buffer      := make([]byte, 3)
    buffer[0]    = TELNET_IAC
    buffer[1]    = cmd
    buffer[2]    = telopt
    me.FilterRaw(buffer, out)
}    
   
// Parse a subnegotiation buffer for a naws event
func (me * Telnet) SubnegotiateNAWS(buffer []byte,  )
    // Some clients, like Gnome-Mud can't even get this right. Grrr!
    // XXx continue here
    if buffer.nil? || buffer.empty? || buffer.size != 4
      monolog.Info("Bad NAWS negotiation: #{buffer}")
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
      me.type   = type
      me.value  = value
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
  case me.sb_telopt
  when TELNET_TELOPT_COMPRESS2
    # received COMPRESS2 begin marker, setup our zlib box and
    # start handling the compressed stream if it's not already.
    me.compress = true
    send_event(:compress, me.compress)
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
    send_event(:subnegotiate, me.sb_telopt, buffer)
    return false
  end
end


  
  def process_byte(byte) 
    # p "process_byte, #{me.state} #{byte}"
    case me.state
    # regular data
    when :data
      if byte == TELNET_IAC
        # receive buffered bytes as data and go to IAC state if it's notempty
        send_event(:data, me.buffer) unless me.buffer.empty?
        me.buffer = ""
        me.state = :iac
      else
        me.buffer << byte
      end
    # IAC received before
    when :iac
      case byte
      # subnegotiation
      when TELNET_SB
        me.state = :sb
      # negotiation commands
      when TELNET_WILL
        me.state = :will
      when TELNET_WONT
        me.state = :wont
      when TELNET_DO
        me.state = :do
      when TELNET_DONT
        me.state = :dont
      # IAC escaping 
      when TELNET_IAC
        me.buffer << TELNET_IAC.chr
        send_event(:data, me.buffer) unless me.buffer.empty?
        me.buffer = ""
        me.state = :data
      # some other command
      else
        send_event(:iac, byte)
        me.state = :data
      end

    # negotiation received before
    when :will, :wont, :do, :dont
      do_negotiate(byte)
      me.state = :data
    # subnegotiation started, determine option to subnegotiate
    when :sb
      me.sb_telopt = byte
      me.state     = :sb_data
    # subnegotiation data, buffer bytes until the end request 
    when :sb_data
      # IAC command in subnegotiation -- either IAC SE or IAC IAC
      if (byte == TELNET_IAC)
        me.state = :sb_data_iac
      elsif (me.sb_telopt == TELNET_TELOPT_COMPRESS && byte == TELNET_WILL)
        # MCCPv1 defined an invalid subnegotiation sequence (IAC SB 85 WILL SE) 
        # to start compression. Catch and discard this case, only support 
        # MMCPv2.
        me.state = data
      else 
        me.buffer << byte
      end

    # IAC received inside a subnegotiation
    when :sb_data_iac
      case byte
        # end subnegotiation
        when TELNET_SE
          me.state = :data
          # process subnegotiation
          compress = do_subnegotiate(me.buffer)
          # if compression was negotiated, the rest of the stream is compressed
          # and processing it requires decompressing it. Return true to signal 
          # this.
          me.buffer = ""
          return true if compress
        # escaped IAC byte
        when TELNET_IAC
        # push IAC into buffer */
          me.buffer << byte
          me.state = :sb_data
        # something else -- protocol error.  attempt to process
        # content in subnegotiation buffer, then evaluate the
        # given command as an IAC code.
        else
          log_error("Unexpected byte after IAC inside SB: %d", byte)
          me.state = :iac
          # subnegotiate with the buffer anyway, even though it's an error
          compress = do_subnegotiate(me.buffer)
          # if compression was negotiated, the rest of the stream is compressed
          # and processing it requires decompressing it. Return true to signal 
          # this.
          me.buffer = ""
          return true if compress
        end
    when :data  
      # buffer any other bytes
      me.buffer << byte
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
    send_event(:data, me.buffer) unless me.buffer.empty?
    me.buffer = ""
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
    me.compress = true
  end
  
  # send formatted data
  def telnet_raw_printf(fmt, *args)
    buf   = sprintf(fmt, *args)
    telnet_send(buf)
  end

  CRLF  = "\r\n"
  CRNUL = "\r\0"
  
  # send formatted data with \r and \n translation in addition to IAC IAC 
  def telnet_printf(fmt, *args)
    buf   = sprintf(fmt, *args)
    buf.gsub!("\r", CRNUL)
    buf.gsub!("\n", CRLF)
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



