package world

import "math/rand"
import "crypto/sha1"
import "fmt"



const MAKE_SALT_AID string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789."


func MakeSalt() string {    
    c1 := MAKE_SALT_AID[rand.Intn(len(MAKE_SALT_AID))]
    c2 := MAKE_SALT_AID[rand.Intn(len(MAKE_SALT_AID))]
    res := string( []byte{ c1, c2 })
    return res
}

func WoeCryptPassword(password string, salt string) string {
    if len(salt) < 1 {
        salt = MakeSalt()
    }
    to_hash := salt + password
    return salt + fmt.Sprintf("%x", sha1.Sum([]byte(to_hash)))
}


func WoeCryptChallenge(hash, trypass string) bool {
    salt := hash[0:2]
    try  := WoeCryptPassword(trypass, salt)
    return try == hash
}




