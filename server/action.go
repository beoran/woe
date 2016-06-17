package server

import "bytes"
import "errors"
import "regexp"
// import "github.com/beoran/woe/telnet"
import "github.com/beoran/woe/world"
import "github.com/beoran/woe/monolog"


/* Actiondata are the params passed to an Actin, neatly wrapped in 
 * a struct */
type ActionData struct {
    Client * Client
    Server * Server
    World  * world.World
    Account* world.Account
    Action * Action
    Command  []byte
    Rest     []byte
    Argv     [][]byte
}

/* A handler for an action. */
type ActionHandler func(data * ActionData) (err error) 

/* Actions that  a client can perform on a server or in 
 * the server's world. */
type Action struct {
    Name        string
    Privilege   world.Privilege
    Handler     ActionHandler
}


var ActionMap map[string] Action

func AddAction(name string, privilege world.Privilege, handler ActionHandler) {
    monolog.Info("Adding new action %s with privilege %d", name, privilege)
    action := Action{name, privilege, handler}
    ActionMap[name] = action
}

func doShout(data * ActionData) (err error) {
      data.Server.Broadcast("Client said %s\r\n", data.Rest)
      return nil  
}

func doShutdown(data * ActionData) (err error) {    
    data.Server.Broadcast("Shutting down server NOW!\n")
    data.Server.Restart();
    return nil
}

func doRestart(data * ActionData) (err error) {
    data.Server.Broadcast("Restarting server NOW!\n")
    data.Server.Restart();
    return nil
}

func doQuit(data * ActionData) (err error) {  
    data.Client.Printf("Byebye!\n")
    data.Client.Disconnect()
    return nil
}

func doEnableLog(data * ActionData) (err error) {  
    // strings. string(data.Rest)
    return nil
}


func ParseCommand(command []byte, data * ActionData) (err error) {
    /* strip any leading blanks  */
    trimmed    := bytes.TrimLeft(command, " \t")
    re         := regexp.MustCompile("[^ \t,]+")
    parts      := re.FindAll(command, -1)
    
    bytes.SplitN(trimmed, []byte(" \t,"), 2)  
    
    if len(parts) < 1 {
        data.Command = nil
        return errors.New("Come again?")
    }
    data.Command = parts[0]
    if len(parts) > 1 { 
        data.Rest    = parts[1]
        data.Argv    = parts
    } else {
        data.Rest    = nil
        data.Argv    = nil
    }
        
    return nil
} 

func init() {
    ActionMap = make(map[string] Action)
    AddAction("/shutdown"   , world.PRIVILEGE_LORD, doShutdown)
    AddAction("/restart"    , world.PRIVILEGE_LORD, doRestart)
    AddAction("/quit"       , world.PRIVILEGE_ZERO, doQuit)
}

func (client * Client) ProcessCommand(command []byte) {
    ad := &ActionData{client, client.GetServer(), 
        client.GetWorld(), client.GetAccount(), nil, nil, nil, nil }
    _ = ad
    err := ParseCommand(command, ad);
    if err != nil {
        client.Printf("%s", err)
        return
    }
    
    action, ok := ActionMap[string(ad.Command)]
    ad.Action = &action
    
    if ad.Action == nil || (!ok) {
        client.Printf("Unknown command %s.", ad.Command)
        return
    }
    // Check if sufficient rights to perform the action
    if (ad.Action.Privilege > client.GetAccount().Privilege) {
        client.Printf("You lack the privilege to %s (%d vs %d).", 
        ad.Command, ad.Action.Privilege, client.GetAccount().Privilege)
        return
    }
    
    // Finally run action
    ad.Action.Handler(ad)
} 

