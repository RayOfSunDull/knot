package knot

import (
	"os"
	"os/exec"
	"path/filepath"
	"errors"
    "strings"
)

type PlatformDirs struct {
    ConfigDir string
    TempDir string
    BinDir string
}

type Platform interface {
    GetPlatformDirs() (PlatformDirs, error)
    GetPythonCommand() (string, error)
}

type Linux struct {}

func (linux Linux) GetPlatformDirs() (PlatformDirs, error) {
    homeDir, err := os.UserHomeDir();
    if err != nil {
        return PlatformDirs{}, err
    }

    configDir := filepath.Join(homeDir, ".config")
    tempDir := "/tmp"

    binDir := filepath.Join(homeDir, ".local/bin")
    preferredBinDir := filepath.Join(homeDir, "bin")
    // we prefer $HOME/bin to $HOME/.local/bin, so if
    // it's in $PATH we will change binDir to that
    path := os.Getenv("PATH")
    pathElements := strings.Split(path, ":")
    for _, pathElement := range pathElements {
        if pathElement != preferredBinDir {
            continue
        }
        binDir = preferredBinDir
        break
    }

    return PlatformDirs{
        ConfigDir: configDir,
        TempDir: tempDir,
        BinDir: binDir}, nil
}

func (linux Linux) GetPythonCommand() (string, error) {
	python3, err := exec.LookPath("python3")
	if err == nil {
        return string(python3), nil
    }

	python3NotFound := errors.New("unable to find Python 3")

	python, err := exec.LookPath("python")
	if err != nil { return "", python3NotFound }

	getPythonVersion := exec.Command(
		"python", "-c", "'import sys; print(sys.version_info[0])'")
	
	version, err := getPythonVersion.Output()
	if err != nil {
		return "", errors.New("unable to determine Python version")
	}

	if string(version) != "3" {
        return "", python3NotFound
    }

	return string(python), nil
}
