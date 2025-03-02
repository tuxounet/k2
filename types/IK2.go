package types

type IK2Body = interface{}

type IK2[TBody IK2Body] struct {
	K2 IK2Root[TBody] `yaml:"k2"`
}

type IK2Root[TBody IK2Body] struct {
	Metadata IK2Metadata `yaml:"metadata"`
	Body     TBody       `yaml:"body"`
}

type IK2Metadata struct {
	ID      string `yaml:"id"`
	Kind    string `yaml:"kind"`
	Version string `yaml:"version"`
	Path    string `yaml:"path"`
	Folder  string `yaml:"folder"`
}
