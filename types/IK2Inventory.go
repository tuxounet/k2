package types

type IK2Inventory struct {
	K2 IK2InventoryRoot `yaml:"k2"`
}

type IK2InventoryRoot struct {
	Metadata IK2Metadata      `yaml:"metadata"`
	Body     IK2InventoryBody `yaml:"body"`
}

type IK2InventoryBody struct {
	Folders struct {
		Ignore    []string `yaml:"ignore"`
		Applies   []string `yaml:"applies"`
		Templates []string `yaml:"templates"`
	} `yaml:"folders"`
	Vars map[string]string `yaml:"vars"`
}
