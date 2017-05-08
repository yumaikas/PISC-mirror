package main

import "gopkg.in/sorcix/irc.v2"
import "fmt"

var ModIRCKit = PISCModule{
	Author:    "Andrew Owen",
	Name:      "IRCKit",
	License:   "MIT",
	DocString: "A wrapper around IRC, built to make it easy to write IRC bots",
	Load:      loadIRCKit,
}

type ircConn irc.Conn
type ircMessage irc.Message

var InvalidIRCMessage = fmt.Errorf("IRC message formmated incorrectly")

func (conn *ircConn) close(m *machine) error {
	return ((*irc.Conn)(conn)).Close()
}

func (conn *ircConn) write(m *machine) error {
	str := m.popValue().String()
	msg := irc.ParseMessage(str)

	fmt.Println(msg)
	if msg == nil {
		return InvalidIRCMessage
	}

	err := (*irc.Conn)(conn).Encode(msg)
	if err != nil {
		return err
	}
	return nil
}

func (_msg ircMessage) getMessageCommand(m *machine) error {
	msg := irc.Message(_msg)
	m.pushValue(String(msg.Command))
	return nil
}

func (_msg ircMessage) getMessagePrefixIsUserLike(m *machine) error {
	msg := irc.Message(_msg)
	m.pushValue(Boolean(msg.IsHostmask()))
	return nil
}

func (_msg ircMessage) getMessagePrefixIsUserLike(m *machine) error {
	msg := irc.Message(_msg)
	m.pushValue(Boolean(msg.IsServer()))
	return nil
}

func (conn *ircConn) readMessage(m *machine) error {
	_msg, err := (*irc.Conn)(conn).Decode()
	if err != nil {
		return err
	}

	msg := ircMessage(_msg)

	msgDict := Dict{
		"command": GoFunc(msg.getMessageCommand),
		"is-userlike": GoFunc(msg.getMessagePrefixIsUserLike),
		"is-serverlike": GoFunc(msg.g)
	}

	m.pushValue(msgDict)
	return nil
}

func (conn *ircConn) readMessageString(m *machine) error {
	msg, err := (*irc.Conn)(conn).Decode()
	if err != nil {
		return err
	}
	str := String(msg.Bytes())
	m.pushValue(str)
	return nil
}

// TODO: Load IRCKit here, using IRCX library
func loadIRCKit(m *machine) error {

	m.addGoWord("irc-dial", "( addr-str -- conn )", func(m *machine) error {
		addr := m.popValue().String()
		conn, err := irc.Dial(addr)
		if err != nil {
			return err
		}
		stackConn := (*ircConn)(conn)

		connDict := Dict{
			"close":               GoFunc(stackConn.close),
			"send-message":        GoFunc(stackConn.write),
			"recieve-message-str": GoFunc(stackConn.readMessageString),
		}

		m.pushValue(connDict)
		return nil
	})
	return nil
}
