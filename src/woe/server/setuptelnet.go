package server

import "github.com/beoran/woe/monolog"
import t "github.com/beoran/woe/telnet"
import "github.com/beoran/woe/telnet"
import "strings"
import "strconv"



/* This file contains telnet setup helpers for the client. */


// generic negotiation

func (me * Client) SetupNegotiate(millis int, command byte, option byte, yes_event telnet.EventType, no_event telnet.EventType) (bool, telnet.Event) {
    me.telnet.TelnetSendNegotiate(command, option)
    tev, timeout, close := me.TryReadEvent(millis)

    if tev == nil || timeout || close {
        monolog.Info("Timeout or close in TryReadEvent")
        return false, nil
    }
    
    if telnet.IsEventType(tev, no_event) {
        monolog.Info("Negative event no_event %v %v", tev, no_event)
        return false, tev
    }
    
    if !telnet.IsEventType(tev, yes_event) {
        monolog.Info("Unexpected event yes_event %v %v", tev, yes_event)
        return false, tev
    }
    
    return true, tev
}

  
// Negotiate COMPRESS2 support
func (me * Client) SetupCompress2() telnet.Event {
    ok, tev := me.SetupNegotiate(1000, t.TELNET_WILL, t.TELNET_TELOPT_COMPRESS2, t.TELNET_DO_EVENT, t.TELNET_DONT_EVENT)
    if (!ok) {
        return tev
    } 
     
    me.telnet.TelnetBeginCompress2()
    monolog.Info("Client #{@id} started COMPRESS2 compression")
    me.info.compress2 = true
    return tev
}


// Negotiate NAWS (window size) support
func (me * Client) SetupNAWS() telnet.Event {  
    ok, tev := me.SetupNegotiate(1000, t.TELNET_DO, t.TELNET_TELOPT_NAWS, t.TELNET_WILL_EVENT, t.TELNET_WONT_EVENT)
    if (!ok) {
        return tev
    } 
    
    tev2, _, _ := me.TryReadEvent(1000)
    if (tev2 == nil) || (!telnet.IsEventType(tev2, t.TELNET_NAWS_EVENT)) {
        return tev2
    }
    
    nawsevent := tev.(telnet.NAWSEvent)
    me.info.w = nawsevent.W
    me.info.h = nawsevent.H
    monolog.Info("Client %d window size #{%d}x#{%d}", me.id, me.info.w, me.info.h) 
    me.info.naws     = true
    return nil
}
 
func (me * Client) SetupMSSP() telnet.Event {
    ok, tev := me.SetupNegotiate(1000, t.TELNET_WILL, t.TELNET_TELOPT_MSSP, t.TELNET_DO_EVENT, t.TELNET_DONT_EVENT)
    if (!ok) {
        return tev
    } 
    me.telnet.TelnetSendMSSP(MSSP)
    monolog.Info("Client %d accepts MSSP", me.id) 
    me.info.mssp = true
    return nil
}
 
// Check for MXP (html-like) support (but don't implement it yet)
func (me * Client) SetupMXP() telnet.Event { 

    ok, tev := me.SetupNegotiate(1000, t.TELNET_DO, t.TELNET_TELOPT_MXP, t.TELNET_WILL_EVENT, t.TELNET_WONT_EVENT)
    if (!ok) {
        return tev
    } 
    monolog.Info("Client %d accepts MXP", me.id) 
    me.info.mxp = true
    return nil
}

// Check for MSP (sound) support (but don't implement it yet)
func (me * Client) SetupMSP() telnet.Event { 

    ok, tev := me.SetupNegotiate(1000, t.TELNET_DO, t.TELNET_TELOPT_MSP, t.TELNET_WILL_EVENT, t.TELNET_WONT_EVENT)
    if (!ok) {
        return tev
    } 
    monolog.Info("Client %d accepts MSP", me.id) 
    me.info.msp = true
    return nil
}

