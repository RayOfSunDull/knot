package knot

import (
	"fmt"
	"os"
	"io"
	"time"
	"os/exec"
	"path/filepath"
	"archive/zip"
	// "github.com/signintech/gopdf"
)


func ExportToPNG(src, dst string) error {
	srcStat, err := os.Stat(src)
	if err != nil { return err }
	srcModTime := srcStat.ModTime()

	localLoc, _ := time.LoadLocation("Local")
	dstModTime := time.Date(0,time.January,0,0,0,0,0,localLoc)

	dstStat, err := os.Stat(dst)
	if err == nil {
		dstModTime = dstStat.ModTime()
	}

	if srcModTime.Before(dstModTime) {
		return nil
	}

	srcReader, err := zip.OpenReader(src)
	if err != nil { return err }
	defer srcReader.Close()

	pngData, err := srcReader.Open("mergedimage.png")
	if err != nil { return err }
	defer pngData.Close()

	pngFile, err := os.Create(dst)
	if err != nil { return err }
	defer pngFile.Close()

	_, err = io.Copy(pngFile, pngData)
	return err
}


func ExportBatch(batchNumber int, pi *ProjectInfo, si *SystemInfo) (string, error) {
	batchPath := filepath.Join(
		pi.ContentDir, GetBatchName(pi, batchNumber))

	batchDir, err := os.ReadDir(batchPath)
	if err != nil { return "", err }

	exportPath := filepath.Join(batchPath, pi.ExportDirName)
	if EnsureDirExists(exportPath) != nil {
		return "", err
	}

	for _, item := range batchDir {
		itemName := item.Name()
		extension := filepath.Ext(itemName)
		if extension != ".kra" || item.IsDir() { 
			continue 
		}

		src := filepath.Join(batchPath, itemName)
		
		dst := filepath.Join(exportPath, 
			ChangeFileExt(itemName, "png"))

		ExportToPNG(src, dst)
	}

	outputPath := filepath.Join(
		batchPath, fmt.Sprintf("%s.pdf", filepath.Base(batchPath)))

	exportArgs := []string{
		si.ExportScript,
		"-o", outputPath,
		"-q", fmt.Sprintf("%v", si.ExportQuality)}

	exportDir, err := os.ReadDir(exportPath)
	if err != nil { return "", err }
	
	pageRegexp := GetPageRegexp(".png")
	for _, item := range exportDir {
		itemName := item.Name()
		
		if item.IsDir() || !pageRegexp.Match([]byte(itemName)) {
			continue
		}
		
		exportArgs = append(exportArgs, filepath.Join(exportPath, itemName))
	}
	
	cmd := exec.Command(si.PythonCommand, exportArgs...)

	return outputPath, cmd.Run()
}