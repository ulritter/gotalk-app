package main

import (
	"crypto/tls"
	"log"
	"net"
	"strings"
	"time"
)

func handleConnection(conn net.Conn, inputChannel chan ClientInput, nl Newline) error {
	buf := make([]byte, BUFSIZE)

	session := &Session{conn}
	n, err := conn.Read(buf)
	if err != nil {
		log.Print(lang.Lookup(locale, "Error reading from buffer")+nl.NewLine(), err)
		return err
	}
	var nick string
	pattern := string(buf[:n])
	if (pattern[0] == CMD_ESCAPE_CHAR) && (pattern[n-1] == CMD_ESCAPE_CHAR) {
		nick = string(buf[1 : n-1])
	} else {
		nick = "J_Doe"
	}

	user := &User{name: nick, session: session, timejoined: time.Now().Format("2006.01.02 15:04:05")}
	inputChannel <- ClientInput{
		user,
		&UserJoinedEvent{},
	}

	for {
		n, err := conn.Read(buf)

		if (buf[0] == CMD_ESCAPE_CHAR) || (err != nil) {
			pattern := strings.Fields(string(buf[:n]))
			if (len(pattern) == 1) && (pattern[0] == (CMD_EXIT)) || (err != nil) {
				log.Printf(lang.Lookup(locale, "End condition, closing connection for:")+" %s"+nl.NewLine(), user.name)
				inputChannel <- ClientInput{
					user,
					&UserLeftEvent{user, lang.Lookup(locale, "Goodbye")},
				}
				return err
			} else if (len(pattern) == 2) && (pattern[0] == (CMD_CHANGENICK)) {
				inputChannel <- ClientInput{
					user,
					&UserChangedNickEvent{user, pattern[1]},
				}
			} else if (len(pattern) == 1) && (pattern[0] == (CMD_LISTUSERS)) {
				inputChannel <- ClientInput{
					user,
					&ListUsersEvent{user},
				}
			}
		} else {
			msg := strings.TrimSpace(string(string(buf[:n])))
			e := ClientInput{user, &MessageEvent{msg}}
			inputChannel <- e
		}

	}
}

func startServer(eventChannel chan ClientInput, config *tls.Config, port string, nl Newline) error {
	log.Printf(lang.Lookup(locale, "Starting server on:")+"%s"+nl.NewLine(), port)
	ln, err := tls.Listen("tcp", port, config)
	if err != nil {
		// handle error
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			log.Print(lang.Lookup(locale, "Error accepting connection")+nl.NewLine(), err)
			continue
		}
		go func() {
			if err := handleConnection(conn, eventChannel, nl); err != nil {
				log.Print(lang.Lookup(locale, "Error handling connection or unexpected client exit")+nl.NewLine(), err)
			}
		}()

	}
}
