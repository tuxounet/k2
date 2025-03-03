package stores

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/tuxounet/k2/libs"
)

func getAllFiles(folder string) ([]string, error) {
	var files []string
	err := filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		files = append(files, path)
		return nil
	})

	fileList := make([]string, 0)

	for _, file := range files {
		relPath, err := filepath.Rel(folder, file)
		if err != nil {
			return nil, err
		}
		if relPath == "." {
			continue
		}
		if strings.HasPrefix(filepath.Join(folder, ".git"), file) {
			continue
		}
		fileName := filepath.Base(relPath)
		excludedFiles := []string{"k2.template.yaml", ".DS_Store"}
		if slices.Contains(excludedFiles, fileName) {
			continue
		}

		fileList = append(fileList, relPath)

	}

	return fileList, err
}

func copyFile(src string, dest string, tplData any) error {

	dir := path.Dir(dest)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	fileName := filepath.Base(dest)
	if fileName == ".gitignore" {
		return nil
	}

	source, err := os.ReadFile(src)
	if err != nil {

		return err
	}

	target, err := libs.RenderTemplate(string(source), tplData)
	if err != nil {
		return err
	}

	err = os.WriteFile(dest, target, os.ModePerm)
	if err != nil {
		return err
	}

	return nil

}

func createGitIgnore(files map[string]string, destinationFolder string) error {

	var ignoreContent []string

	ignoreContent = append(ignoreContent, "!k2.apply.yaml")

	ignorePath := filepath.Join(destinationFolder, ".gitignore")

	for _, dest := range files {
		relPath, err := filepath.Rel(destinationFolder, dest)
		if err != nil {
			return err
		}
		fileName := filepath.Base(relPath)
		if fileName == "k2.apply.yaml" {
			ignoreContent = append(ignoreContent, "!"+relPath)
			continue
		}
		ignoreContent = append(ignoreContent, relPath)
	}

	ignoreContent = slices.Compact(ignoreContent)
	slices.Sort(ignoreContent)

	err := os.WriteFile(ignorePath, []byte(strings.Join(ignoreContent, "\n")), os.ModePerm)

	if err != nil {
		return err
	}
	return nil
}
