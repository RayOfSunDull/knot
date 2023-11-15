package knot

import (
	"fmt"
	"os"
	"os/exec"
	"encoding/json"
	"path/filepath"
	"errors"
	"io"
    "github.com/glycerine/zygomys/v6/zygo"
)

type CommandRunner interface {
    Run(inputs []string) (string, error)
    Start(inputs []string) error
}

type SimpleCommandRunner struct {
    commandName string
}
func NewSimpleCommandRunner(commandName string) *SimpleCommandRunner {
    return &SimpleCommandRunner{ commandName: commandName }
}

func (runner *SimpleCommandRunner) Run(inputs []string) (string, error) {
    cmd := exec.Command(runner.commandName, inputs...)
    output, err := cmd.CombinedOutput()
    return string(output), err
}

func (runner *SimpleCommandRunner) Start(inputs []string) error {
    cmd := exec.Command(runner.commandName, inputs...)
    return cmd.Start()
}

type ZygoCommandRunner struct {
    zygoEnv *zygo.Zlisp
    zygoFunc *zygo.SexpFunction
}
func NewZygoCommandRunner(zygoEnv *zygo.Zlisp, zygoFunc *zygo.SexpFunction) *ZygoCommandRunner {
    return &ZygoCommandRunner{ zygoEnv: zygoEnv, zygoFunc: zygoFunc }
}

func (runner *ZygoCommandRunner) Run(inputs []string) (string, error) {
    sexpInputs := make([]zygo.Sexp, len(inputs))
    for i, input := range inputs {
        sexpInputs[i] = &zygo.SexpStr{ S: input }
    }

    sexp, err := runner.zygoEnv.Apply(runner.zygoFunc, sexpInputs)
    return sexp.SexpString(zygo.NewPrintState()), err
}

func (runner *ZygoCommandRunner) Start(inputs []string) error {
    _, err := runner.Run(inputs)
    return err
}


func CommandRunnerFromZygoEnv(zygoEnv *zygo.Zlisp, sexpName string, defaultRunner CommandRunner) CommandRunner {
    sexp, found := zygoEnv.FindObject(sexpName)
    if !found {
        return defaultRunner
    }

    switch sexp.(type) {
    case *zygo.SexpStr:
        return NewSimpleCommandRunner(string(sexp.(*zygo.SexpStr).S))
    case *zygo.SexpFunction:
        return NewZygoCommandRunner(zygoEnv, sexp.(*zygo.SexpFunction))
    default:
        return defaultRunner
    }
}

func IntFromZygoEnv(zygoEnv *zygo.Zlisp, sexpName string, defaultInt int) int {
    sexp, found := zygoEnv.FindObject(sexpName)
    if !found {
        return defaultInt
    }

    switch sexp.(type) {
    case *zygo.SexpInt:
        return int(sexp.(*zygo.SexpInt).Val)
    default:
        return defaultInt
    }
}

type ConfigInfo struct {
	PDFReader CommandRunner
	FileExplorer CommandRunner
	ExportQuality int
}

func LoadConfigInfo(configFile string) ConfigInfo {
    var configInfo ConfigInfo
    configInfo.PDFReader = NewSimpleCommandRunner("evince")
    configInfo.FileExplorer = NewSimpleCommandRunner("nautilus")
    configInfo.ExportQuality = 100

    zygoEnv := zygo.NewZlisp()

    sexps, err := zygoEnv.ParseFile(configFile)
    if err != nil {
        return configInfo
    }
    err = zygoEnv.LoadExpressions(sexps)
    if err != nil {
        return configInfo
    }
    _, err = zygoEnv.Run()
    if err != nil {
        return configInfo
    }

    configInfo.PDFReader = CommandRunnerFromZygoEnv(
        zygoEnv, "PDFReader", configInfo.PDFReader)
    configInfo.FileExplorer = CommandRunnerFromZygoEnv(
        zygoEnv, "FileExplorer", configInfo.FileExplorer)
    configInfo.ExportQuality = IntFromZygoEnv(
        zygoEnv, "ExportQuality", configInfo.ExportQuality)

    return configInfo
}


type TempConfigInfo struct {
	KnotWD string
}


type SystemInfo struct {
	ConfigInfo
	TempConfigInfo
	ConfigDir string
	ConfigFile string
	TempConfigFile string
	ProjectsFile string
	TemplateDir string
	ExportScript string
	PythonCommand string
}


func GetSystemInfo(platform Platform) (SystemInfo, error) {
    platformDirs, err := platform.GetPlatformDirs()
    if err != nil {
        return SystemInfo{}, err
    }
	tempDir := platformDirs.TempDir
	sysConfigDir := platformDirs.ConfigDir

	tempConfigFile := filepath.Join(tempDir, "knotconfig.json")

	tempConfigBytes, errRead := os.ReadFile(tempConfigFile)

	var tci TempConfigInfo
	errUnmarshal := json.Unmarshal(tempConfigBytes, &tci)

	if errRead != nil || errUnmarshal != nil {
		knotWD, err := os.Getwd()
		if err != nil { return SystemInfo{}, err }

		tci = TempConfigInfo{KnotWD: knotWD}
	}

	configDir := filepath.Join(sysConfigDir, "knot")
	projectsFile := filepath.Join(configDir, "projects.json")
	templateDir := filepath.Join(configDir, "templates")
	configFile := filepath.Join(configDir, "config.zy")
	exportScript := filepath.Join(configDir, "export.py")

    ci := LoadConfigInfo(configFile)

	pythonCommand, err := platform.GetPythonCommand()
	if err != nil { return SystemInfo{}, err}

	return SystemInfo{
		ConfigInfo: ci,
		TempConfigInfo: tci,
		ConfigDir: configDir,
		ConfigFile: configFile,
		TempConfigFile: tempConfigFile,
		ProjectsFile: projectsFile,
		TemplateDir: templateDir,
		ExportScript: exportScript,
		PythonCommand: pythonCommand}, nil
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


func OpenFile(si *SystemInfo, file string, open bool) error {
	if !open { return nil }

	extension := filepath.Ext(file)

	switch extension {
	case ".kra":
		cmd := exec.Command("krita", file)
		return cmd.Start()
	case ".pdf":
        return si.PDFReader.Start([]string{file})
	case "": // directory
        return si.FileExplorer.Start([]string{file})
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
