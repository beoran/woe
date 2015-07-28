package server


import (
    // "fmt"
    "net"
    "time"
    // "errors"
    // "io"
    "github.com/beoran/woe/monolog"
    "github.com/beoran/woe/telnet"
)

/* Specific properties of a client. */
type ClientInfo struct {
    w           int
    h           int
    mtts        int
    naws        bool
    compress2   bool
    mssp        bool
    zmp         bool
    msp         bool
    msdp        bool
    mxp         bool
    ttype       bool
    terminals []string
    terminal    string
}

type Client struct {
    server * Server
    id       int
    conn     net.Conn
    alive    bool
    timeout  int
    datachan chan []byte
    errchan  chan error
    timechan chan time.Time 
    telnet * telnet.Telnet
    info     ClientInfo
}


func NewClient(server * Server, id int, conn net.Conn) * Client {
    datachan := make (chan []byte, 1024)
    errchan  := make (chan error, 1)
    timechan := make (chan time.Time, 32)
    telnet   := telnet.New()
    info     := ClientInfo{-1, -1, 0, false, false, false, false, false, false, false, false, nil, "none"}
    return &Client{server, id, conn, true, -1, datachan, errchan, timechan, telnet, info}
}

func (me * Client) Close() {
    me.conn.Close()
    me.alive = false
}

/** Goroutine that does the actual reading of input data, and sends it to the 
 * needed channels. */    
func (me * Client) ServeRead() {
    for (me.alive) { 
        buffer  := make([]byte, 1024, 1024)
        read , err := me.conn.Read(buffer);
        if err != nil {
            me.errchan <- err
            return
        }
        // reply will be stored in me.telnet.Events
        me.telnet.ProcessBytes(buffer[:read])
    }
}

/* Goroutine that sends any data that must be sent through the Telnet protocol 
 * (that any data, really) other to the connected client.
 */
 func (me * Client) ServeWrite() {
     for (me.alive) {
        select {
            case data := <- me.telnet.ToClient:
            monolog.Debug("Will send to client: %v", data)
            me.conn.Write(data)
        }
    }
}


func (me * Client) TryReadEvent(millis int) (event telnet.Event, timeout bool, close bool) {
    select {
        case event := <- me.telnet.Events:
            return event, false, false
                       
        case err  := <- me.errchan:
            monolog.Info("Connection closed: %s\n", err)
            me.Close()
            return nil, false, true
            
        case _ = <- time.Tick(time.Millisecond * time.Duration(millis)):
            return nil, true, false
    }
}

func (me * Client) TryRead(millis int) (data []byte, timeout bool, close bool) {
    
    for (me.alive) { 
        event, timeout, close := me.TryReadEvent(millis)
        if event == nil && (timeout || close) {
            return nil, timeout, close
        }
        switch event := event.(type) {
            case * telnet.DataEvent:
                monolog.Debug("Telnet data event %T : %d.", event, len(event.Data))
                return event.Data, false, false
            default:
                monolog.Info("Ignoring telnet event %T : %v for now.", event, event)
        }
    }
    
    return nil, false, true
}


func (me * Client) Serve() (err error) {
    // buffer := make([]byte, 1024, 1024)
    go me.ServeWrite()
    go me.ServeRead()
    me.SetupTelnet()
    
    for (me.alive) {
        
        
        data, _, _ := me.TryRead(3000)
        
        if data == nil {
           // me.telnet.TelnetPrintf("Too late!\r\n")
        } else {
            me.server.Broadcast(string(data))
        }
        
    }
    return nil
}

func (me * Client) IsAlive() bool {
    return me.alive
}

func (me * Client) WriteString(str string) {
    me.conn.Write([]byte(str))
}
