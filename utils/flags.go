package knot

import (
	"flag"
	// "path/filepath"
)

type Flags struct {
	SilentMode bool
	ContentDirName string
	ContentName string
	InitDirName string
	NextBatch bool
	SpecifiedBatch int
	NextPage bool
	SpecifiedPage int
	ExportLatestBatch bool
	ExportSpecifiedBatch int
	ExportDirName string
	TemplateName string
	DeregisterProject string
	OpenProject string
	OpenBatch int
	ListProjects bool
	PrintWD bool
	SetWD string
}

func GetFlags() Flags {
	silentModePtr := flag.Bool("s", false, "silent mode; disable automatic opening of files")

	contentDirNamePtr := flag.String("cd", "", "name the directory of the content files. If none is specified, they will be dumped at the top level of the project directory")

	contentNamePtr := flag.String("c", "", "name of the content files. If none is specified, the name of the project directory will be used")

	initDirNamePtr := flag.String("i", "", "initialise new project directory with the given name. By default, it will be created in $PWD")
	
	nextBatchPtr := flag.Bool("b", false, "create the next batch of notes")

	specifiedBatchPtr := flag.Int("sb", -1, "create a new batch of notes with a specified batch number")

	nextPagePtr := flag.Bool("p", false, "create the next batch of notes")

	specifiedPagePtr := flag.Int("sp", -1, "create a new page in a  batch of notes with a specified batch number")

	exportLatestBatchPtr := flag.Bool("e", false, "export the latest batch to pdf")

	exportSpecifiedBatchPtr := flag.Int("se", -1, "export a batch with specified batch number to pdf")

	exportDirNamePtr := flag.String("ed", "export", "the subdirectory in each batch where all pages will be exported to pngs")

	templateNamePtr := flag.String("t", "default", "the template used for initialising the new project directory")

	deregisterProjectPtr := flag.String("d", "", "deregister: remove current project from projects list")

	openProjectPtr := flag.String("o", "", "open the latest batch of a given project")

	openBatchPtr := flag.Int("ob", -1, "open the batch with the given number in krita. Ignores silent mode")

	listProjectsPtr := flag.Bool("l", false, "list all registered projects")

	printWD := flag.Bool("pwd", false, "print the current knot working directory")

	setWD := flag.String("wd", "", "set the current knot working directory")

	flag.Parse()

	return Flags{
		SilentMode: 			*silentModePtr,
		ContentDirName: 		*contentDirNamePtr,
		ContentName: 			*contentNamePtr,
		InitDirName: 			*initDirNamePtr,
		NextBatch: 				*nextBatchPtr,
		SpecifiedBatch: 		*specifiedBatchPtr,
		NextPage: 				*nextPagePtr,
		SpecifiedPage: 			*specifiedPagePtr,
		ExportLatestBatch: 		*exportLatestBatchPtr,
		ExportSpecifiedBatch: 	*exportSpecifiedBatchPtr,
		ExportDirName:			*exportDirNamePtr,
		TemplateName:			*templateNamePtr,
		DeregisterProject:		*deregisterProjectPtr,
		OpenProject:			*openProjectPtr,
		OpenBatch:				*openBatchPtr,
		ListProjects:			*listProjectsPtr,
		PrintWD:				*printWD,
		SetWD:					*setWD}
}