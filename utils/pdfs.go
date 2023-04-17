package knot

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"github.com/signintech/gopdf"
)


func ExportToPNG(src, dst string) error {
	dstStat, err := os.Stat(dst)
	if err == nil && dstStat.IsDir() {
		dst = filepath.Join(dst, 
			ChangeFileExt(filepath.Base(src), "png"))
	}
	cmd := exec.Command(
		"krita", src, "--export", 
		"--export-filename", dst)

	return cmd.Run()
}


func ExportBatch(batchNumber int, pi *ProjectInfo) (string, error) {
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

	exportDir, err := os.ReadDir(exportPath)
	if err != nil { return "", err }

	pageSizePtr := gopdf.PageSizeA4
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *pageSizePtr})
	
	for _, item := range exportDir {
		itemName := item.Name()
		extension := filepath.Ext(itemName)
		if extension != ".png" || item.IsDir() { 
			continue 
		}

		imgPath := filepath.Join(exportPath, itemName)
		
		pdf.AddPage()
		pdf.Image(imgPath, 0, 0, pageSizePtr)
	}

	outputPath := filepath.Join(
		batchPath, fmt.Sprintf("%s.pdf", filepath.Base(batchPath)))
	
	pdf.WritePdf(outputPath)

	return outputPath, nil
}


// legacy functions (they require external dependencies)


func ConvertToPDF(src string) error {
	dst := ChangeFileExt(src, "pdf")

	cmd := exec.Command("convert", src, dst)
	// requires imagemagick
	return cmd.Run()
}


func MergePDFs(srcFiles []string, dst string) error {
	GSArgs := []string{
		"-dBATCH", "-dNOPAUSE", "-q",
		"-sDEVICE=pdfwrite",
		"-dPDFSETTINGS=/prepress",
		fmt.Sprintf("-sOutputFile=%s", dst)}
	
	completeGSArgs := append(GSArgs, srcFiles...)

	cmd := exec.Command("gs", completeGSArgs...)
	// requires ghostscript
	err := cmd.Run()

	for _, file := range srcFiles {
		if os.Remove(file) != nil {
			fmt.Printf("could not remove file <%s>\n", file)
		}
	}

	return err
}


func LegacyExportBatch(batchNumber int, pi *ProjectInfo) (string, error) {
	batchPath := filepath.Join(
		pi.ContentDir, GetBatchName(pi, batchNumber))

	batchDir, err := os.ReadDir(batchPath)
	if err != nil { return "", err }

	exportPath := filepath.Join(batchPath, pi.ExportDirName)
	if EnsureDirExists(exportPath) != nil {
		return "", err
	}

	nPages := 0
	for _, item := range batchDir {
		itemName := item.Name()
		extension := filepath.Ext(itemName)
		if extension != ".kra" || item.IsDir(){ 
			continue 
		}

		src := filepath.Join(batchPath, itemName)
		
		dst := filepath.Join(exportPath, 
			ChangeFileExt(itemName, "png"))

		ExportToPNG(src, dst)
		ConvertToPDF(dst)
		nPages += 1
	}

	pageRegexp := GetPageRegexp(".pdf")
	pages := make([]string, 0, nPages)

	exportDir, err := os.ReadDir(exportPath)
	if err != nil { return "", err }

	for _, item := range exportDir {
		if item.IsDir() { continue }
		itemName := item.Name()

		if pageRegexp.Match([]byte(itemName)) {
			pages = append(pages, filepath.Join(exportPath, itemName))
		}
	}

	outputPath := filepath.Join(
		batchPath, fmt.Sprintf("%s.pdf", filepath.Base(batchPath)))
	err = MergePDFs(pages, outputPath)
	if err != nil { return "", err }
	
	return outputPath, nil
}