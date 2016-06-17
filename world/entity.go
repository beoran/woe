package world
import "strings"
import "os"
import "sort"
import "encoding/xml"
import "github.com/beoran/woe/sitef"




// Anything inside the WOE World can be identified by a unique short string 
// description, the Label
type Labeled interface {
    Label() string // Returns a unique label of the thing.
}

type Typed interface {
    Type() string // Returns a string description of the type of the thing. 
}


// An entity is anything that can exist in a World
type Entity struct {
    ID                  string      `xml:"id,attr"`
    Name                string      `xml:"name,attr"`
    Short               string      `xml:"short,attr"`
    Long                string
    Aliases           []string    
    // Privilege level needed to use/interact with/enter/... this Entity
    Privilege           Privilege
}


func EntityNameToID(kind string, name string) string {
    return kind + "_" + strings.ToLower(name)
}


func (me * Entity) InitKind(kind string, name string, 
    privilege Privilege) (* Entity) {
    if me == nil {
        return me
    }
    me.ID       = EntityNameToID(kind, name) 
    me.Name     = name
    me.Short    = name
    me.Long     = name
    me.Privilege= privilege
    return me
}

// Devious polymorphism for Entities...
type Entitylike interface {
    AsEntity() * Entity
}

// A little trick to need less setters and getters
func (me * Entity) AsEntity() (*Entity) {
    return me
}

func (me Entity) Label() string {
    return string(me.ID)
}


type EntitySlice []Entitylike
type EntityMap map[string]Entitylike

type EntityLookup struct {
    slice EntitySlice
    table EntityMap  
}


// By is the type of a "less" function that defines the ordering of its Planet arguments.
type EntityLikeBy func(p1, p2 Entitylike) bool

// planetSorter joins a By function and a slice of Planets to be sorted.
type EntityLikeSorter struct {
    Entities EntitySlice
    By       EntityLikeBy
}


func (me EntitySlice) By(compare func (p1, p2 Entitylike) bool) *EntityLikeSorter {
    return &EntityLikeSorter { me,  EntityLikeBy(compare) };
}

// Len is part of sort.Interface.
func (s * EntityLikeSorter) Len() int {
    return len(s.Entities)
}

// Swap is part of sort.Interface.
func (s *EntityLikeSorter) Swap(i, j int) {
    s.Entities[i], s.Entities[j] = s.Entities[j], s.Entities[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *EntityLikeSorter) Less(i, j int) bool {
    return s.By(s.Entities[i], s.Entities[j])
}

func (me * EntityLikeSorter) Sort() {
    sort.Sort(me)
} 

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by EntityLikeBy) Sort(entities EntitySlice) {
    ps := &EntityLikeSorter{
        Entities: entities,
        By:      by, // The Sort method's receiver is the function (closure) that defines the sort order.
    }
    sort.Sort(ps)
}


func (me EntityLookup) Get(index int) Entitylike {
    return me.slice[index]
}

func (me EntityLookup) Len() int {
    return len(me.table)
}

func (me * EntityLookup) Add(ent Entitylike) int {
    me.slice = append(me.slice, ent)
    me.slice.By(func (e1, e2 Entitylike) bool {
            return (e1.AsEntity().Label() < e2.AsEntity().Label()) 
        }).Sort()
    res := len(me.slice) - 1 
    if me.table == nil {
        me.table = make(EntityMap)
    }
    me.table[ent.AsEntity().Label()] = ent
    return res 
}

func (me EntityLookup) Lookup(id string) (Entitylike) {
    res, ok := me.table[id]
    if !ok {
        return nil
    }
    return res    
}

func (me EntityLookup) Remove(ent Entitylike) bool {
    key := ent.AsEntity().Label()
    delete(me.table, key)
    return true
}

type EntityLookupEachFunc func (id string, ent Entitylike, args...interface{}) Entitylike

func (me EntityLookup) Each(lambda EntityLookupEachFunc, args...interface{}) Entitylike {
    for k, v := range (me.table) {
        res := lambda(k, v, args...) 
        if res != nil {
            return res
        }
    } 
    return nil
}


type EntityList interface {
    Get(int) Entitylike
    Len() int
    Add(Entitylike) int
    Lookup(string) (Entitylike)
    Remove(Entitylike) bool
}

// Interface 
type Savable interface {
    Labeled
    Typed
}

type Loadable interface {    
    Typed
}

func SaveSavable(dirname string, savable Savable) (err error) {
    path := SavePathFor(dirname, savable.Type(), savable.Label())
    
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    enc := xml.NewEncoder(file)
    enc.Indent(" ", "  ")
    return enc.Encode(savable)
}

func LoadLoadable(dirname string, nameid string, result Loadable) (Loadable) {
    path := SavePathFor(dirname, result.Type(), nameid)
    
    file, err := os.Open(path)
    if err != nil {
        return nil
    }
    dec := xml.NewDecoder(file)    
    err = dec.Decode(result)
    if err != nil { return nil }
    return result
}

// A list of Identifier items mapped to their ID's 
type LabeledLister interface {
    Get(string)     Labeled
    Put(Labeled)
    Size()          int
    Index(int)      Labeled
    PutIndex(int)
}

type LabeledList struct {
    byList        []Labeled
    byLabel       map[string] Labeled
}

func NewLabeledList() * LabeledList {
    byname := make(map[string] Labeled)
    return &LabeledList{nil, byname}
}

func (me * LabeledList) Get(id string) Labeled {
    val, ok := me.byLabel[id]
    if !ok { return nil }
    return val
}

func (me * LabeledList) Index(index int) Labeled {
    if index < 0 { return nil } 
    if index > len(me.byList) { return nil }
    val := me.byList[index]
    return val
}


// Save an entity to a sitef record.
func (me * Entity) SaveSitef(rec * sitef.Record) (err error) {
    rec.Put("id", me.ID)
    rec.Put("name", me.Name)
    rec.Put("short", me.Short)
    rec.Put("long",  me.Long)
    return nil
}

// Load an entity from a sitef record.
func (me * Entity) LoadSitef(rec sitef.Record) (err error) {
    me.ID       = rec.Get("id")
    me.Name     = rec.Get("name")
    me.Short    = rec.Get("short")
    me.Long     = rec.Get("long")
    return nil
}


func (me Entity) AskName() string {
    return me.Name
}

func (me Entity) AskShort() string {
    return me.Short
}

func (me Entity) AskLong() string {
    return me.Long
}

func (me Entity) AskPrivilege() Privilege {
    return me.Privilege
}





