package stores

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/tuxounet/k2/libs"
	"github.com/tuxounet/k2/types"

	"gopkg.in/yaml.v3"
)

func (t *TemplatingStore) resolveTemplateGit(hash string) (*types.IK2Template, error) {

	libs.WriteOutputf("resolve template git %s\n", hash)

	var gitRef *types.IK2TemplateRef
	for _, ref := range t.plan.Refs {
		if ref.Hash() == hash {
			gitRef = &ref
			break
		}
	}
	if gitRef == nil {
		return nil, libs.WriteErrorf("template not found: %s", hash)
	}

	refsFolder := filepath.Join(t.plan.inventory.InventoryDir, libs.RefsDir)
	if _, err := os.Stat(refsFolder); os.IsNotExist(err) {
		err := os.MkdirAll(refsFolder, os.ModePerm)
		if err != nil {
			return nil, libs.WriteErrorf("failed to create inventory folder: %w", err)
		}
	}

	templateGitFolder := filepath.Join(refsFolder, hash, ".git")

	targetCloneFolder := filepath.Join(refsFolder, hash)
	if _, err := os.Stat(templateGitFolder); os.IsNotExist(err) {
		cmd := exec.Command("git", "clone", "--single-branch", "--branch", gitRef.Params["branch"], gitRef.Params["repository"], hash)
		cmd.Dir = refsFolder
		output, err := cmd.CombinedOutput()
		if err != nil {
			return nil, libs.WriteErrorf("failed to clone repository: %w, %s", err, output)
		}

		libs.WriteOutputf("output: %s\n", output)
	} else {
		cmd := exec.Command("git", "pull")
		cmd.Dir = targetCloneFolder
		if err := cmd.Run(); err != nil {
			return nil, libs.WriteErrorf("failed to pull repository: %w", err)
		}
	}

	templatePath := filepath.Join(refsFolder, hash, gitRef.Params["path"])
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return nil, libs.WriteErrorf("template not found: %s", templatePath)
	}

	body, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, libs.WriteErrorf("failed to read template: %w", err)
	}

	var template types.IK2Template
	err = yaml.Unmarshal(body, &template)
	if err != nil {
		return nil, libs.WriteErrorf("failed to unmarshal template: %w", err)
	}

	template.K2.Metadata.Folder = filepath.Dir(templatePath)
	template.K2.Metadata.Path = templatePath

	return &template, nil
}
