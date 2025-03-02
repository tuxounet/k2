package stores

import (
	"fmt"
	"path/filepath"

	"github.com/tuxounet/k2/types"
)

type Inventory struct {
	InventoryDir        string
	InventoryKey        string
	files               *FileStore
	inventoryDefinition *types.IK2Inventory
}

func NewInventory(inventoryPath string) (*Inventory, error) {

	inventoryDir := filepath.Dir(inventoryPath)
	if inventoryDir == "" {
		inventoryDir = "."
	}
	fullInventoryDir, err := filepath.Abs(inventoryDir)
	if err != nil {

		return nil, err

	}
	inventoryKey := filepath.Base(inventoryPath)

	files := NewFileStore(fullInventoryDir)
	instance := &Inventory{
		InventoryDir: fullInventoryDir,
		InventoryKey: inventoryKey,
		files:        files,
	}

	inventoryDefinition, err := files.GetAsInventory(inventoryKey)
	if err != nil {
		return nil, err
	}
	instance.inventoryDefinition = inventoryDefinition

	return instance, nil
}

func (i *Inventory) Plan() (*ActionPlan, error) {

	result := NewActionPlan(i)
	appliesFounds, err := i.files.Scan(i.inventoryDefinition.K2.Body.Folders.Applies)
	if err != nil {
		return nil, fmt.Errorf("error scanning applies: %w", err)
	}
	fmt.Printf("Found Apply: %d\n", len(appliesFounds))

	if len(appliesFounds) == 0 {
		fmt.Printf("No applies found\n")
		return result, nil
	}

	templateRefs := make([]types.IK2TemplateRef, 0)
	applies := make([]*types.IK2TemplateApply, 0)
	for _, apply := range appliesFounds {
		templateApply, err := i.files.GetAsTemplateApply(apply)
		if err != nil {
			return nil, fmt.Errorf("error getting template apply: %w", err)
		}
		applies = append(applies, templateApply)
		templateRefs = append(templateRefs, templateApply.K2.Body.Template)

		result.AddEntity(templateApply.K2.Metadata, templateApply)

	}

	templatesFounds, err := i.files.Scan(i.inventoryDefinition.K2.Body.Folders.Templates)
	if err != nil {
		return nil, fmt.Errorf("error scanning templates: %w", err)
	}

	inventoryTemplates := make([]*types.IK2Template, 0)
	for _, template := range templatesFounds {
		template, err := i.files.GetAsTemplate(template)
		if err != nil {
			return nil, fmt.Errorf("error getting template: %w", err)
		}
		inventoryTemplates = append(inventoryTemplates, template)
		result.AddEntity(template.K2.Metadata, template)
	}

	result.Refs = templateRefs

	for _, templateRef := range templateRefs {
		switch templateRef.Source {
		case types.K2TemplateRefSourceInventory:
			found := false
			for _, template := range inventoryTemplates {
				if template.K2.Metadata.ID == templateRef.Params["id"] {
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("template not found: %s", templateRef.Params["id"])
			}
			result.AddTask(ActionTask{
				Type: ActionTaskTypeLocalResolve,
				Params: map[string]interface{}{
					"hash":   templateRef.Hash(),
					"params": fmt.Sprintf("%v", templateRef.Params),
				}})

		case types.K2TemplateRefSourceGit:

			result.AddTask(ActionTask{
				Type: ActionTaskTypeGitResolve,
				Params: map[string]interface{}{
					"hash":   templateRef.Hash(),
					"params": fmt.Sprintf("%v", templateRef.Params),
				},
			})

		default:
			return nil, fmt.Errorf("unknown template source: %s", templateRef.Source)
		}

	}

	for _, apply := range applies {

		templateRef := apply.K2.Body.Template.Hash()
		result.AddEntity(apply.K2.Metadata, apply)
		result.AddTask(ActionTask{
			Type: ActionTaskTypeApply,
			Params: map[string]interface{}{
				"id":  apply.K2.Metadata.ID,
				"ref": templateRef,
			},
		})

	}

	result.Dedup()

	return result, nil
}

func (i *Inventory) Apply(plan *ActionPlan) error {

	err := plan.Apply()
	if err != nil {
		return fmt.Errorf("error applying plan: %w", err)
	}
	return nil

}

func (i *Inventory) Destroy(plan *ActionPlan) error {

	err := plan.Destroy()
	if err != nil {
		return fmt.Errorf("error destroying plan: %w", err)
	}
	return nil

}
