package world

import (
	"fmt"
	"math"
	"strings"

	"github.com/beoran/woe/monolog"
	"github.com/beoran/woe/sitef"
)

/* Aptitudes of a being, species or profession */
type Aptitudes struct {
	Skills     []BeingSkill
	Arts       []BeingArt
	Techniques []BeingTechnique
	Exploits   []BeingExploit
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
	Arts float64
	// Technique multiplier in %. If zero techniques cannot be used.
	Techniques float64
	// Learning speed in %
	Learning float64
	// How much of the being is mechanical.
	// Affects healing arts and medication, and the ability to install parts.
	Mechanical float64
	// Level of Corruption
	Corruption float64
	// Level of "unliving", i,e, previously alive, matter in the being
	Unlife float64
}

func NewKin(id string, name string) *Kin {
	res := new(Kin)
	res.ID = id
	res.Name = name
	return res
}

/* Kin modifier display as a string */
func (me Kin) ToKinModifiers() string {
	return fmt.Sprintf("ARTS : %4.2f TECHS: %4.2f LEARN: %4.2f\n"+
		"MECHA: %4.2f CORRU: %4.2f UNLIF: %4.2f",
		me.Arts, me.Techniques, me.Learning, me.Mechanical,
		me.Corruption, me.Unlife)
}

/* Help for a kin */
func (me Kin) AskLong() string {
	return me.Long + "\n\nTalent modifiers for " + me.Name +
		" kin:\n" + me.ToTalents() + "\n\nKin modifiers for " + me.Name +
		" kin:\n" + me.ToKinModifiers() + "\n"
}

/* Job of a being */
type Job struct {
	Entity
	// Talent modifiers of the profession
	Talents
	// Vitals modifiers of the profession
	Vitals
	// Map of skills that this job starts with and their levels.
	Skills map[string]int
	// Same for arts and techniques and exploits
	Arts       map[string]int
	Techniques map[string]int
	Exploits   map[string]int
	// if a player can choose this or not
	Playable bool
}

/* Help for a job */
func (me Job) AskLong() string {
	return me.Long + "\n\nTalent modifiers for " + me.Name +
		" job:\n" + me.ToTalents() + "\n"
}

/* Gender of a being */
type Gender struct {
	Entity
	// Talent modifiers of the gender
	Talents
	// Vitals modifiers of the gender
	Vitals
}

/* Help for a gender */
func (me Gender) AskLong() string {
	return me.Long + "\n\nGender modifiers for " + me.Name +
		" gender:\n" + me.ToTalents() + "\n"
}

/* Struct, list and map of all genders in WOE. */
var Genders = struct {
	Female     Gender
	Male       Gender
	Intersex   Gender
	Genderless Gender
}{
	Gender{
		Entity: Entity{ID: "gender_female", Name: "Female",
			Short: "Female gender",
			Long:  "No matter the day and age, there are still plenty of beings who are female. Females are slighly more agile and charismatic than those who have another gender.",
		},
		Talents: Talents{Agility: 1, Charisma: 1},
	},

	Gender{
		Entity: Entity{ID: "gender_male", Name: "Male",
			Short: "Male gender",
			Long:  "No matter the day and age, there are still plenty of beings who are male. Males are slighly more strong and studious than those who have another gender.",
		},
		Talents: Talents{Strength: 1, Intelligence: 1},
	},

	Gender{
		Entity: Entity{ID: "gender_intersex", Name: "Intersex",
			Short: "Intersexed",
			Long:  "Not every being can be clearly defined as being male or female. Sometimes, certain beings end up characteristics of both. Intersexed are slighly more dexterous and wise than those who have another gender.",
		},
		Talents: Talents{Dexterity: 1, Wisdom: 1},
	},

	Gender{
		Entity: Entity{ID: "gender_none", Name: "Genderless",
			Short: "No gender",
			Long:  "Some beings lack reproductive sytems and are therefore genderless. Genderless are slighly more tough and emotionally balanced than those who have another gender.",
		},
		Talents: Talents{Toughness: 1, Emotion: 1},
	},
}

func EntitylikeToGender(me Entitylike) *Gender {
	v, ok := me.(*Gender)
	if ok {
		return v
	} else {
		return nil
	}
}

type EntitylikeSlice []Entitylike

var GenderList = []Gender{Genders.Female, Genders.Male,
	Genders.Intersex, Genders.Genderless}

var GenderEntityList = EntitylikeSlice{&Genders.Female, &Genders.Male,
	&Genders.Intersex, &Genders.Genderless}

var GenderMap = map[string]*Gender{
	Genders.Female.ID:     &Genders.Female,
	Genders.Male.ID:       &Genders.Male,
	Genders.Intersex.ID:   &Genders.Intersex,
	Genders.Genderless.ID: &Genders.Genderless,
}

type EntityIterator interface {
	Each(cb func(Entitylike) Entitylike) Entitylike
}

