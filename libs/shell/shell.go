package shell

import (
	"fmt"
	"io/ioutil"
	"os"
	"pisc"
	"strconv"
)

var ModShellUtils = pisc.Module{
	Author:    "Andrew Owen",
	Name:      "ShellUtils",
	License:   "MIT",
	DocString: `A set of fucntions that are used to provide some shell-like functionality`,
	Load:      loadShellWords,
}

func _listFilesAt(m *pisc.Machine) error {
	dirPath := m.PopValue().String()
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
	}
	arr := make([]pisc.StackEntry, len(files))
	for i, f := range files {
		arr[i] = fileInfoToDict(f)
	}
	m.PushValue(&pisc.Vector{Elements: arr})
	return nil
}

func _listFiles(m *pisc.Machine) error {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
	}
	arr := make([]pisc.StackEntry, len(files))
	for i, f := range files {
		arr[i] = fileInfoToDict(f)
	}
	m.PushValue(&pisc.Vector{Elements: arr})
	return nil
}

func _stat(m *pisc.Machine) error {
	path := m.PopValue().String()
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	m.PushValue(fileInfoToDict(info))
	return nil
}

func _pwd(m *pisc.Machine) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	m.PushValue(pisc.String(dir))
	return nil
}

func _envGet(m *pisc.Machine) error {
	key := m.PopValue().String()
	m.PushValue(pisc.String(os.Getenv(key)))
	return nil
}

func _envSet(m *pisc.Machine) error {
	val := m.PopValue().String()
	key := m.PopValue().String()
	return os.Setenv(key, val)
}

func _cd(m *pisc.Machine) error {
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
}

func loadShellWords(m *pisc.Machine) error {
	m.AddGoWord("list-files-at", "( path -- files )", _listFilesAt)
	m.AddGoWord("list-files", "( -- files ) ", _listFiles)
	m.AddGoWord("stat", "( filepath -- info )", _stat)
	m.AddGoWord("pwd", "( -- workingdir )", _pwd)
	m.AddGoWord("env-get", " ( key -- envVal ) ", _envGet)
	m.AddGoWord("env-set", " ( key value -- ) ", _envSet)
	m.AddGoWord("cd", "( new-dir -- ) ", _cd)

	return m.ImportPISCAsset("stdlib/shell.pisc")
}

func fileInfoToDict(info os.FileInfo) pisc.Dict {
	dict := make(pisc.Dict)
	dict["name"] = pisc.String(info.Name())
	dict["size"] = pisc.Integer(info.Size())
	// dict["is-dir"] = Boolean(info.IsDir())
	dict["mode"] = pisc.String(info.Mode().String())
	dict["timestamp"] = pisc.String(strconv.FormatInt(info.ModTime().Unix(), 10))
	dict["__type"] = pisc.String("inode")
	dict["__ordering"] = &pisc.Vector{
		Elements: []pisc.StackEntry{
			pisc.String("name"),
			pisc.String("mode"),
			pisc.String("size"),
			pisc.String("timestamp"),
		},
	}
	return dict
}
