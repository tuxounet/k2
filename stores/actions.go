package stores

import (
	"fmt"
	"path/filepath"

	"github.com/tuxounet/k2/libs"
	"github.com/tuxounet/k2/types"
)

type ActionTaskType string

const (
	ActionTaskTypeLocalResolve ActionTaskType = "local-resolve"
	ActionTaskTypeGitResolve   ActionTaskType = "git-resolve"
	ActionTaskTypeApply        ActionTaskType = "apply"
)

type ActionTask struct {
	Type   ActionTaskType         `yaml:"type"`
	Params map[string]interface{} `yaml:"params"`
}

type ActionPlan struct {
	inventory *Inventory                    `yaml:"-"`
	Tasks     []ActionTask                  `yaml:"actions"`
	Entities  map[string]any                `yaml:"entities"`
	Refs      []types.IK2TemplateRef        `yaml:"refs"`
	Templates map[string]*types.IK2Template `yaml:"-"`
}

func NewActionPlan(inventory *Inventory) *ActionPlan {
	return &ActionPlan{
		inventory: inventory,
		Tasks:     make([]ActionTask, 0),
		Entities:  make(map[string]any, 0),
		Templates: make(map[string]*types.IK2Template, 0),
	}
}
func (ap *ActionPlan) AddEntity(meta types.IK2Metadata, body any) {
	ap.Entities[meta.ID] = body
}

func (ap *ActionPlan) AddTask(task ActionTask) {
	ap.Tasks = append(ap.Tasks, task)
}

func (ap *ActionPlan) Dedup() {

	uniqueActions := make(map[string]ActionTask, 0)
	deduped := make([]ActionTask, 0)

	for _, action := range ap.Tasks {
		actionKey := fmt.Sprintf("%s-%v", action.Type, action.Params)
		_, ok := uniqueActions[actionKey]
		if !ok {
			uniqueActions[actionKey] = action
			deduped = append(deduped, action)
		}

	}

	ap.Tasks = deduped

}

func (ap *ActionPlan) Apply() error {
	templateStore := NewTemplatingStore(ap)

	for _, task := range ap.Tasks {

		switch task.Type {
		case ActionTaskTypeLocalResolve:
			hash := task.Params["hash"].(string)
			tpl, err := templateStore.resolveTemplateInventory(hash)
			if err != nil {
				return err
			}
			ap.Templates[hash] = tpl
		case ActionTaskTypeGitResolve:
			hash := task.Params["hash"].(string)
			tpl, err := templateStore.resolveTemplateGit(hash)

			if err != nil {
				return err
			}
			ap.Templates[hash] = tpl
		case ActionTaskTypeApply:
			id := task.Params["id"].(string)
			ref := task.Params["ref"].(string)

			ok, err := templateStore.ApplyTemplate(id, ref, true)
			if err != nil {
				return err
			}
			if !ok {
				return libs.WriteErrorf("error applying template: %s\n", id)
			}

		default:
			return libs.WriteErrorf("unknown action type: %s\n", task.Type)
		}
	}

	return nil

}
func (ap *ActionPlan) Destroy() error {
	templateStore := NewTemplatingStore(ap)

	for _, task := range ap.Tasks {

		switch task.Type {
		case ActionTaskTypeGitResolve:
			hash := task.Params["hash"].(string)
			refsFolder := filepath.Join(ap.inventory.InventoryDir, libs.RefsDir)
			templateFolder := filepath.Join(refsFolder, hash)
			err := templateStore.destroyTemplateRef(templateFolder)
			if err != nil {
				return err
			}

		case ActionTaskTypeApply:
			id := task.Params["id"].(string)

			err := templateStore.DestroyTemplate(id)
			if err != nil {
				return err
			}

		}
	}

	err := templateStore.cleanupEmptyDirs(ap.inventory.InventoryDir)
	if err != nil {
		return err
	}

	return nil

}
func (ap *ActionPlan) GetEntityAsTemplate(id string) (*types.IK2Template, error) {

	ret, ok := ap.Entities[id]
	if !ok {
		return nil, libs.WriteErrorf("entity not found: %s\n", id)
	}

	return ret.(*types.IK2Template), nil

}

func (ap *ActionPlan) GetEntityAsTemplateApply(id string) (*types.IK2TemplateApply, error) {

	ret, ok := ap.Entities[id]
	if !ok {
		return nil, libs.WriteErrorf("entity not found: %s\n", id)
	}

	return ret.(*types.IK2TemplateApply), nil

}
