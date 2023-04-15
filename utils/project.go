package knot

import (
	"fmt"
	"os"
	"encoding/json"
	"path/filepath"
	"errors"
	"regexp"
)


type ProjectInfo struct {
	ProjectDir string
	ContentDir string
	ContentName string
	ExportDirName string
	TemplateName string
}


type Projects map[string]ProjectInfo


func (projects *Projects) Save(file string) error {
	destination, err := os.Create(file)
	if err != nil { return err }
	defer destination.Close()

	infoAsBytes, err := json.MarshalIndent(*projects, "", "\t")
	if err != nil { return err }

	_, err = destination.Write(infoAsBytes)
	return err
}


func GetProjects(file string) (Projects, error) {
	var result Projects

	fileBytes, err := os.ReadFile(file)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(fileBytes, &result)
	return result, err
}


func GetExistingProjectInfo(file string, projectName string) (ProjectInfo, error) {
	projects, err := GetProjects(file)
	if err != nil {
		return ProjectInfo{}, err
	}

	result, ok := projects[projectName]
	if ok {
		err = nil
	} else {
		err = errors.New(fmt.Sprintf(
			"no project called <%s> in project list", projectName))
	}

	return result, err
}


func GetProjectInfo(flags *Flags, si *SystemInfo, projects *Projects) (ProjectInfo, error) {
	projectName := flags.InitDirName
	if projectName == "" { // see if $PWD is inside a project
		projectsByDir := ArrangeProjectsByDir(projects)
		return FindFirstParentProjectInfo(
			si.KnotWD, projects, &projectsByDir)
	}

	projectDir := filepath.Join(si.KnotWD, projectName)

	contentDir := filepath.Join(projectDir, flags.ContentDirName)
	contentName := flags.ContentName
	if contentName == "" {
		contentName = projectName
	}

	return ProjectInfo{
		ProjectDir: projectDir,
		ContentDir: contentDir,
		ContentName: contentName,
		ExportDirName: flags.ExportDirName,
		TemplateName: flags.TemplateName}, nil
}


func GetContentRegexp(name string) *regexp.Regexp {
	result, _ := regexp.Compile(fmt.Sprintf(
		"^%s-[0-9]+", name))
	return result
}


func GetPageRegexp(extension string) *regexp.Regexp {
	result, _ := regexp.Compile(fmt.Sprintf(
		"^page-[0-9]+%s$", extension))
	return result
}


func NumberOfMatches(dirPath string, re *regexp.Regexp) (int, error) {
	dir, err := os.ReadDir(dirPath)
	if err != nil { return 0, err }

	var result int = 0
	for _, item := range dir {
		if re.Match([]byte(item.Name())) {
			result += 1
		}
	}
	return result, nil
}


func ArrangeProjectsByDir(projects *Projects) map[string]string {
	result := make(map[string]string)
	for projectName, projectInfo := range *projects {
		result[projectInfo.ProjectDir] = projectName
	}
	return result
}


func FindFirstParentProjectInfo(wd string, projects *Projects, projectsByDir *map[string]string) (ProjectInfo, error) {
	if projectName, ok := (*projectsByDir)[wd]; ok {
		result, _ := (*projects)[projectName]
		return result, nil
	}
	wdParent := filepath.Dir(wd)
	if wd == wdParent {
		return ProjectInfo{}, errors.New(
			"working directory is not part of a project")
	}
	return FindFirstParentProjectInfo(wdParent, projects, projectsByDir)
}