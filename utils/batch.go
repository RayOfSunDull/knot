package knot

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func GetBatchName(pi *ProjectInfo, batchNumber int) string {
	return fmt.Sprintf("%s-%d", pi.ContentName, batchNumber)
}


func GetBatchDir(pi *ProjectInfo, batchNumber int) string {
	return filepath.Join(
		pi.ContentDir, GetBatchName(pi, batchNumber))
}


func GetPageName(pageNumber int) string {
	return fmt.Sprintf("page-%d.kra", pageNumber)
}


func MakeBatch(templatePath string, pi *ProjectInfo, batchNumber int, open bool) error {
	templateBatchDir := filepath.Join(templatePath, "batch")

	newBatchDir := filepath.Join(
		pi.ContentDir, GetBatchName(pi, batchNumber))
	
	if err := CopyDir(templateBatchDir, newBatchDir); err != nil {
		return err
	}
	
	page0 := filepath.Join(newBatchDir, "page-0.kra")
	_, err := MoveFile(
		filepath.Join(newBatchDir, "page.kra"),
		page0)

	OpenFile(page0, open)
	return err
}


func CreateProject(templatePath string, si *SystemInfo, pi *ProjectInfo, open bool) error {
	if _, err := os.Stat(pi.ProjectDir); err == nil {
		fmt.Printf("directory <%s> already exists. Assuming you simply want to register it instead of creating a new project\n", pi.ProjectDir)
		return nil
	}

	if err := EnsureDirExists(pi.ProjectDir); err != nil {
		return err
	}
	if err := EnsureDirExists(pi.ContentDir); err != nil {
		return err
	}

	return MakeBatch(templatePath, pi, 0, open)
}


func MakePage(templatePath string, pi *ProjectInfo, batchNumber int, open bool) error {
	batchDir := GetBatchDir(pi, batchNumber)
		
	pageNumber, err := NumberOfMatches(
		batchDir, GetPageRegexp(".kra"))
	if err != nil { return err }
	newPage := filepath.Join(batchDir, GetPageName(pageNumber))

	_, err = CopyFile(
		filepath.Join(templatePath, "batch", "page.kra"),
		newPage)
	if err != nil { return err }
	
	OpenFile(newPage, open)

	return nil
}


func OpenKraFilesInBatch(pi *ProjectInfo, batchNumber int, open bool) error {
	if !open { return nil }
	
	batchPath := GetBatchDir(pi, batchNumber)

	batchDir, err := os.ReadDir(batchPath)
	if err != nil { return err }

	pageRegexp := GetPageRegexp(".kra")
	
	pages := make([]string, 0, 10)
	pages = append(pages, "krita")

	for _, item := range batchDir {
		if item.IsDir() { continue }
		itemName := item.Name()

		if pageRegexp.Match([]byte(itemName)) {
			pages = append(pages, filepath.Join(batchPath, itemName))
		}
	}

	cmd := exec.Command("nohup", pages...)
	return cmd.Start()
}
