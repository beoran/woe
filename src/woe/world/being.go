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

/* Kind of a being or "Kin" for short*/
type Kin struct {
    Entity
    // Talent modifiers of the species
    Talents
    // Vitals modifiers of the species
    Vitals
    Aptitudes 
    // Arts multiplier in %. If zero arts cannot be used.
    Arts        float64
    // Technique multiplier in %. If zero techniques cannot be used.
    Techniques  float64
    // Learning speed in % 
    Learning    float64
    // How much of the being is mechanical. 
    // Affects healing arts and medication, and the ability to install parts.
    Mechanical  float64
    // Level of Corruption
    Corruption  float64
    // If players can choose this or not
    Playable    bool
}


func NewKin(id ID, name string) * Kin {
    res := new(Kin)
    res.ID = id;
    res.Name = name
    return res
}

/* Job of a being */
type Job struct {
    Entity
    // Talent modifiers of the profession
    Talents
    // Vitals modifiers of the profession
    Vitals
    // Map of skills that this job starts with and their levels.
    Skills      map[ID] int
    // Same for arts and techniques and exploits
    Arts        map[ID] int
    Techniques  map[ID] int
    Exploits    map[ID] int
    // if a player can choose this or not    
    Playable    bool
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
    Talents : Talents { Toughness: 1, Emotion: 1 }, 
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



/* All Kins of  WOE */
var KinList = [] Kin {  
    {   Entity: Entity {
            ID: "kin_human", Name: "Human",
            Short: "The primordial conscious beings on Earth.",
            Long: 
`Humans are the primordial kind of conscious beings on Earth. 
They excel at nothing in particular, but are fast learners.`,      
        },
        // No talents because humans get no talent bonuses
        // No stats either because no bonuses there either.
        Arts : 1.0,
        Playable : true,
    },
    {   Entity: Entity {
            ID: "kin_neosa", Name: "Neosa",
            Short: "Nimble beings, skilled in the Arts.",
            Long: 
`Neosa are descendents of humans genetically modified to be lite, agile 
and skilled with the Numen Arts. They are less tough and strong, and less adept
with techniques`,
        },
        
        // AGI+1 EMO+1 STR-1 TOU-1 
        Talents : Talents { Strength : -1, Toughness: -1, 
            Agility : 1, Emotion : 1, },
        Arts        : 1.2,
        Techniques  : 0.8,
        Playable    : true,
    },
    {   Entity: Entity {
            ID: "kin_mantu", Name: "Mantu",
            Short: "Hardy, stocky beings, skilled in the Techniques.",
            Long: 
`Mantu are descendents of humans genetically modified to be hardy, stocky 
and skilled with Techniques. They are somewhat less agine and less adept with 
Numen Arts than humans.`,
        },
        // STR+1 1 TOU+1  AGI-1 EMO-1 
        Talents : Talents { Strength : +1, Toughness: +1, 
            Agility : -1, Emotion : -1, },
        
        Arts : 0.8,
        Techniques : 1.2,
        Playable : true,
    },
    {   Entity: Entity {
            ID: "kin_cyborg", Name: "Cyborg",
            Short: "Human enhanced with robotic parts. ",
            Long: 
`Cyborgs are humans who either through neccesity or through their own will, 
have been enhanced with robotic parts. They are far more skilled with 
Techniques, and can install some Parts, but their Nummen arts are only half 
as effective. They are partially mechanical and healing arts and medication 
is not as effective on them, but they can be repaired.`,
        },
        // STR+1 1 TOU+1 DEX+1 INT+1 
        Talents : Talents { Strength : +1, Toughness: +1, 
            Dexterity : +1, Intelligence: +1, },
        Arts : 0.5,
        Techniques : 1.5,
        Mechanical : 0.5,
        Playable : true,
    },
    {   Entity: Entity {
            ID: "kin_android", Name: "Android",
            Short: "Human shaped biomchanical robot at the service of humans. ",
            Long: 
`Androids are conscious human shaped robots with the imperative to serve humans.  
Highly effective with Techniques, and can install many Parts, but cannot use 
any Nummen arts. Since thay are not alive, they technically cannot die.`,
        },
        // STR+1 1 TOU+1 DEX+1 INT+1 
        Talents : Talents { Strength : +2, Toughness: +2, 
            Dexterity : +2, Intelligence: +2, },
        Arts : 0.0,
        Techniques : 2.0,
        Mechanical : 1.0,
        Playable : true,
    },
    
    {   Entity: Entity {
            ID: "kin_maverick", Name: "Maverick",
            Short: "Human shaped biomechanical robot running wild. ",
            Long: 
`Mavericks are androids in which the imperative to serve humans has 
been destroyed or disabled.  Highly effective with Techniques, and can install 
many Parts, but cannot use any Nummen arts. Since thay are not alive, they 
technically cannot die.  They are feared by Humans and hated by Androids.`,
        },
        // STR+1 1 TOU+1 DEX+1 INT+1 
        Talents : Talents { Strength : +3, Toughness: +3, 
            Dexterity : +2, Intelligence: +2, Charisma: -2 },
        Arts : 0.0,
        Techniques : 2.0,
        Mechanical : 1.0,
        Playable : false,
    },
    
    
    {   Entity: Entity {
            ID: "kin_robot", Name: "Robot",
            Short: "Non conscious mechanical robot.",
            Long: 
`In the wars of the past many robots were built for offense or defense. 
Unfortunately, they are self repairing and often even able to replicate 
if they find suitable materials. No wonder they are still prowling
the Earth millennia after.`,
        },
        // STR+1 1 TOU+1 DEX+1 INT+1 
        Talents : Talents { Strength : +4, Toughness: +4, 
            Dexterity : +2, Intelligence: +2, Charisma: -4},
        Arts : 0.0,
        Techniques : 2.0,
        Mechanical: 1.0,
        Playable : false,
    },
    
    {   Entity: Entity {
            ID: "kin_drone", Name: "Drone",
            Short: "Flying combat robot. ",
            Long: 
`Out of control robots are a pain, out of control flying robots even more so!
They might be less though than normal robots, but they move extremely quickly.
`,
        },
        Talents : Talents { Strength : +2, Toughness: +2, 
            Agility: +4, Dexterity : +2, Intelligence: +2, Charisma: -4},
        Arts : 0.0,
        Techniques : 2.0,
        Mechanical : 1.0,
        Playable : false,

    },
    {   Entity: Entity {
            ID: "kin_turret", Name: "Turret",
            Short: "Immobile automated defense system. ",
            Long: 
`
The ancients would set up robotic defense system to guard certain areas.
These defense systems might be immobile, but they are deadly accurate.
Furthermore they are extremely resillient and self repairing. 
No wonder they are still actve after all these years.
`,
        },
        Talents : Talents { Strength : +2, Toughness: +4, 
            Agility: -4, Dexterity : +4, Intelligence: +4, Charisma: -4},
        Arts : 0.0,
        Techniques : 2.0,
        Mechanical : 1.0,
        Playable : false,
    },
    
    {   Entity: Entity {
            ID: "kin_beast", Name: "Beast",
            Short: "Beast that prowls the wild. ",
            Long: 
`Due to the damage to the ecosystem, beasts have rapidly evolved in the last 
60000 years. As a result, most all of them, even the plant eaters, are ferocious 
and aggressive, to protect themselves and their offspring from Humans.`,
        },
        Talents : Talents { Strength : +2, Toughness: +2, 
            Agility : +1, Intelligence: -5, },
    },

    {   Entity: Entity {
            ID: "kin_bird", Name: "Bird",
            Short: "Flying being that prowls the wild. ",
            Long: 
`Beasts can be dangerous, flying beasts are all the more so! 
They might be less resillient, but all the more agile.`,
        },
        Talents : Talents { Strength : +1, Toughness: +1, 
            Agility : +3, Intelligence: -5, },
    },

    {   Entity: Entity {
            ID: "kin_fish", Name: "Fish",
            Short: "Fish like being that swims the seas or rivers. ",
            Long: 
`Now you know why you were always told not to swim in rivers or seas.
Fish dart through the water, attacking with their razor sharp teeth.
`,
        },
        Talents : Talents { Strength : +3, Toughness: +2, 
            Agility : +2, Intelligence: -5, },
    },

    {   Entity: Entity {
            ID: "kin_amphibian", Name: "Amphibian",
            Short: "Being that lives both on land and in the water. ",
            Long: 
`Covered with a slimy skin and often toxic, these being can not only swim 
but also purse you on land.`,
        },
        Talents : Talents { Strength : +1, Toughness: +2, 
            Agility : +1, Dexterity: +1, Intelligence: -5, },
    },

    {   Entity: Entity {
            ID: "kin_reptile", Name: "Reptile",
            Short: "Scaly creepy crawling beasts.",
            Long: 
`Reptiles have been around for a long time, and it looks they will be around for 
a long time still. They may be slow, especially in colder weather, nevertheless
they remain dangerous.
`,
        },
        Talents : Talents { Strength : +2, Toughness: +2, 
            Agility: -1, Intelligence: -5, },
    },
    {   Entity: Entity {
            ID: "kin_crustacean", Name: "Custacian",
            Short: "Beast protected by a tough shell",
            Long: 
`You might find it hard to inflict any damage to these well armoured beings.
Their shells protect them against damage and allow them to live in the water
and as well on the land.`,
        },
        Talents : Talents { Strength : +2, Toughness: +4, 
                            Intelligence: -5,  },
    },
    {   Entity: Entity {
            ID: "kin_insect", Name: "Insect",
            Short: "Beast with articulated legs and bodies.",
            Long: 
`The climate of Earth hs shifted dramaticaly over the last 60000 years,
and as a result, larger creepy crawlers became more successful.
`,
        },
        Talents : Talents { Strength : +2, Toughness: +4,  Intelligence: -5, },
    },


    {   Entity: Entity {
            ID: "kin_aquatic", Name: "beast",
            Short: "Aquatic beast. ",
            Long: 
`Whether in the rivers or the deep seas, these soft bodies creatures 
are often toxic and quite dangerous.`,
        },
        Talents : Talents { Strength : +2, Dexterity: +2, Intelligence: -5, },
    },
    {   Entity: Entity {
            ID: "kin_corrupted", Name: "corrupted",
            Short: "Beast corupted by Omen. ",
            Long: 
`Some animals became corrupted by Omen. As a result, they became much stronger 
and more resillient. Fortunately, the Omen is weak against certain Numen arts.
Beware, their attacks might be contageous...`,
        },
        Talents : Talents { Strength : +4, Toughness: +4, 
            Agility : +1, Intelligence: -3, Wisdom: -5 },
    },
    {   Entity: Entity {
            ID: "kin_deceased", Name: "Deceased",
            Short: "Deceased biological being animated by Omen. ",
            Long: 
`Some living beings become corrupted by Omen to the point that they remain
animated even after their biological bodies have already stopped functioning.
Such beings are termed the Deceased. They are resillient, strong and cunning. 
Fortunately, the Omen is weak against certain Numen arts. But beware, their 
attacks might be contageous...
`,
        },
        Talents : Talents { Strength : +2, Toughness: +2, 
            Agility : +1, Intelligence: -1, Wisdom: -7 },
    },
    

}



/* All jobs of WOE 
 * agent        officer     guardian
 * worker       brawler     builder
 * hunter       gunsman     rogue
 * explorer     ranger      rebel
 * tinker       engineer    wrecker 
 * homemaker    musician    trader
 * scholar      scientist   hacker
 * medic        cleric      artist
 * esper        dilettante  (not playable) Danger
 * 
 * 
 * hunter scholar esper worker medic agent officer cleric guardian ranger wrecker engineer tinker scientist 
 * 
 * 
 * 

Officer                 STR + 1
Worker                  TOU + 1
Engineer                DEX + 1
Hunter                  AGI + 1
Scholar                 INT + 1
Doctor                  WIS +1
Cleric                  CHA + 1

 * 
 * 
*/
var JobList = [] Job {  
    {   Entity: Entity {
            ID: "job_hunter", Name: "Hunter",
            Short: "Hunter who chases beasts and mavericks.",
            Long: 
`Hunters protect human settlements from Beasts and Mavericks, and try to keep their numbers down.
`,      
        },
        Talents : Talents { Dexterity: +2 },            
        Skills : map[ID] int{ "skill_guns": 10,  },
        Playable: true,
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
    Emotion         int     `xml:"Talents>EMO,omitempty"`
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
    Kin    
    Job
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
    return (me.Dexterity * 2 + me.Emotion) / 3
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
      return (me.Emotion * 2 + me.Dexterity) / 3
}

// Generates a prompt for use with the being/character
func (me * Being) ToPrompt() string {
    if me.Emotion > 0 {
        return fmt.Sprintf("HP:%s MP:%s JP:%s LP:%s", me.HP.TNM(), me.MP.TNM(), me.JP.TNM, me.LP.TNM())
    } else {
        return fmt.Sprintf("HP:%s MP:%s LP:%s", me.HP.TNM(), me.MP.TNM(), me.LP.TNM())
    }
}


// Generates an overview of the essentials of the being as a string.
func (me * Being) ToEssentials() string {
    return fmt.Sprintf("%s lvl %d %s %s %s", me.Name, me.Level, me.Gender.Name, me.Kin.Name, me.Job.Name)
}

// Generates an overview of the physical talents of the being as a string.
func (me * Being) ToBodyTalents() string {
    return fmt.Sprintf("STR: %3d    TOU: %3d    AGI: %3d    DEX: %3d", me.Strength, me.Toughness, me.Agility, me.Dexterity)
}

// Generates an overview of the mental talents of the being as a string.
func (me * Being) ToMindTalents() string {
    return fmt.Sprintf("INT: %3d    WIS: %3d    CHA: %3d    EMO: %3d", me.Intelligence, me.Wisdom, me.Charisma, me.Emotion)
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






