package pisc

import "gopkg.in/sorcix/irc.v2"
import "fmt"

var ModIRCKit = Module{
	Author:    "Andrew Owen",
	Name:      "IRCKit",
	License:   "MIT",
	DocString: "A wrapper around IRC, built to make it easy to write IRC bots",
	Load:      loadIRCKit,
}

type ircConn irc.Conn
type ircMessage irc.Message

var InvalidIRCMessage = fmt.Errorf("IRC message formmated incorrectly")

func (conn *ircConn) close(m *Machine) error {
	return ((*irc.Conn)(conn)).Close()
}

func (conn *ircConn) write(m *Machine) error {
	str := m.PopValue().String()
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

func (_msg ircMessage) getMessageCommand(m *Machine) error {
	msg := irc.Message(_msg)
	m.PushValue(String(msg.Command))
	return nil
}

func (_msg ircMessage) getMessagePrefixIsUserLike(m *Machine) error {
	msg := irc.Message(_msg)
	m.PushValue(Boolean(msg.IsHostmask()))
	return nil
}

func (_msg ircMessage) getMessagePrefixIsServerLike(m *Machine) error {
	msg := irc.Message(_msg)
	m.PushValue(Boolean(msg.IsServer()))
	return nil
}

func (_msg ircMessage) getName(m *Machine) error {
	msg := irc.Message(_msg)
	m.PushValue(String(msg.Name))
	return nil
}

func (_msg ircMessage) getParams(m *Machine) error {
	msg := irc.Message(_msg)
	params := make(Array, len(msg.Params))
	for i := 0; i < len(msg.Params); i++ {
		params[i] = String(msg.Params[i])
	}
	m.PushValue(Array(params))
	return nil
}

func (conn *ircConn) readMessage(m *Machine) error {
	_msg, err := (*irc.Conn)(conn).Decode()
	if err != nil {
		return err
	}

	msg := (*ircMessage)(_msg)
	msgDict := Dict{
		"command":       GoFunc(msg.getMessageCommand),
		"name":          GoFunc(msg.getName),
		"is-userlike":   GoFunc(msg.getMessagePrefixIsUserLike),
		"is-serverlike": GoFunc(msg.getMessagePrefixIsServerLike),
		"params":        GoFunc(msg.getParams),
	}

	m.PushValue(msgDict)
	return nil
}

func (conn *ircConn) readMessageString(m *Machine) error {
	msg, err := (*irc.Conn)(conn).Decode()
	if err != nil {
		return err
	}
	str := String(msg.Bytes())
	m.PushValue(str)
	return nil
}

// TODO: Load IRCKit here, using IRCX library
func loadIRCKit(m *Machine) error {

	m.AddGoWord("irc-dial", "( addr-str -- conn )", func(m *Machine) error {
		addr := m.PopValue().String()
		conn, err := irc.Dial(addr)
		if err != nil {
			return err
		}
		stackConn := (*ircConn)(conn)

		connDict := Dict{
			"close":               GoFunc(stackConn.close),
			"send-message":        GoFunc(stackConn.write),
			"recieve-message-str": GoFunc(stackConn.readMessageString),
			"recieve-message":     GoFunc(stackConn.readMessage),
		}

		m.PushValue(connDict)
		return nil
	})
	return nil
}
