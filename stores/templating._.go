package stores

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/tuxounet/k2/libs"
)

type TemplatingStore struct {
	plan *ActionPlan
}

func NewTemplatingStore(plan *ActionPlan) *TemplatingStore {
	return &TemplatingStore{
		plan: plan,
	}
}

func (t *TemplatingStore) ApplyTemplate(templateApplyId string, templateHash string, produceGitIgnore bool) (bool, error) {
	fmt.Println("apply template", templateApplyId)
	template, ok := t.plan.Templates[templateHash]
	if !ok {
		return false, fmt.Errorf("template not found: %s", templateHash)

	}

	apply, err := t.plan.GetEntityAsTemplateApply(templateApplyId)
	if err != nil {
		return false, err
	}

	firstTime, err := t.isFirstTimeApply(templateApplyId)
	if err != nil {
		return false, err
	}
	if firstTime {
		err = template.ExecuteBootstrap(apply)
		if err != nil {
			return false, err
		}
		err = apply.ExecuteBootstrap()
		if err != nil {
			return false, err
		}
	}

	err = template.ExecutePre(apply)
	if err != nil {
		return false, err
	}
	err = apply.ExecutePre()
	if err != nil {
		return false, err
	}

	templateSourceFolder := template.K2.Metadata.Folder

	allTemplateFiles, err := getAllFiles(templateSourceFolder)
	if err != nil {
		return false, err
	}

	copyMap := make(map[string]string)

	for _, file := range allTemplateFiles {

		sourcePath := filepath.Join(templateSourceFolder, file)

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return false, err
		}
		if fileInfo.IsDir() {
			continue
		}

		targetpath := filepath.Join(apply.K2.Metadata.Folder, file)

		copyMap[sourcePath] = targetpath

	}

	for source, destination := range copyMap {
		fmt.Println("copy", source, destination)

		err = copyFile(source, destination, libs.MergeMaps(template.K2.Body.Parameters, apply.K2.Body.Vars))
		if err != nil {
			return false, err
		}
	}

	err = createGitIgnore(copyMap, apply.K2.Metadata.Folder)
	if err != nil {
		return false, err
	}

	err = template.ExecutePost(apply)
	if err != nil {
		return false, err
	}
	err = apply.ExecutePost()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (t *TemplatingStore) DestroyTemplate(templateApplyId string) error {
	fmt.Println("destroy template", templateApplyId)

	apply, err := t.plan.GetEntityAsTemplateApply(templateApplyId)
	if err != nil {
		return err
	}

	err = apply.ExecuteNuke()
	if err != nil {
		return err
	}

	folder := apply.K2.Metadata.Folder
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		return nil
	}

	gitIgnoreFile := filepath.Join(folder, ".gitignore")

	if _, err := os.Stat(gitIgnoreFile); os.IsNotExist(err) {
		return nil
	}

	files, err := os.ReadFile(gitIgnoreFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(files), "\n")
	lines = append(lines, ".gitignore")
	for _, line := range lines {
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "!") {
			continue
		}

		file := filepath.Join(folder, line)
		if _, err := os.Stat(file); os.IsNotExist(err) {
			continue
		}

		stat, err := os.Stat(file)
		if err != nil {
			return err
		}

		if stat.IsDir() {
			err = os.RemoveAll(file)
			if err != nil {
				return err
			}
		} else {
			err = os.Remove(file)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *TemplatingStore) isFirstTimeApply(templateApplyId string) (bool, error) {
	apply, err := t.plan.GetEntityAsTemplateApply(templateApplyId)
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(apply.K2.Metadata.Folder); os.IsNotExist(err) {
		return true, nil
	}

	gitIgnoreFile := filepath.Join(apply.K2.Metadata.Folder, ".gitignore")

	if _, err := os.Stat(gitIgnoreFile); os.IsNotExist(err) {
		return true, nil
	}

	//Read the .gitignore file
	fileContent, err := os.ReadFile(gitIgnoreFile)
	if err != nil {
		return false, err
	}

	present := slices.Contains(strings.Split(string(fileContent), "\n"), "!k2.apply")
	return present, nil

}

func (t *TemplatingStore) cleanupEmptyDirs(folder string) error {
	var folders []string
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			folders = append(folders, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	sort.Slice(folders, func(i, j int) bool {
		return len(folders[i]) > len(folders[j])
	})

	for _, folder := range folders {

		sub, err := os.ReadDir(folder)
		if err != nil {
			return err
		}
		if len(sub) > 0 {
			continue
		}

		err = os.Remove(folder)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TemplatingStore) destroyTemplateRef(folder string) error {
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		return nil
	}

	err := os.RemoveAll(folder)
	if err != nil {
		return err
	}

	return nil
}
