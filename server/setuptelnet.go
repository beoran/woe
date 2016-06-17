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
    
    evtype := telnet.EventTypeOf(tev)
    
    if evtype == no_event {
        monolog.Info("Negative event no_event %v %v %v", tev, evtype, no_event)
        return false, tev
    }
    
    if evtype != yes_event {
        monolog.Info("Unexpected event yes_event %v %v %v", tev, evtype, yes_event)
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
    
    nawsevent, ok := tev2.(*telnet.NAWSEvent)
    if ok {
        me.info.w = nawsevent.W
        me.info.h = nawsevent.H
        monolog.Info("Client %d window size #{%d}x#{%d}", me.id, me.info.w, me.info.h) 
        me.info.naws     = true
    }
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
    monolog.Debug("Client %d supports terminals? %s %v", me.id, name, me.info.terminals)
    for index := range me.info.terminals {
        if (me.info.terminals[index] == name) {
            return true
        }
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
        GET_TTYPE: for index := 0 ; index < 3 ; index++ {
            tev2, _, _ = me.TryReadEvent(1000)
            etyp := telnet.EventTypeOf(tev2)
            monolog.Info("Waiting for TTYPE: %T %v %d", tev2, tev2, etyp)
            if tev2 != nil && telnet.IsEventType(tev2, t.TELNET_TTYPE_EVENT) {
                monolog.Info("TTYPE received: %T %v %d", tev2, tev2, etyp)
                break GET_TTYPE
            }
        }
        
        if tev2 == nil || !telnet.IsEventType(tev2, t.TELNET_TTYPE_EVENT) {
            etyp := telnet.EventTypeOf(tev2)
            monolog.Warning("Received no TTYPE: %T %v %d", tev2, tev2, etyp)
            return tev2
        }
        
        ttypeevent := tev2.(*telnet.TTypeEvent)
        now = ttypeevent.Name
        if (!me.HasTerminal(now)) {
            me.info.terminals = append(me.info.terminals, now)
        }
        me.info.terminal = now
    }
    
    monolog.Info("Client %d supports terminals %v", me.id, me.info.terminals)
    monolog.Info("Client %d active terminal %v", me.id, me.info.terminal)

    //  MTTS support
    for i := range me.info.terminals {
        term := me.info.terminals[i]
        monolog.Info("Checking MTTS support: %s", term)
        if strings.HasPrefix(term, "MTTS ") {
            // it's an mtts terminal
            strnum := strings.TrimPrefix(term, "MTTS ")
            num, err := strconv.Atoi(strnum)
            if err == nil {
                me.info.mtts = num
                monolog.Info("Client %d supports mtts %d", me.id, me.info.mtts)                
            } else {
                monolog.Warning("Client %d could not parse mtts %s %v", me.id, strnum, err)
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
    // me.ColorTest()
}



