package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

var tmpDir = ".tmp"
var emulatorDir = "emulator"
var sdDir = "sdcard"

type Environment struct {
	BasePath string `json:"base_path"`
	SDSize   string `json:"sd_size"`
	Emulator string `json:"emulator"`
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

func (e *Environment) SDCardName() string {
	return fmt.Sprintf("sdcard-%s", e.SDSize)
}

func (e *Environment) Save() error {
	bytes, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("failed to marshal environment (%w)", err)
	}

	err = os.WriteFile(path.Join(e.BasePath, "environment.json"), bytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to save environment (%w)", err)
	}
	return nil
}

func LoadEnvironment(basePath string) (*Environment, error) {
	bytes, err := os.Open(path.Join(basePath, "environment.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to open environment file (%w)", err)
	}
	defer bytes.Close()

	env := &Environment{
		BasePath: basePath,
	}
	err = json.NewDecoder(bytes).Decode(env)
	if err != nil {
		return nil, fmt.Errorf("failed to parse environment file (%w)", err)
	}
	return env, nil
}
