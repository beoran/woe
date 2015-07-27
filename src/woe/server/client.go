package server


import (
    // "fmt"
    "net"
    "time"
    // "errors"
    // "io"
    "github.com/beoran/woe/monolog"
)

type Client struct {
    server * Server
    id       int
    conn     net.Conn
    telnet   interface{}
    alive    bool
    timeout  int
    datachan chan []byte
    errchan  chan error
    timechan chan time.Time 
}


func NewClient(server * Server, id int, conn net.Conn) * Client {
    datachan := make (chan []byte, 1024)
    errchan  := make (chan error, 1)
    timechan := make (chan time.Time, 32)
    return &Client{server, id, conn, nil, true, -1, datachan, errchan, timechan}
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
        _ , err := me.conn.Read(buffer);
        if err != nil {
            me.errchan <- err
            return
        }
        me.datachan <- buffer
    }
}


func (me * Client) TryRead(millis int) (data [] byte, timeout bool, close bool) {
    select {
        case data := <- me.datachan:
            return data, false, false
                       
        case err  := <- me.errchan:
            monolog.Info("Connection closed: %s\n", err)
            me.Close()
            return nil, false, true
            
        case _ = <- time.Tick(time.Millisecond * time.Duration(millis)):
            return nil, true, false
    }
}


func (me * Client) Serve() (err error) {
    // buffer := make([]byte, 1024, 1024)
    go me.ServeRead()
    for (me.alive) {
        data, _, _ := me.TryRead(3000)
        
        if data == nil {
            me.conn.Write([]byte("Too late!\r\n"))
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
