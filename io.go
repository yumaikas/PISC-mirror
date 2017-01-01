package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type PISCReader interface {
	io.RuneReader
	io.ByteReader
	io.Reader
	ReadString(delim byte) (string, error)
}

func makeReader(reader PISCReader) Dict {
	file := make(Dict)
	EOF := false
	file["read-byte"] = GoFunc(func(m *machine) error {
		b, err := reader.ReadByte()
		if err == io.EOF {
			EOF = true
			err = nil
		}
		if err != nil {
			return err
		}
		m.pushValue(Integer(int(b)))
		return nil
	})
	file["read-rune"] = GoFunc(func(m *machine) error {
		ch, _, err := reader.ReadRune()
		if err == io.EOF {
			EOF = true
			err = nil
		}
		if err != nil {
			return err
		}
		m.pushValue(String(string(ch)))
		return nil
	})
	file["read-line"] = GoFunc(func(m *machine) error {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			EOF = true
			err = nil
		}
		if err != nil {
			return err
		}
		// Deal with \r on windows
		m.pushValue(String(strings.TrimRight(str, "\r\n")))
		return nil
	})
	file["EOF"] = GoFunc(func(m *machine) error {
		m.pushValue(Boolean(EOF))
		return nil
	})
	return file
}

func (m *machine) loadIOWords() error {
	m.predefinedWords["import"] = GoWord(func(m *machine) error {
		fileName := m.popValue().(String)
		data, err := ioutil.ReadFile(string(fileName))
		if err != nil {
			return err
		}
		// Reading in the data
		code := &codeList{
			idx:  0,
			code: string(data),
			codePosition: codePosition{
				source: "file:" + string(fileName),
			},
		}
		err = m.execute(code)
		if err != nil {
			return err
		}
		return nil
	})

	m.predefinedWords["filepath>string"] = GoWord(func(m *machine) error {
		fileName := m.popValue().(String)
		data, err := ioutil.ReadFile(string(fileName))
		if err != nil {
			return err
		}
		m.pushValue(String(string(data)))
		return nil
	})

	m.predefinedWords["open-file-writer"] = GoWord(func(m *machine) error {
		fileName := m.popValue().(String)
		goFile, err := os.Create(string(fileName))
		if err != nil {
			return err
		}
		err = goFile.Chmod(os.ModePerm | 0644)
		if err != nil {
			return err
		}
		fileWriter := bufio.NewWriter(goFile)
		var file = Dict(make(map[string]stackEntry))
		file["close"] = GoFunc(func(m *machine) error {
			return goFile.Close()
		})
		file["write-line"] = GoFunc(func(m *machine) error {
			str := m.popValue().String()
			_, err := fileWriter.WriteString(str + "\n")
			fileWriter.Flush()
			return err
		})

		file["write-string"] = GoFunc(func(m *machine) error {
			str := m.popValue().String()
			_, err := fileWriter.WriteString(str)
			fileWriter.Flush()
			return err
		})
		m.pushValue(file)

		return nil

	})

	m.predefinedWords["open-file-reader"] = GoWord(func(m *machine) error {
		var fileName = m.popValue().(String)
		// var file = Dict(make(map[string]stackEntry))
		goFile, err := os.Open(string(fileName))
		if err != nil {
			return err
		}
		var reader = bufio.NewReader(goFile)
		file := makeReader(reader)
		file["close"] = GoFunc(func(m *machine) error {
			return goFile.Close()
		})
		m.pushValue(file)
		/*
			toExec := &codeList{
				idx:    0,
				code:   quot,
				spaces: make([]string, 0),
			}

			executeWordsOnMachine(m, toExec)
		*/
		return nil
	})

	m.predefinedWords["priv_puts"] = NilWord(func(m *machine) {
		data := m.popValue().(String)
		fmt.Println(string(data))
	})
	return nil
}