// Check for MSDP (two way MSSP) support (but don't implement it yet)
func (me * Client) SetupMSDP() telnet.Event { 

    ok, tev := me.SetupNegotiate(1000, t.TELNET_WILL, t.TELNET_TELOPT_MSDP, t.TELNET_DO_EVENT, t.TELNET_DONT_EVENT)
    if (!ok) {
        return tev
    } 
    monolog.Info("Client %d accepts MSDP", me.id) 
    me.info.msdp = true
    return nil
}

func (me * Client) HasTerminal(name string) bool {
    for index := range me.info.terminals {
        return me.info.terminals[index] == name
    }
    return false
}




// Negotiate MTTS/TTYPE (TERMINAL TYPE) support
func (me * Client)  SetupTType() telnet.Event {
    me.info.terminals = nil
    ok, tev := me.SetupNegotiate(1000, t.TELNET_DO, t.TELNET_TELOPT_TTYPE, t.TELNET_WILL_EVENT, t.TELNET_WONT_EVENT)
    if (!ok) {
        return tev
    }
        
    var last string = "none"
    var now  string = ""
    
    for last != now {
        last = now
        me.telnet.TelnetTTypeSend()
        var tev2 telnet.Event = nil
        // Some clients (like KildClient, but not TinTin or telnet), 
        // insist on spamming useless NUL characters
        // here... So we have to retry a few times to get a ttype_is
        // throwing away any undesirable junk in between.
        for index := 0 ; index < 3 ; index++ {
            tev2, _, _ := me.TryReadEvent(1000)
        
            if tev2 != nil && telnet.IsEventType(tev2, t.TELNET_TTYPE_EVENT) {
                break
            }
        }
        
        if tev2 == nil || !telnet.IsEventType(tev2, t.TELNET_TTYPE_EVENT) {
            return tev2
        }
        
        ttypeevent := tev.(*telnet.TTypeEvent)
        now = ttypeevent.Name
        if (!me.HasTerminal(now)) {
            me.info.terminals = append(me.info.terminals, now)
        }
        me.info.terminal = now
    }
    
    monolog.Info("Client %d supports terminals %v", me.id, me.info.terminals)
    //  MTTS support
    for i := range me.info.terminals {
        term := me.info.terminals[i]
        if strings.HasPrefix(term, "MTTS ") {
            // it's an mtts terminal
            strnum := strings.TrimPrefix(term, "MTTS ")
            num, err := strconv.Atoi(strnum)
            if err != nil {
                me.info.mtts = num
                monolog.Info("Client %d supports mtts %d", me.id, me.info.mtts)                
            }
        }
    }
    me.info.ttype = true
    return nil
}

func (me * Client) SetupTelnet() {
    for {
      tev, _, _ := me.TryReadEvent(500)
      if tev != nil {
        monolog.Info("Client %d telnet setup received: %v", me.id, tev)
      } else {
        monolog.Info("Client %d no telnet setup received", me.id)
        break
      }
    }
    me.SetupMSSP()
    // me.SetupCompress2
    me.SetupNAWS()
    me.SetupTType()
    me.SetupMXP()
    me.SetupMSP()
    me.SetupMSDP()
    // color_test
}


