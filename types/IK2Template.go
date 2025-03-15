package types

import (
	"strings"

	"github.com/tuxounet/k2/libs"
)

type IK2Template struct {
	K2 IK2TemplateRoot `yaml:"k2"`
}

type IK2TemplateRoot struct {
	Metadata IK2Metadata     `yaml:"metadata"`
	Body     IK2TemplateBody `yaml:"body"`
}

type IK2TemplateBody struct {
	Name       string            `yaml:"name"`
	Parameters map[string]string `yaml:"parameters"`
	Scripts    struct {
		Bootstrap []string `yaml:"bootstrap"`
		Pre       []string `yaml:"pre"`
		Post      []string `yaml:"post"`
		Nuke      []string `yaml:"nuke"`
	} `yaml:"scripts"`
}

func (t *IK2Template) ExecutePre(target *IK2TemplateApply) error {
	return t.executeScript(target, t.K2.Body.Scripts.Pre)

}

func (t *IK2Template) ExecutePost(target *IK2TemplateApply) error {
	return t.executeScript(target, t.K2.Body.Scripts.Post)

}

func (t *IK2Template) ExecuteNuke(target *IK2TemplateApply) error {
	return t.executeScript(target, t.K2.Body.Scripts.Nuke)

}

func (t *IK2Template) ExecuteBootstrap(target *IK2TemplateApply) error {
	return t.executeScript(target, t.K2.Body.Scripts.Bootstrap)

}

func (t *IK2Template) executeScript(target *IK2TemplateApply, script []string) error {
	if len(script) == 0 {
		return nil
	}

	libs.WriteOutputf("template execute script: %v\n", script)

	for _, line := range script {
		line := strings.TrimSpace(line)
		if line == "" {
			continue
		}

		err := libs.ExecCommand(line, target.K2.Metadata.Folder, libs.MergeMaps(t.K2.Body.Parameters, target.K2.Body.Vars))
		if err != nil {
			return libs.WriteErrorf("error executing script: %w", err)
		}
	}

	return nil
}
