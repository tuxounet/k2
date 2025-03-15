package stores

import (
	"github.com/tuxounet/k2/libs"
	"github.com/tuxounet/k2/types"
)

func (t *TemplatingStore) resolveTemplateInventory(hash string) (*types.IK2Template, error) {

	libs.WriteOutputf("resolve template inventory %s\n", hash)

	refs := t.plan.Refs

	for _, ref := range refs {
		if ref.Hash() == hash {
			inventoryId := ref.Params["id"]
			return t.plan.GetEntityAsTemplate(inventoryId)
		}

	}

	return nil, libs.WriteErrorf("template not found: %s", hash)
}
