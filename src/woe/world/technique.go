package world


/* Techniques, arts, exploits and and crafts are special or specific abilities 
 * that fall under a generic skill. Crafts are not stored in a separate 
 * static list. Rather the crafting requirements are stored in the item list 
 * which is loaded from disk. Beings then have a list of the ID's of items they 
 * have learned to craft.
 */
 
type Technique struct {
    Entity
    Kind        ID
    Effect      ID
    Level       int
    Cost        int
    Skill       ID
    skill     * Skill
    onUse       func (me * Technique, caster * Being, targets ...*Being) (bool)
}


type BeingTechnique struct {
    being       * Being
    technique   * Technique
    Being         ID
    Technique     ID
}

/* An exploit is a special technique that can only be used a few times 
 * per day.
 * 
 * Maximum amount of uses (integer math!)
 * EXPLOIT_RANK = 1 + (EXPLOIT_LEVEL / 10)
 * Uses.Max = 1 + (SKILL_LEVEL / (10 * EXPLOIT_RANK))
 * 
 */
type Exploit    Technique


type BeingExploit struct {
    being       * Being
    exploit     * Exploit
    Being         ID
    Exploit       ID
    // How many times the exploit may be used.
    Uses          Vital
}



func CraftCookOmlet(caster * Being, targets ...*Being) bool {
    return true
} 


/* Generic rage technique performer. */
func DoRageTechnique(me * Technique, caster * Being, targets ...*Being) bool {    
    return true
}
      

/*
 * 
 * Requires techniques to be useful: 
 * Rage, Sleight, Robotech, Stealth, Medical, 
 */


var TechniqueList = []Technique {
    /* Rage techiques, need more later */ 
    { 
      Skill  : "skill_rage",  
      Entity : Entity { ID: "tech_bite", Name: "Bite", 
      Short  : "A beast's bite attack.", },      
      Effect : "crush",
      Level  : 1,
      Cost   : 1,
      onUse  : func (me * Technique, caster * Being, targets ...*Being) bool {
            return true
     }, 
      
    },
    
    { 
      Skill  : "skill_rage",  
      Entity : Entity { ID: "tech_claw", Name: "Claw",
      Short  : "A beast's claw attack.", },      
      Effect : "cut",
      Level  : 1,
      Cost   : 1,
    },
            
    { 
      Skill  : "skill_rage",  
      Entity : Entity { ID: "tech_sting", Name: "Sting",
      Short  : "A beast's sting attack.", },      
      Effect : "pierce",
      Level  : 1,
      Cost   : 1,
      onUse  : DoRageTechnique,
    },
    
    { 
      Skill  : "skill_horn",  
      Entity : Entity { ID: "tech_sting", Name: "Horn",
      Short  : "A beast's horn attack.", },      
      Effect : "pierce",
      Level  : 1,
      Cost   : 1,
    },
    
    /* Sleight techniques, need more later  */
    { 
      Skill  : "skill_sleight",  
      Entity : Entity { ID: "tech_distract", Name: "Distract",
      Short  : "Distract an opponent, makes you likely to steal for one turn.", },
      Effect : "distract",
      Level  : 1,
      Cost   : 20,
    },
    
    { 
      Skill  : "skill_sleight",  
      Entity : Entity { ID: "tech_steal", Name: "Steal",
      Short  : "Attempt to steal an item from an opponent.", },
      Effect : "steal",
      Level  : 1,
      Cost   : 10,
    },
    
    /* Robotic techniques need more later. */
    { 
      Skill  : "skill_robotic",  
      Entity : Entity { ID: "tech_spark", Name: "Spark",
      Short  : "Minor shock damage to one opponent", },
      Effect : "shock",
      Level  : 1,
      Cost   : 10,
    },
    
    { 
      Skill  : "skill_robotic",  
      Entity : Entity { ID: "tech_drill", Name: "Drill",
      Short  : "Minor piercing damage to one opponent", },
      Effect : "pierce",
      Level  : 1,
      Cost   : 10,
    },
    
    { 
      Skill  : "skill_robotic",  
      Entity : Entity { ID: "tech_saw", Name: "Saw",
      Short  : "Minor cutting damage to one opponent", },
      Effect : "cut",
      Level  : 1,
      Cost   : 10,
    },
    
    { 
      Skill  : "skill_robotic",  
      Entity : Entity { ID: "tech_stomp", Name: "Stomp",
      Short  : "Minor crushing damage to one opponent", },
      Effect : "crush",
      Level  : 1,
      Cost   : 10,
    },
    
    /* Stealth skills */    
    { 
      Skill  : "skill_stealth",  
      Entity : Entity { ID: "tech_tiptoe", Name: "Tiptoe",
      Short  : "Walk more silently for a short time", },
      Effect : "stealth",
      Level  : 1,
      Cost   : 20,
    },

    { 
      Skill  : "skill_stealth",  
      Entity : Entity { ID: "tech_hide", Name: "Hide",
      Short  : "Attempt to hide, giving you a chance to ambush your opponent.", },
      Effect : "hide",
      Level  : 1,
      Cost   : 40,
    },
    
    /* Medical skills */
    { 
      Skill  : "skill_medical",  
      Entity : Entity { ID: "tech_diagnose", Name: "Diagnose",
      Short  : "Get detailed infomation on the status of one being.", },
      Effect : "diagnose",
      Level  : 1,
      Cost   : 20,
    },
    
    { 
      Skill  : "skill_medical",
      Entity : Entity { ID: "tech_first_aid", Name: "First Aid",
      Short  : "Slightly heals HP, restores one LP of one being.", },
      Effect : "heal",
      Level  : 1,
      Cost   : 20,
    },
    
}









