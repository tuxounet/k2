package types

type IK2Stack struct {
	Version string       `yaml:"version"`
	Stack   IK2StackBody `yaml:"stack"`
}

type IK2StackBody struct {
	Description string            `yaml:"description"`
	Extends     string            `yaml:"extends,omitempty"`
	Env         map[string]string `yaml:"env"`
	Layers      []IK2StackLayer   `yaml:"layers"`
}

type IK2StackLayer struct {
	Layer string            `yaml:"layer"`
	Plan  string            `yaml:"plan"`
	Env   map[string]string `yaml:"env"`
}
