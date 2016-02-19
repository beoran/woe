package world

import "fmt"
// import "strconv"
import "github.com/beoran/woe/monolog"
import "github.com/beoran/woe/sitef"

/* Characters may exist "outside" of any world, because, at login of a client,
 * the Character must be loaded or created before it enters a World. */

type Character struct {
    Being       
    AccountName string
    account   * Account
}


func NewCharacter(being Being, accountname string, account * Account) (*Character) {
    return &Character{being, accountname, account}
}


// Save a character as a sitef record.
func (me * Character) SaveSitef(rec sitef.Record) (err error) {
    rec["accountname"]  = me.AccountName
    // TODO: saving: me.Being.SaveSitef(rec)
    return nil
}

// Load a character from a sitef record.
func (me * Character) LoadSitef(rec sitef.Record) (err error) {
    me.AccountName = rec["accountname"] 
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
    character.AccountName   = record["AccountName"]
    account, err           := DefaultWorld.LoadAccount(dirname, character.AccountName);
    if err != nil  {
        return nil, err
    } 
    
    if account == nil {
        return nil, fmt.Errorf("Cound not load account %s for character %s", 
            character.AccountName, character.Name)
    }
    
    character.account = account
    
    return character, nil
}
