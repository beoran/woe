package world

const (
    TALENT_NONE ID = "TALENT_NONE"
    TALENT_STR  ID = "STR"
    TALENT_TOU  ID = "TOU"
    TALENT_AGI  ID = "AGI"
    TALENT_DEX  ID = "DEX"
    TALENT_INT  ID = "INT"
    TALENT_WIS  ID = "WIS"
    TALENT_CHA  ID = "CHA"
    TALENT_EMO  ID = "EMO"
)
    

type Skill struct {
    Entity
    Kind              ID
    Talent            ID
    derived           func (m * Being)  int
}

type BeingSkill struct {
    being       * Being
    skill       * Skill
    talent      * int
    Experience    int
    Next          int
    Level         int
}

type JobSkill struct {
    Experience    int
    Level         int
    
    
}

/* Requires Arts to be useful:
 * Arcane Tone Moon Sun Tree Flame Snow Bolt Omen
 * 
 * Requires techniques to be useful: 
 * Rage, Sleight, Robotech, Stealth, Medical, 
 * 
 * Don't strictly need techniques but they are desirable:
 * Shield, Lore, Hacking, Barter, Social, Bravery, Survival, Knife, Acrobatics, Travelling, Science, Engineering
 * 
 * Optionally could use some techniques:
 * Sword, Cannon, Fist, Mauls, Toiling,  Gun, Polearm,  Explosives, 
 * 
 * Other's don't really require techniques. These are armor wearing skills, 
 * and crafting skills (crafting info is stored in the item info in stead of in a technique).
 * 
 */


