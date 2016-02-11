package server

import (
    "log"
  //  "io"
    "net"
  //  "errors"
    "os"
    "math/rand"
    "time"
    "fmt"
    "path/filepath"
    "github.com/beoran/woe/monolog"
    "github.com/beoran/woe/world"
)

var MSSP map[string] string

const MAX_CLIENTS = 1000

func init() {
     MSSP = map[string] string {
          "NAME"        : "Workers Of Eruta",
          "UPTIME"      : string(time.Now().Unix()),
          "PLAYERS"     : "0",
          "CRAWL DELAY" : "0",
          "CODEBASE"    : "WOE",
          "CONTACT"     : "beoran@gmail.com",
          "CREATED"     : "2015",
           "ICON"       : "None",
          "LANGUAGE"    : "English",
          "LOCATION"    : "USA",
          "MINIMUM AGE" : "18",
          "WEBSITE"     : "beoran.net",
          "FAMILY"      : "Custom",
          "GENRE"       : "Science Fiction",
          "GAMEPLAY"    : "Adventure",
          "STATUS"      : "Alpha",
          "GAMESYSTEM"  : "Custom",
          "INTERMUD"    : "",
          "SUBGENRE"    : "None",
          "AREAS"       : "0",
          "HELPFILES"   : "0",
          "MOBILES"     : "0",
          "OBJECTS"     : "0",
          "ROOMS"       : "1",
          "CLASSES"     : "0",
          "LEVELS"      : "0",
          "RACES"       : "3",
          "SKILLS"      : "900",
          "ANSI"        : "1",
          "MCCP"        : "1",
          "MCP"         : "0",
          "MSDP"        : "0",
          "MSP"         : "0",
          "MXP"         : "0",
          "PUEBLO"      : "0",
          "UTF-8"       : "1",
          "VT100"       : "1",
          "XTERM 255 COLORS" : "1",
          "PAY TO PLAY"      : "0",
          "PAY FOR PERKS"    : "0",
          "HIRING BUILDERS"  : "0",
          "HIRING CODERS"    : "0" }  
}




type Server struct {
    address               string
    listener              net.Listener
    logger              * log.Logger
    logfile             * os.File
    clients map[int]    * Client 
    tickers map[string] * Ticker
    alive                 bool
    World               * world.World
}


type Ticker struct {
    * time.Ticker
    Server        * Server
    Name            string
    Milliseconds    int
    callback        func(me * Ticker, t time.Time) (stop bool)
}


const DEFAULT_MOTD =

`
###############################
#       Workers Of Eruta      # 
###############################

`


func (me * Server) SetupWorld() error {
    /*
    me.World, _ = world.LoadWorld(me.DataPath(), "WOE")
    if me.World == nil {
        monolog.Info("Creating new default world...")
        me.World = world.NewWorld("WOE", DEFAULT_MOTD)
        err := me.World.Save(me.DataPath())
        if err != nil {
            monolog.Error("Could not save world: %v", err)
            return err
        } else {
            monolog.Info("Saved default world.")
        }
    }
    */
    return nil
}


func NewServer(address string) (server * Server, err error) {
    listener, err := net.Listen("tcp", address);
    if (err != nil) { 
        io.Printf("")
        return nil, err
    }
    
    logfile, err := os.OpenFile("log.woe", os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0660)
    if (err != nil) {
        return nil, err
    }
    
    logger := log.New(logfile, "woe", log.Llongfile | log.LstdFlags)
    clients := make(map[int] * Client)
    tickers := make(map[string] * Ticker)

    server = &Server{address, listener, logger, logfile, clients, tickers, true, nil}
    err = server.SetupWorld()
    server.AddDefaultTickers()
    
    return server, err
}

func (me * Server) Close() {
    me.logfile.Close();
}


func NewTicker(server * Server, name string, milliseconds int, callback func (me * Ticker, t time.Time) bool) (* Ticker) {
    ticker := time.NewTicker(time.Millisecond * time.Duration(milliseconds))
    return &Ticker {ticker, server, name, milliseconds, callback}
}


func (me * Ticker) Run() {
    OUTER: 
    for me.Server.alive {
        for tick := range me.C {
            if (!me.callback(me, tick)) {
                break OUTER;
            }
        }
    }
}


func (me * Server) RemoveTicker(name string) {
    ticker, have := me.tickers[name]
    if (!have) {
        return
    }    
    ticker.Stop()
    delete(me.tickers, name)
}

func (me * Server) StopTicker(name string) {
    ticker, have := me.tickers[name]
    if (!have) {
        return
    }    
    ticker.Stop();
}



func (me * Server) AddTicker(name string, milliseconds int, callback func (me * Ticker, t time.Time) bool) (* Ticker) {
    _, have := me.tickers[name]
    
    if have {
        me.RemoveTicker(name)
    }
        
    ticker := NewTicker(me, name, milliseconds, callback)
    me.tickers[name] = ticker
    go ticker.Run();
    
    return ticker
}


func onWeatherTicker (me * Ticker, t time.Time) bool {
    monolog.Info("Weather Ticker tick tock.")
    return true
}


func (me * Server) AddDefaultTickers() {
    me.AddTicker("weather", 10000, onWeatherTicker)    
}

func (me * Server) handleDisconnectedClients() {
    for { 
        time.Sleep(1)
        for id, client := range me.clients {
            if (!client.IsAlive()) {
                monolog.Info("Client %d has disconnected.", client.id)
                client.Close()
                delete(me.clients, id);
            }
        }   
    }
}

func (me * Server) findFreeID() (id int, err error) {
    for id = 0 ; id < MAX_CLIENTS ; id++ {
        client, have := me.clients[id]
        if (!have) || (client == nil) {
            return id, nil
        }
    }
    return -1, fmt.Errorf("Too many clients!");
}

func (me * Server) onConnect(conn net.Conn) (err error) {
    id, err := me.findFreeID()
    if err != nil {
        monolog.Info("Refusing connection for %s: too many clients. ", conn.RemoteAddr().String())
        conn.Close()
        return nil
    }
    monolog.Info("New client connected from %s, id %d. ", conn.RemoteAddr().String(), id)
    client := NewClient(me, id, conn)
    me.clients[id] = client
    return client.Serve()
}

func (me * Server) Serve() (err error) { 
    // Setup random seed here, or whatever
    rand.Seed(time.Now().UTC().UnixNano())
    
    go me.handleDisconnectedClients()
    
    for (me.alive) {
        conn, err := me.listener.Accept()
        if err != nil {
            return err
        }
        go me.onConnect(conn)
    }
    return nil
}


func (me * Server) Broadcast(message string) {
    for _, client := range me.clients {
        if (client.IsAlive()) {
            client.WriteString(message)
        }
    }       
}


// Returns the data path of the server
func (me * Server) DataPath() string {
    // 
    cwd, err := os.Getwd();
    if  err != nil {
        cwd = "."
    }
    
    return filepath.Join(cwd, "data", "var")
}

// Returns the script path of the server
func (me * Server) ScriptPath() string {
    // 
    cwd, err := os.Getwd();
    if err != nil {
        cwd = "."
    }
    
    return filepath.Join(cwd, "data", "script")
}




