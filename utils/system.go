package knot

import (
	"fmt"
	"os"
	"os/exec"
	"encoding/json"
	"path/filepath"
	"errors"
	"io"
)


type TempConfigInfo struct {
	KnotWD string
}


type SystemInfo struct {
	WD string
	ConfigDir string
	TempConfigFile string
	ProjectsFile string
	TemplateDir string
}


func GetSystemInfo() (SystemInfo, error) {
	var wd string

	tempConfigFile := filepath.Join("/tmp", "knotconfig.json")

	tempConfigBytes, errRead := os.ReadFile(tempConfigFile)

	var tci TempConfigInfo
	errUnmarshal := json.Unmarshal(tempConfigBytes, &tci)

	if errRead == nil && errUnmarshal == nil {
		wd = tci.KnotWD
	} else {
		pwd, err := os.Getwd()
		if err != nil { return SystemInfo{}, err }

		wd = pwd
	}

	homeDir, err := os.UserHomeDir()
	if err != nil { return SystemInfo{}, err }

	configDir := filepath.Join(homeDir, ".config", "knot")
	projectsFile := filepath.Join(configDir, "projects.json")
	templateDir := filepath.Join(configDir, "templates")

	return SystemInfo{
		WD: wd,
		ConfigDir: configDir,
		TempConfigFile: tempConfigFile,
		ProjectsFile: projectsFile,
		TemplateDir: templateDir}, nil
}


func SetTempConfigInfo(si *SystemInfo, tci *TempConfigInfo) error {
	tempConfigInfoBytes, err := json.MarshalIndent(*tci, "", "\t")
	if err != nil { return err }

	tempConfigFileName := si.TempConfigFile

	_, err = os.Stat(tempConfigFileName)
	if err == nil { 
		err = os.Remove(tempConfigFileName)
		if err != nil { return err }
	}

	tempConfigFile, err := os.Create(tempConfigFileName)
	if err != nil { return err }
	defer tempConfigFile.Close()

	_, err = tempConfigFile.Write(tempConfigInfoBytes)
	return err
}


func SetTempKnotWD(si *SystemInfo, knotWD string) error {
	absKnotWD, err := filepath.Abs(knotWD)
	if err != nil { return err }

	return SetTempConfigInfo(
		si, &TempConfigInfo{KnotWD: absKnotWD})
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