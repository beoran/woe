package world

import "github.com/beoran/woe/sitef"
import "github.com/beoran/woe/monolog"
// import "fmt"
import "errors"



type Direction  string

type Exit struct {
    Direction
    ToRoomID    int
    toRoom    * Room
}

type Room struct {
    Entity
    Exits   map[Direction]Exit
}

// Load a room from a sitef file.
func LoadRoom(dirname string, id string) (room * Room, err error) {
    
    path := SavePathFor(dirname, "room", id)
    
    records, err := sitef.ParseFilename(path)
    if err != nil {
        return nil, err
    }
    
    if len(records) < 1 {
        return nil, errors.New("No room found!")
    }
    
    record := records[0]
    monolog.Info("Loading Room record: %s %v", path, record)
    
    room = new(Room)
    room.Entity.LoadSitef(*record)
    /*
    account.Name            = record.Get("name")
    account.Hash            = record.Get("hash")
    account.Algo            = record.Get("algo")
    account.Email           = record.Get("email")
    account.Points          = record.GetIntDefault("points", 0)
    account.Privilege       = Privilege(record.GetIntDefault("privilege", 
                                int(PRIVILEGE_NORMAL)))
    
    var nchars int
    nchars                  = record.GetIntDefault("characters", 0)
    _ = nchars    
    */
    monolog.Info("Loaded Room: %s %v", path, room)
    return room, nil
}
