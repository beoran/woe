package world

type DamageKind string
const (
    DAMAGE_CUT          DamageKind = "cut"
    DAMAGE_CRUSH        DamageKind = "crush"
    DAMAGE_PIERCE       DamageKind = "pierce"
    DAMAGE_HEAT         DamageKind = "heat"
    DAMAGE_COLD         DamageKind = "cold"
    DAMAGE_SHOCK        DamageKind = "shock"
    DAMAGE_TOXIC        DamageKind = "toxic"
    DAMAGE_LASER        DamageKind = "laser"
    DAMAGE_BLAST        DamageKind = "blast"
    DAMAGE_TONE         DamageKind = "tone"
    DAMAGE_CORRUPTION   DamageKind = "corruption"
    DAMAGE_ARCANE       DamageKind = "arcane"
    DAMAGE_HEAL         DamageKind = "heal"
    DAMAGE_REPAIR       DamageKind = "repair"
)
    
    



type ItemKind string

const (
    ITEM_MEDICINE   ItemKind = "medicine"
    
    ITEM_CAP        ItemKind = "cap"
    ITEM_RIBBON     ItemKind = "ribbon"    
    ITEM_HAT        ItemKind = "hat"
    ITEM_SCARF      ItemKind = "scarf"
    ITEM_CIRCLET    ItemKind = "circlet"
    ITEM_HEADGEAR   ItemKind = "headgear"
    ITEM_CROWN      ItemKind = "crown"
    ITEM_HELMET     ItemKind = "helmet"
    
    ITEM_CAPE       ItemKind = "cape"
    ITEM_COAT       ItemKind = "coat"
    ITEM_ROBE       ItemKind = "robe"
    ITEM_VEST       ItemKind = "vest"
    ITEM_CHEST      ItemKind = "chest"
    ITEM_SUIT       ItemKind = "suit"
    ITEM_ARMOR      ItemKind = "armor"

    ITEM_SANDAL     ItemKind = "sandal"
    ITEM_SHOE       ItemKind = "shoe"
    ITEM_BOOT       ItemKind = "boot"
    

 
 //not sure...   
    ITEM_PANTS      ItemKind = "pants"
    ITEM_SKIRT      ItemKind = "skirt"
    ITEM_GREAVES    ItemKind = "greaves"  
    ITEM_RING       ItemKind = "ring"
    ITEM_BRACELET   ItemKind = "bracelet"
    ITEM_ARMLET     ItemKind = "armlet"
    ITEM_SLEEVE     ItemKind = "sleeve"
    ITEM_PAULDRON   ItemKind = "pauldron"
    ITEM_GAUNTLET   ItemKind = "gauntlet"
    ITEM_NECKLACE   ItemKind = "necklace"
    ITEM_EARRING    ItemKind = "earring"
    ITEM_LAMP       ItemKind = "lamp"
    ITEM_FLASHLIGHT ItemKind = "flashlight"
  // sure again
   
    ITEM_SWORD      ItemKind = "sword"
    ITEM_TWOHANDER  ItemKind = "twohander"
    ITEM_KNIFE      ItemKind = "knife"
    ITEM_DAGGER     ItemKind = "dagger"
    ITEM_GLOVE      ItemKind = "glove"
    ITEM_CLAW       ItemKind = "claw"
    ITEM_WAND       ItemKind = "wand"
    ITEM_STAFF      ItemKind = "staff"
    ITEM_FOCUS      ItemKind = "focus"
    
    // Not sure if these will be implemented.
    
    ITEM_AXE        ItemKind = "axe"
    ITEM_MAUL       ItemKind = "maul"
    ITEM_SPEAR      ItemKind = "spear"
    ITEM_NAGINATA   ItemKind = "naginata"
    ITEM_BOW        ItemKind = "bow"
    ITEM_CROSSBOW   ItemKind = "crossbow"
    ITEM_ARROW      ItemKind = "arrow"
    ITEM_BOLT       ItemKind = "bolt"

    // But these below are planned.
    
    ITEM_NEEDLER    ItemKind = "needler"
    ITEM_HANDGUN    ItemKind = "handgun"
    ITEM_LASERGUN   ItemKind = "lasergun"
    ITEM_MACHINEGUN ItemKind = "machinegun"
    ITEM_CANNON     ItemKind = "cannon"
    ITEM_BAZOOKA    ItemKind = "bazooka"
    
    ITEM_NEEDLE     ItemKind = "needle"
    ITEM_BULLET     ItemKind = "bullet"
    ITEM_CELL       ItemKind = "cell"
    ITEM_ROCKET     ItemKind = "rocket"
    
    ITEM_EXPLOSIVE  ItemKind = "explosive"
    ITEM_GRENADE    ItemKind = "grenade"
    ITEM_REPAIR     ItemKind = "repair"
    ITEM_ORE        ItemKind = "ore"
    ITEM_INGOT      ItemKind = "ingot"
    ITEM_METAL      ItemKind = "metal"
    ITEM_PLANT      ItemKind = "plant"
    ITEM_FRUIT      ItemKind = "fruit"
    ITEM_WOOD       ItemKind = "wood"
    ITEM_FOOD       ItemKind = "food"
    ITEM_TRAP       ItemKind = "trap"
    ITEM_HANDCUFFS  ItemKind = "handcuffs"
    ITEM_CHEMICAL   ItemKind = "chemical"
    ITEM_FISH       ItemKind = "fish"
    ITEM_MEAT       ItemKind = "meat"
    ITEM_HIDE       ItemKind = "hide"
    ITEM_LEATHER    ItemKind = "leather"
    ITEM_FIBER      ItemKind = "fiber"
    ITEM_CLOTH      ItemKind = "cloth"
    ITEM_CERAMIC    ItemKind = "ceramic"
    ITEM_POLYMER    ItemKind = "polymer"
    
    // android parts
    // AGI, STR, CHA 
    ITEM_LEGPARTS   ItemKind = "legparts"
    // DEX, STR, CHA
    ITEM_ARMPARTS   ItemKind = "armparts"
    // TOU, CHA
    ITEM_BODYPARTS  ItemKind = "bodyparts"
    // INT, WIS, CHA
    ITEM_HEADPARTS  ItemKind = "headparts"
    
    
    ITEM_           ItemKind = ""
)

