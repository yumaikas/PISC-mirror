package pisc

// "strings"

var ModDictionaryCore = Module{
	Author:    "Andrew Owen",
	Name:      "DictionaryCore",
	License:   "MIT", // TODO: Clarify here
	Load:      loadDictMod,
	DocString: "Words for manipulating dictionaries",
	// Possible: indicate PISC files to be loaded?
}

var ErrMissingKey = Error{
	message: "Dictionary was missing key",
}

func loadDictMod(m *Machine) error {
	return loadDictWords(m)
}

func _buildDict(m *Machine) error {
	dict := make(map[string]StackEntry)
	m.PushValue(Dict(dict))
	return nil
}

func _dictHasKey(m *Machine) error {
	key := m.PopValue().(String).String()
	dict := m.PopValue().(Dict)
	_, ok := dict[key]
	m.PushValue(Boolean(ok))
	return nil
}

func _dictSet(m *Machine) error {
	key := m.PopValue().(String).String()
	value := m.PopValue()
	dict := m.PopValue().(Dict)
	dict[string(key)] = value
	return nil
}

func _dictPush(m *Machine) error {
	key := m.PopValue().(String).String()
	value := m.PopValue()
	dict := m.Values[len(m.Values)-1].(Dict)
	dict[string(key)] = value
	return nil
}

func getMissingKeyErr(key string) Error {
	return Error{
		message: "Dictionary was missing key:" + key,
	}
}

func _dictGet(m *Machine) error {
	key := m.PopValue().(String).String()
	dict := m.PopValue().(Dict)
	if val, found := dict[key]; found {
		m.PushValue(val)
		return nil
	} else {
		return getMissingKeyErr(key)
	}
}

func _dictKeys(m *Machine) error {
	dict := m.PopValue().(Dict)

	keyArr := make([]StackEntry, dict.Length())

	var i int = 0
	for k, _ := range dict {
		keyArr[i] = String(k)
		i++
	}
	m.PushValue(&Vector{Elements: keyArr})
	return nil
}

func _dictGetRand(m *Machine) error {
	dict := m.PopValue().(Dict)
	// Rely on random key ordering
	for k, v := range dict {
		m.PushValue(String(k))
		m.PushValue(v)
		break
	}
	return nil
}

// dict quot -- ...
func _dictEachKey(m *Machine) error {
	quot := m.PopValue().(*Quotation)
	dict := m.PopValue().(Dict)

	for k, _ := range dict {
		m.PushValue(String(k))
		err := m.CallQuote(quot)
		if IsLoopError(err) && LoopShouldEnd(err) {
			return nil
		}
		if IsLoopError(err) && !LoopShouldEnd(err) {
			continue
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func loadDictWords(m *Machine) error {

	m.AddGoWordWithStack("<dict>",
		"( -- dict )",
		"Place an empty dictionary on the stack",
		_buildDict)
	m.AddGoWordWithStack("dict-has-key?",
		"( dict key -- has-key? )",
		"Check to see if a dictionary has a given key",
		_dictHasKey)

	m.AddGoWordWithStack(
		"dict-set",
		"( dict value key -- )",
		"Set the entry in dict at key to value",
		_dictSet)
	m.AddGoWordWithStack(
		"dict-push",
		"( dict value key -- dict )",
		"Set the entry in dict at key to value, leaving the dict on the stack",
		_dictPush)

	m.AddGoWordWithStack("dict-get",
		"( dict key -- value|error? )",
		"Get the key, or error if it doesn't exist",
		_dictGet)
	m.AddGoWordWithStack("dict-keys",
		"( dict -- { keys } )",
		"Loads all the keys for a dictionary into an array",
		_dictKeys)
	m.AddGoWordWithStack("dict-get-rand",
		"( dict -- key value )",
		"Get a random key/value pair from this dictionary",
		_dictGetRand)
	m.AddGoWordWithStack("dict-each-key",
		"(dict quot -- .. )",
		"Run a function over each key in the dictionary",
		_dictEachKey)

	return m.ImportPISCAsset("stdlib/dicts.pisc")
}
