package server

/* This file contains dialog helpers for the client. */

// import "github.com/beoran/woe/monolog"
import t "github.com/beoran/woe/telnet"
import "github.com/beoran/woe/telnet"
import "github.com/beoran/woe/world"
import "github.com/beoran/woe/monolog"
import "bytes"
import "regexp"
// import "fmt"
import "strconv"
// import "strings"

const NEW_CHARACTER_PRICE = 4
  
// Switches to "password" mode.
func (me * Client) PasswordMode() telnet.Event {
// The server sends "IAC WILL ECHO", meaning "I, the server, will do any 
// echoing from now on." The client should acknowledge this with an IAC DO 
// ECHO, and then stop putting echoed text in the input buffer. 
// It should also do whatever is appropriate for password entry to the input 
// box thing - for example, it might * it out. Text entered in server-echoes 
// mode should also not be placed any command history.
// don't use the Q state machne for echos
    me.telnet.TelnetSendBytes(t.TELNET_IAC, t.TELNET_WILL, t.TELNET_TELOPT_ECHO)
    tev, _, _:= me.TryReadEvent(100)
    if tev != nil && !telnet.IsEventType(tev, t.TELNET_DO_EVENT) { 
        return tev
    }
    return nil
}

// Switches to "normal, or non-password mode.
func (me * Client) NormalMode() telnet.Event {
// When the server wants the client to start local echoing again, it s}s 
// "IAC WONT ECHO" - the client must respond to this with "IAC DONT ECHO".
// Again don't use Q state machine.   
    me.telnet.TelnetSendBytes(t.TELNET_IAC, t.TELNET_WONT, t.TELNET_TELOPT_ECHO)
    tev, _, _ := me.TryReadEvent(100)
    if tev != nil && !telnet.IsEventType(tev, t.TELNET_DONT_EVENT) { 
        return tev
    }
    return nil
}

func (me * Client) Printf(format string, args ...interface{}) {
    me.telnet.TelnetPrintf(format, args...)
}

func (me * Client) ColorTest() {
    me.Printf("\033[1mBold\033[0m\r\n")
    me.Printf("\033[3mItalic\033[0m\r\n")
    me.Printf("\033[4mUnderline\033[0m\r\n")
    for fg := 30; fg < 38; fg++ {
        me.Printf("\033[%dmForeground Color %d\033[0m\r\n", fg, fg)
        me.Printf("\033[1;%dmBold Foreground Color %d\033[0m\r\n", fg, fg)
    }
    
    for bg := 40; bg < 48; bg++ {
        me.Printf("\033[%dmBackground Color %d\033[0m\r\n", bg, bg)
        me.Printf("\033[1;%dmBold Background Color %d\033[0m\r\n", bg, bg)
    }    
}


// Blockingly reads a single command from the client
func (me * Client) ReadCommand() (something []byte) {
    something = nil
    for something == nil { 
        something, _, _ = me.TryRead(-1)
        if something != nil {
            something = bytes.TrimRight(something, "\r\n")
            return something
        }
    } 
    return nil       
}

func (me * Client) AskSomething(prompt string, re string, nomatch_prompt string, noecho bool) (something []byte) {
    something = nil
    
    if noecho {
      me.PasswordMode()
    }

    for something == nil || len(something) == 0 { 
        me.Printf("%s:", prompt)
        something, _, _ = me.TryRead(-1)
        if something != nil {
            something = bytes.TrimRight(something, "\r\n")
            if len(re) > 0 {
                ok, _ := regexp.Match(re, something)
                if !ok {
                    me.Printf("\n%s\n", nomatch_prompt)
                    something = nil
                }
            }
        }
    }
    
    if noecho {
      me.NormalMode()
      me.Printf("\n")
    }
        
    return something
}

func (me * Client) AskYesNo(prompt string) bool {
    res := me.AskSomething(prompt + " (y/n)","[ynYN]", "Please answer y or n.", false)
    if res[0] == 'Y'|| res[0] == 'y' { 
        return true
    } else {
        return false
    }    
}

func (me * Client) AskEntityListOnce(heading string, prompt string, noecho bool, elist world.EntitylikeSlice) (result world.Entitylike) { 
    list := elist.FilterPrivilege(me.account.Privilege)
    me.Printf("\n%s\n\n",heading)
    for i, v := range(list) {
        e := v.AsEntity()
        me.Printf("[%d] %s: %s\n", i+1, e.Name, e.Short)
    }
    me.Printf("\n")
    aid := me.AskSomething(prompt, "", "", false);
    iresp, err := strconv.Atoi(string(aid))
    if err != nil { /* Try name. */
        e := list.FindName(string(aid))
        if e != nil {
            return e
        } else {
            me.Printf("Name not found in list. Please choose a number or name from the list above.\n")
        }
    } else if (iresp>0) && (iresp<=len(list)) { /* In range. */
        return list[iresp-1]
    } else {
        me.Printf("Please choose a number or name from the list above.\n")
    }
    return nil
}
    