type EquipWhere string

const (
    EQUIP_NONE      EquipWhere = "none"
    EQUIP_HEAD      EquipWhere = "head"
    EQUIP_TORSO     EquipWhere = "torso"
    EQUIP_OFFHAND   EquipWhere = "offhand"
    EQUIP_DOMINANT  EquipWhere = "dominant"
    EQUIP_AMMO      EquipWhere = "ammo"    
    EQUIP_FEET      EquipWhere = "feet"
    EQUIP_FOCUS     EquipWhere = "focus"
    EQUIP_PHONE     EquipWhere = "phone"
    
    EQUIP_GLOVES    EquipWhere = "gloves"
    EQUIP_NECK      EquipWhere = "neck"
    EQUIP_LEGS      EquipWhere = "legs"
    EQUIP_ARMS      EquipWhere = "arms"
    EQUIP_RIGHTRING EquipWhere = "rightring"
    EQUIP_LEFTRING  EquipWhere = "leftring"
    EQUIP_LIGHT     EquipWhere = "light"
    EQUIP_          EquipWhere = ""    
)

type Item struct {
    Entity
    Quality       int
    Price         int
    Kind          ItemKind
    Damage        DamageKind
    // equipment location,  "none" if not equippable
    Equippable    EquipWhere
    // Level of crafing skill needed to craft this, or of harvesting skill to harvest this, 
    // or of mining skill to mine this. Negative if cannot be crafted nor harvested, nor mined.    
    Level         int
    // Id's of ingredients to craft this item. Empty if it cannot be crafted.
    Ingredients []ID
    // Id of item this item can be upgraded/enhanced to. empty or "none" if cannot be upgraded.
    Upgrade       ID
    // ID of item this item can degrade into. empty or "none" if cannot be degraded.
    Degrade       ID
    // ID of technique/art/item to craft this item teaches when used, empty or none if it teaches nothing.
    // If it'ss a skill, the XP of teaching is determine by the Quality of the item.   
    Teaches       ID
}

type ItemPointer struct {
    ID     ID
    item * Item
}




