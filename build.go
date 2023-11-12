package main

import (
	"fmt"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
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


func main() {
	hostSystem := system{
		osys: stringToOS(runtime.GOOS), 
		arch: stringToArch(runtime.GOARCH)}

	noBuild := flag.Bool("no-build", false, "if true, the binaries will not be built from source")
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
	
	if ! *noBuild {
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
			}
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
		var err error

		switch hostSystem.osys {
		case LINUX:
			linux := knot.Linux{}
            platformDirs, err := linux.GetPlatformDirs()
			if err != nil { log.Fatal(err) }

			configDir = platformDirs.ConfigDir
			configDir = filepath.Join(configDir, "knot")

            binDir = platformDirs.BinDir
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
		
		exportInstall := filepath.Join(configDir, "export.py")
		if ! *quiet {
			fmt.Printf("copying export script to %s\n", exportInstall)
		}
		_, err = knot.CopyFile(
			filepath.Join(utilsDir, "export.py"),
			exportInstall)
		if err != nil { log.Fatal(err) }
	}
}
