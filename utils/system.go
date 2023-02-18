package knot

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"errors"
	"strings"
	"io"
)

func GetEnvironmentVariables() map[string]string {
	environment := os.Environ()

	result = make(map[string]string, len(environment))

	for idx, keyValueString := range environment {
		keyValuePair := strings.Split(keyValueString, "=")

		key, value = keyValuePair[0], keyValuePair[1]

		result[key] = value
	}

	return result
}

func GetKnotWD() (string, error) {
	envVar := GetEnvironmentVariables()

	knotwd, ok := envVar["KNOTWD"]
	if ok { return knotwd, nil }

	pwd, err = os.Getwd()
	if err != nil { return "", err }

	return pwd, nil
}


type SystemInfo struct {
	Wd string
	ConfigDir string
	ProjectsFile string
	TemplateDir string
}

func GetSystemInfo() (SystemInfo, error) {
	wd, err := os.Getwd()
	if err != nil { return SystemInfo{}, err }

	homeDir, err := os.UserHomeDir()
	if err != nil { return SystemInfo{}, err }

	configDir := filepath.Join(homeDir, ".config", "knot")
	projectsFile := filepath.Join(configDir, "projects.json")
	templateDir := filepath.Join(configDir, "templates")

	return SystemInfo{
		Wd: wd,
		ConfigDir: configDir,
		ProjectsFile: projectsFile,
		TemplateDir: templateDir}, nil
}


func OpenFile(file string, open bool) error {
	if !open { return nil }

	extension := filepath.Ext(file)

	switch extension {
	case ".kra":
		cmd := exec.Command("nohup", "krita", file)
		return cmd.Start()
	case ".pdf":
		cmd := exec.Command("nohup", "evince", file)
		return cmd.Start()
	case "": // directory
		cmd := exec.Command("nautilus", file)
		return cmd.Start()
	default:
		return errors.New(fmt.Sprintf(
			"extension %s is unsupported", extension))
	}
}

func CreateFile(path string) error {
	file, ok := os.Create(path)
	file.Close()
	return ok
}

func CopyFile(src, dst string) (int64, error) {
	srcStat, err := os.Stat(src)
	if err != nil { return 0, err }

	if !srcStat.Mode().IsRegular() {
		return 0, errors.New(fmt.Sprintf(
			"%s is not a regular file", src))
	}

	source, err := os.Open(src)
	defer source.Close()
	if err != nil { return 0, err }

	dstStat, err := os.Stat(dst)
	if err == nil && dstStat.IsDir() {
		dst = filepath.Join(dst, filepath.Base(src))
	}

	destination, err := os.Create(dst)
	if err != nil { return 0, err }
	defer destination.Close()

	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func CopyDir(src, dst string) error {
	srcName := filepath.Base(src)
	dstDir := dst

	_, ok := os.Stat(dstDir)
	if ok == nil { // if the dir exists, create a new dir inside it
		dstDir = filepath.Join(dst, srcName)
		_, ok = os.Stat(dstDir)
		if ok == nil {
			return errors.New(fmt.Sprintf("directory <%s> already exists", dstDir))
		}
	}

	os.Mkdir(dstDir, 0750)

	srcDir, ok := os.ReadDir(src)
	if ok != nil { return ok }

	for _, item := range srcDir {
		itemName :=  item.Name()
		destination := filepath.Join(dstDir, itemName)
		source := filepath.Join(src, itemName)
		
		if item.IsDir() {
			ok := CopyDir(source, destination)
			if ok != nil { return ok }
		} else {
			_, ok := CopyFile(source, destination)
			if ok != nil { return ok }
		}
	} 
	return nil
}

func MoveFile(src, dst string) (int, error) {
	srcStat, err := os.Stat(src)
	if err != nil { return 0, err }
	if !srcStat.Mode().IsRegular() {
		return 0, errors.New(fmt.Sprintf(
			"%s is not a regular file", src))
	}

	sourceBytes, err := os.ReadFile(src)
	if err != nil { return 0, err }
	
	if os.Remove(src) != nil { 
		return 0, err 
	}

	if dstStat, err := os.Stat(dst); err == nil && dstStat.IsDir() {
		dst = filepath.Join(dst, filepath.Base(src))
	}

	destination, err := os.Create(dst)
	if err != nil { return 0, err }
	defer destination.Close()

	return destination.Write(sourceBytes)
}

func EnsureDirExists(dirName string) error {
	if _, err := os.Stat(dirName); err != nil {
		err = os.MkdirAll(dirName, os.ModePerm)
		return err
	}
	return nil
}

func FileWithoutExt(fileName string) string {
	return fileName[:len(fileName) - len(filepath.Ext(fileName))]
}

func ChangeFileExt(fileName string, newExtension string) string {
	return fmt.Sprintf(
		"%s.%s", FileWithoutExt(fileName), newExtension)
}