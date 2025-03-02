package types

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
	} `yaml:"scripts"`
}
