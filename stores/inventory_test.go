package stores

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// copySamplesDir creates a temporary copy of the samples directory for testing
func copySamplesDir(t *testing.T) string {
	t.Helper()

	// Find the project root (where samples/ lives)
	projectRoot := findProjectRoot(t)
	samplesDir := filepath.Join(projectRoot, "samples")

	tmpDir := t.TempDir()
	err := copyDirRecursive(samplesDir, tmpDir)
	if err != nil {
		t.Fatalf("failed to copy samples dir: %v", err)
	}
	return tmpDir
}

func findProjectRoot(t *testing.T) string {
	t.Helper()
	// Walk up from current working directory to find go.mod
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find project root")
		}
		dir = parent
	}
}

func copyDirRecursive(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(destPath, data, info.Mode())
	})
}

func TestInventory_NewInventory_FromSamples(t *testing.T) {
	samplesDir := copySamplesDir(t)
	inventoryPath := filepath.Join(samplesDir, "k2.inventory.yaml")

	inv, err := NewInventory(inventoryPath)
	if err != nil {
		t.Fatalf("failed to create inventory: %v", err)
	}
	if inv == nil {
		t.Fatal("expected non-nil inventory")
	}
	if inv.InventoryDir != samplesDir {
		t.Fatalf("expected inventory dir '%s', got '%s'", samplesDir, inv.InventoryDir)
	}
	if inv.InventoryKey != "k2.inventory.yaml" {
		t.Fatalf("expected key 'k2.inventory.yaml', got '%s'", inv.InventoryKey)
	}
	if inv.inventoryDefinition.K2.Metadata.ID != "k2.cli.sample.inventory" {
		t.Fatalf("expected id 'k2.cli.sample.inventory', got '%s'", inv.inventoryDefinition.K2.Metadata.ID)
	}
}

