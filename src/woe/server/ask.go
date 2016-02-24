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
// import "strconv"


  
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
  

const LOGIN_RE = "^[A-Za-z][A-Za-z0-9]+$"

func (me * Client) AskLogin() []byte {
    return me.AskSomething("Login", LOGIN_RE, "Login must consist of a letter followed by letters or numbers.", false)
}

const EMAIL_RE = "@"

func (me * Client) AskEmail() []byte {
    return me.AskSomething("E-mail", EMAIL_RE, "Email must have at least an @ in there somewhere.", false)
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
    /*
    if bytes.HasPrefix(command, []byte("/quit")) {
      me.Printf("Byebye!\n")
      me.alive = false
    } else if bytes.HasPrefix(command, []byte("/shutdown")) {
      me.server.Broadcast("Shutting down server NOW!\n")
      me.server.Shutdown();
    } else if bytes.HasPrefix(command, []byte("/restart")) {
      me.server.Broadcast("Restarting down server NOW!\n")
      me.server.Restart();
    } else {
      me.server.Broadcast("Client %d said %s\r\n", me.id, command)  
    }
    */
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
    
    me.account, err = me.server.World.LoadAccount(me.server.DataPath(), string(login))    
    if err != nil {
        monolog.Warning("Could not load account %s: %v", login, err)  
    }
    if me.account != nil {
      return me.ExistingAccountDialog()
    } else {
      return me.NewAccountDialog(string(login))
    }
}
 
func (me * Client) CharacterDialog() bool {
    login  := me.AskLogin()
    if login == nil { return false }
    var err error
    me.account, err = world.LoadAccount(me.server.DataPath(), string(login))    
    if err != nil {
        monolog.Warning("Could not load account %s: %v", login, err)  
    }
    if me.account != nil {
      return me.ExistingAccountDialog()
    } else {
      return me.NewAccountDialog(string(login))
    }
}