func (me EntitylikeSlice) Each(cb func(Entitylike) Entitylike) Entitylike {
	for i := 0; i < len(me); i++ {
		res := cb(me[i])
		if res != nil {
			return res
		}
	}
	return nil
}

func (me EntitylikeSlice) Filter(cb func(Entitylike) Entitylike) EntitylikeSlice {
	result := make(EntitylikeSlice, 0)
	for i := 0; i < len(me); i++ {
		res := cb(me[i])
		if res != nil {
			result = append(result, res)
		}
	}
	return result
}

/* Finds the name irrespecful of the case */
func (me EntitylikeSlice) FindName(name string) Entitylike {
	return me.Each(func(e Entitylike) Entitylike {
		if strings.ToLower(e.AsEntity().Name) == strings.ToLower(name) {
			return e
		} else {
			return nil
		}
	})
}

/* Finds the ID  */
func (me EntitylikeSlice) FindID(id string) Entitylike {
	if id == "" {
		return nil
	}

	return me.Each(func(e Entitylike) Entitylike {
		if strings.ToLower(e.AsEntity().ID) == id {
			return e
		} else {
			return nil
		}
	})
}

/* Filters the list by privilege level (only those allowed by the level are retained) */
func (me EntitylikeSlice) FilterPrivilege(privilege Privilege) EntitylikeSlice {
	return me.Filter(func(e Entitylike) Entitylike {
		if e.AsEntity().Privilege <= privilege {
			return e
		} else {
			return nil
		}
	})
}