var SkillList = []Skill {
    { Entity : Entity { ID: "skill_sword", Name: "Sword", 
      Short: "Fighting with swords and katanas.", },
      Talent : TALENT_STR,
    },
    { Entity : Entity { ID: "skill_cannon", Name: "Cannon", 
      Short: "Use and shooting of cannons, launchers and bazookas", },
      Talent : TALENT_STR,
    },
    { Entity : Entity { ID: "skill_medium_gear", Name: "Medium Gear", 
      Short: "Use of medium gear such as vests, gear and shoes.", },  
      Talent : TALENT_STR,
    },
    { Entity : Entity { ID: "skill_weaponsmith", Name: "Weaponsmith", 
      Short: "Crafting of hand to hand style weapons.", },  
      Talent : TALENT_STR,
    },
    { Entity : Entity { ID: "skill_mining", Name: "Mining", 
      Short: "Mining of minerals.", },  
      Talent : TALENT_STR,
    },
    { Entity : Entity { ID: "skill_rage", Name: "Rage",
      Short: "Techniques of beasts.", },
      Talent : TALENT_STR,
    },
 

    { Entity : Entity { ID: "skill_fist", Name: "Fist", 
      Short: "Martial arts and unarmed fighting, use of gloves.", },
      Talent : TALENT_TOU,
    },
    { Entity : Entity { ID: "skill_mauls", Name: "Mauls", 
      Short: "Fighting with heavy meelee weapons such as axes and hammers", },
      Talent : TALENT_TOU,
    },
    { Entity : Entity { ID: "skill_heavy_gear", Name: "Heavy Gear", 
      Short: "Use of heavy gear such as armor suits, helmets and boots.", },  
      Talent : TALENT_TOU,
    },
    { Entity : Entity { ID: "skill_armorer", Name: "Armorer", 
      Short: "Crafing of heavy equipment and armor.", },
      Talent : TALENT_TOU,
    },
    { Entity : Entity { ID: "skill_toiling", Name: "Toiling", 
      Short: "Physical exertion and carrying heavy loads.", },
      Talent : TALENT_TOU,
    },    
    { Entity : Entity { ID: "skill_survival", Name: "Survival", 
      Short: "Fishing and hunting skills, recovering biological ingredients.", },
      Talent : TALENT_TOU,
    },


    { Entity : Entity { ID: "skill_knife", Name: "Knife",
      Short: "Fighting with knives, daggers.", },
      Talent : TALENT_DEX,
    },
    { Entity : Entity { ID: "skill_sleight", Name: "Sleight", 
      Short: "Disarming of traps, picking of locks and pickpocketing from beasts.", },
      Talent : TALENT_DEX,
    },
    { Entity : Entity { ID: "skill_gun", Name: "Gun", 
      Short: "Use and shooting of guns and bowguns", },
      Talent : TALENT_DEX,
    },
    { Entity : Entity { ID: "skill_gunsmith", Name: "Gunsmith", 
      Short: "Crafting of guns , cannons and other ranged weapons.", },
      Talent : TALENT_DEX,
    },
    { Entity : Entity { ID: "skill_outfitter", Name: "Outfitter", 
      Short: "Crafting of medium weight equipment.", },
      Talent : TALENT_DEX,
    },
    { Entity : Entity { ID: "skill_robotic", Name: "Robotic", 
      Short: "Special techniques of robots and androids.", },
      Talent : TALENT_DEX,
    },



    { Entity : Entity { ID: "skill_polearm", Name: "Polearm", 
      Short: "Fighting with pole arms such as spears, lances and naginatas.", },
      Talent : TALENT_AGI,
    },

    { Entity : Entity { ID: "skill_shield", Name: "Shield",
      Short: "Use of bucklers and shields and shield relates techniques.", },  
      Talent : TALENT_AGI,
    },
    { Entity : Entity { ID: "skill_light_gear", Name: "Light Gear", 
      Short: "Use of light gear such as robes, caps and sandals.", },
      Talent : TALENT_AGI,
    },
    { Entity : Entity { ID: "skill_acrobatics", Name: "Acrobatics", 
      Short: "Evading attacks, running away from danger.", },
      Talent : TALENT_AGI,
    },
    { Entity : Entity { ID: "skill_stealth", Name: "Stealth", 
      Short: "Stealthy movement and suprising opponents.", },
      Talent : TALENT_AGI,
    },
    { Entity : Entity { ID: "skill_traveling", Name: "Traveling", 
      Short: "Efficient traveling and skills to reduce MP efficiency", },
      Talent : TALENT_AGI,
    },
 
 
    { Entity : Entity { ID: "skill_explosives", Name: "Explosives", 
      Short: "Crafting, use and disarming of bombs and explosives.", },
      Talent : TALENT_INT,
    },    
    { Entity : Entity { ID: "skill_science", Name: "Science", 
      Short: "Crafting of healing items and poisons and using them.", },
      Talent : TALENT_INT,
    },
    { Entity : Entity { ID: "skill_lore", Name: "Lore", 
      Short: "General and historic knowledge.", },
      Talent : TALENT_INT,
    },
    { Entity : Entity { ID: "skill_medical", Name: "Medical", 
      Short: "Healing techniques and use of medication.", },
      Talent : TALENT_INT,
    },
    { Entity : Entity { ID: "skill_engineering", Name: "Engineering", 
      Short: "Crafting of android parts, repair and construction of machines", },
      Talent : TALENT_INT,
    },    
    { Entity : Entity { ID: "skill_hacking", Name: "Hacking", 
      Short: "Knowledge of computers, boosting of machines", },
      Talent : TALENT_INT,
    },

        
    { Entity : Entity { ID: "skill_barter", Name: "Barter", 
      Short: "Trading and bribing", },
      Talent : TALENT_CHA,
    },
    { Entity : Entity { ID: "skill_social", Name: "Social", 
      Short: "Talking to NPC's and convincing them. Supportive group techniques and leadership skills.", },
      Talent : TALENT_CHA,
    },
    { Entity : Entity { ID: "skill_arcane", Name: "Arcane Arts", 
      Short: "Extraordinary Numen Arts.", },  
      Talent : TALENT_CHA,
    },
    { Entity : Entity { ID: "skill_tone", Name: "Tone Arts", 
      Short: "Numen arts that control sound and buffing effects.", },  
      Talent : TALENT_CHA,
    },
    { Entity : Entity { ID: "skill_tailor", Name: "Tailor", 
      Short: "Crafting of light clothes and equipment.", },  
      Talent : TALENT_CHA,
    },
    { Entity : Entity { ID: "skill_cooking", Name: "Cooking",
      Short: "Crafting of food, for temporary boosts.", },  
      Talent : TALENT_CHA,
    },


    { Entity : Entity { ID: "skill_earth", Name: "Earth Arts",
      Short: "Numen arts that control earth and explosive effects.", },  
      Talent : TALENT_WIS,
    },
    { Entity : Entity { ID: "skill_sun", Name: "Sun Arts", 
      Short: "Numen arts that control light and purification.", },  
      Talent : TALENT_WIS,
    },
    { Entity : Entity { ID: "skill_tree", Name: "Tree Arts", 
      Short: "Numen arts that control living beings, nature and healing.", },  
      Talent : TALENT_WIS,
    },
    { Entity : Entity { ID: "skill_bravery", Name: "Bravery", 
      Short: "Courage and mental resistance.", },  
      Talent : TALENT_WIS,
    },
    { Entity : Entity { ID: "skill_jeweler", Name: "Jeweler", 
      Short: "Crafting of jewels mundane and imbued with Numen.", },  
      Talent : TALENT_WIS,
    },
    { Entity : Entity { ID: "skill_artistic", Name: "Artistic", 
      Short: "Graphical and plastic arts to craft artworks.", },  
      Talent : TALENT_WIS,
    },
    
    { Entity : Entity { ID: "skill_staff", Name: "Staff", 
      Short: "Fighting with staves and wands.", },
      Talent : TALENT_EMO,
    },
    { Entity : Entity { ID: "skill_flame", Name: "Flame Arts", 
      Short: "Numen arts that control fire and heat.", },  
      Talent : TALENT_EMO,
    },
    { Entity : Entity { ID: "skill_snow", Name: "Snow Arts", 
      Short: "Numen arts that control water, ice and cold.", },  
      Talent : TALENT_EMO,
    },
    { Entity : Entity { ID: "skill_bolt", Name: "Bolt Arts", 
      Short: "Numen arts that control electricity, wind and air.", },  
      Talent : TALENT_EMO,
    },
    { Entity : Entity { ID: "skill_numensmith", Name: "Numensmith",
      Short: "Crafting of wands, rods and other items enhanced with Numen.", },  
      Talent : TALENT_EMO,
    },
    { Entity : Entity { ID: "skill_omen", Name: "Omen Arts",
      Short: "Arts that use Omen in stead of Numen, by Corrupted beings", },
      Talent : TALENT_EMO,
    },

    
}


/*
 *  Other ideas: 
 *  ESS, AGI
 * 
 *  
 * Cooking -- Food for various bufffs 
 * 
 */

