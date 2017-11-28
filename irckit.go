package pisc

import (
	"fmt"
	"time"

	"gopkg.in/sorcix/irc.v2"
)

var ModIRCKit = Module{
	Author:    "Andrew Owen",
	Name:      "IRCKit",
	License:   "MIT",
	DocString: "A wrapper around IRC, built to make it easy to write IRC bots",
	Load:      loadIRCKit,
}

var NL = "\n"

type ircConn irc.Conn
type ircMessage irc.Message

var InvalidIRCMessage = fmt.Errorf("IRC message formmated incorrectly")

func (conn *ircConn) close(m *Machine) error {
	return ((*irc.Conn)(conn)).Close()
}

func (conn *ircConn) write(m *Machine) error {
	str := m.PopValue().String()
	msg := irc.ParseMessage(str)

	if msg == nil {
		return InvalidIRCMessage
	}

	err := (*irc.Conn)(conn).Encode(msg)
	go func() { fmt.Println(msg) }()
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
	params := make([]StackEntry, len(msg.Params))
	for i := 0; i < len(msg.Params); i++ {
		params[i] = String(msg.Params[i])
	}
	m.PushValue(&Vector{Elements: params})
	return nil
}

func (_msg ircMessage) msgRaw(m *Machine) error {
	msg := irc.Message(_msg)
	m.PushValue(String(string(msg.Bytes())))
	return nil
}

