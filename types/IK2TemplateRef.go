package types

import (
	"crypto/sha256"
	"fmt"
)

type IK2TemplateRefSource string

const (
	K2TemplateRefSourceInventory IK2TemplateRefSource = "inventory"
	K2TemplateRefSourceGit       IK2TemplateRefSource = "git"
)

type IK2TemplateRef struct {
	// The name of the template
	Source IK2TemplateRefSource `yaml:"source"`
	Params map[string]string    `yaml:"params"`
}

func (t *IK2TemplateRef) Hash() string {
	value := fmt.Sprintf("%s-%v", t.Source, t.Params)
	sha256 := sha256.New()
	sha256.Write([]byte(value))
	return fmt.Sprintf("%x", sha256.Sum(nil))

}
