package world

import (
    "fmt"
    "strings"
)

/* Aptitudes of a being, species or profession */
type Aptitudes struct {
    Skills      []BeingSkill
    Arts        []BeingArt
    Techniques  []BeingTechnique
    Exploits    []BeingExploit
}

/* Kind of a being*/
type BeingKind struct {
    Entity
    // Talent modifiers of the species
    Talents
    // Vitals modifiers of the species
    Vitals
    Aptitudes 
    // Arts multiplier in %. If zero arts cannot be used.
    Arts        float64
    // If players can choose this or not
    Playable    bool
}


func NewBeingKind(id ID, name string) * BeingKind {
    res := new(BeingKind)
    res.ID = id;
    res.Name = name
    return res
}

/* Profession of a being */
type Profession struct {
    Entity
    // Talent modifiers of the profession
    Talents
    // Vitals modifiers of the profession
    Vitals
    Aptitudes    
}


/* Gender of a being */
type Gender struct {
    Entity
    // Talent modifiers of the gender
    Talents
    // Vitals modifiers of the gender
    Vitals
}

/* Struct, list and map of all genders in WOE. */
var Genders = struct {
    Female          Gender
    Male            Gender
    Intersex        Gender 
    Genderless      Gender
}{
    Gender {
    Entity : Entity{ ID: "female", Name: "Female"},
    Talents : Talents { Agility: 1, Charisma: 1 }, 
    }, 

    Gender { 
    Entity : Entity{ ID: "male", Name: "Male"},
    Talents : Talents { Strength: 1, Intelligence: 1 },
    }, 

    Gender {
    Entity : Entity{ ID: "intersex", Name: "Intersex"},
    Talents : Talents { Dexterity: 1, Wisdom: 1 }, 
    },

    Gender {
    Entity : Entity{ ID: "genderless", Name: "Genderless"},
    Talents : Talents { Toughness: 1, Wisdom: 1 }, 
    },
    
}


var GenderList  = []*Gender{&Genders.Female, &Genders.Male, 
    &Genders.Intersex, &Genders.Genderless }

var GenderMap = map[ID]*Gender {
    Genders.Female.ID       : &Genders.Female,
    Genders.Male.ID         : &Genders.Male,
    Genders.Intersex.ID     : &Genders.Intersex,
    Genders.Genderless.ID   : &Genders.Genderless,
}



/* All BeingKinds of  WOE */
var BeingKindList = [] BeingKind {  
    {   Entity: Entity {
            ID: "kin_", Name: "Human",
            Short: "The primordial conscious beings on Earth.",
            Long: 
            `Humans are the primordial kind of conscious beings on Earth. The excel at nothing in particular, but are fast learners.`,      
        },
        // No talents because humans get no talent bonuses
        // No stats either because no bonuses there either.
        Arts : 1.0,
        Playable : true,
    },
}











type LabeledPointer struct {
    ID ID
    labeled * Labeled
}

type GenderPointer struct {
    ID ID
    gender * Gender
}

//}

/* Vital statistic of a Being. */
type Vital struct {
    Now int `xml:"N,attr"`
    Max int `xml:"X,attr"`
}

/* Report a vital statistic as a Now/Max string */
func (me * Vital) ToNowMax() string {
    return fmt.Sprintf("%4d/%4d", me.Now , me.Max)
}

// alias of the above, since I'm lazy at times
func (me * Vital) TNM() string {
    return me.ToNowMax()
}


/* Report a vital statistic as a rounded percentage */
func (me * Vital) ToPercentage() string {
    percentage := (me.Now * 100) / me.Max 
    return fmt.Sprintf("%d", percentage)
}

/* Report a vital statistic as a bar of characters */
func (me * Vital) ToBar(full string, empty string, length int) string {
    numfull := (me.Now * length) / me.Max
    numempty := length - numfull
    return strings.Repeat(empty, numempty) + strings.Repeat(full, numfull)  
}


type Talents struct {
    Strength        int     `xml:"Talents>STR,omitempty"`
    Toughness       int     `xml:"Talents>TOU,omitempty"`
    Agility         int     `xml:"Talents>AGI,omitempty"`
    Dexterity       int     `xml:"Talents>DEX,omitempty"`
    Intelligence    int     `xml:"Talents>INT,omitempty"`
    Wisdom          int     `xml:"Talents>WIS,omitempty"`
    Charisma        int     `xml:"Talents>CHA,omitempty"`
    Essence         int     `xml:"Talents>ESS,omitempty"`
}
  

