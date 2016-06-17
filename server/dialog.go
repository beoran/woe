// dialog.go
package server

/* This file contains dialogs for the client. The dialog helpers are in ask.go. */

import "github.com/beoran/woe/monolog"
import "github.com/beoran/woe/world"

const NEW_CHARACTER_PRICE = 4

func (me *Client) ExistingAccountDialog() bool {
	pass := me.AskPassword()
	for pass == nil {
		me.Printf("Password may not be empty!\n")
		pass = me.AskPassword()
	}

	if !me.account.Challenge(string(pass)) {
		me.Printf("Password not correct!\n")
		me.Printf("Disconnecting!\n")
		return false
	}
	return true
}

func (me *Client) NewAccountDialog(login string) bool {
	for me.account == nil {
		me.Printf("\nWelcome, %s! Creating new account...\n", login)
		pass1 := me.AskPassword()

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
		if email == nil {
			return false
		}

		me.account = world.NewAccount(login, string(pass1), string(email), 7)
		err := me.account.Save(me.server.DataPath())

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

func (me *Client) AccountDialog() bool {
	login := me.AskLogin()
	if login == nil {
		return false
	}
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

func (me *Client) NewCharacterDialog() bool {
	noconfirm := true
	extra := TrivialAskOptionList{TrivialAskOption("Cancel")}

	me.Printf("New character:\n")
	charname := me.AskCharacterName()

	existing, aname, _ := world.LoadCharacterByName(me.server.DataPath(), string(charname))

	for existing != nil {
		if aname == me.account.Name {
			me.Printf("You already have a character with a similar name!\n")
		} else {
			me.Printf("That character name is already taken by someone else.\n")
		}
		charname := me.AskCharacterName()
		existing, aname, _ = world.LoadCharacterByName(me.server.DataPath(), string(charname))
	}

	kinres := me.AskOptionListExtra("Please choose the kin of this character", "Kin?> ", false, noconfirm, KinListAsker(world.KinList), extra)

	if sopt, ok := kinres.(TrivialAskOption); ok {
		if string(sopt) == "Cancel" {
			me.Printf("Character creation canceled.\n")
			return true
		} else {
			return true
		}
	}

	kin := kinres.(world.Kin)

	genres := me.AskOptionListExtra("Please choose the gender of this character", "Gender?> ", false, noconfirm, GenderListAsker(world.GenderList), extra)
	if sopt, ok := kinres.(TrivialAskOption); ok {
		if string(sopt) == "Cancel" {
			me.Printf("Character creation canceled.\n")
			return true
		} else {
			return true
		}
	}

	gender := genres.(world.Gender)

	jobres := me.AskOptionListExtra("Please choose the job of this character", "Gender?> ", false, noconfirm, JobListAsker(world.JobList), extra)
	if sopt, ok := kinres.(TrivialAskOption); ok {
		if string(sopt) == "Cancel" {
			me.Printf("Character creation canceled.\n")
			return true
		} else {
			return true
		}
	}

	job := jobres.(world.Job)

	character := world.NewCharacter(me.account,
		string(charname), &kin, &gender, &job)

	me.Printf("%s", character.Being.ToStatus())

	ok := me.AskYesNo("Is this character ok?")

	if !ok {
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

func (me *Client) DeleteCharacterDialog() bool {
	extra := []AskOption{TrivialAskOption("Cancel"), TrivialAskOption("Disconnect")}

	els := me.AccountCharacterList()
	els = append(els, extra...)
	result := me.AskOptionList("Character to delete?",
		"Character?>", false, false, els)

	if alt, ok := result.(TrivialAskOption); ok {
		if string(alt) == "Disconnect" {
			me.Printf("Disconnecting")
			return false
		} else if string(alt) == "Cancel" {
			me.Printf("Canceled")
			return true
		} else {
			monolog.Warning("Internal error, unhandled case.")
			return true
		}
	}

	character := result.(*world.Character)
	/* A character that is deleted gives NEW_CHARACTER_PRICE +
	 * level / (NEW_CHARACTER_PRICE * 2) points, but only after the delete. */
	np := NEW_CHARACTER_PRICE + character.Level/(NEW_CHARACTER_PRICE*2)
	me.account.DeleteCharacter(me.server.DataPath(), character)
	me.account.Points += np

	return true
}

func (me *Client) AccountCharacterList() AskOptionSlice {
	els := make(AskOptionSlice, 0, 16)
	for i := 0; i < me.account.NumCharacters(); i++ {
		chara := me.account.GetCharacter(i)
		els = append(els, chara)
	}
	return els
}

func (me *Client) ChooseCharacterDialog() bool {
	extra := []AskOption{
		SimpleAskOption{"New", "Create New character",
			"Create a new character. This option costs 4 points.",
			world.PRIVILEGE_ZERO},
		SimpleAskOption{"Disconnect", "Disconnect from server",
			"Disconnect your client from this server.",
			world.PRIVILEGE_ZERO},
		SimpleAskOption{"Delete", "Delete character",
			"Delete a character. A character that has been deleted cannot be reinstated. You will receive point bonuses for deleting your characters that depend on their level.",
			world.PRIVILEGE_ZERO},
	}

	var pchara *world.Character = nil

	for pchara == nil {
		els := me.AccountCharacterList()
		els = append(els, extra...)
		result := me.AskOptionList("Choose a character?", "Character?>", false, true, els)
		switch opt := result.(type) {
		case SimpleAskOption:
			if opt.Name == "New" {
				if me.account.Points >= NEW_CHARACTER_PRICE {
					if !me.NewCharacterDialog() {
						return false
					}
				} else {
					me.Printf("Sorry, you have no points left to make new characters!\n")
				}
			} else if opt.Name == "Disconnect" {
				me.Printf("Disconnecting\n")
				return false
			} else if opt.Name == "Delete" {
				if !me.DeleteCharacterDialog() {
					return false
				}
			} else {
				me.Printf("Internal error, alt not valid: %v.", opt)
			}
		case *world.Character:
			pchara = opt
		default:
			me.Printf("What???")
		}
		me.Printf("You have %d points left.\n", me.account.Points)
	}

	me.character = pchara
	me.Printf("%s\n", me.character.Being.ToStatus())
	me.Printf("Welcome, %s!\n", me.character.Name)

	return true
}

func (me *Client) CharacterDialog() bool {
	me.Printf("You have %d remaining points.\n", me.account.Points)
	for me.account.NumCharacters() < 1 {
		me.Printf("You have no characters yet!\n")
		if me.account.Points >= NEW_CHARACTER_PRICE {
			if !me.NewCharacterDialog() {
				return false
			}
		} else {
			me.Printf("Sorry, you have no characters, and no points left to make new characters!\n")
			me.Printf("Please contact the staff of WOE if you think this is a mistake.\n")
			me.Printf("Disconnecting!\n")
			return false
		}
	}

	return me.ChooseCharacterDialog()
}
