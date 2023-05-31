package knot

import (
	"os"
	"os/exec"
	"path/filepath"
	"errors"
)

type Platform interface {
	GetConfigDir() (string, error)
	GetTempDir() (string, error)
	GetBinDir() (string, error)
	GetPythonCommand() (string, error)
}

type Unix struct {}

func (unix Unix) GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil { return "", err }

	return filepath.Join(home, ".config"), nil
}

func (unix Unix) GetTempDir() (string, error) {
	return "/tmp", nil
}

func (unix Unix) GetBinDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil { return "", err }

	return filepath.Join(home, "bin"), nil
}

func (unix Unix) GetPythonCommand() (string, error) {
	python3, err := exec.LookPath("python3")
	if err == nil { return string(python3), nil }

	python3NotFound := errors.New("unable to find Python 3")

	python, err := exec.LookPath("python")
	if err != nil { return "", python3NotFound }

	getPythonVersion := exec.Command(
		"python", "-c", "'import sys; print(sys.version_info[0])'")
	
	version, err := getPythonVersion.Output()
	if err != nil {
		return "", errors.New("unable to determine Python version")
	}

	if string(version) != "3" { return "", python3NotFound }

	return string(python), nil
}