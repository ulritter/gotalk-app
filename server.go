package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"
)

// read from connection, recognize request types and pass appropriate event types to the session handler (serverDialog())
func (a *application) handleConnection(conn net.Conn, inputChannel chan ClientInput) error {
	buf := make([]byte, BUFSIZE)

	session := &Session{conn}
	n, err := conn.Read(buf)
	if err != nil {
		a.logger.Print(a.lang.Lookup(a.config.locale, "Error reading from buffer")+a.config.newline, err)
		return err
	}
	msg := Message{}
	msg.Body = nil
	msg.UnmarshalMSG(buf[:n])

	if (msg.Action != ACTION_INIT) || (len(msg.Body) != 2) {
		return fmt.Errorf(a.lang.Lookup(a.config.locale, "Wrong connection initialization message."))
	} else {
		// expecting format {ACTION_INIT, [{<nickname>}, {<revision level>}]}
		if msg.Body[1] != REVISION {
			sendMessage(conn, ACTION_REVISION, []string{REVISION})
			return fmt.Errorf(a.lang.Lookup(a.config.locale,
				"Connection request from ")+conn.RemoteAddr().(*net.TCPAddr).IP.String()+a.lang.Lookup(a.config.locale,
				" rejected. ")+a.lang.Lookup(a.config.locale,
				"Wrong client revision level. Should be: ")+" %s"+a.lang.Lookup(a.config.locale, ", actual: ")+"%s", REVISION, msg.Body[1])
		}
	}
	user := &User{name: msg.Body[0], session: session, timejoined: time.Now().Format("2006.01.02 15:04:05")}
	inputChannel <- ClientInput{
		user,
		&UserJoinedEvent{},
	}

	for {
		n, err1 := conn.Read(buf)
		if err1 != nil {
			a.logger.Printf(a.lang.Lookup(a.config.locale, "End condition, closing connection for:")+" %s"+a.config.newline, user.name)
			inputChannel <- ClientInput{
				user,
				&UserLeftEvent{user, a.lang.Lookup(a.config.locale, "Goodbye")},
			}
			return err1
		}

		msg.Action = ""
		msg.Body = nil
		err2 := msg.UnmarshalMSG(buf[:n])

		if err2 != nil {
			a.logger.Printf(a.lang.Lookup(a.config.locale, "Warning: Corrupt JSON Message from: ")+" %s"+a.config.newline, user.name)
			a.logger.Println(err2)
		}

		if msg.Action == ACTION_EXIT {
			a.logger.Printf(a.lang.Lookup(a.config.locale, "End condition, closing connection for:")+" %s"+a.config.newline, user.name)
			//echo exit condition for organized client shutdown
			sendMessage(conn, ACTION_EXIT, []string{""})
			inputChannel <- ClientInput{
				user,
				&UserLeftEvent{user, a.lang.Lookup(a.config.locale, "Goodbye")},
			}
			return err1
		}

		switch msg.Action {
		case ACTION_CHANGENICK:
			if len(msg.Body) == 1 {
				inputChannel <- ClientInput{
					user,
					&UserChangedNickEvent{user, msg.Body[0]},
				}
			}
		case ACTION_LISTUSERS:
			if msg.Action == ACTION_LISTUSERS {
				inputChannel <- ClientInput{
					user,
					&ListUsersEvent{user},
				}
			}
		case ACTION_SENDMESSAGE:
			if len(msg.Body) == 1 {
				sendmsg := strings.TrimSpace(msg.Body[0])
				e := ClientInput{user, &MessageEvent{sendmsg}}
				inputChannel <- e
			}
		default:
		}

	}
}

// this function is called by main() in the case the app needs to operate as server
// wait for connections and start a handler for each connection
func (a *application) startServer(eventChannel chan ClientInput, config *tls.Config, port string) error {
	a.logger.Printf(a.lang.Lookup(a.config.locale, "Starting server on port ")+"%s"+a.config.newline, port)
	ln, err := tls.Listen("tcp", port, config)
	if err != nil {
		// handle error
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			a.logger.Print(a.lang.Lookup(a.config.locale, "Error accepting connection")+a.config.newline, err)
			continue
		}
		go func() {
			if err := a.handleConnection(conn, eventChannel); err != nil {
				a.logger.Print(a.lang.Lookup(a.config.locale, "Error handling connection or unexpected client exit")+a.config.newline, err)
			}
		}()

	}
}
