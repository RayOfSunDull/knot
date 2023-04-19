package knot

import (
	"fmt"
	"os"
	"io"
	"time"
	"os/exec"
	"path/filepath"
	"archive/zip"
	"github.com/signintech/gopdf"
)

type CompressionLevel int

const (
	DEFAULT CompressionLevel = iota
	PREPRESS
	EBOOK
)


func IntToCompressionLevel(n int) CompressionLevel {
	if n >= 0 && n <= 2 { return CompressionLevel(n) }
	return CompressionLevel(0)
}


func CompressionLevelToGS(cl CompressionLevel) string {
	switch cl {
	case DEFAULT:
		return "/default"
	case PREPRESS:
		return "/prepress"
	case EBOOK:
		return "/ebook"
	}
	return ""
}


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

	pngFile, err := os.Create(dst)
	if err != nil { return err }

	_, err = io.Copy(pngFile, pngData)
	return err
}


// requires ghostscript
func CompressPDF(fileName string, compressionLevel int) error {
	compressedFileName := fmt.Sprintf("%s.pdf", fileName)

	cmd := exec.Command(
		"gs", "-dBATCH", "-dNOPAUSE", "-q",
		"-sDEVICE=pdfwrite",
		fmt.Sprintf("-dPDFSETTINGS=%s", CompressionLevelToGS(
			IntToCompressionLevel(compressionLevel))),
		fmt.Sprintf("-sOutputFile=%s", compressedFileName),
		fileName)

	err := cmd.Run()
	if err != nil { return err }

	err = os.Remove(fileName)
	if err != nil { return err }
	
	err = os.Rename(compressedFileName, fileName)

	return err
}


func ExportBatch(batchNumber int, pi *ProjectInfo, compressionLevel int) (string, error) {
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
	
	pageRegexp := GetPageRegexp(".png")
	exportDir, err := os.ReadDir(exportPath)
	if err != nil { return "", err }

	pageSizePtr := gopdf.PageSizeA4
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *pageSizePtr})
	
	for _, item := range exportDir {
		itemName := item.Name()
		
		if item.IsDir() || !pageRegexp.Match([]byte(itemName)) {
			continue
		}

		imgPath := filepath.Join(exportPath, itemName)
		
		pdf.AddPage()
		pdf.Image(imgPath, 0, 0, pageSizePtr)
	}

	outputPath := filepath.Join(
		batchPath, fmt.Sprintf("%s.pdf", filepath.Base(batchPath)))
	
	pdf.WritePdf(outputPath)

	if compressionLevel != 0 { 
		CompressPDF(outputPath, compressionLevel)
	}

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