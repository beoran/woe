package world

import "path/filepath"
import "github.com/beoran/woe/sitef"
import "github.com/beoran/woe/monolog"
import "fmt"
import "errors"

type Privilege int

const (
    PRIVILEGE_ZERO        = Privilege(iota * 100)
    PRIVILEGE_NORMAL
    PRIVILEGE_MASTER
    PRIVILEGE_LORD
    PRIVILEGE_IMPLEMENTOR
)


type Named struct {
    Name string
}


type Account struct {
    Name              string
    Hash              string
    Algo              string
    Email             string
    Points            int
    Privilege         Privilege
    CharacterNames  []string
    characters      [] * Character
}

func SavePathForXML(dirname string, typename string, name string) string {
    return filepath.Join(dirname, typename, name + ".xml")
} 

func SavePathFor(dirname string, typename string, name string) string {
    return filepath.Join(dirname, typename, name + ".sitef")
}



func NewAccount(name string, pass string, email string, points int) (*Account) {
    return &Account{name, pass, "plain", email, points, PRIVILEGE_NORMAL, nil, nil}
}

// Password Challenge for an account.
func (me * Account) Challenge(challenge string) bool {
    if me.Algo == "plain" {
        return me.Hash == challenge
    }
    // XXX implement encryption later
    return false
}


// Add a character to an account.
func (me * Account) AddCharacter(chara * Character) {
    me.characters = append(me.characters, chara)
}


// Save an account as a sitef file.
func (me * Account) Save(dirname string) (err error) {
    path := SavePathFor(dirname, "account", me.Name)
    
    rec                := sitef.NewRecord()
    rec.Put("name",         me.Name)
    rec.Put("hash",         me.Hash)
    rec.Put("algo",         me.Algo)
    rec.Put("email",        me.Email)
    rec.PutInt("points",    me.Points)
    rec.PutInt("privilege", int(me.Privilege))
    rec.PutInt("characters",len(me.characters))
    for i, chara   := range me.characters {
        key        := fmt.Sprintf("characters[%d]", i)
        rec.Put(key, chara.ID)
        
    }
    monolog.Debug("Saving Acccount record: %s %v", path, rec)
    return sitef.SaveRecord(path, *rec)
}

// Load an account from a sitef file.
func LoadAccount(dirname string, name string) (account *Account, err error) {
    
    path := SavePathFor(dirname, "account", name)
    
    records, err := sitef.ParseFilename(path)
    if err != nil {
        return nil, err
    }
    
    if len(records) < 1 {
        return nil, errors.New("No record found!")
    }
    
    record := records[0]
    monolog.Info("Loading Account record: %s %v", path, record)
    
    account = new(Account)
    account.Name            = record.Get("name")
    account.Hash            = record.Get("hash")
    account.Algo            = record.Get("algo")
    account.Email           = record.Get("email")
    account.Points          = record.GetIntDefault("points", 0)
    account.Privilege       = Privilege(record.GetIntDefault("privilege", 
                                int(PRIVILEGE_NORMAL)))
    
    nchars                 := record.GetIntDefault("characters", 0)
    account.characters      = make([] * Character, 0, nchars)
    monolog.Info("Try to load %d characters:\n", nchars)
    for index := 0 ; index < nchars ; index ++ {

        chid := record.GetArrayIndex("characters", index)
        monolog.Info("Loading character: %d %s\n", index, chid)
        
        ch, err := account.LoadCharacter(dirname, chid);
        if err != nil {
            monolog.Error("Could not load character %s: %s", chid, err.Error())
            // return nil, err
        } else {
            account.characters = append(account.characters, ch)
        } 
    }
    
    
    /* Todo: load characters here... */    
    monolog.Info("Loaded Account: %s %v", path, account)
    return account, nil
}

 
func (me * Account) NumCharacters() int {
    return len(me.characters)
} 

func (me * Account) GetCharacter(index int) (* Character) {
    return me.characters[index]
} 

func (me * Account) FindCharacter(character * Character) (index int) {
    for k, c := range me.characters {
        if c == character  {
            return k
        }
    }
    return -1;
} 

// Delete a character from this account.
func (me * Account) DeleteCharacter(dirname string, character * Character) bool {
    
    if i:= me.FindCharacter(character) ; i < 0 {
        monolog.Warning("Could not find character: %v %d", character, i)
        return false;  
    } else {
        copy(me.characters[i:], me.characters[i+1:])
        newlen := len(me.characters) - 1 
        me.characters[newlen] = nil
        me.characters = me.characters[:newlen]
    }
    /// Save self so the deletion is correctly recorded.
    me.Save(dirname)
    
    return character.Delete(dirname)
} 


func (me * Account) CharacterEntitylikeSlice() EntitylikeSlice {
    els := make(EntitylikeSlice, 0, 16)
    for i:= 0 ; i < me.NumCharacters(); i++ {
        chara := me.GetCharacter(i)
        els = append(els, chara)
    }
    return els
}





