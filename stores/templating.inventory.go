package stores

import (
	"fmt"
	"k2/types"
)

func (t *TemplatingStore) resolveTemplateInventory(hash string) (*types.IK2Template, error) {

	fmt.Println("resolve template invenntory", hash)

	refs := t.plan.Refs

	for _, ref := range refs {
		if ref.Hash() == hash {
			inventoryId := ref.Params["id"]
			return t.plan.GetEntityAsTemplate(inventoryId)
		}

	}

	return nil, fmt.Errorf("template not found: %s", hash)
}
