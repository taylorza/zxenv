package engine

import "path"

var tmpDir = ".tmp"
var emulatorDir = "emulator"
var sdDir = "sdcard"

type Environment struct {
	BasePath string
}

func (e *Environment) TempPath() string {
	return path.Join(e.BasePath, tmpDir)
}

func (e *Environment) EmulatorPath() string {
	return path.Join(e.BasePath, emulatorDir)
}

func (e *Environment) SDPath() string {
	return path.Join(e.BasePath, sdDir)
}