func (conn *ircConn) readMessage(m *Machine) error {
	_msg, err := (*irc.Conn)(conn).Decode()
	if err != nil {
		return err
	}

	buf := _msg.Bytes()

	go func() { fmt.Println("READ:", string(buf)) }()
	msg := (*ircMessage)(_msg)
	msgDict := Dict{
		"command":       GoFunc(msg.getMessageCommand),
		"name":          GoFunc(msg.getName),
		"is-userlike":   GoFunc(msg.getMessagePrefixIsUserLike),
		"is-serverlike": GoFunc(msg.getMessagePrefixIsServerLike),
		"params":        GoFunc(msg.getParams),
		"raw-str":       GoFunc(msg.msgRaw),
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
	go func() { fmt.Println("READ:", string(str)) }()
	m.PushValue(str)
	return nil
}

func ircDial(m *Machine) error {
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
}

func attachIRC_GLOBALS(m *Machine, globals Dict) {
	m.AddGoWordWithStack(
		"IRC_GLOBALS",
		" ( -- dict ) ",
		"Puts the IRC_GLOBALS dict on the stack, allowing you to save and store information out of it.",
		func(inner_m *Machine) error {
			inner_m.PushValue(globals)
			return nil
		})
}

func buildIRCEvalVM() (*Machine, error) {
	ircVM := &Machine{
		Values:               make([]StackEntry, 0),
		DefinedWords:         make(map[string]*CodeQuotation),
		DefinedStackComments: make(map[string]string),
		PredefinedWords:      make(map[string]GoWord),
		PrefixWords:          make(map[string]*CodeQuotation),
		HelpDocs:             make(map[string]string),
		DispatchBudget:       2000,
		IsBudgeted:           true,
	}
	err := ircVM.LoadModules(
		ModLoopCore,
		ModLocalsCore,
		ModDictionaryCore,
		ModStringsCore,
		ModMathCore,
		ModBoolCore,
		ModVectorCore,
		ModSymbolCore,
		ModRandomCore,
		ModPISCCore)
	if err != nil {
		return nil, err
	}

	vmGlobals := make(Dict)
	attachIRC_GLOBALS(ircVM, vmGlobals)
	err = ircVM.ExecuteString("seed-rand-time", CodePosition{Source: "irckit.go"})
	if err != nil {
		return nil, err
	}
	return ircVM, nil
}

func (ircVM *Machine) ircRestart(m *Machine) error {
	ircVM.ExecuteString("IRC_GLOBALS", CodePosition{Source: "irckit.go"})
	vmGlobals, ok := ircVM.PopValue().(Dict)
	if !ok {
		//
		fmt.Println("Somehow IRC_GLOBALS got corrupted, so it's getting reset")
		vmGlobals = make(Dict)
	}

	vm, err := buildIRCEvalVM()
	if err != nil {
		time.Sleep(time.Minute * 1)
		return err
	}
	attachIRC_GLOBALS(vm, vmGlobals)
	ircVM = vm
	return nil
}

func safeEval(m *Machine, code string) (err error) {
	defer func() {
		pErr := recover()
		if pErr != nil {
			err = fmt.Errorf("Error while running irc eval %v", pErr)
		}
	}()
	err = m.ExecuteString(code, CodePosition{Source: "IRC Connection"})
	return
}

func (ircVM *Machine) evalOnVM(m *Machine) error {
	code := string(m.PopValue().(String))
	err := safeEval(ircVM, code)
	// If there is an error, put it on the stack
	if err != nil {
		m.PushValue(&Vector{
			Elements: []StackEntry{String(err.Error())},
		})
	} else {
		m.PushValue(&Vector{Elements: ircVM.Values})
	}

	// Reset the intermediate VM state.
	ircVM.Values = make([]StackEntry, 0)
	ircVM.NumDispatches = 0
	return nil
}

func stackIRCEvalVM(m *Machine) error {
	vm, err := buildIRCEvalVM()
	if err != nil {
		return err
	}
	vmDict := Dict{
		"restart": GoFunc(vm.ircRestart),
		"eval":    GoFunc(vm.evalOnVM),
	}
	m.PushValue(vmDict)
	return nil
}

func loadIRCKit(m *Machine) error {
	m.AppendToHelpTopic("irc-message",
		"An irc-message is a dict which supports the following calls (assuming a irc-message in $msg)"+NL+
			"`$msg .command`. ( -- command-str ) Retrives the command for the message, for example, PRIVMSG or JOIN"+NL+
			"`$msg .name`. Gets the nick or server name that sent the message"+NL+
			"`$msg .is-userlike`. ( -- userlike? ) Returns a boolean indicating if the name looks like a username"+NL+
			"`$msg .is-serverlike`. ( -- serverlike? ) Returns a boolean indicating if the name looks like a server name"+NL+
			"`$msg .params`. ( -- params-vec ) Returns a boolean indicating if the name looks like a server name"+NL+
			"`$msg .raw-str. ( -- str ) Gets the raw string of bytes that formed the underlying IRC message`"+NL+
			"")
	m.AddGoWordWithStack(
		"irc-dial",
		"( addr-str -- conn )",
		"Returns a an irc connnection, which supports the folloing calls (assuming a var $conn):`"+NL+
			"`$conn .close` ( -- ) which closes the underlying connnection"+NL+NL+

			"`\"/PRIVMSG #channel message\" $conn .send-message` ( message-str -- ) which takes an IRC command in string form,"+
			"parses it, and attemps to send it on the connection behind `$conn`"+NL+
			"`$conn .recieve-message-str` ( -- message-str ) which waits for a message on the connection, and then pulls the string form of that message out and puts it on the stack"+NL+
			"`$conn .recieve-message` ( -- message-dict ) which waits for a message on `$conn`, and then puts an @irc-message object on the stack"+NL+
			"",
		ircDial)
	m.AddGoWordWithStack(
		"<irc-vm>",
		"( -- vm )",
		"Builds a PISC vm that has been harded for IRC. AN IRC VM supports the following calls:"+NL+
			"`$vm .eval` ( code -- result-vec ) Takes code, evaluates on $vm, and pushes a vector that contains the results "+
			"of running the code. Example:  "+NL+
			"`\"1 2 +\" $vm .eval` => `{ 3 }`"+NL+
			"`$vm .restart` "+NL+
			""+NL+
			"",
		stackIRCEvalVM)
	return nil
}

// IDEAS:
/*

:ON JOIN ( %{ .nick .channel } -- )
    ->>nick :nick
    ->channel :channel
;

*/
