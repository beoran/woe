package world

// An art is a special case of technique that requires Essence, 
// and consumes JP in stead of MP
type Art Technique



type BeingArt struct {
    being       * Being
    art         * Art
    Experience    int
    Level         int
}



/* Arcane Tone Earth Sun Tree Flame Snow Bolt Omen */

var ArtList = []Art {
    /* Arcane arts, need more later */ 
    { 
      Skill  : "skill_arcane",  
      Entity : Entity { ID: "art_scry", Name: "Scry", 
      Short  : "Sense mobiles in the current zone", },
      Effect : "arcane",
      Level  : 1,
      Cost   : 10,
    },

    { 
      Skill  : "skill_arcane",  
      Entity : Entity { ID: "art_melt", Name: "Melt", 
      Short  : "A piercing art attack aganst one being", },      
      Effect : "arcane",
      Level  : 1,
      Cost   : 30,
    },
    
    { 
      Skill  : "skill_tone",  
      Entity : Entity { ID: "art_hope_popsong", Name: "Hope Popsong", 
      Short  : "An encouraging song. Restores some MP.", },      
      Effect : "restore mp",
      Level  : 1,
      Cost   : 30,
    },
    
    { 
      Skill  : "skill_tone",  
      Entity : Entity { ID: "art_scream", Name: "Scream", 
      Short  : "Scream that cause sound damage and may scare weaker foes away.", },
      Effect : "tone",
      Level  : 1,
      Cost   : 30,
    },
    
    { 
      Skill  : "skill_earth",  
      Entity : Entity { ID: "art_sand", Name: "Sand", 
      Short  : "Darkens the sight of a foe with sand, causing them to miss more often.", },
      Effect : "cause blind",
      Level  : 1,
      Cost   : 30,
    },
    
    { 
      Skill  : "skill_earth",  
      Entity : Entity { ID: "art_boom", Name: "Boom", 
      Short  : "Small explosion attack against one being.", },      
      Effect : "blast",
      Level  : 1,
      Cost   : 30,
    },
    
    
    { 
      Skill  : "skill_sun",  
      Entity : Entity { ID: "art_flash", Name: "Flash", 
      Short  : "Blinds foes and purifies corrupted beings.", },      
      Effect : "laser",
      Level  : 1,
      Cost   : 10,
    },
    
    { 
      Skill  : "skill_sun",
      Entity : Entity { ID: "art_sunbeam", Name: "Sunbeam", 
      Short  : "Small laser attack against one foe.", },      
      Effect : "laser",
      Level  : 1,
      Cost   : 30,
    },
    
    { 
      Skill  : "skill_tree",  
      Entity : Entity { ID: "art_poison", Name: "poison", 
      Short  : "Small toxic attack against one opponent.", },      
      Effect : "toxic",
      Level  : 1,
      Cost   : 30,
    },
    { 
      Skill  : "skill_tree",
      Entity : Entity { ID: "art_invigorate", Name: "Invigorate", 
      Short  : "Slighly heals the HP and MP of one's self.", },      
      Effect : "heal",
      Level  : 1,
      Cost   : 30,
    },
    
    { 
      Skill  : "skill_",  
      Entity : Entity { ID: "art_", Name: "", 
      Short  : "", },      
      Effect : "",
      Level  : 1,
      Cost   : 30,
    },
    { 
      Skill  : "skill_",  
      Entity : Entity { ID: "art_", Name: "", 
      Short  : "", },      
      Effect : "",
      Level  : 1,
      Cost   : 30,
    },
    
    { 
      Skill  : "skill_",  
      Entity : Entity { ID: "art_", Name: "", 
      Short  : "", },      
      Effect : "",
      Level  : 1,
      Cost   : 30,
    },
    { 
      Skill  : "skill_",  
      Entity : Entity { ID: "art_", Name: "", 
      Short  : "", },      
      Effect : "",
      Level  : 1,
      Cost   : 30,
    },
    
    { 
      Skill  : "skill_",  
      Entity : Entity { ID: "art_", Name: "", 
      Short  : "", },      
      Effect : "",
      Level  : 1,
      Cost   : 30,
    },
    { 
      Skill  : "skill_",  
      Entity : Entity { ID: "art_", Name: "", 
      Short  : "", },      
      Effect : "",
      Level  : 1,
      Cost   : 30,
    },
    
    { 
      Skill  : "skill_",  
      Entity : Entity { ID: "art_", Name: "", 
      Short  : "", },      
      Effect : "",
      Level  : 1,
      Cost   : 30,
    },
    { 
      Skill  : "skill_",  
      Entity : Entity { ID: "art_", Name: "", 
      Short  : "", },      
      Effect : "",
      Level  : 1,
      Cost   : 30,
    },
    
    { 
      Skill  : "skill_",  
      Entity : Entity { ID: "art_", Name: "", 
      Short  : "", },      
      Effect : "",
      Level  : 1,
      Cost   : 30,
    },
    
    { 
      Skill  : "skill_",  
      Entity : Entity { ID: "art_", Name: "", 
      Short  : "", },      
      Effect : "",
      Level  : 1,
      Cost   : 30,
    },
    
    
    
    
    
    
    
        
     
}


 
