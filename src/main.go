package main

import (
	"fmt"
	"path/filepath"
	"knot/utils"
	"log"
)


func main() {
	flags := knot.GetFlags()

	systemInfo, err := knot.GetSystemInfo()
	if err != nil { log.Fatal(err) }

	projects, err := knot.GetProjects(systemInfo.ProjectsFile)
	if err != nil { log.Fatal(err) }

	projectInfo, errProjectInfo := knot.GetProjectInfo(
		&flags, &systemInfo, &projects)

	if errProjectInfo != nil && false { fmt.Println(errProjectInfo) }

	projectName := filepath.Base(projectInfo.ProjectDir)

	templatePath := filepath.Join(
		systemInfo.TemplateDir, projectInfo.TemplateName)
	
	contentRegexp := knot.GetContentRegexp(
		projectInfo.ContentName)

	open := !flags.SilentMode

	if flags.InitDirName != "" {
		err := knot.CreateProject(
			templatePath, &systemInfo, &projectInfo, open)
		if err != nil { log.Fatal(err) }

		projects[projectName] = projectInfo
	}

	if flags.NextBatch {
		batchNumber, err := knot.NumberOfMatches(
			projectInfo.ContentDir, contentRegexp)
		if err != nil { log.Fatal(err) }

		knot.MakeBatch(templatePath, &projectInfo, batchNumber, open)
	}

	if flags.SpecifiedBatch >= 0 {
		knot.MakeBatch(
			templatePath, &projectInfo, flags.SpecifiedBatch, open)
	}

	if flags.NextPage {
		latestBatch, err := knot.NumberOfMatches(
			projectInfo.ContentDir, contentRegexp)
		if err != nil { log.Fatal(err) }
		latestBatch -= 1
		
		knot.MakePage(templatePath, &projectInfo, latestBatch, open)
	}

	if flags.SpecifiedPage >= 0 {
		knot.MakePage(
			templatePath, &projectInfo, flags.SpecifiedPage, open)
	}

	if flags.ExportLatestBatch {
		latestBatch, err := knot.NumberOfMatches(
			projectInfo.ContentDir, contentRegexp)
		if err != nil { log.Fatal(err) }
		latestBatch -= 1

		output, err := knot.ExportBatch(latestBatch, &projectInfo)
		if err != nil { log.Fatal(err) }

		knot.OpenFile(output, open)
	}

	if flags.ExportSpecifiedBatch >= 0 {
		output, err := knot.ExportBatch(
			flags.ExportSpecifiedBatch, &projectInfo)
		if err != nil { log.Fatal(err) }

		knot.OpenFile(output, open)
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
		info, err := knot.GetExistingProjectInfo(
			systemInfo.ProjectsFile, flags.OpenProject)
		if err != nil { log.Fatal(err) }
		knot.OpenFile(info.ProjectDir, true)

		latestBatch, err := knot.NumberOfMatches(
			info.ContentDir, 
			knot.GetContentRegexp(info.ContentName))
		if err != nil { log.Fatal(err) }
		latestBatch -= 1
		
		err = knot.OpenKraFilesInBatch(
			&info, latestBatch, open)
		
		if err != nil { log.Fatal(err) }
	}

	if flags.OpenBatch >= 0 {
		err := knot.OpenKraFilesInBatch(
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
		fmt.Println(systemInfo.WD)
	}

	if flags.SetWD != "" {
		absWD, err := filepath.Abs(flags.SetWD)

		err = knot.SetTempKnotWD(&systemInfo, absWD)

		if err != nil { log.Fatal(err) }
	}

	if err = projects.Save(systemInfo.ProjectsFile); err != nil {
		log.Fatal(err)
	}
}