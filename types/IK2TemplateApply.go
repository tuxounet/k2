package types

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
	} `yaml:"scripts"`
}
