package stores

import (
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func executeScript(data interface{}, stage string, destinationFolder string) error {
	// Implement the script execution logic here
	return nil
}

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
		excludedFiles := []string{".gitignore", "k2.template.yaml", ".DS_Store"}
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

	source, err := os.ReadFile(src)
	if err != nil {

		return err
	}

	tpl, err := template.New("template").Funcs(sprig.FuncMap()).Parse(string(source))
	if err != nil {
		return err
	}

	var outBuffer strings.Builder
	outIO := io.MultiWriter(&outBuffer)

	err = tpl.Execute(outIO, tplData)
	if err != nil {
		return err
	}

	err = os.WriteFile(dest, []byte(outBuffer.String()), os.ModePerm)
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

	err := os.WriteFile(ignorePath, []byte(strings.Join(ignoreContent, "\n")), os.ModePerm)

	if err != nil {
		return err
	}
	return nil
}