func (me * Client) AskEntityList(heading string, prompt string, noecho bool, list world.EntitylikeSlice) (result world.Entitylike) {     
    for {
        result = me.AskEntityListOnce(heading, prompt, noecho, list)
        if result != nil {
            e := result.AsEntity()
            me.Printf("\n%s: %s\n\n%s\n\n", e.Name, e.Short, e.Long)
            if noecho || me.AskYesNo("Confirm?") {
                return result
            }
        }
    }
}
    

const LOGIN_RE = "^[A-Za-z]+$"

func (me * Client) AskLogin() []byte {
    return me.AskSomething("Login", LOGIN_RE, "Login must consist of a letter followed by letters or numbers.", false)
}


const EMAIL_RE = "@"

func (me * Client) AskEmail() []byte {
    return me.AskSomething("E-mail", EMAIL_RE, "Email must have at least an @ in there somewhere.", false)
}

func (me * Client) AskCharacterName() []byte {
    return me.AskSomething("Character Name", LOGIN_RE, "Character name consisst of letters only.", false)
}

func (me * Client) AskPassword() []byte {
    return me.AskSomething("Password", "", "", true)
}

func (me * Client) AskRepeatPassword() []byte {
    return me.AskSomething("Repeat Password", "", "", true)
}

func (me * Client) HandleCommand() {
    command := me.ReadCommand()
    me.ProcessCommand(command)
}
 
func (me * Client) ExistingAccountDialog() bool {
    pass  := me.AskPassword()
    for pass == nil {
        me.Printf("Password may not be empty!\n")        
        pass  = me.AskPassword()
    }
    
    if !me.account.Challenge(string(pass)) {
        me.Printf("Password not correct!\n")
        me.Printf("Disconnecting!\n")
        return false
    }    
    return true
}

func (me * Client) NewAccountDialog(login string) bool {
    for me.account == nil {    
      me.Printf("\nWelcome, %s! Creating new account...\n", login)
      pass1  := me.AskPassword()
      
      if pass1 == nil { 
          return false
      }
      
      pass2 := me.AskRepeatPassword()
      
      if pass1 == nil { 
          return false
      }
      
      if string(pass1) != string(pass2) {
        me.Printf("\nPasswords do not match! Please try again!\n")
        continue
      }
      
      email := me.AskEmail()
      if email == nil { return false  }
      
      me.account = world.NewAccount(login, string(pass1), string(email), 7)
      err      := me.account.Save(me.server.DataPath())
      
      if err != nil {      
        monolog.Error("Could not save account %s: %v", login, err)  
        me.Printf("\nFailed to save your account!\nPlease contact a WOE administrator!\n")
        return false
      }
      
      monolog.Info("Created new account %s", login)  
      me.Printf("\nSaved your account.\n")
      return true
    }
    return false
}
  
func (me * Client) AccountDialog() bool {
    login  := me.AskLogin()
    if login == nil { return false }
    var err error
    
    if me.server.World.GetAccount(string(login)) != nil {
        me.Printf("Account already logged in!\n")
        me.Printf("Disconnecting!\n")
        return false 
    }
    
    me.account, err = me.server.World.LoadAccount(string(login))    
    if err != nil {
        monolog.Warning("Could not load account %s: %v", login, err)  
    }
    if me.account != nil {
      return me.ExistingAccountDialog()
    } else {
      return me.NewAccountDialog(string(login))
    }
}

func (me * Client) NewCharacterDialog() bool {
    me.Printf("New character:\n")
    charname := me.AskCharacterName()
    
    kin := me.AskEntityList("Please choose the kin of this character", "Kin of character? ", false, world.KinEntityList)
    me.Printf("%s %v\n", charname, kin)
     
    gender := me.AskEntityList("Please choose the gender of this character", "Gender? ", false, world.GenderList)
    me.Printf("%s %v\n", charname, gender)

    job := me.AskEntityList("Please choose the job of this character", "Job? ", false, world.JobEntityList)
    me.Printf("%s %v\n", charname, job)
    
    character := world.NewCharacter(me.account, 
                    string(charname), kin, gender, job)
    
    me.Printf("%s", character.Being.ToStatus());
    
    ok := me.AskYesNo("Is this character ok?")
    
    if (!ok) {
        me.Printf("Character creation canceled.\n")
        return true
    }
    
    me.account.AddCharacter(character)
    me.account.Points -= NEW_CHARACTER_PRICE 
    me.account.Save(me.server.DataPath())
    character.Save(me.server.DataPath())
    me.Printf("Character %s saved.\n", character.Being.Name)
    

    return true
}    

 
func (me * Client) CharacterDialog() bool {
    me.Printf("You have %d remaining points.\n", me.account.Points)
    for me.account.NumCharacters() < 1 {
        me.Printf("You have no characters yet!\n")
        if (me.account.Points > 0) {
            me.NewCharacterDialog();
        } else {
            me.Printf("Sorry, you have no points left to make new characters!\n")
            me.Printf("Please contact the staff of WOE if you think this is a mistake.\n")
            me.Printf("Disconnecting!\n")
            return false 
        }
    }
    return true
}
 

