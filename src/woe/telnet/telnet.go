package telnet

// import "bytes"
import "io"
import "strings"
import "fmt"
import "compress/zlib"
import "github.com/beoran/woe/monolog"


// This Telnet struct implements a subset of the Telnet protocol.

// Telnet states
type TelnetState int

const (
    data_state TelnetState  = iota
    iac_state               = iota
    will_state              = iota
    wont_state              = iota
    do_state                = iota
    dont_state              = iota
    sb_state                = iota
    sb_data_state           = iota
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


type SubnegotiateEvent struct {
    Telopt    byte
    Buffer [] byte
}

func (me SubnegotiateEvent) isEvent() {}


type IACEvent struct {
    Telopt    byte
}

func (me IACEvent) isEvent() {}


type CompressEvent struct {
    Compress  bool
}

func (me CompressEvent) isEvent() {}


//Storage for environment values
type Environment struct {
    Type byte
    Value string
}


type EnvironmentEvent struct {
    Telopt    byte
    Vars      [] Environment
}

func (me EnvironmentEvent) isEvent() {}


type MSSPEvent struct {
    Telopt    byte
    Vars      map[string] string
}

func (me MSSPEvent) isEvent() {}


type ZMPEvent struct {
    Vars      []string
}

func (me ZMPEvent) isEvent() {}

type WillEvent struct {
    Telopt byte
}

func (me WillEvent) isEvent() {}

type WontEvent struct {
    Telopt byte
}

func (me WontEvent) isEvent() {}


type DoEvent struct {
    Telopt byte
}

func (me DoEvent) isEvent() {}

type DontEvent struct {
    Telopt byte
}

func (me DontEvent) isEvent() {}


// Telnet event type constants
type EventType int

const (
    TELNET_DATA_EVENT           EventType =  iota
    TELNET_NAWS_EVENT           EventType =  iota
    TELNET_TTYPE_EVENT          EventType =  iota
    TELNET_SUBNEGOTIATE_EVENT   EventType =  iota
    TELNET_IAC_EVENT            EventType =  iota
    TELNET_COMPRESS_EVENT       EventType =  iota
    TELNET_ENVIRONMENT_EVENT    EventType =  iota
    TELNET_MSSP_EVENT           EventType =  iota
    TELNET_ZMP_EVENT            EventType =  iota
    TELNET_WILL_EVENT           EventType =  iota
    TELNET_WONT_EVENT           EventType =  iota
    TELNET_DO_EVENT             EventType =  iota
    TELNET_DONT_EVENT           EventType =  iota
    TELNET_UNKNOWN_EVENT        EventType =  iota
)


/* Returns the numerical event type of an event. Useful for direct comparison. */
func EventTypeOf(event Event) EventType {
    switch event.(type) {
        case DataEvent, *DataEvent:
            return TELNET_DATA_EVENT
        case NAWSEvent, *NAWSEvent:
            return TELNET_NAWS_EVENT
        case TTypeEvent, *TTypeEvent:
            return TELNET_TTYPE_EVENT
        case SubnegotiateEvent, *SubnegotiateEvent:
            return TELNET_SUBNEGOTIATE_EVENT
        case IACEvent, *IACEvent:
            return TELNET_IAC_EVENT
        case CompressEvent, *CompressEvent:
            return TELNET_COMPRESS_EVENT
        case EnvironmentEvent, *EnvironmentEvent:
            return TELNET_ENVIRONMENT_EVENT
        case MSSPEvent, *MSSPEvent:
            return TELNET_MSSP_EVENT
        case ZMPEvent, *ZMPEvent:
            return TELNET_ZMP_EVENT
        case WillEvent, *WillEvent:
            return TELNET_WILL_EVENT
        case WontEvent, *WontEvent:
            return TELNET_WONT_EVENT
        case DoEvent, *DoEvent:
            return TELNET_DO_EVENT
        case DontEvent, *DontEvent:
            return TELNET_DONT_EVENT
        default:
            monolog.Error("Unknown event type %T %v", event, event)
            return TELNET_UNKNOWN_EVENT
    }
}

// Returns true if the event is of the given type, or false if not
func IsEventType(event Event, typ EventType) bool {
    return EventTypeOf(event) == typ;
}


type EventChannel chan(Event)


type Telopt struct {
    telopt byte
    us     byte
    him    byte
}
    
type Telnet struct { 
  Events            EventChannel
  ToClient          chan([]byte)
  telopts map[byte] Telopt 
  state             TelnetState 
  compress          bool
  zwriter           zlib.Writer
  zreader           io.ReadCloser
  buffer          []byte
  sb_telopt         byte
}

func New() (telnet * Telnet) {
    
    events     := make(EventChannel, 64)
    toclient   := make(chan([]byte), 64)
    telopts    := make (map[byte] Telopt)
    state      := data_state
    compress   := false
    var zwriter zlib.Writer
    var zreader io.ReadCloser
    var buffer []byte = nil
    sb_telopt  := byte(0)
    telnet      = &Telnet { events, toclient, telopts, state, compress, zwriter, zreader, buffer, sb_telopt }
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
func (me * Telnet) SendRaw(in []byte) {
    // XXX Handle compression here later
    me.ToClient <- in
} 

// Filters text, escaping IAC bytes. 
func (me * Telnet) SendEscaped(in []byte) {
    buffer := make([]byte, len(in) * 2, len(in) * 2) 
    outdex := 0
    /* Double IAC characters to escape them. */
    for index := 0; index < len(in) ; index++ {
        now := in[index]
        if now == TELNET_IAC {
            buffer[outdex] = TELNET_IAC; 
            outdex++;    
        }
        buffer[outdex] = now;
        outdex++;
    }
    me.SendRaw(buffer)
} 

// Send negotiation bytes
func (me * Telnet) SendNegotiate(cmd byte, telopt byte) {
    buffer      := make([]byte, 3)
    buffer[0]    = TELNET_IAC
    buffer[1]    = cmd
    buffer[2]    = telopt
    me.SendRaw(buffer)
}

func (me * Telnet) SendEvent(event Event) {
    me.Events <- event
}
   
// Parse a subnegotiation buffer for a naws event
func (me * Telnet) SubnegotiateNAWS(buffer []byte) {
    // Some clients, like Gnome-Mud can't even get this right. Grrr!
    if buffer == nil || len(buffer) != 4 {
      monolog.Warning("Bad NAWS negotiation: #{buffer}")
      return
    }
    var w int   = (int(buffer[0]) << 8) + int(buffer[1])
    var h int   = (int(buffer[2]) << 8) + int(buffer[3])
    me.SendEvent(&NAWSEvent{w, h})
}

// process an ENVIRON/NEW-ENVIRON subnegotiation buffer
func (me * Telnet) SubnegotiateEnviron(buffer []byte) {
    var vars []Environment
    var cmd []byte
    fb   := buffer[0]  
    // First byte must be a valid command 
    if fb != TELNET_ENVIRON_SEND && fb != TELNET_ENVIRON_IS && fb != TELNET_ENVIRON_INFO {
      monolog.Warning("telopt environment subneg command not valid")
    }
    
    cmd = append(cmd, fb)   
    
    if len(buffer) == 1 { 
      me.SendEvent(&EnvironmentEvent{fb, vars})
      return
    }
        
    // Second byte must be VAR or USERVAR, if present
    sb := buffer[1]
    if sb != TELNET_ENVIRON_VAR && fb != TELNET_ENVIRON_USERVAR {
      monolog.Warning("telopt environment subneg missing variable type")
      return
    }
    
    // ensure last byte is not an escape byte (makes parsing later easier) 
    lb := buffer[len(buffer) - 1]
    if lb == TELNET_ENVIRON_ESC {
      monolog.Warning("telopt environment subneg ends with ESC")
      return
    }

/* XXX : not implemented yet
    var variable * Environment = nil
    index           := 1
    escape          := false
    
    for index := 1 ; index < len(buffer) ; index++ {
      c := buffer[index]  
      switch c {
        case TELNET_ENVIRON_VAR: 
            fallthrough
        case TELNET_ENVIRON_VALUE:
            fallthrough
        case TELNET_ENVIRON_USERVAR:
            if escape {
                escape = false
                variable.Value  = append(variable.Value, c)
            } else if (variable != nil) {
                vars            = append(vars, variable)
                variable        = new(Environment)
                variable.Type   = c
            } else {
                variable        = new(Environment)
                variable.Type   = c
            }
      case TELNET_ENVIRON_ESC:
        escape = true
      default:
        variable.Value = append(variable.Value, c)
      }
    }
    // Finally send event
    me.SendEvent(&EnvironmentEvent{fb, vars})
*/
}


const (
    MSTATE_NONE = 0
    MSTATE_VAR  = 1
    MSTATE_VAL  = 2
)

// process an MSSP subnegotiation buffer
func (me * Telnet) SubnegotiateMSSP(buffer []byte) {
    if len(buffer) < 1 {
        return
    }
  
    fb    := buffer[0]  
    // first byte must be a valid command
    if fb != TELNET_MSSP_VAR {
        monolog.Warning("telopt MSSP subneg data not valid")
        return
    }
  
    variables := make(map[string] string)
    var variable []byte
    var value []byte
    mstate := MSTATE_NONE
    
    for index := 0 ; index <  len(buffer) ; index ++ {
        c     := buffer[index]
        
        switch c {
            case TELNET_MSSP_VAR:
            mstate = MSTATE_VAR
            if mstate == MSTATE_VAR {
                variables[string(variable)] = string(value)
                variable = nil
                value    = nil
            }
            case TELNET_MSSP_VAL:
                mstate = MSTATE_VAL
            default:
                if mstate == MSTATE_VAL {
                    variable = append(variable, c)
                } else {
                    value = append(value, c)
                }  
        }
    }
    me.SendEvent(&MSSPEvent{fb, variables})
}


// Parse ZMP command subnegotiation buffers 
func (me * Telnet) SubnegotiateZMP(buffer []byte) {
  var vars []string
  var variable []byte
  var b byte
  for index := 0 ; index < len(buffer) ; index++ {
      b = buffer[index]
      if b == 0 {
        vars     = append(vars, string(variable))
        variable = nil
      } else {
        variable = append(variable, b)
      }  
  }
  me.SendEvent(&ZMPEvent{vars})
}

// parse TERMINAL-TYPE command subnegotiation buffers
func (me * Telnet) SubnegotiateTType(buffer []byte) {
  // make sure request is not empty
  if len(buffer) == 0 {
    monolog.Warning("Incomplete TERMINAL-TYPE request");
    return 
  }
  
  fb    := buffer[0]
  if fb != TELNET_TTYPE_IS && fb != TELNET_TTYPE_SEND {
    monolog.Warning("TERMINAL-TYPE request has invalid type %d (%v)", fb, buffer)
    return
  }
  
  term := string(buffer[1:])
  me.SendEvent(&TTypeEvent{fb, term})
}


// process a subnegotiation buffer; returns true if the current buffer
// must be aborted and reprocessed due to COMPRESS2 being activated
func (me * Telnet) DoSubnegotiate(buffer []byte) bool {
    switch me.sb_telopt {
        case TELNET_TELOPT_COMPRESS2:
        // received COMPRESS2 begin marker, setup our zlib box and
        // start handling the compressed stream if it's not already.
        me.compress = true
        me.SendEvent(&CompressEvent{me.compress})
        return true
        // specially handled subnegotiation telopt types
        case TELNET_TELOPT_TTYPE:
            me.SubnegotiateTType(buffer)
        case TELNET_TELOPT_ENVIRON:
            me.SubnegotiateEnviron(buffer)
        case TELNET_TELOPT_NEW_ENVIRON:
            me.SubnegotiateEnviron(buffer)
        case TELNET_TELOPT_MSSP:
            me.SubnegotiateMSSP(buffer)
        case TELNET_TELOPT_NAWS:
            me.SubnegotiateNAWS(buffer)
        case TELNET_TELOPT_ZMP:
            me.SubnegotiateZMP(buffer)
        default:    
            // Send catch all subnegotiation event
            me.SendEvent(&SubnegotiateEvent{me.sb_telopt, buffer})
    }
    return false
}

func (me * Telnet) DoNegotiate(state TelnetState, telopt byte) bool {
    switch me.state {
        case will_state:
            me.SendEvent(&WillEvent{telopt})
        case wont_state:
            me.SendEvent(&WontEvent{telopt})
        case do_state:
            me.SendEvent(&DoEvent{telopt})
        case dont_state:
            me.SendEvent(&DontEvent{telopt})
        default:
            monolog.Warning("State not vvalid in  telnet negotiation.")
    }
    me.state = data_state
    return false
}

// Send the current buffer as a DataEvent if it's not empty
// Also empties the buffer if it wasn't emmpty
func (me * Telnet) maybeSendDataEventAndEmptyBuffer() {
    if (me.buffer != nil) && (len(me.buffer) > 0) {
        me.SendEvent(&DataEvent{me.buffer})
        me.buffer = nil
    }
}

// Append a byte to the data buffer
// Also empties the buffer if it wasn't emmpty
func (me * Telnet) appendByte(bin byte) {
    monolog.Debug("Appending to telnet buffer: %d %d", len(me.buffer), cap(me.buffer))
    me.buffer = append(me.buffer, bin)
}

// Process a byte in the data state 
func (me * Telnet) dataStateProcessByte(bin byte) bool {
    if bin == TELNET_IAC {
        // receive buffered bytes as data and go to IAC state if it's notempty
        me.maybeSendDataEventAndEmptyBuffer()
        me.state = iac_state
    } else {
        me.appendByte(bin)
    }
    return false
}

// Process a byte in the IAC state 
func (me * Telnet) iacStateProcessByte(bin byte) bool {
    switch bin {
      // subnegotiation
      case TELNET_SB:
        me.state = sb_state
      // negotiation commands
      case TELNET_WILL:
        me.state = will_state
      case TELNET_WONT:
        me.state = wont_state
      case TELNET_DO:
        me.state = do_state
      case TELNET_DONT:
        me.state = dont_state
      // IAC escaping
      case TELNET_IAC:
        me.appendByte(TELNET_IAC)
        me.maybeSendDataEventAndEmptyBuffer()
        me.state = data_state
      // some other command
      default:
        me.SendEvent(IACEvent { bin })
        me.state = data_state
    }
    return false      
}


// Process a byte in the subnegotiation data state 
func (me * Telnet) sbdataStateProcessByte(bin byte) bool {
    // IAC command in subnegotiation -- either IAC SE or IAC IAC
    if (bin == TELNET_IAC)  {
        me.state = sb_data_iac_state
    } else if me.sb_telopt == TELNET_TELOPT_COMPRESS &&  bin == TELNET_WILL {
        // MCCPv1 defined an invalid subnegotiation sequence (IAC SB 85 WILL SE) 
        // to start compression. Catch and discard this case, only support 
        // MMCPv2.
        me.state = data_state
    } else {
        me.appendByte(bin)
    }
    return false
}

// Process a byte in the IAC received when processing subnegotiation data state 
func (me * Telnet) sbdataiacStateProcessByte(bin byte) bool {
    switch bin { 
        //end subnegotiation
        case TELNET_SE:
        me.state = data_state
        // process subnegotiation
        compress := me.DoSubnegotiate(me.buffer)
        // if compression was negotiated, the rest of the stream is compressed
        // and processing it requires decompressing it. Return true to signal 
        // this.
        me.buffer = nil
        if compress {
            return true 
        }
            
        // escaped IAC byte
        case TELNET_IAC:
        // push IAC into buffer
        me.appendByte(bin)
        me.state = sb_data_state
        // something else -- protocol error.  attempt to process
        // content in subnegotiation buffer, then evaluate the
        // given command as an IAC code.
        default:
        monolog.Warning("Unexpected byte after IAC inside SB: %d", bin)
        me.state = iac_state
        // subnegotiate with the buffer anyway, even though it's an error
        compress := me.DoSubnegotiate(me.buffer)
        // if compression was negotiated, the rest of the stream is compressed
        // and processing it requires decompressing it. Return true to signal 
        // this.
        me.buffer = nil
        if compress {
            return true 
        }
    }
    return false    
}


// Process a single byte received from the client 
func (me * Telnet) ProcessByte(bin byte) bool {
    monolog.Debug("ProcessByte %d %d", bin, me.state)
    switch me.state {
    // regular data
        case data_state:
        return me.dataStateProcessByte(bin)
    // IAC received before
        case iac_state:
        return me.iacStateProcessByte(bin)
        case will_state, wont_state, do_state, dont_state:
        return me.DoNegotiate(me.state, bin)
        // subnegotiation started, determine option to subnegotiate
        case sb_state:
        me.sb_telopt = bin      
        me.state     = sb_data_state
        // subnegotiation data, buffer bytes until the end request 
        case sb_data_state:
        return me.sbdataStateProcessByte(bin)
        // IAC received inside a subnegotiation
        case sb_data_iac_state:
        return me.sbdataiacStateProcessByte(bin)
        default:
            //  programing error, shouldn't happen
            panic("Error in telnet state machine!")        
    }
    // return false to signal compression needn't start
    return false
}
 
// Process multiple bytes received from the client
func (me * Telnet) ProcessBytes(bytes []byte) {
    for index := 0 ; index < len(bytes) ; {
        bin := bytes[index]
        compress := me.ProcessByte(bin)
        if compress {
            // paper over this for a while... 
            // new_bytes = Zlib.inflate(arr.pack('c*')) rescue nil
            // if new_bytes
            //arr = new_bytes.bytes.to_a
        }
        index ++
    }
    me.maybeSendDataEventAndEmptyBuffer()
}

  
// Call this when the server receives data from the client
func (me * Telnet) TelnetReceive(data []byte) {
// the COMPRESS2 protocol seems to be half-duplex in that only 
// the server's data stream is compressed (unless maybe if the client
// is asked to also compress with a DO command ?)
    me.ProcessBytes(data)
}

// Send a bytes array (raw) to the client
func (me * Telnet) TelnetSendBytes(bytes ...byte) {
    me.SendRaw(bytes)
}

// Send an iac command 
func (me * Telnet) TelnetSendIac(cmd byte) {
    me.TelnetSendBytes(TELNET_IAC, cmd)
}

// Send negotiation. Currently rfc1143 is not implemented, so beware of 
// server client loops. The simplest way to avoid those is to never answer any 
// client requests, only send server requests.
func (me * Telnet) TelnetSendNegotiate(cmd byte, telopt byte) {
    me.TelnetSendBytes(TELNET_IAC, cmd, telopt)
}
        
// Send non-command data (escapes IAC bytes)
func (me * Telnet) TelnetSend(buffer []byte) {
    me.SendEscaped(buffer)
}
  
// send subnegotiation header
func (me * Telnet) TelnetBeginSubnegotiation(telopt byte) {
    me.TelnetSendBytes(TELNET_IAC, TELNET_SB, telopt)
}


// send subnegotiation ending
func (me * Telnet) TelnetEndSubnegotiation() {
    me.TelnetSendBytes(TELNET_IAC, TELNET_SE)
}

// Send complete subnegotiation
func (me * Telnet) TelnetSubnegotiation(telopt byte, buffer []byte) {
    me.TelnetBeginSubnegotiation(telopt)
    if buffer != nil {
        me.TelnetSend(buffer) 
    }
    me.TelnetEndSubnegotiation()
}
  
// Ask client to start accepting compress2 compression
func (me * Telnet) TelnetBeginCompress2() {
    me.TelnetSendBytes(TELNET_IAC, TELNET_SB, TELNET_TELOPT_COMPRESS2, TELNET_IAC, TELNET_SE);
    me.compress = true
}

// Send formatted data to the client
func (me * Telnet) TelnetRawPrintf(format string, args ...interface{}) {
    buf  := fmt.Sprintf(format, args...)
    me.TelnetSend([]byte(buf))
}

const CRLF  = "\r\n"
const CRNUL = "\r\000"
  
// send formatted data with \r and \n translation in addition to IAC IAC 
// escaping
func (me * Telnet) TelnetPrintf(format string, args ...interface{}) {
    buf  := fmt.Sprintf(format, args...)
    buf   = strings.Replace(buf, "\r", CRNUL, -1)
    buf   = strings.Replace(buf, "\n", CRLF , -1)
    me.TelnetSend([]byte(buf))
}

// NEW-ENVIRON subnegotation
func (me * Telnet) TelnetNewenviron(cmd []byte) {
    me.TelnetSubnegotiation(TELNET_TELOPT_NEW_ENVIRON, cmd)
}

// send TERMINAL-TYPE SEND command
func (me * Telnet)  TelnetTTypeSend() {
    me.TelnetSendBytes(TELNET_IAC, TELNET_SB, TELNET_TELOPT_TTYPE, TELNET_TTYPE_SEND, TELNET_IAC, TELNET_SE)
}

// send TERMINAL-TYPE IS command 
func (me * Telnet)  TelnetTTypeIS(ttype string) {
    me.TelnetSendBytes(TELNET_IAC, TELNET_SB, TELNET_TELOPT_TTYPE, TELNET_TTYPE_IS)
    me.TelnetSend([]byte(ttype))
}

// send MSSP data
func (me * Telnet) TelnetSendMSSP(mssp map[string] string) {
    var buf []byte 
    for key, val := range mssp { 
      buf = append(buf, TELNET_MSSP_VAR)
      buf = append(buf, key...)
      buf = append(buf, TELNET_MSSP_VAL)
      buf = append(buf, val...)
    }
    me.TelnetSubnegotiation(TELNET_TELOPT_MSSP, buf)
}



