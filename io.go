package pisc

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var ModIOCore = Module{
	Author:    "Andrew Owen",
	Name:      "IOCore",
	License:   "MIT",
	DocString: "Some File I/O bits",
	Load:      loadIOCore,
}

/*
// Windows only right now!
func getch(m *Machine) error {
	char := C.getch()
	m.PushValue(Integer(char))
	return nil
}
*/

func loadIOCore(m *Machine) error {
	m.AddGoWord("import", "( file-path -- )", GoWord(importPISC))
	// m.AddGoWord("getkey", "( -- keyval )", GoWord(getch))
	m.AddGoWord("ESC",
		"( -- ESC-char ) Emits the terminal escape char, for use in terminal escape codes",
		GoWord(func(m *Machine) error {
			m.PushValue(String("\x1B"))
			return nil
		}))
	// TODO: Consider deleting this later, if it isn't used.
	m.AddGoWord("import-asset", "( file-path -- )", GoWord(importAssetPISC))

	m.AddGoWord("get-str-at-path", "( path -- contents )", GoWord(func(m *Machine) error {
		fileName := m.PopValue().(String)
		data, err := ioutil.ReadFile(string(fileName))
		if err != nil {
			return err
		}
		m.PushValue(String(string(data)))
		return nil
	}))

	m.AddGoWord("save-str-to-path", "( str path -- )", GoWord(func(m *Machine) error {
		fileName := m.PopValue().(String)
		data := m.PopValue().String()
		return ioutil.WriteFile(string(fileName), []byte(data), os.FileMode(0644))
	}))

	m.PredefinedWords["open-file-writer"] = GoWord(func(m *Machine) error {
		fileName := m.PopValue().(String)
		goFile, err := os.Create(string(fileName))
		if err != nil {
			return err
		}
		// err = goFile.Chmod(os.ModePerm | 0644)
		if err != nil {
			return err
		}
		fileWriter := bufio.NewWriter(goFile)
		var file = Dict(make(map[string]StackEntry))
		file["close"] = GoFunc(func(m *Machine) error {
			return goFile.Close()
		})
		file["write-line"] = GoFunc(func(m *Machine) error {
			str := m.PopValue().String()
			_, err := fileWriter.WriteString(str + "\n")
			fileWriter.Flush()
			return err
		})

		file["write-string"] = GoFunc(func(m *Machine) error {
			str := m.PopValue().String()
			_, err := fileWriter.WriteString(str)
			fileWriter.Flush()
			return err
		})
		m.PushValue(file)

		return m.ImportPISCAsset("stdlib/io.pisc")

	})

	m.PredefinedWords["open-file-reader"] = GoWord(func(m *Machine) error {
		var fileName = m.PopValue().(String)
		// var file = Dict(make(map[string]StackEntry))
		goFile, err := os.Open(string(fileName))
		if err != nil {
			return err
		}
		var reader = bufio.NewReader(goFile)
		file := MakeReader(reader)
		file["close"] = GoFunc(func(m *Machine) error {
			return goFile.Close()
		})
		m.PushValue(file)
		return nil
	})

	m.PredefinedWords["priv_puts"] = NilWord(func(m *Machine) {
		data := m.PopValue().(String)
		fmt.Print(string(data))
	})
	return nil
}

type PISCReader interface {
	io.RuneReader
	io.ByteReader
	io.Reader
	ReadString(delim byte) (string, error)
}

func MakeReader(reader PISCReader) Dict {
	file := make(Dict)
	EOF := false
	file["read-byte"] = GoFunc(func(m *Machine) error {
		b, err := reader.ReadByte()
		if err == io.EOF {
			EOF = true
			err = nil
		}
		if err != nil {
			return err
		}
		m.PushValue(Integer(int(b)))
		return nil
	})
	file["read-rune"] = GoFunc(func(m *Machine) error {
		ch, _, err := reader.ReadRune()
		if err == io.EOF {
			EOF = true
			err = nil
		}
		if err != nil {
			return err
		}
		m.PushValue(String(string(ch)))
		return nil
	})
	file["read-line"] = GoFunc(func(m *Machine) error {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			EOF = true
			err = nil
		}
		if err != nil {
			return err
		}
		// Deal with \r on windows
		m.PushValue(String(strings.TrimRight(str, "\r\n")))
		return nil
	})
	file["EOF"] = GoFunc(func(m *Machine) error {
		m.PushValue(Boolean(EOF))
		return nil
	})
	return file
}

func importPISC(m *Machine) error {
	fileName := m.PopValue().(String)
	data, err := ioutil.ReadFile(string(fileName))
	if err != nil {
		return err
	}
	// Reading in the data
	code, err := stringToQuotation(string(data), CodePosition{
		Source: "file:" + string(fileName),
	})
	if err != nil {
		return err
	}
	err = m.execute(code)
	if err != nil {
		return err
	}
	return nil
}

func (m *Machine) ImportPISCAsset(assetkey string) error {
	data, err := Asset(string(assetkey))
	if err != nil {
		return err
	}
	// Reading in the data
	code, err := stringToQuotation(string(data), CodePosition{Source: "file:" + string(assetkey)})
	if err != nil {
		return err
	}
	err = m.execute(code)
	if err != nil {
		return err
	}
	return nil
}

func importAssetPISC(m *Machine) error {
	fileName := m.PopValue().(String).String()
	return m.ImportPISCAsset(fileName)
}
