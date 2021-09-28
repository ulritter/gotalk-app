package main

import (
	"fmt"
	"log"
	"time"
)

func handleServerDialog(clientInputChannel <-chan ClientInput, nl Newline) {
	room := &Room{}
	for input := range clientInputChannel {
		switch event := input.event.(type) {
		case *MessageEvent:
			currentTime := time.Now()
			log.Printf(lang.Lookup(locale, "Received Message at")+" %s "+lang.Lookup(locale, "from")+" [%s]: %s"+nl.NewLine(), currentTime.Format("2006.01.02 15:04:05"), input.user.name, event.msg)
			for _, user := range room.users {
				user.session.WriteString(fmt.Sprintf("[%s]: %s", input.user.name, event.msg))
			}
		case *UserJoinedEvent:
			log.Print("User joined: "+nl.NewLine(), input.user.name)
			room.users = append(room.users, input.user)
			input.user.session.WriteStatus(fmt.Sprintf(lang.Lookup(locale, "Welcome")+" %s", input.user.name))
			for _, user := range room.users {
				if user != input.user {
					user.session.WriteStatus(fmt.Sprintf(nl.NewLine()+"%s "+lang.Lookup(locale, "entered the room"), input.user.name))
				}
			}
		case *UserLeftEvent:
			log.Printf(lang.Lookup(locale, "User left:")+" %s %s"+nl.NewLine(), event.msg, event.user.name)
			var users []*User
			for _, user := range room.users {
				if user != input.user {
					user.session.WriteStatus(fmt.Sprintf(nl.NewLine()+"%s "+lang.Lookup(locale, "left the room"), input.user.name))
					users = append(users, user)
				}
			}
			room.users = users

		case *UserChangedNickEvent:
			log.Printf(lang.Lookup(locale, "User")+" %s "+lang.Lookup(locale, "has changed his nick to:")+" %s"+nl.NewLine(), event.user.name, event.nick)
			for _, user := range room.users {
				if user != input.user {
					user.session.WriteStatus(fmt.Sprintf(nl.NewLine()+"[%s] "+lang.Lookup(locale, "has changed his nick to:")+" [%s]", event.user.name, event.nick))

				} else {
					user.session.WriteStatus(fmt.Sprintf(nl.NewLine()+lang.Lookup(locale, "You have changed your nick\nfrom")+" [%s] "+lang.Lookup(locale, "to")+" [%s]", event.user.name, event.nick))
				}
			}
			input.user.name = event.nick
		case *ListUsersEvent:
			log.Printf(lang.Lookup(locale, "User")+" %s "+lang.Lookup(locale, "has requested user list")+nl.NewLine(), input.user.name)
			input.user.session.WriteStatus(nl.NewLine() + lang.Lookup(locale, "User list:") + nl.NewLine() + lang.Lookup(locale, "=========="))

			for _, user := range room.users {
				input.user.session.WriteStatus(fmt.Sprintf("[%s] - "+lang.Lookup(locale, "joined at")+" [%s]", user.name, user.timejoined))
			}
		}
	}
}