/* All Kins of  WOE */
var KinList = []Kin{
	Kin{
		Entity: Entity{
			ID: "kin_human", Name: "Human",
			Short: "The primordial conscious beings on Earth.",
			Long: `Humans are the primordial kind of conscious beings on Earth. 
They excel at nothing in particular, but are fast learners.`,
		},
		// No talents because humans get no talent bonuses
		// No stats either because no bonuses there either.
		Arts:       1.0,
		Techniques: 1.0,
		Learning:   1.5,
	},
	{
		Entity: Entity{
			ID: "kin_neosa", Name: "Neosa",
			Short: "Nimble beings, skilled in the Arts.",
			Long: `Neosa are descendents of humans genetically modified to be lite, agile 
and skilled with the Numen Arts. They are less tough and strong, and less adept
with techniques. They can be recognized by their long and pointy ears.`,
		},

		// AGI+1 EMO+1 STR-1 TOU-1
		Talents: Talents{Strength: -1, Toughness: -1,
			Agility: 1, Emotion: 1},
		Arts:       1.2,
		Techniques: 0.8,
		Learning:   1.0,
	},
	{
		Entity: Entity{
			ID: "kin_mantu", Name: "Mantu",
			Short: "Hardy, stocky, furry beings, skilled in the Techniques.",
			Long: `Mantu are descendents of humans genetically modified to be hardy, stocky 
and skilled with Techniques. They are somewhat less agine and less adept with 
Numen Arts than humans. They have a soft fur which covers their whole body.`,
		},
		// STR+1 1 TOU+1  AGI-1 EMO-1
		Talents: Talents{Strength: +1, Toughness: +1,
			Agility: -1, Emotion: -1},

		Arts:       0.8,
		Techniques: 1.2,
		Learning:   1.0,
	},
	{
		Entity: Entity{
			ID: "kin_cyborg", Name: "Cyborg",
			Short: "Human enhanced with robotic parts. ",
			Long: `Cyborgs are humans who either through neccesity or through their own will, 
have been enhanced with robotic parts. They are far more skilled with 
Techniques, and can install some Parts, but their Nummen arts are only half 
as effective. They are partially mechanical and healing arts and medication 
is not as effective on them, but they can be repaired.`,
		},
		// STR+1 1 TOU+1 DEX+1 INT+1
		Talents: Talents{Strength: +1, Toughness: +1,
			Dexterity: +1, Intelligence: +1},
		Arts:       0.5,
		Techniques: 1.5,
		Mechanical: 0.5,
		Learning:   1.1,
	},
	{
		Entity: Entity{
			ID: "kin_android", Name: "Android",
			Short: "Human shaped biomchanical robot at the service of humans. ",
			Long: `Androids are conscious human shaped robots with the imperative to serve humans.  
Highly effective with Techniques, and can install many Parts, but cannot use 
any Numen arts. Since thay are not alive, they technically cannot die.`,
		},
		// STR+1 1 TOU+1 DEX+1 INT+1
		Talents: Talents{Strength: +2, Toughness: +2,
			Dexterity: +2, Intelligence: +2},
		Arts:       0.0,
		Techniques: 2.0,
		Mechanical: 1.0,
		Learning:   1.0,
	},

	{
		Entity: Entity{
			ID: "kin_maverick", Name: "Maverick",
			Short: "Human shaped biomechanical robot running wild. ",
			Long: `Mavericks are androids in which the imperative to serve humans has 
been destroyed or disabled.  Highly effective with Techniques, and can install 
many Parts, but cannot use any Numen arts. Since thay are not alive, they 
technically cannot die.  They are feared by Humans and hated by Androids.`,
			Privilege: PRIVILEGE_IMPLEMENTOR,
		},
		// STR+1 1 TOU+1 DEX+1 INT+1
		Talents: Talents{Strength: +3, Toughness: +3,
			Dexterity: +2, Intelligence: +2, Charisma: -2},
		Arts:       0.0,
		Techniques: 2.0,
		Mechanical: 1.0,
		Learning:   1.0,
	},

	{
		Entity: Entity{
			ID: "kin_robot", Name: "Robot",
			Short: "Non conscious mechanical robot.",
			Long: `In the wars of the past many robots were built for offense or defense. 
Unfortunately, they are self repairing and often even able to replicate 
if they find suitable materials. No wonder they are still prowling
the Earth millennia after.`,
			Privilege: PRIVILEGE_IMPLEMENTOR,
		},
		// STR+1 1 TOU+1 DEX+1 INT+1
		Talents: Talents{Strength: +4, Toughness: +4,
			Dexterity: +2, Intelligence: +2, Charisma: -4},
		Arts:       0.0,
		Techniques: 2.0,
		Mechanical: 1.0,
		Learning:   1.0,
	},

	{
		Entity: Entity{
			ID: "kin_drone", Name: "Drone",
			Short: "Flying combat robot. ",
			Long: `Out of control robots are a pain, out of control flying robots even more so!
They might be less though than normal robots, but they move extremely quickly.
`,
			Privilege: PRIVILEGE_IMPLEMENTOR,
		},
		Talents: Talents{Strength: +2, Toughness: +2,
			Agility: +4, Dexterity: +2, Intelligence: +2, Charisma: -4},
		Arts:       0.0,
		Techniques: 2.0,
		Mechanical: 1.0,
		Learning:   1.0,
	},

	{
		Entity: Entity{
			ID: "kin_turret", Name: "Turret",
			Short: "Immobile automated defense system. ",
			Long: `
The ancients would set up robotic defense system to guard certain areas.
These defense systems might be immobile, but they are deadly accurate.
Furthermore they are extremely resillient and self repairing. 
No wonder they are still active after all these years.
`,
			Privilege: PRIVILEGE_IMPLEMENTOR,
		},
		Talents: Talents{Strength: +2, Toughness: +4,
			Agility: -4, Dexterity: +4, Intelligence: +4, Charisma: -4},
		Arts:       0.0,
		Techniques: 2.0,
		Mechanical: 1.0,
		Learning:   1.0,
	},

	{
		Entity: Entity{
			ID: "kin_beast", Name: "Beast",
			Short: "Beast that prowls the wild. ",
			Long: `Due to the damage to the ecosystem, beasts have rapidly evolved in the last 
60000 years. As a result, most all of them, even the plant eaters, are ferocious 
and aggressive, to protect themselves and their offspring from Humans.`,
			Privilege: PRIVILEGE_IMPLEMENTOR,
		},
		Talents: Talents{Strength: +2, Toughness: +2,
			Agility: +1, Intelligence: -5},
		Arts:       1.0,
		Techniques: 1.0,
		Learning:   1.0,
	},

	{
		Entity: Entity{
			ID: "kin_bird", Name: "Bird",
			Short: "Flying being that prowls the wild. ",
			Long: `Beasts can be dangerous, flying beasts are all the more so! 
They might be less resillient, but all the more agile.`,
			Privilege: PRIVILEGE_IMPLEMENTOR,
		},
		Talents: Talents{Strength: +1, Toughness: +1,
			Agility: +3, Intelligence: -5},
		Arts:       1.0,
		Techniques: 1.0,
		Learning:   1.0,
	},

	{
		Entity: Entity{
			ID: "kin_fish", Name: "Fish",
			Short: "Fish like being that swims the seas or rivers. ",
			Long: `Now you know why you were always told not to swim in rivers or seas.
Fish dart through the water, attacking with their razor sharp teeth.
`,
			Privilege: PRIVILEGE_IMPLEMENTOR,
		},
		Talents: Talents{Strength: +3, Toughness: +2,
			Agility: +2, Intelligence: -5},
		Arts:       1.0,
		Techniques: 1.0,
		Learning:   1.0,
	},

	{
		Entity: Entity{
			ID: "kin_amphibian", Name: "Amphibian",
			Short: "Being that lives both on land and in the water. ",
			Long: `Covered with a slimy skin and often toxic, these being can not only swim 
but also purse you on land.`,
			Privilege: PRIVILEGE_IMPLEMENTOR,
		},
		Talents: Talents{Strength: +1, Toughness: +2,
			Agility: +1, Dexterity: +1, Intelligence: -5},
		Arts:       1.0,
		Techniques: 1.0,
		Learning:   1.0,
	},

	{
		Entity: Entity{
			ID: "kin_reptile", Name: "Reptile",
			Short: "Scaly creepy crawling beasts.",
			Long: `Reptiles have been around for a long time, and it looks they will be around for 
a long time still. They may be slow, especially in colder weather, nevertheless
they remain dangerous.
`,
			Privilege: PRIVILEGE_IMPLEMENTOR,
		},
		Talents: Talents{Strength: +2, Toughness: +2,
			Agility: -1, Intelligence: -5},
		Arts:       1.0,
		Techniques: 1.0,
		Learning:   1.0,
	},

	{
		Entity: Entity{
			ID: "kin_crustacean", Name: "Custacian",
			Short: "Beast protected by a tough shell",
			Long: `You might find it hard to inflict any damage to these well armoured beings.
Their shells protect them against damage and allow them to live in the water
and as well on the land.`,
			Privilege: PRIVILEGE_IMPLEMENTOR,
		},
		Talents: Talents{Strength: +2, Toughness: +4,
			Intelligence: -5},
		Arts:       1.0,
		Techniques: 1.0,
		Learning:   1.0,
	},

	{
		Entity: Entity{
			ID: "kin_insect", Name: "Insect",
			Short: "Beast with articulated legs and bodies.",
			Long: `The climate of Earth hs shifted dramaticaly over the last 60000 years,
and as a result, larger creepy crawlers became more successful.
`,
			Privilege: PRIVILEGE_IMPLEMENTOR,
		},
		Talents:    Talents{Strength: +2, Toughness: +4, Intelligence: -5},
		Arts:       1.0,
		Techniques: 1.0,
		Learning:   1.0,
	},

	{
		Entity: Entity{
			ID: "kin_aquatic", Name: "beast",
			Short: "Aquatic beast. ",
			Long: `Whether in the rivers or the deep seas, these soft bodies creatures 
are often toxic and quite dangerous.`,
			Privilege: PRIVILEGE_IMPLEMENTOR,
		},
		Talents:    Talents{Strength: +2, Dexterity: +2, Intelligence: -5},
		Arts:       1.0,
		Techniques: 1.0,
		Learning:   1.0,
	},

	{
		Entity: Entity{
			ID: "kin_corrupted", Name: "corrupted",
			Short: "Beast corupted by Omen. ",
			Long: `Some animals became corrupted by Omen. As a result, they became much stronger 
and more resillient. Fortunately, the Omen is weak against certain Numen arts.
Beware, their attacks might be contageous...`,
			Privilege: PRIVILEGE_IMPLEMENTOR,
		},
		Talents: Talents{Strength: +4, Toughness: +4,
			Agility: +1, Intelligence: -3, Wisdom: -5},
		Arts:       1.0,
		Techniques: 1.0,
		Corruption: 1.0,
		Learning:   1.0,
	},

	{
		Entity: Entity{
			ID: "kin_deceased", Name: "Deceased",
			Short: "Deceased biological being animated by Omen. ",
			Long: `Some living beings become corrupted by Omen to the point that they remain
animated even after their biological bodies have already stopped functioning.
Such beings are termed the Deceased. They are resillient, strong and cunning. 
Fortunately, the Omen is weak against certain Numen arts. But beware, their 
attacks might be contageous...
`,
			Privilege: PRIVILEGE_IMPLEMENTOR,
		},
		Talents: Talents{Strength: +2, Toughness: +2,
			Agility: +1, Intelligence: -1, Wisdom: -7},
		Arts:       1.0,
		Techniques: 1.0,
		Learning:   1.0,
		Corruption: 1.0,
		Unlife:     1.0,
	},
}

