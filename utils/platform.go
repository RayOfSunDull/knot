package knot

import (
	"os"
	"path/filepath"
)

type Platform interface {
	GetConfigDir() (string, error)
	GetTempDir() (string, error)
	GetBinDir() (string, error)
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