func TestInventory_NewInventory_InvalidPath(t *testing.T) {
	_, err := NewInventory("/nonexistent/path/k2.inventory.yaml")
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestInventory_Plan_FromSamples(t *testing.T) {
	samplesDir := copySamplesDir(t)
	inventoryPath := filepath.Join(samplesDir, "k2.inventory.yaml")

	inv, err := NewInventory(inventoryPath)
	if err != nil {
		t.Fatalf("failed to create inventory: %v", err)
	}

	plan, err := inv.Plan()
	if err != nil {
		t.Fatalf("failed to plan: %v", err)
	}

	if plan == nil {
		t.Fatal("expected non-nil plan")
	}

	// The samples have multiple applies and templates
	if len(plan.Tasks) == 0 {
		t.Fatal("expected tasks in plan")
	}
	if len(plan.Entities) == 0 {
		t.Fatal("expected entities in plan")
	}
	if len(plan.Refs) == 0 {
		t.Fatal("expected refs in plan")
	}

	// Check task types present
	hasLocalResolve := false
	hasGitResolve := false
	hasApply := false
	for _, task := range plan.Tasks {
		switch task.Type {
		case ActionTaskTypeLocalResolve:
			hasLocalResolve = true
		case ActionTaskTypeGitResolve:
			hasGitResolve = true
		case ActionTaskTypeApply:
			hasApply = true
		}
	}

	if !hasLocalResolve {
		t.Fatal("expected local-resolve tasks in plan")
	}
	if !hasGitResolve {
		t.Fatal("expected git-resolve tasks in plan")
	}
	if !hasApply {
		t.Fatal("expected apply tasks in plan")
	}
}

func TestInventory_Plan_InventoryEntities(t *testing.T) {
	samplesDir := copySamplesDir(t)
	inventoryPath := filepath.Join(samplesDir, "k2.inventory.yaml")

	inv, err := NewInventory(inventoryPath)
	if err != nil {
		t.Fatalf("failed to create inventory: %v", err)
	}

	plan, err := inv.Plan()
	if err != nil {
		t.Fatalf("failed to plan: %v", err)
	}

	// Check that expected entity IDs are present
	expectedEntities := []string{
		"k2.cli.sample.services.product1.component1",
		"k2.cli.sample.services.product1.component2",
		"k2.cli.sample.services.product1.component2bis",
		"k2.cli.sample.services.product2.componentFromGit1",
		"k2.cli.sample.services.product2.withCustomFiles",
		"k2.cli.sample.templates.kind1",
		"k2.cli.sample.templates.kind2",
		"k2.cli.sample.templates.with-placeholder",
	}

	for _, expectedID := range expectedEntities {
		if _, ok := plan.Entities[expectedID]; !ok {
			t.Fatalf("expected entity '%s' not found in plan", expectedID)
		}
	}
}

func TestInventory_Plan_TasksAreDeduped(t *testing.T) {
	samplesDir := copySamplesDir(t)
	inventoryPath := filepath.Join(samplesDir, "k2.inventory.yaml")

	inv, err := NewInventory(inventoryPath)
	if err != nil {
		t.Fatalf("failed to create inventory: %v", err)
	}

	plan, err := inv.Plan()
	if err != nil {
		t.Fatalf("failed to plan: %v", err)
	}

	// Verify no duplicate tasks
	seen := make(map[string]bool)
	for _, task := range plan.Tasks {
		key := strings.Join([]string{string(task.Type), taskParamsKey(task.Params)}, "-")
		if seen[key] {
			t.Fatalf("duplicate task found: %s", key)
		}
		seen[key] = true
	}
}

func taskParamsKey(params map[string]interface{}) string {
	parts := make([]string, 0)
	for k, v := range params {
		parts = append(parts, k+"="+strings.TrimSpace(strings.ReplaceAll(v.(string), "\n", "")))
	}
	return strings.Join(parts, ",")
}

func TestInventory_Plan_LocalOnlyInventory(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a minimal inventory with only local templates
	invContent := `k2:
  metadata:
    id: test-local
    kind: inventory
  body:
    folders:
      ignore: []
      templates:
        - templates/**/k2.template.yaml
      applies:
        - services/**/k2.apply.yaml
    vars:
      title: test
`
	os.WriteFile(filepath.Join(tmpDir, "k2.inventory.yaml"), []byte(invContent), 0644)

	// Template
	tplDir := filepath.Join(tmpDir, "templates", "simple")
	os.MkdirAll(tplDir, 0755)
	tplContent := `k2:
  metadata:
    id: simple-tpl
    kind: template
  body:
    name: simple
    parameters:
      name: simple
`
	os.WriteFile(filepath.Join(tplDir, "k2.template.yaml"), []byte(tplContent), 0644)
	os.WriteFile(filepath.Join(tplDir, "README.md"), []byte("# {{ .name }}"), 0644)

	// Apply
	applyDir := filepath.Join(tmpDir, "services", "comp1")
	os.MkdirAll(applyDir, 0755)
	applyContent := `k2:
  metadata:
    id: comp1-apply
    kind: template-apply
  body:
    template:
      source: inventory
      params:
        id: simple-tpl
    vars:
      name: mycomp
`
	os.WriteFile(filepath.Join(applyDir, "k2.apply.yaml"), []byte(applyContent), 0644)

	inv, err := NewInventory(filepath.Join(tmpDir, "k2.inventory.yaml"))
	if err != nil {
		t.Fatalf("failed to create inventory: %v", err)
	}

	plan, err := inv.Plan()
	if err != nil {
		t.Fatalf("failed to plan: %v", err)
	}

	if len(plan.Tasks) != 2 {
		t.Fatalf("expected 2 tasks (1 local-resolve + 1 apply), got %d", len(plan.Tasks))
	}

	if plan.Tasks[0].Type != ActionTaskTypeLocalResolve {
		t.Fatalf("expected first task to be local-resolve, got %s", plan.Tasks[0].Type)
	}
	if plan.Tasks[1].Type != ActionTaskTypeApply {
		t.Fatalf("expected second task to be apply, got %s", plan.Tasks[1].Type)
	}
}

func TestInventory_ApplyAndDestroy_LocalOnly(t *testing.T) {
	tmpDir := t.TempDir()

	// Inventory
	invContent := `k2:
  metadata:
    id: test-ad
    kind: inventory
  body:
    folders:
      ignore: []
      templates:
        - templates/**/k2.template.yaml
      applies:
        - services/**/k2.apply.yaml
    vars:
      title: test
`
	os.WriteFile(filepath.Join(tmpDir, "k2.inventory.yaml"), []byte(invContent), 0644)

	// Template
	tplDir := filepath.Join(tmpDir, "templates", "simple")
	os.MkdirAll(tplDir, 0755)
	os.WriteFile(filepath.Join(tplDir, "k2.template.yaml"), []byte(`k2:
  metadata:
    id: simple-tpl
    kind: template
  body:
    name: simple
    parameters:
      name: simple
      description: A simple template
`), 0644)
	os.WriteFile(filepath.Join(tplDir, "README.md"), []byte("# {{ .name }}\n\n{{ .description }}"), 0644)

	// Apply
	applyDir := filepath.Join(tmpDir, "services", "comp1")
	os.MkdirAll(applyDir, 0755)
	os.WriteFile(filepath.Join(applyDir, "k2.apply.yaml"), []byte(`k2:
  metadata:
    id: comp1-apply
    kind: template-apply
  body:
    template:
      source: inventory
      params:
        id: simple-tpl
    vars:
      name: mycomponent
      description: My great component
`), 0644)

	inv, err := NewInventory(filepath.Join(tmpDir, "k2.inventory.yaml"))
	if err != nil {
		t.Fatalf("failed to create inventory: %v", err)
	}

	plan, err := inv.Plan()
	if err != nil {
		t.Fatalf("failed to plan: %v", err)
	}

	// Apply
	err = inv.Apply(plan)
	if err != nil {
		t.Fatalf("failed to apply: %v", err)
	}

	// Check README was generated
	readmePath := filepath.Join(applyDir, "README.md")
	content, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("README.md not generated: %v", err)
	}
	if !strings.Contains(string(content), "mycomponent") {
		t.Fatalf("expected README to contain 'mycomponent', got: %s", string(content))
	}
	if !strings.Contains(string(content), "My great component") {
		t.Fatalf("expected README to contain 'My great component', got: %s", string(content))
	}

	// Check .gitignore was generated
	gitignore, err := os.ReadFile(filepath.Join(applyDir, ".gitignore"))
	if err != nil {
		t.Fatalf(".gitignore not generated: %v", err)
	}
	if !strings.Contains(string(gitignore), "README.md") {
		t.Fatalf("expected .gitignore to reference README.md, got: %s", string(gitignore))
	}

	// Now plan again for destroy
	plan2, err := inv.Plan()
	if err != nil {
		t.Fatalf("failed to plan for destroy: %v", err)
	}

	// Destroy
	err = inv.Destroy(plan2)
	if err != nil {
		t.Fatalf("failed to destroy: %v", err)
	}

	// README should be removed
	if _, err := os.Stat(readmePath); !os.IsNotExist(err) {
		t.Fatal("expected README.md to be removed after destroy")
	}
}

