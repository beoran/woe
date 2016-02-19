package server

// import "github.com/beoran/woe/telnet"
import "github.com/beoran/woe/world"
import "github.com/beoran/woe/monolog"


/* Actiondata are the params passed to an Actin, neatly wrapped in 
 * a struct */
type ActionData struct {
    Client * Client
    Server * Server
    World  * world.World
    Action * Action
    Command  string
    Argv   []string
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

func doShutdown(data * ActionData) (err error) {
    return nil
}

func doRestart(data * ActionData) (err error) {
    return nil
}


func doQuit(data * ActionData) (err error) {    
    return nil
}

func ParseCommand(command string, data * ActionData) {
    data.Command = command
} 

func init() {
    ActionMap = make(map[string] Action)
    AddAction("/shutdown"   , world.PRIVILEGE_LORD, doShutdown)
    AddAction("/restart"    , world.PRIVILEGE_LORD, doRestart)
    AddAction("/quit"       , world.PRIVILEGE_ZERO, doQuit)
}

