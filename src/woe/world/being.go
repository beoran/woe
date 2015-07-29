package world

import (
    "fmt"
    "strings"
)

type BeingKind int

const (
    BEING_KIND_NONE BeingKind = iota
    BEING_KIND_CHARACTER
    BEING_KIND_MANTUH
    BEING_KIND_NEOMAN
    BEING_KIND_HUMAN    
    BEING_KIND_CYBORG
    BEING_KIND_ANDROID    
    BEING_KIND_NONCHARACTER
    BEING_KIND_MAVERICK
    BEING_KIND_ROBOT
    BEING_KIND_BEAST
    BEING_KIND_SLIME
    BEING_KIND_BIRD
    BEING_KIND_REPTILE
    BEING_KIND_FISH
    BEING_KIND_CORRUPT    
    BEING_KIND_NONHUMAN
)

func (me BeingKind) ToString() string {
    switch me {
        case BEING_KIND_NONE:       return "None"
        case BEING_KIND_CHARACTER:  return "Character"
        case BEING_KIND_MANTUH:     return "Mantuh"
        case BEING_KIND_NEOMAN:     return "Neoman"
        case BEING_KIND_HUMAN:      return "Human"    
        case BEING_KIND_CYBORG:     return "Cyborg"
        case BEING_KIND_ANDROID:    return "Android"
        case BEING_KIND_NONCHARACTER: return "Non Character"
        case BEING_KIND_MAVERICK:   return "Maverick"
        case BEING_KIND_ROBOT:      return "Robot"
        case BEING_KIND_BEAST:      return "Beast"
        case BEING_KIND_SLIME:      return "Slime"
        case BEING_KIND_BIRD:       return "Bird"
        case BEING_KIND_REPTILE:    return "Reptile"
        case BEING_KIND_FISH:       return "Fish"
        case BEING_KIND_CORRUPT:    return "Corrupted"
        default: return ""
    }
    return ""
}


type BeingProfession int

const (
    BEING_PROFESSION_NONE       BeingProfession = iota
    BEING_PROFESSION_OFFICER
    BEING_PROFESSION_WORKER
    BEING_PROFESSION_ENGINEER
    BEING_PROFESSION_HUNTER
    BEING_PROFESSION_SCHOLAR
    BEING_PROFESSION_MEDIC
    BEING_PROFESSION_CLERIC  
    BEING_PROFESSION_ROGUE
)

func (me BeingProfession) ToString() string {
    return ""
}


/* Vital statistic of a Being. */
type Vital struct {
    Now int
    Max int
}

/* Report a vital statistic as a Now/Max string */
func (me * Vital) ToNowMax() string {
    return fmt.Sprintf("%d/%d", me.Now , me.Max)
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



type Being struct {
    Entity
    
    // Essentials
    Kind            BeingKind
    Profession      BeingProfession
    Level           int
    
    
    // Talents
    Strength        int
    Toughness       int
    Agility         int
    Dexterity       int
    Intelligence    int
    Wisdom          int
    Charisma        int
    Essence         int
        
    // Vitals
    HP              Vital
    MP              Vital
    JP              Vital
    LP              Vital
        
    // Equipment values
    Offense         int
    Protection      int
    Block           int
    Rapidity        int
    Yield           int

    // Skills array
    // Skills       []Skill
       
    // Arts array
    // Arts         []Art
    // Affects array
    // Affects      []Affect
       
    // Equipment
    // Equipment
    
    // Inventory
    // Inventory 
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
    return fmt.Sprintf("%s %d %s %s", me.Name, me.Level, me.Kind, me.Profession)
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



