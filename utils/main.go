package knot

import (
	"fmt"
	"path/filepath"
	"log"
)


func Main(platform Platform) {
	flags := GetFlags()

	systemInfo, err := GetSystemInfo(platform)
	if err != nil { log.Fatal(err) }

	projects, err := GetProjects(systemInfo.ProjectsFile)
	if err != nil { log.Fatal(err) }

	projectInfo, errProjectInfo := GetProjectInfo(
		&flags, &systemInfo, &projects)

	if errProjectInfo != nil && false { fmt.Println(errProjectInfo) }

	projectName := filepath.Base(projectInfo.ProjectDir)

	templatePath := filepath.Join(
		systemInfo.TemplateDir, projectInfo.TemplateName)
	
	contentRegexp := GetContentRegexp(
		projectInfo.ContentName)

	open := !flags.SilentMode

	if flags.InitDirName != "" {
		err := CreateProject(
			templatePath, &systemInfo, &projectInfo, open)
		if err != nil { log.Fatal(err) }

		projects[projectName] = projectInfo

		err = SetTempKnotWD(&systemInfo, flags.InitDirName)
		if err != nil { log.Fatal(err) }
	}

	if flags.NextBatch {
		batchNumber, err := NumberOfMatches(
			projectInfo.ContentDir, contentRegexp)
		if err != nil { log.Fatal(err) }

		MakeBatch(
			templatePath, &systemInfo, &projectInfo, 
			batchNumber, open)
	}

	if flags.SpecifiedBatch >= 0 {
		MakeBatch(
			templatePath, &systemInfo, &projectInfo, 
			flags.SpecifiedBatch, open)
	}

	if flags.NextPage {
		latestBatch, err := NumberOfMatches(
			projectInfo.ContentDir, contentRegexp)
		if err != nil { log.Fatal(err) }
		latestBatch -= 1
		
		MakePage(
			templatePath, &systemInfo, &projectInfo, 
			latestBatch, open)
	}

	if flags.SpecifiedPage >= 0 {
		MakePage(
			templatePath, &systemInfo, &projectInfo, 
			flags.SpecifiedPage, open)
	}

	if flags.ExportLatestBatch {
		latestBatch, err := NumberOfMatches(
			projectInfo.ContentDir, contentRegexp)
		if err != nil { log.Fatal(err) }
		latestBatch -= 1
		
		var output string
		
		output, err = ExportBatch(
			latestBatch, &projectInfo, &systemInfo)
		if err != nil { log.Fatal(err) }

		OpenFile(&systemInfo, output, open)
	}

	if flags.ExportSpecifiedBatch >= 0 {
		var output string

		output, err = ExportBatch(
			flags.ExportSpecifiedBatch, &projectInfo, &systemInfo)
		if err != nil { log.Fatal(err) }

		OpenFile(&systemInfo, output, open)
	}

	if flags.DeregisterProject != "" {
		deregisteredProjectName := filepath.Base(
			flags.DeregisterProject)
		deregisteredProjectInfo, ok := projects[deregisteredProjectName]
		if ok {
			delete(projects, deregisteredProjectName)

			fmt.Println(fmt.Sprintf(
				"deregistered project <%s>, found in <%s>",
				deregisteredProjectName, 
				deregisteredProjectInfo.ProjectDir))
		}
	}

	if flags.OpenProject != "" {
		info, err := GetExistingProjectInfo(
			systemInfo.ProjectsFile, flags.OpenProject)
		if err != nil { log.Fatal(err) }
		OpenFile(&systemInfo, info.ProjectDir, true)

		latestBatch, err := NumberOfMatches(
			info.ContentDir, 
			GetContentRegexp(info.ContentName))
		if err != nil { log.Fatal(err) }
		latestBatch -= 1
		
		err = OpenKraFilesInBatch(
			&info, latestBatch, open)
		if err != nil { log.Fatal(err) }

		err = SetTempKnotWD(&systemInfo, info.ProjectDir)
		if err != nil { log.Fatal(err) }
	}

	if flags.OpenBatch >= 0 {
		err := OpenKraFilesInBatch(
			&projectInfo, flags.OpenBatch, true)
		if err != nil { log.Fatal(err) }
	}

	if flags.ListProjects {
		fmt.Printf("registered projects:\n")
		for name, info := range projects {
			fmt.Printf("\t project <%s> in <%s>\n", name, info.ProjectDir)
		}
	}

	if flags.PrintWD {
		fmt.Println(systemInfo.KnotWD)
	}

	if flags.SetWD != "" {
		err := SetTempKnotWD(&systemInfo, flags.SetWD)
		if err != nil { log.Fatal(err) }
	}

	if err = projects.Save(systemInfo.ProjectsFile); err != nil {
		log.Fatal(err)
	}
}