func TestInventory_ApplyAndDestroy_MultipleComponents(t *testing.T) {
	tmpDir := t.TempDir()

	// Inventory
	os.WriteFile(filepath.Join(tmpDir, "k2.inventory.yaml"), []byte(`k2:
  metadata:
    id: test-multi
    kind: inventory
  body:
    folders:
      ignore: []
      templates:
        - templates/**/k2.template.yaml
      applies:
        - services/**/k2.apply.yaml
    vars:
      title: multi test
`), 0644)

	// Template kind1
	tplDir1 := filepath.Join(tmpDir, "templates", "kind1")
	os.MkdirAll(tplDir1, 0755)
	os.WriteFile(filepath.Join(tplDir1, "k2.template.yaml"), []byte(`k2:
  metadata:
    id: tpl-kind1
    kind: template
  body:
    name: kind1
    parameters:
      name: kind1
`), 0644)
	os.WriteFile(filepath.Join(tplDir1, "README.md"), []byte("# {{ .name }}"), 0644)

	// Template kind2
	tplDir2 := filepath.Join(tmpDir, "templates", "kind2")
	os.MkdirAll(tplDir2, 0755)
	os.WriteFile(filepath.Join(tplDir2, "k2.template.yaml"), []byte(`k2:
  metadata:
    id: tpl-kind2
    kind: template
  body:
    name: kind2
    parameters:
      name: kind2
`), 0644)
	os.WriteFile(filepath.Join(tplDir2, "README.md"), []byte("## {{ .name }}"), 0644)

	// Apply 1 -> kind1
	applyDir1 := filepath.Join(tmpDir, "services", "svc1")
	os.MkdirAll(applyDir1, 0755)
	os.WriteFile(filepath.Join(applyDir1, "k2.apply.yaml"), []byte(`k2:
  metadata:
    id: svc1-apply
    kind: template-apply
  body:
    template:
      source: inventory
      params:
        id: tpl-kind1
    vars:
      name: service1
`), 0644)

	// Apply 2 -> kind2
	applyDir2 := filepath.Join(tmpDir, "services", "svc2")
	os.MkdirAll(applyDir2, 0755)
	os.WriteFile(filepath.Join(applyDir2, "k2.apply.yaml"), []byte(`k2:
  metadata:
    id: svc2-apply
    kind: template-apply
  body:
    template:
      source: inventory
      params:
        id: tpl-kind2
    vars:
      name: service2
`), 0644)

	// Apply 3 -> kind1 (same template, different apply)
	applyDir3 := filepath.Join(tmpDir, "services", "svc3")
	os.MkdirAll(applyDir3, 0755)
	os.WriteFile(filepath.Join(applyDir3, "k2.apply.yaml"), []byte(`k2:
  metadata:
    id: svc3-apply
    kind: template-apply
  body:
    template:
      source: inventory
      params:
        id: tpl-kind1
    vars:
      name: service3
`), 0644)

	inv, err := NewInventory(filepath.Join(tmpDir, "k2.inventory.yaml"))
	if err != nil {
		t.Fatalf("failed to create inventory: %v", err)
	}

	plan, err := inv.Plan()
	if err != nil {
		t.Fatalf("failed to plan: %v", err)
	}

	err = inv.Apply(plan)
	if err != nil {
		t.Fatalf("failed to apply: %v", err)
	}

	// Check service1 got kind1 template
	c1, _ := os.ReadFile(filepath.Join(applyDir1, "README.md"))
	if !strings.Contains(string(c1), "service1") {
		t.Fatalf("svc1 README: expected 'service1', got: %s", string(c1))
	}

	// Check service2 got kind2 template
	c2, _ := os.ReadFile(filepath.Join(applyDir2, "README.md"))
	if !strings.Contains(string(c2), "service2") {
		t.Fatalf("svc2 README: expected 'service2', got: %s", string(c2))
	}

	// Check service3 got kind1 template
	c3, _ := os.ReadFile(filepath.Join(applyDir3, "README.md"))
	if !strings.Contains(string(c3), "service3") {
		t.Fatalf("svc3 README: expected 'service3', got: %s", string(c3))
	}

	// Destroy
	plan2, _ := inv.Plan()
	err = inv.Destroy(plan2)
	if err != nil {
		t.Fatalf("failed to destroy: %v", err)
	}

	// All READMEs should be gone
	for _, dir := range []string{applyDir1, applyDir2, applyDir3} {
		if _, err := os.Stat(filepath.Join(dir, "README.md")); !os.IsNotExist(err) {
			t.Fatalf("expected README.md to be removed in %s", dir)
		}
	}
}