/*  
  # Switches to "password" mode.
  def password_mode
    # The server sends "IAC WILL ECHO", meaning "I, the server, will do any 
    # echoing from now on." The client should acknowledge this with an IAC DO 
    # ECHO, and then stop putting echoed text in the input buffer. 
    # It should also do whatever is appropriate for password entry to the input 
    # box thing - for example, it might * it out. Text entered in server-echoes 
    # mode should also not be placed any command history.
    # don't use the Q state machne for echos
    @telnet.telnet_send_bytes(TELNET_IAC, TELNET_WILL, TELNET_TELOPT_ECHO)
    tev = wait_for_input(0.1)
    return tev if tev && tev.type != :do
    return nil
  end

  # Switches to "normal, or non-password mode.
  def normal_mode
    # When the server wants the client to start local echoing again, it sends 
    # "IAC WONT ECHO" - the client must respond to this with "IAC DONT ECHO".
    # Again don't use Q state machine.   
    @telnet.telnet_send_bytes(TELNET_IAC, TELNET_WONT, TELNET_TELOPT_ECHO)
    tev = wait_for_input(0.1)
    return tev if tev && tev.type != :dont
    return nil
  end
  
  def color_test
    self.write("\e[1mBold\e[0m\r\n")
    self.write("\e[3mItalic\e[0m\r\n")
    self.write("\e[4mUnderline\e[0m\r\n")
    30.upto(37) do | fg |
      self.write("\e[#{fg}mForeground Color #{fg}\e[0m\r\n")
      self.write("\e[1;#{fg}mBold Foreground Color #{fg}\e[0m\r\n")
    end  
    40.upto(47) do | bg |
      self.write("\e[#{bg}mBackground Color #{bg}\e[0m\r\n")
      self.write("\e[1;#{bg}mBold Background Color #{bg}\e[0m\r\n")
    end    
  end
  
  def setup_telnet
    loop do
      tev = wait_for_input(0.5)
      if tev
        p "setup_telnet", tev
      else
        p "no telnet setup received..."
        break
      end
    end
    setup_mssp
    setup_compress2
    setup_naws
    setup_ttype
    setup_mxp
    setup_msp
    setup_msdp
    # color_test
    
    
    #p "mssp ev #{tev}"
    # @telnet.telnet_send_negotiate(TELNET_WILL, TELNET_TELOPT_MSSP)        
    # tev = wait_for_input(0.5)
    # p "mssp ev #{tev}"
    
    # @telnet.telnet_ttype_send
    
    
  end
 
  LOGIN_RE = /\A[A-Za-z][A-Za-z0-9]*\Z/
  
  def ask_something(prompt, re, nomatch_prompt, noecho=false)
    something = nil
    
    if noecho
      password_mode
    end

    while  something.nil? || something.empty? 
      write("#{prompt}:")
      something = wait_for_command
      if something
          something.chomp!
        if re && something !~ re
          write("\r\n#{nomatch_prompt}\r\n")
          something = nil
        end
      end
    end
    
    if noecho
      normal_mode
    end
    
    something.chomp!
    return something
  end
  
  
  
  def ask_login
    return ask_something("Login", LOGIN_RE, "Login must consist of a letter followed by letters or numbers.")
  end

  EMAIL_RE = /@/

  def ask_email
    return ask_something("E-mail", EMAIL_RE, "Email must have at least an @ in there somewhere.")
  end


  def ask_password(prompt = "Password")
    return ask_something(prompt, nil, "", true) 
  end
  
  def handle_command
    order = wait_for_command
    case order
    when "/quit"
      write("Byebye!\r\n")
      @busy = false
    else
      @server.broadcast("#{@account.id} said #{order}\r\n")
    end
  end
  
  def existing_account_dialog
    pass  = ask_password
    return false unless pass
    unless @account.challenge?(pass)
      printf("Password not correct!\n")
      return false
    end
    return true
  end
  
  def new_account_dialog(login)
    while !@account 
      printf("\nWelcome, %s! Creating new account...\n", login)
      pass1  = ask_password
      return false unless pass1
      pass2 = ask_password("Repeat Password")
      return false unless pass2
      if pass1 != pass2
        printf("\nPasswords do not match! Please try again!\n")
        next
      end
      email = ask_email
      return false unless email
      @account = Woe::Account.new(:id => login, :email => email )
      @account.password   = pass1
      @account.woe_points = 7
      unless @account.save_one
        printf("\nFailed to save your account! Please contact a WOE administrator!\n")
        return false
      end
      printf("\nSaved your account.\n")
      return true
    end
  end
  
  def account_dialog
    login  = ask_login
    return false unless login
    @account = Account.fetch(login)
    if @account
      return existing_account_dialog
    else
      return new_account_dialog(login)
    end
  end
 
*/
