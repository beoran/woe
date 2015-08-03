package world
import "os"
import "encoding/xml"




// Anything insdie the WOE World can be identified by a unique short string 
// description, the Label
type Labeled interface {
    Label() string // Returns a unique label of the thing.
}

type Typed interface {
    Type() string // Returns a tring description of the type of the thing. 
}


// An entity is anything that can exist in a World
type Entity struct {
    ID                  ID          `xml:"id,attr"`
    Name                string      `xml:"name,attr"`
    Short               string      `xml:"short,attr"`
    Long                string
    Aliases           []string          
}

func (me * Entity) Label() string {
    return string(me.ID)
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
    Get(ID)         Labeled
    Put(Labeled)
    Size()          int
    Index(int)      Labeled
    PutIndex(int)
}

type LabeledList struct {
    byList        []Labeled
    byLabel       map[ID] Labeled
}

func NewLabeledList() * LabeledList {
    byname := make(map[ID] Labeled)
    return &LabeledList{nil, byname}
}

func (me * LabeledList) Get(id ID) Labeled {
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








