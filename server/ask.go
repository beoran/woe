package server

/* This file contains dialog helpers for the client. */

// import "github.com/beoran/woe/monolog"
import t "github.com/beoran/woe/telnet"
import "github.com/beoran/woe/telnet"
import "github.com/beoran/woe/world"

// import "github.com/beoran/woe/monolog"
import "bytes"
import "strings"
import "regexp"

// import "fmt"
import "strconv"

// import "strings"

// Switches to "password" mode.
func (me *Client) PasswordMode() telnet.Event {
	// The server sends "IAC WILL ECHO", meaning "I, the server, will do any
	// echoing from now on." The client should acknowledge this with an IAC DO
	// ECHO, and then stop putting echoed text in the input buffer.
	// It should also do whatever is appropriate for password entry to the input
	// box thing - for example, it might * it out. Text entered in server-echoes
	// mode should also not be placed any command history.
	// don't use the Q state machne for echos
	me.telnet.TelnetSendBytes(t.TELNET_IAC, t.TELNET_WILL, t.TELNET_TELOPT_ECHO)
	tev, _, _ := me.TryReadEvent(100)
	if tev != nil && !telnet.IsEventType(tev, t.TELNET_DO_EVENT) {
		return tev
	}
	return nil
}

// Switches to "normal, or non-password mode.
func (me *Client) NormalMode() telnet.Event {
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

func (me *Client) Printf(format string, args ...interface{}) {
	me.telnet.TelnetPrintf(format, args...)
}

func (me *Client) ColorTest() {
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
func (me *Client) ReadCommand() (something []byte) {
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

func (me *Client) AskSomething(prompt string, re string, nomatch_prompt string, noecho bool) (something []byte) {
	something = nil

	if noecho {
		me.PasswordMode()
	}

	for something == nil || len(something) == 0 {
		me.Printf("%s", prompt)
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

func (me *Client) AskYesNo(prompt string) bool {
	res := me.AskSomething(prompt+" (y/n)", "[ynYN]", "Please answer y or n.", false)
	if res[0] == 'Y' || res[0] == 'y' {
		return true
	} else {
		return false
	}
}

// Interface for an item in a list of options.
type AskOption interface {
	// Name of the option, also used to compare input
	AskName() string
	// Short description of the option, shown after name
	AskShort() string
	// Long description, displayed if "help name" is requested
	AskLong() string
	// Accound privilege required or the option to be selectable.
	AskPrivilege() world.Privilege
}

type AskOptionList interface {
	AskOptionListLen() int
	AskOptionListGet(index int) AskOption
}

type TrivialAskOption string

func (me TrivialAskOption) AskName() string {
	return string(me)
}

func (me TrivialAskOption) AskLong() string {
	return ""
}

func (me TrivialAskOption) AskShort() string {
	return ""
}

func (me TrivialAskOption) AskPrivilege() world.Privilege {
	return world.PRIVILEGE_ZERO
}

type TrivialAskOptionList []TrivialAskOption

func (me TrivialAskOptionList) AskOptionListLen() int {
	return len(me)
}

func (me TrivialAskOptionList) AskOptionListGet(index int) AskOption {
	return me[index]
}

type SimpleAskOption struct {
	Name      string
	Short     string
	Long      string
	Privilege world.Privilege
}

func (me SimpleAskOption) AskName() string {
	return me.Name
}

func (me SimpleAskOption) AskShort() string {
	return me.Short
}

func (me SimpleAskOption) AskLong() string {
	return me.Long
}

func (me SimpleAskOption) AskPrivilege() world.Privilege {
	return me.Privilege
}

type SimpleAskOptionList []SimpleAskOption

func (me SimpleAskOptionList) AskOptionListLen() int {
	return len(me)
}

func (me SimpleAskOptionList) AskOptionListGet(index int) AskOption {
	return me[index]
}

type JobListAsker []world.Job

func (me JobListAsker) AskOptionListLen() int {
	return len(me)
}

func (me JobListAsker) AskOptionListGet(index int) AskOption {
	return me[index]
}

type KinListAsker []world.Kin

func (me KinListAsker) AskOptionListLen() int {
	return len(me)
}

func (me KinListAsker) AskOptionListGet(index int) AskOption {
	return me[index]
}

type GenderListAsker []world.Gender

func (me GenderListAsker) AskOptionListLen() int {
	return len(me)
}

func (me GenderListAsker) AskOptionListGet(index int) AskOption {
	return me[index]
}

type AskOptionSlice []AskOption
type AskOptionFilterFunc func(AskOption) AskOption

func (me AskOptionSlice) AskOptionListLen() int {
	return len(me)
}

func (me AskOptionSlice) AskOptionListGet(index int) AskOption {
	return me[index]
}

func AskOptionListEach(me AskOptionList, cb AskOptionFilterFunc) AskOption {
	for i := 0; i < me.AskOptionListLen(); i++ {
		res := cb(me.AskOptionListGet(i))
		if res != nil {
			return res
		}
	}
	return nil
}

func AskOptionListFilter(me AskOptionList, cb AskOptionFilterFunc) AskOptionList {
	result := make(AskOptionSlice, 0)
	for i := 0; i < me.AskOptionListLen(); i++ {
		res := cb(me.AskOptionListGet(i))
		if res != nil {
			result = append(result, res)
		}
	}
	return result
}

/* Finds the name irrespecful of the case */
func AskOptionListFindName(me AskOptionList, name string) AskOption {
	return AskOptionListEach(me, func(e AskOption) AskOption {
		if strings.ToLower(e.AskName()) == strings.ToLower(name) {
			return e
		} else {
			return nil
		}
	})
}

/* Filters the list by privilege level (only those allowed by the level are retained) */
func AskOptionListFilterPrivilege(me AskOptionList, privilege world.Privilege) AskOptionList {
	return AskOptionListFilter(me, func(e AskOption) AskOption {
		if e.AskPrivilege() <= privilege {
			return e
		} else {
			return nil
		}
	})
}

type MergedAskOptionList struct {
	head AskOptionList
	tail AskOptionList
}

func (me *MergedAskOptionList) AskOptionListLen() int {
	return me.head.AskOptionListLen() + me.tail.AskOptionListLen()
}

func (me *MergedAskOptionList) AskOptionListGet(index int) AskOption {
	headlen := me.head.AskOptionListLen()
	if index < headlen {
		return me.head.AskOptionListGet(index)
	}
	return me.tail.AskOptionListGet(index - headlen)
}

/* Merges two AskOptionLists without copying using the MergedAskOptionList rtype  */
func AskOptionListMerge(me AskOptionList, you AskOptionList) AskOptionList {
	return &MergedAskOptionList{me, you}
}

func (me AskOptionSlice) Each(cb AskOptionFilterFunc) AskOption {
	for i := 0; i < len(me); i++ {
		res := cb(me[i])
		if res != nil {
			return res
		}
	}
	return nil
}

func (me AskOptionSlice) Filter(cb AskOptionFilterFunc) AskOptionSlice {
	result := make(AskOptionSlice, 0)
	for i := 0; i < len(me); i++ {
		res := cb(me[i])
		if res != nil {
			result = append(result, res)
		}
	}
	return result
}

/* Finds the name irrespecful of the case */
func (me AskOptionSlice) FindName(name string) AskOption {
	return me.Each(func(e AskOption) AskOption {
		if strings.ToLower(e.AskName()) == strings.ToLower(name) {
			return e
		} else {
			return nil
		}
	})
}

/* Filters the list by privilege level (only those allowed by the level are retained) */
func (me AskOptionSlice) FilterPrivilege(privilege world.Privilege) AskOptionSlice {
	return me.Filter(func(e AskOption) AskOption {
		if e.AskPrivilege() <= privilege {
			return e
		} else {
			return nil
		}
	})
}

func (me *Client) AskOptionListHelp(alist AskOptionList, input []byte) {
	re := regexp.MustCompile("[^ \t,]+")
	argv := re.FindAll(input, -1)
	if len(argv) < 2 {
		me.Printf("Help usage: help <topic>.\n")
		return
	}
	e := AskOptionListFindName(alist, string(argv[1]))
	if e == nil {
		me.Printf("Cannot find topic %s in list. No help available.\n", string(argv[1]))
	} else {
		al := e.AskLong()
		if al == "" {
			me.Printf("Topic %s found, but help is unavailable.\n", string(argv[1]))
		} else {
			me.Printf("Help on %s:\n%s\n", string(argv[1]), e.AskLong())
		}
	}
}

func (me *Client) AskOptionListOnce(heading string, prompt string, noecho bool, alist AskOptionList) (result AskOption) {
	list := AskOptionListFilterPrivilege(alist, me.account.Privilege)
	me.Printf("\n%s\n\n", heading)
	for i := 0; i < list.AskOptionListLen(); i++ {
		v := list.AskOptionListGet(i)
		sh := v.AskShort()
		if sh == "" {
			me.Printf("[%d] %s\n", i+1, v.AskName())
		} else {
			me.Printf("[%d] %s: %s\n", i+1, v.AskName(), sh)
		}
	}
	me.Printf("\n")
	aid := me.AskSomething(prompt, "", "", false)
	iresp, err := strconv.Atoi(string(aid))
	if err != nil { /* Try by name if not a number. */
		e := AskOptionListFindName(alist, string(aid))
		if e != nil {
			return e
		} else if ok, _ := regexp.Match("help", bytes.ToLower(aid)); ok {
			me.AskOptionListHelp(list, aid)
		} else {
			me.Printf("Name not found in list. Please choose a number or name from the list above. Or type help <option> for help on that option.\n")
		}
	} else if (iresp > 0) && (iresp <= list.AskOptionListLen()) {
		/* In range of list. */
		return list.AskOptionListGet(iresp - 1)
	} else {
		me.Printf("Please choose a number or name from the list above.\n")
	}
	return nil
}

func (me *Client) AskOptionList(
	heading string, prompt string, noecho bool,
	noconfirm bool, list AskOptionList) (result AskOption) {
	for {
		result = me.AskOptionListOnce(heading, prompt, noecho, list)
		if result != nil {
			if noconfirm || me.AskYesNo(heading+"\nConfirm "+result.AskName()+"? ") {
				return result
			}
		}
	}
}

func (me *Client) AskOptionListExtra(heading string,
	prompt string, noecho bool, noconfirm bool, list AskOptionList,
	extra AskOptionList) (result AskOption) {
	xlist := AskOptionListMerge(list, extra)
	return me.AskOptionList(heading, prompt, noecho, noconfirm, xlist)
}

func (me *Client) AskEntityListOnce(
	heading string, prompt string, noecho bool,
	elist world.EntitylikeSlice, extras []string) (result world.Entitylike, alternative string) {
	list := elist.FilterPrivilege(me.account.Privilege)
	me.Printf("\n%s\n\n", heading)
	last := 0
	for i, v := range list {
		e := v.AsEntity()
		me.Printf("[%d] %s: %s\n", i+1, e.Name, e.Short)
		last = i + 1
	}

	if extras != nil {
		for i, v := range extras {
			me.Printf("[%d] %s\n", last+i+1, v)
		}
	}

	me.Printf("\n")
	aid := me.AskSomething(prompt, "", "", false)
	iresp, err := strconv.Atoi(string(aid))
	if err != nil { /* Try by name if not a number. */
		e := list.FindName(string(aid))
		if e != nil {
			return e, ""
		} else {
			if extras != nil {
				for _, v := range extras {
					if strings.ToLower(v) == strings.ToLower(string(aid)) {
						return nil, v
					}
				}
			}
			me.Printf("Name not found in list. Please choose a number or name from the list above.\n")
		}
	} else if (iresp > 0) && (iresp <= len(list)) { /* In range of list. */
		return list[iresp-1], ""
	} else if (extras != nil) && (iresp > last) && (iresp <= last+len(extras)) {
		return nil, extras[iresp-last-1]
	} else {
		me.Printf("Please choose a number or name from the list above.\n")
	}
	return nil, ""
}

func (me *Client) AskEntityList(
	heading string, prompt string, noecho bool,
	noconfirm bool, list world.EntitylikeSlice, extras []string) (result world.Entitylike, alternative string) {
	for {
		result, alternative = me.AskEntityListOnce(heading, prompt, noecho, list, extras)
		if result != nil {
			e := result.AsEntity()
			if !noconfirm {
				me.Printf("\n%s: %s\n\n%s\n\n", e.Name, e.Short, e.Long)
			}
			if noconfirm || me.AskYesNo(heading+"\nConfirm "+e.Name+"? ") {
				return result, ""
			}
		} else if alternative != "" {
			if noconfirm || me.AskYesNo("Confirm "+alternative+" ?") {
				return result, alternative
			}
		}
	}
}

const LOGIN_RE = "^[A-Za-z][A-Za-z0-9]*$"

func (me *Client) AskLogin() []byte {
	return me.AskSomething("Login?>", LOGIN_RE, "Login must consist of letters followed by letters or numbers.", false)
}

const EMAIL_RE = "@"

func (me *Client) AskEmail() []byte {
	return me.AskSomething("E-mail?>", EMAIL_RE, "Email must have at least an @ in there somewhere.", false)
}

const CHARNAME_RE = "^[A-Z][A-Za-z]+$"

func (me *Client) AskCharacterName() []byte {
	return me.AskSomething("Character Name?>", CHARNAME_RE, "Character name consist of a capital letter followed by at least one letter.", false)
}

func (me *Client) AskPassword() []byte {
	return me.AskSomething("Password?>", "", "", true)
}

func (me *Client) AskRepeatPassword() []byte {
	return me.AskSomething("Repeat Password?>", "", "", true)
}

func (me *Client) HandleCommand() {
	command := me.ReadCommand()
	me.ProcessCommand(command)
}
