package pisc

import (
	"time"

	"gopkg.in/sorcix/irc.v2"
)
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
	ircVM.AddGoWord("IRC_GLOBALS", " ( -- dict ) ", func(inner_m *Machine) error {
		inner_m.PushValue(vmGlobals)
		return nil
	})

	err = ircVM.ExecuteString("seed-rand-time", CodePosition{Source: "irckit.go"})
	if err != nil {
		return nil, err
	}
	return ircVM, nil
}

func (ircVM *Machine) ircRestart(m *Machine) error {

	ircVM.ExecuteString("IRC_GLOBALS", CodePosition{Source: "irckit.go"})
	vmGlobals := ircVM.PopValue().(Dict)

	vm, err := buildIRCEvalVM()
	if err != nil {
		time.Sleep(time.Minute * 1)
		return err
	}

	vm.AddGoWord("IRC_GLOBALS", " ( -- dict ) ", func(inner_m *Machine) error {
		inner_m.PushValue(vmGlobals)
		return nil
	})

	ircVM = vm
	return nil
}

func saveEval(m *Machine, code string) (err error) {
	defer func() {
		pErr := recover()
		if pErr != nil {
			err = pErr
		}
	}()
	err = ircVM.ExecuteString(code, CodePosition{Source: "IRC Connection"})
}

func (ircVM *Machine) evalOnVM(m *Machine) error {
	code := string(m.PopValue().(String))
	err := saveEval(ircVM, code)
	// If there is an error, put it on the stack
	if err != nil {
		m.PushValue(&Vector{
			Elements: []StackEntry{String(err.Error())},
		})
	} else {
		m.PushValue(&Vector{Elements: ircVM.Values})
	}
	// Clean up the string
	err = m.ExecuteString(`
		"|" str-join 
		"\n" " " str-replace :output
		$output len 100 > [ ${ $output 0 100 str-substr " (Truncated...)" } ] [ $output ] if
`, CodePosition{Source: "irckit.go"})

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

// TODO: Load IRCKit here, using IRCX library
func loadIRCKit(m *Machine) error {
	m.AddGoWord("irc-dial", "( addr-str -- conn )", ircDial)
	m.AddGoWord("<irc-vm>", "( -- vm )", stackIRCEvalVM)
	return nil
}
