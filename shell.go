package pisc

import (
	"fmt"
	"io/ioutil"
	"os"
)

var ModShellUtils = Module{
	Author:    "Andrew Owen",
	Name:      "ShellUtils",
	License:   "MIT",
	DocString: `A set of fucntions that are used to provide some shell-like functionality`,
	Load:      loadShellWords,
}

func loadShellWords(m *Machine) error {

	m.AddGoWord("list-files-at", "( path -- files )", GoWord(func(m *Machine) error {
		dirPath := m.PopValue().String()
		files, err := ioutil.ReadDir(dirPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: ", err)
		}
		arr := make(Array, len(files))
		for i, f := range files {
			arr[i] = fileInfoToDict(f)
		}
		m.PushValue(arr)
		return nil
	}))

	m.AddGoWord("list-files", "( -- files ) ", GoWord(func(m *Machine) error {
		files, err := ioutil.ReadDir(".")
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: ", err)
		}
		arr := make(Array, len(files))
		for i, f := range files {
			arr[i] = fileInfoToDict(f)
		}
		m.PushValue(arr)
		return nil
	}))

	m.AddGoWord("stat", "( filepath -- info )", GoWord(func(m *Machine) error {
		path := m.PopValue().String()
		info, err := os.Stat(path)
		if err != nil {
			return err
		}

		// Hrm... Is there a way to avoid this double allocation (struct into dict)
		m.PushValue(fileInfoToDict(info))
		return nil

	}))

	m.AddGoWord("pwd", "( -- workingdir )", GoWord(func(m *Machine) error {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		m.PushValue(String(dir))
		return nil
	}))

	m.AddGoWord("env-get", " ( key -- envVal ) ", GoWord(func(m *Machine) error {
		key := m.PopValue().String()
		m.PushValue(String(os.Getenv(key)))
		return nil
	}))
	m.AddGoWord("env-set", " ( key value -- ) ", GoWord(func(m *Machine) error {
		val := m.PopValue().String()
		key := m.PopValue().String()
		return os.Setenv(key, val)
	}))

	//
	m.AddGoWord("cd", "( new-dir -- ) ", GoWord(func(m *Machine) error {
		location := m.PopValue().String()
		if err := os.Chdir(location); err != nil {
			return err
		}
		if dir, err := os.Getwd(); err != nil {
			return err
		} else {
			fmt.Println(dir)
		}
		return nil
	}))

	return m.importPISCAsset("stdlib/shell.pisc")
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
