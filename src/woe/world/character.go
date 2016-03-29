package world

import "fmt"
import "os"
// import "strconv"
import "github.com/beoran/woe/monolog"
import "github.com/beoran/woe/sitef"

/* Characters may exist "outside" of any world, because, at login of a client,
 * the Character must be loaded or created before it enters a World. */

type Character struct {
    Being       
    Account * Account
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

// Save a character into a a sitef record.
func (me * Character) SaveSirec(rec * sitef.Record) (err error) {
    rec.Put("accountname", me.Account.Name)
    me.Being.SaveSitef(rec)
    return nil
}

// Load a character from a sitef record.
func (me * Character) LoadSirec(rec sitef.Record) (err error) {
    aname := rec.Get("accountname")
    account, err := DefaultWorld.LoadAccount(aname)
    if err != nil {
        return err
    } 
    me.Account = account
    me.Being.LoadSitef(rec)
    return nil
}


// Save a character as a sitef file.
func (me * Character) Save(dirname string) (err error) {
    path := SavePathFor(dirname, "character", me.ID)
    
    rec                := sitef.NewRecord()
    me.SaveSirec(rec)
    monolog.Debug("Saving Character record: %s %v", path, rec)
    return sitef.SaveRecord(path, *rec)
}


// Load a character from a sitef file. Does no account checking, but returns the account name.
func LoadCharacter(dirname string, id string) (character *Character, aname string, err error) {
    
    path := SavePathFor(dirname, "character", id)
    
    records, err := sitef.ParseFilename(path)
    if err != nil {
        return nil, "", err
    }
    
    if len(records) < 1 {
        return nil, "",  fmt.Errorf("No sitef record found for %s!", id)
    }
    
    record := records[0]
    monolog.Info("Loading Character record: %s %v", path, record)
    
    character               = new(Character)
    aname                   = record.Get("accountname")
    character.Being.LoadSitef(*record);
    
    return character, aname, nil
}

// Load a character WITH A GIVEN NAME from a sitef file. Does no account checking, but returns the account name.
func LoadCharacterByName(dirname string, name string) (character *Character, aname string, err error) {
    id := EntityNameToID("character", name)
    return LoadCharacter(dirname, id)
}


// Load an character from a sitef file for the given account.
func (account * Account) LoadCharacter(dirname string, id string) (character *Character, err error) {
    
    character, aname, err := LoadCharacter(dirname, id)
    if character == nil {
        return character, err
    }
    
    if aname != account.Name  {
        err := fmt.Errorf("Account doesn't match! %s %s", aname, account.Name)
        monolog.Error("%s", err.Error())
        return nil, err
    }
        
    character.Account = account    
    return character, nil
}



// Deletes the character itself from disk
func (me * Character) Delete(dirname string) bool {
    path := SavePathFor(dirname, "character", me.ID)
    
    if err := os.Remove(path) ; err != nil {
        monolog.Warning("Could not delete character: %v %s: %s", 
            me, path, err.Error())
        return false
    }
    
    me.Account = nil
    monolog.Info("Character deleted: %s", me.ID)    
    return true
} 

func (me * Character) AskLong() string {
    return me.Long + "\n" + me.ToStatus()
}
