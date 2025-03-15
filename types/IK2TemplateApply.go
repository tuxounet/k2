package types

import (
	"github.com/tuxounet/k2/libs"
)

type IK2TemplateApply struct {
	K2 IK2TemplateApplyRoot `yaml:"k2"`
}

type IK2TemplateApplyRoot struct {
	Metadata IK2Metadata          `yaml:"metadata"`
	Body     IK2TemplateApplyBody `yaml:"body"`
}

type IK2TemplateApplyBody struct {
	Template IK2TemplateRef    `yaml:"template"`
	Vars     map[string]string `yaml:"vars"`
	Scripts  struct {
		Bootstrap []string `yaml:"bootstrap"`
		Pre       []string `yaml:"pre"`
		Post      []string `yaml:"post"`
		Nuke      []string `yaml:"nuke"`
	} `yaml:"scripts"`
}

func (t *IK2TemplateApply) ExecutePre() error {
	return t.executeScript(t.K2.Body.Scripts.Pre)

}

func (t *IK2TemplateApply) ExecutePost() error {
	return t.executeScript(t.K2.Body.Scripts.Post)

}

func (t *IK2TemplateApply) ExecuteNuke() error {
	return t.executeScript(t.K2.Body.Scripts.Nuke)

}

func (t *IK2TemplateApply) ExecuteBootstrap() error {
	return t.executeScript(t.K2.Body.Scripts.Bootstrap)

}

func (t *IK2TemplateApply) executeScript(script []string) error {
	if len(script) == 0 {
		return nil
	}
	libs.WriteOutputf("template-apply execute script: %s %v\n", t.K2.Metadata.ID, script)

	for _, line := range script {
		err := libs.ExecCommand(line, t.K2.Metadata.Folder, t.K2.Body.Vars)
		if err != nil {
			return libs.WriteErrorf("error executing script: %s %w", t.K2.Metadata.ID, err)
		}
	}
	return nil
}