func TestInventory_Apply_WithScripts(t *testing.T) {
	tmpDir := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir, "k2.inventory.yaml"), []byte(`k2:
  metadata:
    id: test-scripts
    kind: inventory
  body:
    folders:
      ignore: []
      templates:
        - templates/**/k2.template.yaml
      applies:
        - services/**/k2.apply.yaml
    vars: {}
`), 0644)

	tplDir := filepath.Join(tmpDir, "templates", "scripted")
	os.MkdirAll(tplDir, 0755)
	os.WriteFile(filepath.Join(tplDir, "k2.template.yaml"), []byte(`k2:
  metadata:
    id: scripted-tpl
    kind: template
  body:
    name: scripted
    parameters:
      name: scripted
    scripts:
      bootstrap:
        - echo "template bootstrap {{ .name }}"
      pre:
        - echo "template pre {{ .name }}"
      post:
        - echo "template post {{ .name }}"
`), 0644)
	os.WriteFile(filepath.Join(tplDir, "README.md"), []byte("# {{ .name }}"), 0644)

	applyDir := filepath.Join(tmpDir, "services", "scripted-svc")
	os.MkdirAll(applyDir, 0755)
	os.WriteFile(filepath.Join(applyDir, "k2.apply.yaml"), []byte(`k2:
  metadata:
    id: scripted-apply
    kind: template-apply
  body:
    template:
      source: inventory
      params:
        id: scripted-tpl
    vars:
      name: myservice
    scripts:
      bootstrap:
        - echo "apply bootstrap"
      pre:
        - echo "apply pre"
      post:
        - echo "apply post"
`), 0644)

	inv, err := NewInventory(filepath.Join(tmpDir, "k2.inventory.yaml"))
	if err != nil {
		t.Fatalf("failed to create inventory: %v", err)
	}

	plan, err := inv.Plan()
	if err != nil {
		t.Fatalf("failed to plan: %v", err)
	}

	err = inv.Apply(plan)
	if err != nil {
		t.Fatalf("failed to apply: %v", err)
	}

	// Verify files are generated
	content, err := os.ReadFile(filepath.Join(applyDir, "README.md"))
	if err != nil {
		t.Fatalf("README.md not generated: %v", err)
	}
	if !strings.Contains(string(content), "myservice") {
		t.Fatalf("expected 'myservice', got: %s", string(content))
	}
}

