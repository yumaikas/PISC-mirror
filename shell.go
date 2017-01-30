package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func (m *machine) loadShellWords() {
	m.addGoWord("list-files", "( -- files ) ", GoWord(func(m *machine) error {
		files, err := ioutil.ReadDir(".")
		if err != nil {
			fmt.Println("Error: ", err)
		}
		arr := make(Array, len(files))
		for i, f := range files {
			arr[i] = fileInfoToDict(f)
		}
		m.pushValue(arr)
		return nil
	}))

	m.addGoWord("pwd", "( -- workingdir )", GoWord(func(m *machine) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		m.pushValue(String(dir))
		return nil
	}))

	m.addGoWord("env-get", " ( key -- envVal ) ", GoWord(func(m *machine) error {
		key := m.popValue().String()
		m.pushValue(String(os.Getenv(key)))
		return nil
	}))
	m.addGoWord("env-set", " ( key value -- ) ", GoWord(func(m *machine) error {
		val := m.popValue().String()
		key := m.popValue().String()
		return os.Setenv(key, val)
	}))

	//
	m.predefinedWords["cd"] = GoWord(func(m *machine) error {
		location := m.popValue().String()
		if err := os.Chdir(location); err != nil {
			return err
		}
		if dir, err := os.Getwd(); err != nil {
			return err
		} else {
			fmt.Println(dir)
		}
		return nil
	})
}

func fileInfoToDict(info os.FileInfo) Dict {
	dict := make(Dict)
	dict["name"] = String(info.Name())
	dict["size"] = Integer(info.Size())
	// dict["is-dir"] = Boolean(info.IsDir())
	dict["mode"] = String(info.Mode().String())
	dict["__type"] = String("inode")
	dict["__ordering"] = Array{
		String("name"),
		String("mode"),
		String("size"),
	}
	return dict
}
