package world

import "os"
import "encoding/xml"
import "github.com/beoran/woe/monolog"

/* Ekements of the WOE game world.  
 * Only Zones, Rooms and their Exits, Items, 
 * Mobiles & Characters are saved
 * and loaded from disk. All the rest 
 * is kept statically delared in code for simplocoty.
*/

/* ID used for anything in a world but the world itself and the account. */
type ID string 


type World struct {
    Name                 string
    ZoneIDS         []   ID
    zones           [] * Zone
    CharacterIDS    []   ID
    characters      [] * Character
    MOTD                 string
    
    /* Skills, etc that exist in this world */
    genders         map[ID] *Gender    
    species         map[ID] *BeingKind
    professions     map[ID] *Profession
    skills          map[ID] *Skill
    arts            map[ID] *Art
    techniques      map[ID] *Technique
    exploits        map[ID] *Exploit
    
    /* Botha array and map are needed for serialization. */
    Genders         [] Gender
    Species         [] BeingKind
    Professions     [] Profession
    Skills          [] Skill
    Arts            [] Art
    Techniques      [] Technique
    Exploits        [] Exploit
        
}



func (me * World) AddBeingKind(toadd * BeingKind) {
    me.species[toadd.ID] = toadd
    me.Species = append(me.Species, *toadd)
}

func (me * World) AddProfession(toadd * Profession) {
    me.professions[toadd.ID] = toadd
    me.Professions = append(me.Professions, *toadd)
}

func (me * World) AddSkill(toadd * Skill) {
    me.skills[toadd.ID] = toadd
    me.Skills = append(me.Skills, *toadd)
}

func (me * World) AddArt(toadd * Art) {
    me.arts[toadd.ID] = toadd
    me.Arts = append(me.Arts, *toadd)
}

func (me * World) AddTechnique(toadd * Technique) {
    me.techniques[toadd.ID] = toadd
    me.Techniques = append(me.Techniques, *toadd)
}

func (me * World) AddExploit(toadd * Exploit) {
    me.exploits[toadd.ID] = toadd
    me.Exploits = append(me.Exploits, *toadd)
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
    world.species = make(map[ID] *BeingKind)
    
    world.AddWoeDefaults()
    return world;
}


func HaveID(ids [] ID, id ID) bool {
    for index := 0 ; index < len(ids) ; index++ {
        if ids[index] == id { return true }  
    }
    return false
}

func (me * World) AddZone(zone * Zone) {
    me.zones = append(me.zones, zone)
    if (!HaveID(me.ZoneIDS, zone.ID)) {
        me.ZoneIDS = append(me.ZoneIDS, zone.ID)
    }    
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
    for _, v := range me.Species {me.species[v.ID] = &v }
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



