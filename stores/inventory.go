package stores

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tuxounet/k2/libs"
	"github.com/tuxounet/k2/types"
	"gopkg.in/yaml.v3"
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
		return nil, libs.WriteErrorf("error scanning applies: %w\n", err)
	}
	libs.WriteDetail("found %d applies", len(appliesFounds))

	if len(appliesFounds) == 0 {
		libs.WriteDetail("no applies found")
		return result, nil
	}

	templateRefs := make([]types.IK2TemplateRef, 0)
	applies := make([]*types.IK2TemplateApply, 0)
	for _, apply := range appliesFounds {
		templateApply, err := i.files.GetAsTemplateApply(apply)
		if err != nil {
			return nil, libs.WriteErrorf("error getting template apply: %w", err)
		}
		applies = append(applies, templateApply)
		templateRefs = append(templateRefs, templateApply.K2.Body.Template)

		result.AddEntity(templateApply.K2.Metadata, templateApply)

	}

	templatesFounds, err := i.files.Scan(i.inventoryDefinition.K2.Body.Folders.Templates)
	if err != nil {
		return nil, libs.WriteErrorf("error scanning templates: %w", err)
	}

	inventoryTemplates := make([]*types.IK2Template, 0)
	for _, template := range templatesFounds {
		template, err := i.files.GetAsTemplate(template)
		if err != nil {
			return nil, libs.WriteErrorf("error getting template: %w", err)
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
				return nil, libs.WriteErrorf("template not found: %s", templateRef.Params["id"])
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
			return nil, libs.WriteErrorf("unknown template source: %s", templateRef.Source)
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
		return libs.WriteErrorf("error applying plan: %w", err)
	}
	return nil

}

func (i *Inventory) Destroy(plan *ActionPlan) error {

	err := plan.Destroy()
	if err != nil {
		return libs.WriteErrorf("error destroying plan: %w", err)
	}
	return nil

}

func (i *Inventory) ListStacks() ([]StackInfo, error) {
	stacksFolder := i.inventoryDefinition.K2.Body.Folders.Stacks
	if stacksFolder == "" {
		return nil, fmt.Errorf("no stacks folder defined in inventory")
	}

	stacksDir := filepath.Join(i.InventoryDir, stacksFolder)
	entries, err := os.ReadDir(stacksDir)
	if err != nil {
		return nil, fmt.Errorf("cannot read stacks directory '%s': %w", stacksDir, err)
	}

	var stacks []StackInfo
	for _, e := range entries {
		if e.IsDir() || (!strings.HasSuffix(e.Name(), ".yaml") && !strings.HasSuffix(e.Name(), ".yml")) {
			continue
		}
		stackName := strings.TrimSuffix(strings.TrimSuffix(e.Name(), ".yaml"), ".yml")
		stackFile := filepath.Join(stacksDir, e.Name())

		description := ""
		data, err := os.ReadFile(stackFile)
		if err == nil {
			var def types.IK2Stack
			if yaml.Unmarshal(data, &def) == nil && def.Stack.Description != "" {
				description = strings.TrimSpace(def.Stack.Description)
			}
		}

		layerCount := 0
		if err == nil {
			var def types.IK2Stack
			if yaml.Unmarshal(data, &def) == nil {
				layerCount = len(def.Stack.Layers)
			}
		}

		stacks = append(stacks, StackInfo{
			Name:        stackName,
			Description: description,
			LayerCount:  layerCount,
		})
	}

	return stacks, nil
}

type StackInfo struct {
	Name        string
	Description string
	LayerCount  int
}