func TestInventory_Apply_WithSubdirectoryFiles(t *testing.T) {
	tmpDir := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir, "k2.inventory.yaml"), []byte(`k2:
  metadata:
    id: test-subdir
    kind: inventory
  body:
    folders:
      ignore: []
      templates:
        - templates/**/k2.template.yaml
      applies:
        - services/**/k2.apply.yaml
    vars: {}
`), 0644)

	tplDir := filepath.Join(tmpDir, "templates", "withsub")
	subDir := filepath.Join(tplDir, "sub")
	os.MkdirAll(subDir, 0755)
	os.WriteFile(filepath.Join(tplDir, "k2.template.yaml"), []byte(`k2:
  metadata:
    id: sub-tpl
    kind: template
  body:
    name: withsub
    parameters:
      name: withsub
`), 0644)
	os.WriteFile(filepath.Join(tplDir, "README.md"), []byte("# {{ .name }}"), 0644)
	os.WriteFile(filepath.Join(subDir, "subfile.txt"), []byte("sub: {{ .name }}"), 0644)

	applyDir := filepath.Join(tmpDir, "services", "svc-sub")
	os.MkdirAll(applyDir, 0755)
	os.WriteFile(filepath.Join(applyDir, "k2.apply.yaml"), []byte(`k2:
  metadata:
    id: sub-apply
    kind: template-apply
  body:
    template:
      source: inventory
      params:
        id: sub-tpl
    vars:
      name: myservice
`), 0644)

	inv, err := NewInventory(filepath.Join(tmpDir, "k2.inventory.yaml"))
	if err != nil {
		t.Fatalf("failed to create inventory: %v", err)
	}

	plan, err := inv.Plan()
	if err != nil {
		t.Fatalf("failed to plan: %v", err)
	}

	err = inv.Apply(plan)
	if err != nil {
		t.Fatalf("failed to apply: %v", err)
	}

	// Check main file
	content, _ := os.ReadFile(filepath.Join(applyDir, "README.md"))
	if !strings.Contains(string(content), "myservice") {
		t.Fatalf("expected 'myservice' in README, got: %s", string(content))
	}

	// Check sub file
	subContent, err := os.ReadFile(filepath.Join(applyDir, "sub", "subfile.txt"))
	if err != nil {
		t.Fatalf("subfile.txt not generated: %v", err)
	}
	if !strings.Contains(string(subContent), "myservice") {
		t.Fatalf("expected 'myservice' in subfile, got: %s", string(subContent))
	}
}

