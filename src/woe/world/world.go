package world

import "github.com/beoran/woe/monolog"
import "github.com/beoran/woe/sitef"
import "errors"

/* Elements of the WOE game world.  
 * Only Zones, Rooms and their Exits, Items, 
 * Mobiles & Characters are saved
 * and loaded from disk. All the rest 
 * is kept statically delared in code for simplicity.
*/


type World struct {
    Name                      string
    MOTD                      string
    dirname                   string
    entitymap       map[string] * Entity
    zonemap         map[string] * Zone
    zones                [] * Zone
    charactermap    map[string] * Character
    characters           []   Character
    roommap         map[string] * Room
    rooms                []   Room
    itemmap         map[string] * Item
    items                []   Item
    mobilemap       map[string] * Mobile
    mobiles              []   Mobile
    accounts             [] * Account
    accountmap      map[string] * Account
}



func (me * World) AddWoeDefaults() {
    /*
    me.AddSpecies(NewSpecies("sp_human"  , "Human"))
    me.AddSpecies(NewSpecies("sp_neosa"  , "Neosa"))
    me.AddSpecies(NewSpecies("sp_mantu"  , "Mantu"))
    me.AddSpecies(NewSpecies("sp_cyborg" , "Cyborg"))
    me.AddSpecies(NewSpecies("sp_android", "Android"))
    */
}

func NewWorld(name string, motd string, dirname string) (*World) {
    world := new(World)
    world.Name          = name
    world.MOTD          = motd
    world.dirname       = dirname
    world.accountmap    = make(map[string] * Account)
    world.itemmap       = make(map[string] * Item)
    world.roommap       = make(map[string] * Room)
    world.charactermap  = make(map[string] * Character)

    world.AddWoeDefaults()
    return world;
}


func HaveID(ids [] string, id string) bool {
    for index := 0 ; index < len(ids) ; index++ {
        if ids[index] == id { return true }  
    }
    return false
}

func (me * World) AddEntity(entity * Entity) {
    me.entitymap[entity.ID] = entity;
}

func (me * World) AddZone(zone * Zone) {
    me.zones = append(me.zones, zone)
    me.zonemap[zone.ID] = zone;
    me.AddEntity(&zone.Entity);
}

// Save an account as a sitef file.
func (me * World) Save(dirname string) (err error) {
    path := SavePathFor(dirname, "world", me.Name)
    
    rec                  := sitef.NewRecord()
    rec.Put("name",         me.Name)
    rec.Put("motd",         me.MOTD)
    monolog.Debug("Saving World record: %s %v", path, rec)
    return sitef.SaveRecord(path, *rec)
}

// Load a world from a sitef file.
func LoadWorld(dirname string, name string) (world * World, err error) {
    
    path := SavePathFor(dirname, "world", name)
    
    records, err := sitef.ParseFilename(path)
    if err != nil {
        return nil, err
    }
    
    if len(records) < 1 {
        return nil, errors.New("No record found!")
    }
    
    record := records[0]
    monolog.Info("Loading World record: %s %v", path, record)
    
    world = NewWorld(record.Get("name"), record.Get("motd"), dirname)
    monolog.Info("Loaded World: %s %v", path, world)
    return world, nil
}


// Returns an acccount that has already been loaded or nil if not found
func (me * World) GetAccount(name string) (account * Account) {
    account, ok := me.accountmap[name];
    if !ok {
        return nil
    }
    return account
} 

// Loads an account to be used with this world. Characters will be linked.
// If the account was already loaded, returns that in stead.
func (me * World) LoadAccount(name string) (account *Account, err error) {
    account = me.GetAccount(name)
    if (account != nil) {
        return account, nil
    }
    
    account, err = LoadAccount(me.dirname, name);
    if err != nil {
        return account, err
    }
    me.accountmap[account.Name] = account
    return account, nil
}

// Removes an account from this world by name.
func (me * World) RemoveAccount(name string) {
    _, have := me.accountmap[name]
    if (!have) {
        return
    }    
    delete(me.accountmap, name)
}

// Default world pointer
var DefaultWorld * World


// Returns an item that has already been loaded or nil if not found
func (me * World) GetItem(id string) (item * Item) {
    item, ok := me.itemmap[id]
    if !ok {
        return nil
    }
    return item
} 

// Loads an item to be used with this world. 
// If the item was already loaded, returns that in stead.
func (me * World) LoadItem(id string) (item *Item, err error) {
    item = me.GetItem(id)
    
    if (item != nil) {
        return item, nil
    }
    
    item, err = LoadItem(me.dirname, id);
    if err != nil {
        return item, err
    }
    me.itemmap[item.ID] = item
    return item, nil
}

// Removes an item from this world by ID.
func (me * World) RemoveItem(id string) {
    _, have := me.itemmap[id]
    if (!have) {
        return
    }    
    delete(me.itemmap, id)
}


// Returns a Room that has already been loaded or nil if not found
func (me * World) GetRoom(id string) (room * Room) {
    room, ok := me.roommap[id]
    if !ok {
        return nil
    }
    return room
} 



// Loads a Room to be used with this world. 
// If the room was already loaded, returns that in stead.
func (me * World) LoadRoom(id string) (room *Room, err error) {
    room = me.GetRoom(id)
    
    if (room != nil) {
        return room, nil
    }
    
    room, err = LoadRoom(me.dirname, id);
    if err != nil {
        return room, err
    }
    me.roommap[room.ID] = room
    return room, nil
}

// Removes an item from this world by ID.
func (me * World) RemoveRoom(id string) {
    _, have := me.roommap[id]
    if (!have) {
        return
    }    
    delete(me.roommap, id)
}



