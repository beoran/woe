package world

import "path/filepath"
import "os"
import "encoding/xml"


type Named struct {
    Name string
}


type Account struct {
    Name              string
    PasswordHash      string
    PasswordAlgo      string
    Email             string
    WoePoints         int
    CharacterNames  []string
}

func SavePathFor(dirname string, typename string, name string) string {
    return filepath.Join(dirname, typename, name + ".xml")
} 

func NewAccount(name string, pass string, email string, points int) (*Account) {
    return &Account{name, pass, "plain", email, points, nil}
}

// Password Challenge for an account.
func (me * Account) Challenge(challenge string) bool {
    if me.PasswordAlgo == "plain" {
        return me.PasswordHash == challenge
    }
    // XXX implement encryption later
    return false
}

func (me * Account) Save(dirname string) (err error) {
    path := SavePathFor(dirname, "account", me.Name)
    
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    enc := xml.NewEncoder(file)
    enc.Indent(" ", "  ")
    return enc.Encode(me)
}

func LoadAccount(dirname string, name string) (account *Account, err error) {
    path := SavePathFor(dirname, "account", name)
    
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    dec := xml.NewDecoder(file)    
    account = new(Account)
    err = dec.Decode(account)
    return account, nil
}