func TestInventory_Plan_UnknownTemplateRef(t *testing.T) {
	tmpDir := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir, "k2.inventory.yaml"), []byte(`k2:
  metadata:
    id: test-unknown
    kind: inventory
  body:
    folders:
      ignore: []
      templates:
        - templates/**/k2.template.yaml
      applies:
        - services/**/k2.apply.yaml
    vars: {}
`), 0644)

	// No template defined
	applyDir := filepath.Join(tmpDir, "services", "svc1")
	os.MkdirAll(applyDir, 0755)
	os.WriteFile(filepath.Join(applyDir, "k2.apply.yaml"), []byte(`k2:
  metadata:
    id: bad-apply
    kind: template-apply
  body:
    template:
      source: inventory
      params:
        id: nonexistent-template
    vars:
      name: test
`), 0644)

	inv, err := NewInventory(filepath.Join(tmpDir, "k2.inventory.yaml"))
	if err != nil {
		t.Fatalf("failed to create inventory: %v", err)
	}

	_, err = inv.Plan()
	if err == nil {
		t.Fatal("expected error for unknown template reference")
	}
}

func TestInventory_Plan_NoApplies(t *testing.T) {
	tmpDir := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir, "k2.inventory.yaml"), []byte(`k2:
  metadata:
    id: test-empty
    kind: inventory
  body:
    folders:
      ignore: []
      templates:
        - templates/**/k2.template.yaml
      applies:
        - services/**/k2.apply.yaml
    vars: {}
`), 0644)

	inv, err := NewInventory(filepath.Join(tmpDir, "k2.inventory.yaml"))
	if err != nil {
		t.Fatalf("failed to create inventory: %v", err)
	}

	plan, err := inv.Plan()
	if err != nil {
		t.Fatalf("failed to plan: %v", err)
	}
	if len(plan.Tasks) != 0 {
		t.Fatalf("expected 0 tasks for empty inventory, got %d", len(plan.Tasks))
	}
}

