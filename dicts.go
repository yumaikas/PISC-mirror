package pisc

// "strings"

var ModDictionaryCore = Module{
	Author:  "Andrew Owen",
	Name:    "DictionaryCore",
	License: "MIT", // TODO: Clarify here
	Load:    loadDictMod,
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
	// Peek, since we have no intention of popping here.
	dict := m.PopValue().(Dict)
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

	keyArr := make(Array, dict.Length())

	var i int = 0
	for k, _ := range dict {
		keyArr[i] = String(k)
		i++
	}
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
		err := m.execute(quot.inner)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadDictWords(m *Machine) error {

	m.AddGoWord("<dict>", "( -- dict ) Place an empty dictionary on the stack", _buildDict)
	m.AddGoWord("dict-has-key?", "( dict key -- has-key? )", _dictHasKey)
	m.AddGoWord("dict-set", "( dict value key -- dict )", _dictSet)

	m.AddGoWord("dict-get", "( dict key -- value|error? )", _dictGet)
	m.AddGoWord("dict-keys", "( dict -- { keys }) Puts all the keys for a dictionary in an array", _dictKeys)
	m.AddGoWord("dict-get-rand", "( dict -- key value )", _dictGetRand)
	m.AddGoWord("dict-each-key", "(dict quot -- .. )", _dictEachKey)

	return m.ImportPISCAsset("stdlib/dicts.pisc")
}
