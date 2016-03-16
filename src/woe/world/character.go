package world

import "fmt"
// import "strconv"
import "github.com/beoran/woe/monolog"
import "github.com/beoran/woe/sitef"

/* Characters may exist "outside" of any world, because, at login of a client,
 * the Character must be loaded or created before it enters a World. */

type Character struct {
    Being       
    * Account
}


func NewCharacterFromBeing(being Being, account * Account) (*Character) {
    return &Character{being, account}
}

func (me * Character) Init(account * Account, name string, 
              kin Entitylike, gender Entitylike, job Entitylike) (* Character) {
    me.Account = account
    me.Being.Init("character", name, account.Privilege, kin, gender, job)
    return me
}

func NewCharacter(account * Account, name string, 
              kin Entitylike, gender Entitylike, job Entitylike) (* Character) {
    me := &Character{};
    return me.Init(account, name, kin, gender, job)
}

// Save a character as a sitef record.
func (me * Character) SaveSitef(rec sitef.Record) (err error) {
    rec["accountname"]  = me.Account.Name
    me.Being.SaveSitef(rec)
    
    return nil
}

// Load a character from a sitef record.
func (me * Character) LoadSitef(rec sitef.Record) (err error) {
    aname := rec["accountname"]
    account, err := DefaultWorld.LoadAccount(aname)
    if err != nil {
        return err
    } 
    me.Account = account
    me.Being.LoadSitef(rec)
    // TODO: load being. me.Being.SaveSitef(rec)
    return nil
}


// Save a character as a sitef file.
func (me * Character) Save(dirname string) (err error) {
    path := SavePathFor(dirname, "character", me.Name)
    
    rec                := make(sitef.Record)
    me.SaveSitef(rec)
    monolog.Debug("Saving Character record: %s %v", path, rec)
    return sitef.SaveRecord(path, rec)
}

// Load an character from a sitef file.
func LoadCharacter(dirname string, name string) (character *Character, err error) {
    
    path := SavePathFor(dirname, "character", name)
    
    records, err := sitef.ParseFilename(path)
    if err != nil {
        return nil, err
    }
    
    if len(records) < 1 {
        return nil, fmt.Errorf("No sitef record found for %s!", name)
    }
    
    record := records[0]
    monolog.Info("Loading Account record: %s %v", path, record)
    
    character               = new(Character)
    aname                  := record["AccountName"]
    account, err           := DefaultWorld.LoadAccount(aname);
    if err != nil  {
        return nil, err
    } 
    
    if account == nil {
        return nil, fmt.Errorf("Cound not load account %s for character %s", 
            aname, character.Name)
    }
    
    character.Account = account
    
    return character, nil
}
