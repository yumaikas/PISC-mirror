package main

import "pisc"

var ModBuildCore = pisc.Module{
	Author:    "Andrew Owen",
	Name:      "BuildKit",
	License:   "MIT",
	DocString: `commands for running pisc Buildscripts`,
	Load:      loadBuildWords,
}

func loadBuildWords(m *pisc.Machine) error {
	return nil
}
