package world

import "path/filepath"
import "os"
import "encoding/xml"
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

// Helpers for use with sitef records
func SitefStoreString(rec sitef.Record, key string, val string) {
    rec[key] = val
}

func SitefStoreInt(rec sitef.Record, key string, val int) {
    rec[key] = fmt.Sprintf("%d", val)
}


func SitefStoreArray(rec sitef.Record, key string, val LabeledList) {

}


// Save an account as a sitef file.
func (me * Account) Save(dirname string) (err error) {
    path := SavePathFor(dirname, "account", me.Name)
    
    rec                := make(sitef.Record)
    rec.Put("name",         me.Name)
    rec.Put("hash",         me.Hash)
    rec.Put("algo",         me.Algo)
    rec.Put("email",        me.Email)
    rec.PutInt("points",    me.Points)
    rec.PutInt("privilege", int(me.Privilege))
    rec.PutInt("characters",len(me.characters))
    for i, chara   := range me.characters {
        key        := fmt.Sprintf("characters[%d]", i)
        rec.Put(key, chara.Name)
    }
    monolog.Debug("Saving Acccount record: %s %v", path, rec)
    return sitef.SaveRecord(path, rec)
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
    
    var nchars int
    nchars                  = record.GetIntDefault("characters", 0)
    _ = nchars
    /* Todo: load characters here... */    
    monolog.Info("Loaded Account: %s %v", path, record)
    return account, nil
}

 

func (me * Account) SaveXML(dirname string) (err error) {
    path := SavePathForXML(dirname, "account", me.Name)
    
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    enc := xml.NewEncoder(file)
    enc.Indent(" ", "  ")
    return enc.Encode(me)
}

func LoadAccountXML(dirname string, name string) (account *Account, err error) {
    path := SavePathForXML(dirname, "account", name)
    
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    dec := xml.NewDecoder(file)    
    account = new(Account)
    err = dec.Decode(account)
    return account, nil
}






