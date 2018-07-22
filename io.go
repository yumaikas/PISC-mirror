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
// m.AddGoWord*Change?*("getkey", "( -- keyval )", GoWord(getch))
*/

func _emitESC(m *Machine) error {
	m.PushValue(String("\x1B"))
	return nil
}

func _getStrAtPath(m *Machine) error {
	fileName := m.PopValue().(String)
	data, err := ioutil.ReadFile(string(fileName))
	if err != nil {
		return err
	}
	m.PushValue(String(string(data)))
	return nil
}

func _saveStrToPath(m *Machine) error {
	fileName := m.PopValue().(String)
	data := m.PopValue().String()
	return ioutil.WriteFile(string(fileName), []byte(data), os.FileMode(0644))
}

func _openFileWriter(m *Machine) error {
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
}

func _openFileReader(m *Machine) error {
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
}

func _priv_puts(m *Machine) error {
	data := m.PopValue().(String)
	fmt.Print(string(data))
	return nil
}

func loadIOCore(m *Machine) error {
	NL := "\n"
	m.AddGoWordWithStack(
		"import",
		"( file-path -- )",
		"Loads the PISC file at the given path",
		importPISC)

	m.AddGoWordWithStack("ESC",
		"( -- ESC-char )",
		"Emits the terminal escape char, for use in terminal escape codes",
		_emitESC)

	// TODO: Consider deleting this later, if it isn't used.
	m.AddGoWordWithStack(
		"import-asset",
		"( file-path -- ? )",
		"Load the packed script (mostly used for the standard library)",
		importAssetPISC)

	m.AddGoWordWithStack(
		"get-str-at-path",
		"( path -- contents )",
		"Load the contents at path into a string",
		_getStrAtPath)

	m.AddGoWordWithStack(
		"save-str-to-path",
		"( str path -- )",
		"Save the value in str to the file at path",
		_saveStrToPath)

	m.AppendToHelpTopic("readers",
		`Readers can currently be made from strings or files

They support the following calls (assuming a reader in $reader)

- <code>$reader .read-byte<code> reads a single byte, and puts it atop the stack
- <code>$reader .read-rune<code> reads a single UTF-8 rune
- <code>$reader .read-line<code> reads a single line into a string

`)
	m.AddGoWordWithStack(
		"open-file-writer",
		"( path -- file-writer )",
		"Opens a file-writer that that writes to the supplied path, if can be made"+NL+
			"A file-writer supports 3 calls: "+NL+
			`<code>"str" $writer .write-line<code>, which writes "str\n" to the file`+NL+
			`- <code>"str" $writer .write-string<code>, which writes "str" to the file`+NL+
			`- <code>"str" $writer .write-string<code>, which writes "str" to the file`+NL,
		_openFileWriter)

	m.AddGoWordWithStack(
		"open-file-reader",
		"( path -- file-reader )",
		"Opens a file-readr that reads from the file at the supplied path, if a file exists at the given path"+NL+
			"Support the standard @reader calls, as well as <code>.close<code>"+NL+
			"See also @readers",
		_openFileReader)

	m.AddGoWordWithStack("priv_puts", "( str -- )", "Prints str to STDOUT", _priv_puts)
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
	err = m.execute(code) // Correctly built
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
	err = m.execute(code) // Correctly built
	if err != nil {
		return err
	}
	return nil
}

func importAssetPISC(m *Machine) error {
	fileName := m.PopValue().(String).String()
	return m.ImportPISCAsset(fileName)
}
