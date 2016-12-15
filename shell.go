package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func (m *machine) loadShellWords() {
	m.predefinedWords["ls"] = GoWord(func(m *machine) error {
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
	})

	m.predefinedWords["pwd"] = GoWord(func(m *machine) error {
		if dir, err := os.Getwd(); err != nil {
			return err
		} else {
			m.pushValue(String(dir))
			return nil
		}
	})

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
