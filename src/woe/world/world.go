package world

import "os"
import "encoding/xml"
import "github.com/beoran/woe/monolog"

/* Elements of the WOE game world.  
 * Only Zones, Rooms and their Exits, Items, 
 * Mobiles & Characters are saved
 * and loaded from disk. All the rest 
 * is kept statically delared in code for simplicity.
*/

/* ID used for anything in a world but the world itself and the account. */
type ID string 


type World struct {
    Name                      string
    MOTD                      string
    entitymap       map[ID] * Entity
    zonemap         map[ID] * Zone
    zones                [] * Zone
    charactermap    map[ID] * Character
    characters           []   Character
    roommap         map[ID] * Room
    rooms                []   Room
    itemmap         map[ID] * Item
    items                []   Item
    mobilemap       map[ID] * Mobile
    mobiles              []   Mobile
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

func NewWorld(name string, motd string) (*World) {
    world := new(World)
    world.Name = name
    world.MOTD = motd
    
    world.AddWoeDefaults()
    return world;
}


func HaveID(ids [] ID, id ID) bool {
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

func (me * World) Save(dirname string) (err error) {
    path := SavePathFor(dirname, "world", me.Name)
    
    file, err := os.Create(path)
    if err != nil {
        monolog.Error("Could not load %name: %v", err)
        return err
    }
    enc := xml.NewEncoder(file)
    enc.Indent(" ", "  ")
    res := enc.Encode(me)
    if (res != nil) {
        monolog.Error("Could not save %s: %v", me.Name, err)
    }
    return res
}

func (me * World) onLoad() {
}

func LoadWorld(dirname string, name string) (result *World, err error) {
    path := SavePathFor(dirname, "world", name)
    
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    dec := xml.NewDecoder(file)    
    result = new(World)
    err = dec.Decode(result)
    if err != nil {
        monolog.Error("Could not load %s: %v", name, err)
        panic(err)
    }
    
    result.onLoad()
    return result, nil
}