var KinEntityList EntitylikeSlice

func init() {
	KinEntityList = make(EntitylikeSlice, len(KinList))
	for i := range KinList {
		e := KinList[i]
		monolog.Debug("KinList: %s", e.Name)
		KinEntityList[i] = &e
	}
}

func EntitylikeToKin(me Entitylike) *Kin {
	v, ok := me.(*Kin)
	if ok {
		return v
	} else {
		return nil
	}
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
 * hunter scholar esper worker medic agent officer cleric guardian ranger
 * wrecker engineer tinker scientist
 *
 *
 *
Agent                   STR + 2
Worker                  TOU + 2
Engineer                DEX + 2
Hunter                  AGI + 2
Scholar                 INT + 2
Medic                   WIS + 2
Cleric                  CHA + 2
 *
 *
*/
var JobList = []Job{
	{Entity: Entity{
		ID: "job_agent", Name: "Agent",
		Short: "Agent employed by the government of Eruta.",
		Long: `Agents work for the government of Eruta. They are given all sorts of tasks, but tend to focus on capturing criminals. Physical strength is their forte.
`,
	},
		Talents: Talents{Strength: +2},
		Skills: map[string]int{
			"skill_sword":       10,
			"skill_heavy_gear":  10,
			"skill_acrobatics":  5,
			"skill_cannon":      5,
			"skill_weaponsmith": 5,
		},
		Playable: true,
	},
	{Entity: Entity{
		ID: "job_worker", Name: "Worker",
		Short: "Worker in construction or mining.",
		Long: `Workers take on the more heavy jobs such as construction and mining. They are a tough lot.
`,
	},
		Talents: Talents{Toughness: +2},
		Skills: map[string]int{
			"skill_maul":       10,
			"skill_heavy_gear": 10,
			"skill_toiling":    5,
			"skill_explosives": 5,
			"skill_mining":     5,
		},
		Playable: true,
	},
	{Entity: Entity{
		ID: "job_engineer", Name: "Engineer",
		Short: "Expert in machines and technology.",
		Long: `Engineers are experts in technology. They can construct and repair most any machine. They tend to be higly dexterous.
`,
	},
		Talents: Talents{Toughness: +2},
		Skills: map[string]int{
			"skill_gun":         10,
			"skill_medium_gear": 10,
			"skill_engineering": 5,
			"skill_fist":        5,
			"skill_gunsmith":    5,
		},
		Playable: true,
	},
	{Entity: Entity{
		ID: "job_hunter", Name: "Hunter",
		Short: "Hunter who chases beasts and mavericks.",
		Long: `Hunters protect human settlements from Beasts and Mavericks, and try to keep their numbers down. They excel in Agility.
`,
	},
		Talents: Talents{Agility: +2},
		Skills: map[string]int{
			"skill_polearm":     10,
			"skill_medium_gear": 10,
			"skill_shield":      5,
			"skill_gun":         5,
			"skill_survival":    5,
		},
		Playable: true,
	},
	{Entity: Entity{
		ID: "job_scholar", Name: "Scholar",
		Short: "Scholar who studies science.",
		Long: `Scolars focus in studying science and discovering what was previously unknown. Reknowned for their intelligence.
`,
	},
		Talents: Talents{Intelligence: +2},
		Skills: map[string]int{
			"skill_staff":      10,
			"skill_light_gear": 10,
			"skill_science":    5,
			"skill_lore":       5,
			"skill_artistic":   5,
		},
		Playable: true,
	},
	{Entity: Entity{
		ID: "job_medic", Name: "Medic",
		Short: "Medic who heals the wounded and cures the ill.",
		Long: `Medics focus on healing the wounded and curng the ill. Need wisdom to deal with their patients.
`,
	},
		Talents: Talents{Wisdom: +2},
		Skills: map[string]int{
			"skill_knife":      10,
			"skill_light_gear": 10,
			"skill_medical":    5,
			"skill_bravery":    5,
			"skill_gun":        5,
		},
		Playable: true,
	},
	{Entity: Entity{
		ID: "job_cleric", Name: "Cleric",
		Short: "Clerics tend to the spiritual.",
		Long:  `Clerics tend to the spiritual well being of others in the name of Lord Kei. It is a job that requires high Charisma.`,
	},
		Talents: Talents{Charisma: +2},
		Skills: map[string]int{
			"skill_fist":       10,
			"skill_light_gear": 10,
			"skill_shield":     5,
			"skill_social":     5,
			"skill_arcane":     5,
		},
		Playable: true,
	},
}

var JobEntityList EntitylikeSlice

func init() {
	JobEntityList = make(EntitylikeSlice, len(JobList))
	for i := range JobList {
		e := JobList[i]
		monolog.Debug("JobList: %s", e.Name)
		JobEntityList[i] = &e
	}
}

func EntitylikeToJob(me Entitylike) *Job {
	v, ok := me.(*Job)
	if ok {
		return v
	} else {
		return nil
	}
}

type LabeledPointer struct {
	ID      string
	labeled *Labeled
}

type GenderPointer struct {
	ID     string
	gender *Gender
}

//}

/* Vital statistic of a Being. */
type Vital struct {
	Now int `xml:"N,attr"`
	Max int `xml:"X,attr"`
}

/* Report a vital statistic as a Now/Max string */
func (me *Vital) ToNowMax() string {
	return fmt.Sprintf("%4d/%4d", me.Now, me.Max)
}

// alias of the above, since I'm lazy at times
func (me *Vital) TNM() string {
	return me.ToNowMax()
}

/* Report a vital statistic as a rounded percentage */
func (me *Vital) ToPercentage() string {
	percentage := (me.Now * 100) / me.Max
	return fmt.Sprintf("%d", percentage)
}

/* Report a vital statistic as a bar of characters */
func (me *Vital) ToBar(full string, empty string, length int) string {
	numfull := (me.Now * length) / me.Max
	numempty := length - numfull
	return strings.Repeat(empty, numempty) + strings.Repeat(full, numfull)
}

type Talents struct {
	Strength     int `xml:"Talents>STR,omitempty"`
	Toughness    int `xml:"Talents>TOU,omitempty"`
	Agility      int `xml:"Talents>AGI,omitempty"`
	Dexterity    int `xml:"Talents>DEX,omitempty"`
	Intelligence int `xml:"Talents>INT,omitempty"`
	Wisdom       int `xml:"Talents>WIS,omitempty"`
	Charisma     int `xml:"Talents>CHA,omitempty"`
	Emotion      int `xml:"Talents>EMO,omitempty"`
}

type Vitals struct {
	HP Vital
	MP Vital
	JP Vital
	LP Vital
}

type EquipmentValues struct {
	Offense    int
	Protection int
	Block      int
	Rapidity   int
	Yield      int
}

type Being struct {
	Entity

	// Essentials
	*Gender
	*Kin
	*Job
	Level int

	// A being has talents.
	Talents
	// A being has vitals
	Vitals
	// A being has Equipment values
	EquipmentValues
	// A being has aptitudes
	Aptitudes

	// Skills       map[string]BeingSkill

	// Arts array
	// Arts         []Art
	// Affects array
	// Affects      []Affect

	// Equipment
	Equipment

	// Inventory
	Inventory

	// Location pointer
	Room *Room
}

var BasicTalent Talents = Talents{
	Strength:     10,
	Toughness:    10,
	Agility:      10,
	Dexterity:    10,
	Intelligence: 10,
	Wisdom:       10,
	Charisma:     10,
	Emotion:      10,
}

// Derived stats
func (me *Talents) Force() int {
	return (me.Strength*2 + me.Wisdom) / 3
}

func (me *Talents) Vitality() int {
	return (me.Toughness*2 + me.Charisma) / 3
}

func (me *Talents) Quickness() int {
	return (me.Agility*2 + me.Intelligence) / 3
}

func (me *Talents) Knack() int {
	return (me.Dexterity*2 + me.Emotion) / 3
}

func (me *Talents) Understanding() int {
	return (me.Intelligence*2 + me.Toughness) / 3
}

func (me *Talents) Grace() int {
	return (me.Charisma*2 + me.Agility) / 3
}

func (me *Talents) Zeal() int {
	return (me.Wisdom*2 + me.Strength) / 3
}

func (me *Talents) Numen() int {
	return (me.Emotion*2 + me.Dexterity) / 3
}

/*

Stats of beings:


Talents:
*
* Talents describe the basic constitution of a being.

Strength: explosive physical strength
Toughness: physical resillience
Agility: Speed of motion and bodily balance
Dexterity: Fine motor sills and hand eye coordination
Intelligence: Book smarts and studiousness
Wisdom: Insight and intution
Emotion: Emotional control and insight.
Charisma: Social abilities and appeal to others

Force: Physical and mental vigor.
Vitality: Resistance against damage.
Quickness: Phyisical and mental speed.
Knack: Manual and mental adroitness
Understanding: Deep insight though stubborn work.
Grace: Appeal and delicacy of motion and presence.
Zeal: Mental and religious perspicacity.
Numen: Effectiveness of arts through passion and quick hands.

HP: Hull Power/ Health Power: Protection against wounds. When 0, the being is
    likely to get stunned and vulnerable against LP damage.
MP: Motion Power. Needed for motion and performing techniques.
JP: Junction Power. Needed for Numen arts.
LP: Life Power. Resilience to actual wounds. When 0, the being dies or is
    destroyed.

Offense: Effectiveness of weapon (ranged or melee).
Protection: Effectiveness of armor.
Blocking: Efectiveness of shield or parrying weapon.
Rapidity: Gain or loss of speed due to armor (may be negative).


Calculations:
Life Power      LP      Vitality / 2 + KIN_LP_BONUS
Husk Power      HP      (((Level * Vitality) div 5) + Level*2 + Vitality*2 + 8) * RACE_HP_MUL
Junction Power  JP      ((Level * Numen) div 4) + Level*2 + Numen*2) * RACE_JP_MUL
Motion Power    MP      (((Level * Zeal) div 4) + Level * 2 + Zeal * 2) + 30) * RACE_MP_MUL

XXX: this will be changd to only include Equipment related values, skill and talent are added to the calculation later.
Offense         OFF     Quality of equipped weapon
Protection      PRO     Sum of quality of equipped gear.
Blocking        BLO     (shield quality ) OR ( weapon skill + weapon quality ) if learned weapon's parry technique, otherwise ineffective.
Rapidity        RAP     - weight of armor  - weight of shield - weight of weapon.
Yield           YIE     Numen - interference penalty of gear - interference penalty of shield - interference of weapon + quality of staff if equipped + quality of equipped Focus.


Offense
Protection
Block
Rapidity
Yield

*/

// Generates a prompt for use with the being/character
func (me *Being) ToPrompt() string {
	if (me.Kin != nil) && (me.Kin.Arts > 0.0) {
		return fmt.Sprintf("HP:%s MP:%s LP:%s JP:%s", me.HP.TNM(), me.MP.TNM(), me.LP.TNM(), me.JP.TNM())
	} else {
		return fmt.Sprintf("HP:%s MP:%s LP:%s", me.HP.TNM(), me.MP.TNM(), me.LP.TNM())
	}
}

func (me *Being) GenderName() string {
	if me.Gender == nil {
		return "????"
	}
	return me.Gender.Name
}

func (me *Being) KinName() string {
	if me.Kin == nil {
		return "????"
	}
	return me.Kin.Name
}

func (me *Being) JobName() string {
	if me.Job == nil {
		return "????"
	}
	return me.Job.Name
}

// Generates an overview of the essentials of the being as a string.
func (me *Being) ToEssentials() string {
	return fmt.Sprintf("%s lvl %d %s %s %s", me.Name, me.Level, me.GenderName(), me.KinName(), me.JobName())
}

// Generates an overview of the physical talents of the being as a string.
func (me *Talents) ToBodyTalents() string {
	return fmt.Sprintf("STR: %3d    TOU: %3d    AGI: %3d    DEX: %3d", me.Strength, me.Toughness, me.Agility, me.Dexterity)
}

// Generates an overview of the mental talents of the being as a string.
func (me *Talents) ToMindTalents() string {
	return fmt.Sprintf("INT: %3d    WIS: %3d    CHA: %3d    EMO: %3d", me.Intelligence, me.Wisdom, me.Charisma, me.Emotion)
}

// Generates an overview of the physical derived talents of the being as a string.
func (me *Talents) ToBodyDerived() string {
	return fmt.Sprintf("FOR: %3d    VIT: %3d    QUI: %3d    KNA: %3d", me.Force(), me.Vitality(), me.Quickness(), me.Knack())
}

// Generates an overview of the mental derived talents of the being as a string.
func (me *Talents) ToMindDerived() string {
	return fmt.Sprintf("UND: %3d    GRA: %3d    ZEA: %3d    NUM: %3d", me.Understanding(), me.Grace(), me.Zeal(), me.Numen())
}

// Generates an overview of the derived talents of the being as a string.
func (me *Talents) ToDerived() string {
	status := me.ToBodyDerived()
	status += "\n" + me.ToMindDerived()
	return status
}

// Generates an overview of all talents as a string.
func (me *Talents) ToTalents() string {
	status := me.ToBodyTalents()
	status += "\n" + me.ToMindTalents()
	return status
}

// Generates an overview of the equipment values of the being as a string.
func (me *EquipmentValues) ToEquipmentValues() string {
	return fmt.Sprintf("OFF: %3d    PRO: %3d    BLO: %3d    RAP: %3d    YIE: %3d", me.Offense, me.Protection, me.Block, me.Rapidity, me.Yield)
}

// Generates an overview of the status of the being as a string.
func (me *Being) ToStatus() string {
	status := me.ToEssentials()
	status += "\n" + me.ToTalents()
	status += "\n" + me.ToDerived()
	status += "\n" + me.ToEquipmentValues()
	status += "\n" + me.ToPrompt()
	status += "\n"
	return status
}

func (me *Talents) GrowFrom(from Talents) {
	me.Strength += from.Strength
	me.Toughness += from.Toughness
	me.Agility += from.Agility
	me.Dexterity += from.Dexterity
	me.Intelligence += from.Intelligence
	me.Wisdom += from.Wisdom
	me.Charisma += from.Charisma
	me.Emotion += from.Emotion
}

func (me *Vital) NewMax(max int) {
	oldmax := me.Max
	me.Max = max
	delta := me.Max - oldmax
	me.Now += delta
	if me.Now > me.Max {
		me.Now = me.Max
	}
	if me.Now < 0 {
		me.Now = 0
	}
}

func (me *Being) RecalculateVitals() {
	newhp := (me.Level * me.Vitality() / 5) + me.Level*2 + me.Vitality()*2
	newhp += 8

	me.Vitals.HP.NewMax(newhp)
	me.Vitals.LP.NewMax((me.Vitality()/2 + 4))

	newjp := (me.Level*me.Numen())/4 + me.Level*2 + me.Numen()*2
	newjpf := float64(newjp)
	newmp := (me.Level*me.Zeal())/4 + me.Level*2 + me.Zeal()*2 + 32
	newmpf := float64(newmp)

	if me.Kin != nil {
		newjpf *= me.Kin.Arts
		newmpf *= me.Kin.Techniques
	}

	me.Vitals.MP.NewMax(int(math.Floor(newjpf)))
	me.Vitals.JP.NewMax(int(math.Floor(newmpf)))
}

func (me *Being) Init(kind string, name string, privilege Privilege,
	kin Entitylike, gender Entitylike, job Entitylike) *Being {
	if me == nil {
		return me
	}
	me.Entity.InitKind(kind, name, privilege)

	realkin := EntitylikeToKin(kin)
	realgen := EntitylikeToGender(gender)
	realjob := EntitylikeToJob(job)

	monolog.Info("Init being: Kin: %v", realkin)
	monolog.Info("Init being: Gender: %v", realgen)
	monolog.Info("Init being: Job: %v", realjob)

	me.Kin = realkin
	me.Gender = realgen
	me.Job = realjob

	me.Talents.GrowFrom(BasicTalent)
	me.Talents.GrowFrom(me.Kin.Talents)
	me.Talents.GrowFrom(me.Gender.Talents)
	me.Talents.GrowFrom(me.Job.Talents)

	me.Level = 1
	me.RecalculateVitals()

	return me
}

func NewBeing(kind string, name string, privilege Privilege,
	kin Entitylike, gender Entitylike, job Entitylike) *Being {
	res := &Being{}
	res.Init(kind, name, privilege, kin, gender, job)
	return res
}

func (me *Being) Type() string {
	return "being"
}

func (me *Being) Save(datadir string) {
	SaveSavable(datadir, me)
}

func LoadBeing(datadir string, nameid string) *Being {
	res, _ := LoadLoadable(datadir, nameid, new(Being)).(*Being)
	return res
}

func (me *Talents) SaveSitef(rec *sitef.Record) (err error) {
	rec.PutStruct("", *me)
	return nil
}

func (me *Vitals) SaveSitef(rec *sitef.Record) (err error) {
	rec.PutStruct("", *me)
	return nil
}

func (me *EquipmentValues) SaveSitef(rec *sitef.Record) (err error) {
	rec.PutStruct("", *me)
	return nil
}

func (me *Aptitudes) SaveSitef(rec *sitef.Record) (err error) {
	nskills := len(me.Skills)
	rec.PutInt("skills", nskills)
	for i := 0; i < nskills; i++ {
		rec.PutArrayIndex("skills", i, me.Skills[i].skill.ID)
	}

	return nil
}

func (me *Inventory) SaveSitef(rec *sitef.Record) (err error) {

	return nil
}

// Save a being to a sitef record.
func (me *Being) SaveSitef(rec *sitef.Record) (err error) {
	me.Entity.SaveSitef(rec)
	rec.PutInt("level", me.Level)

	if me.Gender != nil {
		rec.Put("gender", me.Gender.ID)
	}

	if me.Job != nil {
		rec.Put("job", me.Job.ID)
	}

	if me.Kin != nil {
		rec.Put("kin", me.Kin.ID)
	}

	me.Talents.SaveSitef(rec)
	me.Vitals.SaveSitef(rec)
	me.EquipmentValues.SaveSitef(rec)
	me.Aptitudes.SaveSitef(rec)
	me.Inventory.SaveSitef(rec)

	if me.Room != nil {
		rec.Put("room", me.Room.ID)
	}

	return nil
}

func (me *Talents) LoadSitef(rec sitef.Record) (err error) {
	rec.GetStruct("", me)
	return nil
}

func (me *Vitals) LoadSitef(rec sitef.Record) (err error) {
	rec.GetStruct("", me)
	return nil
}

func (me *EquipmentValues) LoadSitef(rec sitef.Record) (err error) {
	rec.GetStruct("", me)
	return nil
}

func (me *Aptitudes) LoadSitef(rec sitef.Record) (err error) {
	//    rec.GetStruct("", *me)
	return nil
}

func (me *Inventory) LoadSitef(rec sitef.Record) (err error) {
	//    rec.GetStruct("", *me)
	return nil
}

// Load a being from a sitef record.
func (me *Being) LoadSitef(rec sitef.Record) (err error) {
	me.Entity.LoadSitef(rec)

	me.Level = rec.GetIntDefault("level", 1)

	me.Gender = EntitylikeToGender(GenderEntityList.FindID(rec.Get("gender")))
	me.Job = EntitylikeToJob(JobEntityList.FindID(rec.Get("job")))
	me.Kin = EntitylikeToKin(KinEntityList.FindID(rec.Get("kin")))

	me.Talents.LoadSitef(rec)
	me.Vitals.LoadSitef(rec)
	me.EquipmentValues.LoadSitef(rec)
	me.Aptitudes.LoadSitef(rec)
	me.Inventory.LoadSitef(rec)

	if rec.Get("room") != "" {
		var err error
		me.Room, err = DefaultWorld.LoadRoom(rec.Get("room"))
		if err != nil {
			monolog.WriteError(err)
			return err
		}
	}
	return nil
}
