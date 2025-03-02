package stores

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

	err = template.ExecuteBootstrap(apply)
	if err != nil {
		return false, err
	}
	err = apply.ExecuteBootstrap()
	if err != nil {
		return false, err
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

		err = copyFile(source, destination, apply.K2.Body.Vars)
		if err != nil {
			return false, err
		}
	}

	err = createGitIgnore(copyMap, apply.K2.Metadata.Folder)
	if err != nil {
		return false, err
	}

	err = apply.ExecutePost()
	if err != nil {
		return false, err
	}
	err = template.ExecutePost(apply)
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
	err = apply.ExecutePost()
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

	if _, err := os.Stat(gitIgnoreFile); os.IsExist(err) {
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
	}

	var folders []string
	err = filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			folders = append(folders, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

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
