package main

import (
	"fmt"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"errors"
	"runtime"
	"log"
	"knot/utils"
)

type operatingSystem int
type architecture int

// other OS's / architectures to be implemented
const(
	LINUX operatingSystem = iota
	// WINDOWS operatingSystem = iota
	// MACOS operatingSystem = iota

	AMD64 architecture = iota
	// ARM64 architecture = iota
)

func stringToOS(osName string) operatingSystem {
	switch osName {
	case "linux":
		return LINUX
	default:
		return LINUX
	}
}

func stringToArch(archName string) architecture {
	switch archName {
	case "amd64":
		return AMD64
	default:
		return AMD64
	}
}

func (osys *operatingSystem) str() string {
	switch *osys {
	case LINUX:
		return "linux"
	default:
		return ""
	}
}

func (arch *architecture) str() string {
	switch *arch {
	case AMD64:
		return "amd64"
	default:
		return ""
	}
}

type system struct {
	osys operatingSystem
	arch architecture
}

func (target *system)use() {
	os.Setenv("GOOS", target.osys.str())
	os.Setenv("GOARCH", target.arch.str())
}

func (target *system) formatBin(name string) string {
	ext := ""
	// if target.osys == WINDOWS {
	// 	ext = ".exe"
	// }
	return fmt.Sprintf(
		"knot_%s_%s_%s%s", name, 
		target.osys.str(), target.arch.str(), ext)
}

func (target *system) buildGo(src,dst string) (string, error) {
	target.use()

	buildCmd := exec.Command("go", "build",
		"-o", dst, src)

	out, err := buildCmd.CombinedOutput()

	return string(out), err
}

func (target *system) buildCython(src,dst,aux string) error {
	cythonCmd := exec.Command("cython",
		src, "-o", aux,
		"--embed", "-3")

	if cythonErr := cythonCmd.Run(); cythonErr != nil {
		return cythonErr
	}
	switch target.osys {
	case LINUX:
		gccCmd := exec.Command("gcc",
			aux,
			"-Wno-deprecated-declarations",
			"-Wl,--copy-dt-needed-entries",
			"-o", dst,
			"-I/usr/include/python3.11",
			"-L./export-venv/lib/python3.11/site-packages",
			"-lpython3")
		
		if gccOut, gccErr := gccCmd.CombinedOutput(); gccErr != nil {
			return errors.New(fmt.Sprintf("%s:\n%s", gccErr, string(gccOut)))
		}

		return os.Remove(aux)
	default:
		return errors.New(fmt.Sprintf(
			"compiling for OS %s has not been implemented",
			target.osys.str()))
	}
}


func main() {
	hostSystem := system{
		osys: stringToOS(runtime.GOOS), 
		arch: stringToArch(runtime.GOARCH)}

	install := flag.Bool("install", false, "whether to install the program after building")
	quiet := flag.Bool("quiet", false, "suppress build logs")
	overwriteTemplates := flag.Bool("overwrite-templates", false, "whether to overwrite the templates of a previous installation, if it exists")

	flag.Parse()

	binaryNames := [1]string{
		"main"}
	targetOSs := [1]operatingSystem{
		LINUX}
	targetArchs := [1]architecture{
		AMD64}

	wd, err := os.Getwd()
	if err != nil { log.Fatal(err) }
	topSrcDir := filepath.Join(wd, "src")
	topBinDir := filepath.Join(wd, "bin")
	utilsDir := filepath.Join(wd, "utils")
	auxDir := filepath.Join(wd, "aux")

	for _, osys := range targetOSs {
		srcDir := filepath.Join(topSrcDir, osys.str())
		binDir := filepath.Join(topBinDir, osys.str())
		for _, arch := range targetArchs {
			target := system{osys: osys, arch: arch}

			for _, binName := range binaryNames {
				if ! *quiet {
					fmt.Printf("building %s for %s %s\n",
						binName, arch.str(), osys.str())
				}
				
				out, err := target.buildGo(
					filepath.Join(srcDir, fmt.Sprintf("%s.go", binName)),
					filepath.Join(binDir, target.formatBin(binName)))
				if err != nil { log.Fatal(out) }
			}

			err = target.buildCython(
				filepath.Join(utilsDir, "export.py"),
				filepath.Join(binDir, target.formatBin("export")),
				filepath.Join(auxDir, "export.c"))
			if err != nil { log.Fatal(err) }
		}
	}

	if *install {
		if ! *quiet {
			fmt.Printf(
				"installing knot for %s %s\n\n",
				hostSystem.arch.str(), hostSystem.osys.str())
		}
		var configDir string
		var binDir string
		exportFileName := "export"

		switch hostSystem.osys {
		case LINUX:
			unix := knot.Unix{}
			configDir, err = unix.GetConfigDir()
			if err != nil { log.Fatal(err) }

			binDir, err = unix.GetBinDir()
			if err != nil { log.Fatal(err) }
		default:
		}

		hostBinDir := filepath.Join(topBinDir, hostSystem.osys.str())

		knotInstall := filepath.Join(binDir, "knot")
		if ! *quiet {
			fmt.Printf("installing main in %s\n", knotInstall)
		}
		_, err = knot.CopyFile(
			filepath.Join(hostBinDir, hostSystem.formatBin("main")),
			knotInstall)
		if err != nil { log.Fatal(err) }
		
		if ! *quiet {
			fmt.Printf("creating %s\n", configDir)
		}
		err = knot.EnsureDirExists(configDir)
		if err != nil { log.Fatal(err) }
		
		templateInstall := filepath.Join(configDir, "templates")
		if *overwriteTemplates {
			os.RemoveAll(templateInstall)
		}
		if _, err = os.Stat(templateInstall); err != nil {
			if ! *quiet {
				fmt.Printf("copying templates to %s\n", templateInstall)
			}
			err = knot.CopyDir(
				filepath.Join(wd, "templates"),
				templateInstall)
			if err != nil { log.Fatal(err) }
		}
		
		exportInstall := filepath.Join(configDir, exportFileName)
		if ! *quiet {
			fmt.Printf("copying export binary to %s\n", exportInstall)
		}
		_, err = knot.CopyFile(
			filepath.Join(hostBinDir, hostSystem.formatBin("export")),
			exportInstall)
		if err != nil { log.Fatal(err) }
	}
}