func TestInventory_Apply_TemplateParametersOverriddenByVars(t *testing.T) {
	tmpDir := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir, "k2.inventory.yaml"), []byte(`k2:
  metadata:
    id: test-override
    kind: inventory
  body:
    folders:
      ignore: []
      templates:
        - templates/**/k2.template.yaml
      applies:
        - services/**/k2.apply.yaml
    vars: {}
`), 0644)

	tplDir := filepath.Join(tmpDir, "templates", "tpl1")
	os.MkdirAll(tplDir, 0755)
	os.WriteFile(filepath.Join(tplDir, "k2.template.yaml"), []byte(`k2:
  metadata:
    id: tpl-override
    kind: template
  body:
    name: tpl1
    parameters:
      name: default-name
      description: default-desc
`), 0644)
	os.WriteFile(filepath.Join(tplDir, "README.md"), []byte("{{ .name }} - {{ .description }}"), 0644)

	applyDir := filepath.Join(tmpDir, "services", "svc-override")
	os.MkdirAll(applyDir, 0755)
	os.WriteFile(filepath.Join(applyDir, "k2.apply.yaml"), []byte(`k2:
  metadata:
    id: apply-override
    kind: template-apply
  body:
    template:
      source: inventory
      params:
        id: tpl-override
    vars:
      name: custom-name
      description: custom-desc
`), 0644)

	inv, err := NewInventory(filepath.Join(tmpDir, "k2.inventory.yaml"))
	if err != nil {
		t.Fatalf("failed to create inventory: %v", err)
	}

	plan, err := inv.Plan()
	if err != nil {
		t.Fatalf("failed to plan: %v", err)
	}

	err = inv.Apply(plan)
	if err != nil {
		t.Fatalf("failed to apply: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(applyDir, "README.md"))
	if err != nil {
		t.Fatalf("README.md not generated: %v", err)
	}
	// Vars from apply should override template parameters (MergeMaps behavior)
	if !strings.Contains(string(content), "custom-name") {
		t.Fatalf("expected 'custom-name' (vars override params), got: %s", string(content))
	}
	if !strings.Contains(string(content), "custom-desc") {
		t.Fatalf("expected 'custom-desc' (vars override params), got: %s", string(content))
	}
}

func TestInventory_Plan_UnknownTemplateSource(t *testing.T) {
	tmpDir := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir, "k2.inventory.yaml"), []byte(`k2:
  metadata:
    id: test-unknown-source
    kind: inventory
  body:
    folders:
      ignore: []
      templates:
        - templates/**/k2.template.yaml
      applies:
        - services/**/k2.apply.yaml
    vars: {}
`), 0644)

	applyDir := filepath.Join(tmpDir, "services", "svc1")
	os.MkdirAll(applyDir, 0755)
	os.WriteFile(filepath.Join(applyDir, "k2.apply.yaml"), []byte(`k2:
  metadata:
    id: bad-source
    kind: template-apply
  body:
    template:
      source: unknown-source
      params:
        id: whatever
    vars:
      name: test
`), 0644)

	inv, err := NewInventory(filepath.Join(tmpDir, "k2.inventory.yaml"))
	if err != nil {
		t.Fatalf("failed to create inventory: %v", err)
	}

	_, err = inv.Plan()
	if err == nil {
		t.Fatal("expected error for unknown template source")
	}
	if !strings.Contains(err.Error(), "unknown template source") {
		t.Fatalf("expected 'unknown template source' error, got: %v", err)
	}
}

func TestInventory_Apply_IdempotentReapply(t *testing.T) {
	tmpDir := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir, "k2.inventory.yaml"), []byte(`k2:
  metadata:
    id: test-reapply
    kind: inventory
  body:
    folders:
      ignore: []
      templates:
        - templates/**/k2.template.yaml
      applies:
        - services/**/k2.apply.yaml
    vars: {}
`), 0644)

	tplDir := filepath.Join(tmpDir, "templates", "simple")
	os.MkdirAll(tplDir, 0755)
	os.WriteFile(filepath.Join(tplDir, "k2.template.yaml"), []byte(`k2:
  metadata:
    id: simple-tpl
    kind: template
  body:
    name: simple
    parameters:
      name: simple
`), 0644)
	os.WriteFile(filepath.Join(tplDir, "README.md"), []byte("# {{ .name }}"), 0644)

	applyDir := filepath.Join(tmpDir, "services", "svc1")
	os.MkdirAll(applyDir, 0755)
	os.WriteFile(filepath.Join(applyDir, "k2.apply.yaml"), []byte(`k2:
  metadata:
    id: svc1-apply
    kind: template-apply
  body:
    template:
      source: inventory
      params:
        id: simple-tpl
    vars:
      name: myservice
`), 0644)

	inv, err := NewInventory(filepath.Join(tmpDir, "k2.inventory.yaml"))
	if err != nil {
		t.Fatalf("failed to create inventory: %v", err)
	}

	// First apply
	plan1, _ := inv.Plan()
	err = inv.Apply(plan1)
	if err != nil {
		t.Fatalf("first apply failed: %v", err)
	}

	// Second apply (re-apply)
	plan2, _ := inv.Plan()
	err = inv.Apply(plan2)
	if err != nil {
		t.Fatalf("second apply (re-apply) failed: %v", err)
	}

	content, _ := os.ReadFile(filepath.Join(applyDir, "README.md"))
	if !strings.Contains(string(content), "myservice") {
		t.Fatalf("content should still be correct after re-apply, got: %s", string(content))
	}
}
