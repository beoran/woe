package server

import (
	// "fmt"
	"net"
	"time"
	// "errors"
	// "io"
	"github.com/beoran/woe/monolog"
	"github.com/beoran/woe/telnet"
	"github.com/beoran/woe/world"
)

/* Specific properties of a client. */
type ClientInfo struct {
	w         int
	h         int
	mtts      int
	naws      bool
	compress2 bool
	mssp      bool
	zmp       bool
	msp       bool
	msdp      bool
	mxp       bool
	ttype     bool
	terminals []string
	terminal  string
}

type Client struct {
	server   *Server
	id       int
	conn     net.Conn
	alive    bool
	timeout  int
	datachan chan []byte
	errchan  chan error
	timechan chan time.Time
	telnet   *telnet.Telnet
	info     ClientInfo

	// Account of client or nil if not yet selected.
	account *world.Account
	// Character client is plaing with or nil if not yet selected.
	character *world.Character
	// Message channels that this client is listening to once fully logged in.
	// Not to be confused with Go channels.
	channels map[string]bool
}

func NewClient(server *Server, id int, conn net.Conn) *Client {
	datachan := make(chan []byte, 1024)
	errchan := make(chan error, 1)
	timechan := make(chan time.Time, 32)
	telnet := telnet.New()
	channels := make(map[string]bool)
	info := ClientInfo{w: -1, h: -1, terminal: "none"}
	return &Client{server, id, conn, true, -1, datachan, errchan, timechan, telnet, info, nil, nil, channels}
}

func (me *Client) Close() {
	me.conn.Close()
	me.alive = false
	if me.account != nil {
		me.server.World.RemoveAccount(me.account.Name)
	}
	me.account = nil
}

/** Goroutine that does the actual reading of input data, and sends it to the
 * needed channels. */
func (me *Client) ServeRead() {
	for me.alive {
		buffer := make([]byte, 1024, 1024)
		read, err := me.conn.Read(buffer)
		if err != nil {
			me.errchan <- err
			return
		}
		// reply will be stored in me.telnet.Events channel
		me.telnet.ProcessBytes(buffer[:read])
	}
}

/* Goroutine that sends any data that must be sent through the Telnet protocol
 * to the connected client.
 */
func (me *Client) ServeWrite() {
	for me.alive {
		select {
		case data := <-me.telnet.ToClient:
			monolog.Log("SERVEWRITE", "Will send to client: %v", data)
			me.conn.Write(data)
		}
	}
}

func (me *Client) TryReadEvent(millis int) (event telnet.Event, timeout bool, close bool) {
	var timerchan <-chan (time.Time)

	if millis >= 0 {
		timerchan = time.Tick(time.Millisecond * time.Duration(millis))
	} else {
		/* If time is negative, block by using a fake time channel that never gets sent anyting */
		timerchan = make(<-chan (time.Time))
	}

	select {
	case event := <-me.telnet.Events:
		return event, false, false

	case err := <-me.errchan:
		monolog.Info("Connection closed: %s\n", err)
		me.Close()
		return nil, false, true

	case _ = <-timerchan:
		return nil, true, false
	}
}

func (me *Client) HandleNAWSEvent(nawsevent *telnet.NAWSEvent) {
	me.info.w = nawsevent.W
	me.info.h = nawsevent.H
	monolog.Info("Client %d window size #{%d}x#{%d}", me.id, me.info.w, me.info.h)
	me.info.naws = true
}

func (me *Client) TryRead(millis int) (data []byte, timeout bool, close bool) {

	for me.alive {
		event, timeout, close := me.TryReadEvent(millis)
		if event == nil && (timeout || close) {
			return nil, timeout, close
		}
		switch event := event.(type) {
		case *telnet.DataEvent:
			monolog.Log("TELNETDATAEVENT", "Telnet data event %T : %d.", event, len(event.Data))
			return event.Data, false, false
		case *telnet.NAWSEvent:
			monolog.Log("TELNETNAWSEVENT", "Telnet NAWS event %T.", event)
			me.HandleNAWSEvent(event)
		default:
			monolog.Info("Ignoring telnet event %T : %v for now.", event, event)
		}
	}

	return nil, false, true
}

func (me *Client) Serve() (err error) {
	// buffer := make([]byte, 1024, 1024)
	go me.ServeWrite()
	go me.ServeRead()
	me.SetupTelnet()
	if me.server.World != nil {
		me.Printf(me.server.World.MOTD)
	}
	if !me.AccountDialog() {
		time.Sleep(3)
		// sleep so output gets flushed, hopefully.
		// Also slow down brute force attacks.
		me.Close()
		return nil
	}

	if !me.CharacterDialog() {
		time.Sleep(3)
		// sleep so output gets flushed, hopefully.
		// Also slow down brute force attacks.
		me.Close()
		return nil
	}

	me.Printf("Welcome, %s\n", me.account.Name)

	for me.alive {
		me.HandleCommand()
		/*
		   data, _, _ := me.TryRead(3000)

		   if data == nil {
		      // me.telnet.TelnetPrintf("Too late!\r\n")
		   } else {
		       me.server.Broadcast(string(data))
		   }
		*/

	}
	return nil
}

func (me *Client) Disconnect() {
	me.alive = false
}

func (me *Client) IsAlive() bool {
	return me.alive
}

func (me *Client) IsLoginFinished() bool {
	return me.IsAlive() && (me.account != nil) && (me.character != nil)
}

func (me *Client) SetChannel(channelname string, value bool) {
	me.channels[channelname] = value
}

func (me *Client) IsListeningToChannel(channelname string) bool {
	res, ok := me.channels[channelname]
	// All chanels are active by default, and must be actively disabled by the client.
	if ok && (!res) {
		return false
	}
	return true
}

func (me *Client) WriteString(str string) {
	me.conn.Write([]byte(str))
}

/** Accessor */
func (me *Client) GetServer() *Server {
	return me.server
}

/** Accessor */
func (me *Client) GetWorld() *world.World {
	return me.server.World
}

/** Accessor */
func (me *Client) GetAccount() *world.Account {
	return me.account
}