type Vitals struct {
    HP              Vital
    MP              Vital
    JP              Vital
    LP              Vital
}

type EquipmentValues struct {
    Offense         int
    Protection      int
    Block           int
    Rapidity        int
    Yield           int
}


type Being struct {
    Entity
    
    // Essentials
    Gender          
    BeingKind    
    Profession
    Level           int
    
    // A being has talents.
    Talents
    // A being has vitals
    Vitals
    // A being has Equipment values
    EquipmentValues
    // A being has aptitudes
    Aptitudes
    
    // Skills array
    // Skills       []Skill
       
    // Arts array
    // Arts         []Art
    // Affects array
    // Affects      []Affect
       
    // Equipment
    // Equipment
    
    // Inventory
    Inventory         Inventory
    
    // Location pointer
    room            * Room
    
}

// Derived stats 
func (me *Being) Force() int {
    return (me.Strength * 2 + me.Wisdom) / 3
}
    
func (me *Being) Vitality() int {
    return (me.Toughness * 2 + me.Charisma) / 3
}

func (me *Being) Quickness() int {
    return (me.Agility * 2 + me.Intelligence) / 3
}

func (me * Being) Knack() int {
    return (me.Dexterity * 2 + me.Essence) / 3
}
    
    
func (me * Being) Understanding() int {
    return (me.Intelligence * 2 + me.Toughness) / 3
}

func (me * Being) Grace() int { 
    return (me.Charisma * 2 + me.Agility) / 3
}
    
func (me * Being) Zeal() int {
    return (me.Wisdom * 2 + me.Strength) / 3
}


func (me * Being) Numen() int {
      return (me.Essence * 2 + me.Dexterity) / 3
}

// Generates a prompt for use with the being/character
func (me * Being) ToPrompt() string {
    if me.Essence > 0 {
        return fmt.Sprintf("HP:%s MP:%s JP:%s LP:%s", me.HP.TNM(), me.MP.TNM(), me.JP.TNM, me.LP.TNM())
    } else {
        return fmt.Sprintf("HP:%s MP:%s LP:%s", me.HP.TNM(), me.MP.TNM(), me.LP.TNM())
    }
}


// Generates an overview of the essentials of the being as a string.
func (me * Being) ToEssentials() string {
    return fmt.Sprintf("%s lvl %d %s %s %s", me.Name, me.Level, me.Gender.Name, me.BeingKind.Name, me.Profession.Name)
}

// Generates an overview of the physical talents of the being as a string.
func (me * Being) ToBodyTalents() string {
    return fmt.Sprintf("STR: %3d    TOU: %3d    AGI: %3d    DEX: %3d", me.Strength, me.Toughness, me.Agility, me.Dexterity)
}

// Generates an overview of the mental talents of the being as a string.
func (me * Being) ToMindTalents() string {
    return fmt.Sprintf("INT: %3d    WIS: %3d    CHA: %3d    ESS: %3d", me.Intelligence, me.Wisdom, me.Charisma, me.Essence)
}

// Generates an overview of the equipment values of the being as a string.
func (me * Being) ToEquipmentValues() string {
    return fmt.Sprintf("OFF: %3d    PRO: %3d    BLO: %3d    RAP: %3d    YIE: %3d", me.Offense, me.Protection, me.Block, me.Rapidity, me.Yield)
}

// Generates an overview of the status of the being as a string.
func (me * Being) ToStatus() string {
    status := me.ToEssentials()
    status += "\n" + me.ToBodyTalents();
    status += "\n" + me.ToMindTalents();
    status += "\n" + me.ToEquipmentValues();
    status += "\n" + me.ToPrompt();
    status += "\n"
      return status
}

func (me *Being) Type() string {
    return "being"
}

func (me *Being) Save(datadir string) {
    SaveSavable(datadir, me)
}

func LoadBeing(datadir string, nameid string) * Being {    
    res, _  := LoadLoadable(datadir, nameid, new(Being)).(*Being)
    return res
}






