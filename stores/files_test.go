package stores

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tuxounet/k2/types"
)

func TestNewFileStore(t *testing.T) {
	tmpDir := t.TempDir()
	fs := NewFileStore(tmpDir)
	if fs.Dir != tmpDir {
		t.Fatalf("expected dir '%s', got '%s'", tmpDir, fs.Dir)
	}
	if len(fs.K2) != 0 {
		t.Fatalf("expected no K2 items, got %d", len(fs.K2))
	}
}

func TestNewFileStore_RelativePath(t *testing.T) {
	fs := NewFileStore("relative")
	if !filepath.IsAbs(fs.Dir) {
		t.Fatalf("expected absolute dir, got '%s'", fs.Dir)
	}
}

func TestFileStore_GetAsInventory(t *testing.T) {
	tmpDir := t.TempDir()
	content := `k2:
  metadata:
    id: test-inventory
    kind: inventory
  body:
    folders:
      ignore: []
      templates:
        - templates/**/k2.*.yaml
      applies:
        - services/**/k2.apply.yaml
    vars:
      title: test
`
	err := os.WriteFile(filepath.Join(tmpDir, "k2.inventory.yaml"), []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	fs := NewFileStore(tmpDir)
	inv, err := fs.GetAsInventory("k2.inventory.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inv.K2.Metadata.ID != "test-inventory" {
		t.Fatalf("expected 'test-inventory', got '%s'", inv.K2.Metadata.ID)
	}
	if inv.K2.Metadata.Kind != "inventory" {
		t.Fatalf("expected 'inventory', got '%s'", inv.K2.Metadata.Kind)
	}
	if len(inv.K2.Body.Folders.Templates) != 1 {
		t.Fatalf("expected 1 template pattern, got %d", len(inv.K2.Body.Folders.Templates))
	}
	if len(inv.K2.Body.Folders.Applies) != 1 {
		t.Fatalf("expected 1 apply pattern, got %d", len(inv.K2.Body.Folders.Applies))
	}
	if inv.K2.Body.Vars["title"] != "test" {
		t.Fatalf("expected title='test', got '%s'", inv.K2.Body.Vars["title"])
	}
	if inv.K2.Metadata.Path != filepath.Join(tmpDir, "k2.inventory.yaml") {
		t.Fatalf("expected path '%s', got '%s'", filepath.Join(tmpDir, "k2.inventory.yaml"), inv.K2.Metadata.Path)
	}
	if inv.K2.Metadata.Folder != tmpDir {
		t.Fatalf("expected folder '%s', got '%s'", tmpDir, inv.K2.Metadata.Folder)
	}
}

func TestFileStore_GetAsInventory_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	fs := NewFileStore(tmpDir)
	_, err := fs.GetAsInventory("nonexistent.yaml")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestFileStore_GetAsTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	content := `k2:
  metadata:
    id: test-template
    kind: template
  body:
    name: kind1
    parameters:
      name: kind1
      description: Template of type kind1
    scripts:
      bootstrap:
        - echo "template boot of {{ .name }}"
      post:
        - echo "template fin of {{ .name }}"
`
	err := os.WriteFile(filepath.Join(tmpDir, "k2.template.yaml"), []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	fs := NewFileStore(tmpDir)
	tpl, err := fs.GetAsTemplate("k2.template.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tpl.K2.Metadata.ID != "test-template" {
		t.Fatalf("expected 'test-template', got '%s'", tpl.K2.Metadata.ID)
	}
	if tpl.K2.Body.Name != "kind1" {
		t.Fatalf("expected 'kind1', got '%s'", tpl.K2.Body.Name)
	}
	if tpl.K2.Body.Parameters["name"] != "kind1" {
		t.Fatalf("expected name='kind1', got '%v'", tpl.K2.Body.Parameters["name"])
	}
	if len(tpl.K2.Body.Scripts.Bootstrap) != 1 {
		t.Fatalf("expected 1 bootstrap script, got %d", len(tpl.K2.Body.Scripts.Bootstrap))
	}
	if len(tpl.K2.Body.Scripts.Post) != 1 {
		t.Fatalf("expected 1 post script, got %d", len(tpl.K2.Body.Scripts.Post))
	}
}

func TestFileStore_GetAsTemplateApply(t *testing.T) {
	tmpDir := t.TempDir()
	content := `k2:
  metadata:
    id: test-apply
    kind: template-apply
  body:
    template:
      source: inventory
      params:
        id: test-template
    vars:
      name: mycomp
      description: Component description
    scripts:
      bootstrap:
        - echo "boot {{ .name }}"
      post:
        - echo "fin"
`
	err := os.WriteFile(filepath.Join(tmpDir, "k2.apply.yaml"), []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	fs := NewFileStore(tmpDir)
	apply, err := fs.GetAsTemplateApply("k2.apply.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if apply.K2.Metadata.ID != "test-apply" {
		t.Fatalf("expected 'test-apply', got '%s'", apply.K2.Metadata.ID)
	}
	if apply.K2.Body.Template.Source != types.K2TemplateRefSourceInventory {
		t.Fatalf("expected 'inventory', got '%s'", apply.K2.Body.Template.Source)
	}
	if apply.K2.Body.Template.Params["id"] != "test-template" {
		t.Fatalf("expected template id 'test-template', got '%s'", apply.K2.Body.Template.Params["id"])
	}
	if apply.K2.Body.Vars["name"] != "mycomp" {
		t.Fatalf("expected name='mycomp', got '%v'", apply.K2.Body.Vars["name"])
	}
	if len(apply.K2.Body.Scripts.Bootstrap) != 1 {
		t.Fatalf("expected 1 bootstrap script, got %d", len(apply.K2.Body.Scripts.Bootstrap))
	}
}

func TestFileStore_GetAsTemplateApply_GitSource(t *testing.T) {
	tmpDir := t.TempDir()
	content := `k2:
  metadata:
    id: test-git-apply
    kind: template-apply
  body:
    template:
      source: git
      params:
        repository: https://github.com/tuxounet/k2.git
        branch: main
        path: samples/templates/fromGit1/k2.template.yaml
    vars:
      name: component
`
	err := os.WriteFile(filepath.Join(tmpDir, "k2.apply.yaml"), []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	fs := NewFileStore(tmpDir)
	apply, err := fs.GetAsTemplateApply("k2.apply.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if apply.K2.Body.Template.Source != types.K2TemplateRefSourceGit {
		t.Fatalf("expected 'git', got '%s'", apply.K2.Body.Template.Source)
	}
	if apply.K2.Body.Template.Params["repository"] != "https://github.com/tuxounet/k2.git" {
		t.Fatalf("unexpected repository: %s", apply.K2.Body.Template.Params["repository"])
	}
}

func TestFileStore_GetKey(t *testing.T) {
	tmpDir := t.TempDir()
	content := `k2:
  metadata:
    id: generic-item
    kind: something
  body:
    field: value
`
	err := os.WriteFile(filepath.Join(tmpDir, "test.yaml"), []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	fs := NewFileStore(tmpDir)
	item, err := fs.GetKey("test.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.K2.Metadata.ID != "generic-item" {
		t.Fatalf("expected 'generic-item', got '%s'", item.K2.Metadata.ID)
	}
}

func TestFileStore_GetKey_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	err := os.WriteFile(filepath.Join(tmpDir, "bad.yaml"), []byte("not: [valid: yaml: {{"), 0644)
	if err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	fs := NewFileStore(tmpDir)
	_, err = fs.GetKey("bad.yaml")
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestFileStore_Scan(t *testing.T) {
	tmpDir := t.TempDir()

	// Create directory structure
	dirs := []string{
		"services/product1/component1",
		"services/product1/component2",
		"templates/kind1",
	}
	for _, d := range dirs {
		err := os.MkdirAll(filepath.Join(tmpDir, d), 0755)
		if err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
	}

	// Create files
	files := map[string]string{
		"services/product1/component1/k2.apply.yaml": "k2:\n  metadata:\n    id: c1\n",
		"services/product1/component2/k2.apply.yaml": "k2:\n  metadata:\n    id: c2\n",
		"templates/kind1/k2.template.yaml":           "k2:\n  metadata:\n    id: t1\n",
		"other.txt":                                  "not a yaml",
	}
	for name, content := range files {
		err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644)
		if err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
	}

	fs := NewFileStore(tmpDir)

	// Test scanning for applies
	results, err := fs.Scan([]string{"services/**/k2.apply.yaml"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 applies, got %d: %v", len(results), results)
	}

	// Test scanning for templates
	results, err = fs.Scan([]string{"templates/**/k2.*.yaml"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 template, got %d: %v", len(results), results)
	}
}

func TestFileStore_Scan_NoMatch(t *testing.T) {
	tmpDir := t.TempDir()
	err := os.WriteFile(filepath.Join(tmpDir, "other.txt"), []byte("text"), 0644)
	if err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	fs := NewFileStore(tmpDir)
	results, err := fs.Scan([]string{"**/*.yaml"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestFileStore_Scan_MultiplePatterns(t *testing.T) {
	tmpDir := t.TempDir()

	os.MkdirAll(filepath.Join(tmpDir, "a"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "a", "k2.apply.yaml"), []byte("data"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "a", "k2.apply.yml"), []byte("data"), 0644)

	fs := NewFileStore(tmpDir)
	results, err := fs.Scan([]string{"a/k2.apply.yaml", "a/k2.apply.yml"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d: %v", len(results), results)
	}
}

func TestFileStore_Scan_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	fs := NewFileStore(tmpDir)
	results, err := fs.Scan([]string{"**/*.yaml"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}
