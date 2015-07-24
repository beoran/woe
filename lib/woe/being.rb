module Woe
  module Being
    include Serdes
    include Monolog

    serdes_reader   :id
    serdes_reader   :account
    
    # Essentials
    serdes_reader   :kind         # AKA "race"
    serdes_reader   :level
    serdes_reader   :profession   # AKA "class"
    
    # Talents
    serdes_reader   :strength
    serdes_reader   :toughness
    serdes_reader   :agility
    serdes_reader   :dexterity
    serdes_reader   :intelligence
    serdes_reader   :wisdom
    serdes_reader   :charisma
    serdes_reader   :essence
    
    # Vitals
    serdes_reader   :hp
    serdes_reader   :hp_max
    serdes_reader   :mp
    serdes_reader   :mp_max
    serdes_reader   :jp
    serdes_reader   :jp_max
    serdes_reader   :lp
    serdes_reader   :lp_max
        
    # Equipment values
    serdes_reader   :offense
    serdes_reader   :protection
    serdes_reader   :block
    serdes_reader   :rapidity
    serdes_reader   :yield


    # Skills array
    serdes_reader   :skills
    
    # Affects array
    serdes_reader   :affects
    
    # Derived stats 
    def force
      return (strength * 2 + wisdom) / 3
    end
    
    def vitality
      return (toughness * 2 + charisma) / 3
    end
    
    def quickness
      return (agilty * 2 + intelligence) / 3
    end
    
    def knack
      return (dexterity * 2 + essence) / 3
    end
    
    def understanding
      return (intelligence * 2 + toughness) / 3
    end
    
    def grace 
      return (charisma * 2 + agility) / 3
    end
    
    def zeal
      return (wisdom * 2 + strength) / 3
    end
    
    def numen
      return (essence * 2 + DEX) / 3
    end

    # Generates a prompt for use with the being
    def to_prompt
      if essence > 0
        "HP:#{hp}/#{hp_max} MP:#{mp}/#{mp_max} JP:#{jp}/#{jp_max} LP:#{lp}/#{lp_max}>"
      else
        "HP:#{hp}/#{hp_max} MP:#{mp}/#{mp_max} LP:#{lp}/#{lp_max}>"
      end
    end
  
    # Generates an overview of the status of the being.
    def to_status
      status <<- STATUS_END
      ********************************
        Status
        
        #{id} level #{level} #{kind} #{profession}
        
        Talents
        
        STR #{strength}
        TOU #{toughness}
        AGI #{agility}
        DEX #{dexterity}
        INT #{intelligence}
        WIS #{wisdom}
        CHA #{charisma}
        ESS #{essence}
        
        Equipment
        
        OFF #{offense}
        PRO #{protection}
        BLO #{block}
        RAP #{rapidity}
        YIE #{self.yield}
        
        "HP:#{hp}/#{hp_max} MP:#{mp}/#{mp_max} JP:#{jp}/#{jp_max} LP:#{lp}/#{lp_max}>"

      ********************************
      STATUS_END
      return status
    end
    
    # Sets up the being based on the talents
    def setup_on_birth
      
    end 
    
  
  end
end
    

