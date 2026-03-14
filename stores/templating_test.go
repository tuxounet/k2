package stores

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tuxounet/k2/types"
)

func TestNewTemplatingStore(t *testing.T) {
	plan := NewActionPlan(nil)
	ts := NewTemplatingStore(plan)
	if ts == nil {
		t.Fatal("expected non-nil TemplatingStore")
	}
	if ts.plan != plan {
		t.Fatal("plan reference mismatch")
	}
}

func TestTemplatingStore_ResolveTemplateInventory(t *testing.T) {
	plan := NewActionPlan(nil)

	tpl := &types.IK2Template{
		K2: types.IK2TemplateRoot{
			Metadata: types.IK2Metadata{ID: "my-template"},
			Body: types.IK2TemplateBody{
				Name:       "kind1",
				Parameters: map[string]any{"name": "kind1"},
			},
		},
	}
	plan.AddEntity(types.IK2Metadata{ID: "my-template"}, tpl)

	ref := types.IK2TemplateRef{
		Source: types.K2TemplateRefSourceInventory,
		Params: map[string]string{"id": "my-template"},
	}
	plan.Refs = []types.IK2TemplateRef{ref}

	ts := NewTemplatingStore(plan)
	result, err := ts.resolveTemplateInventory(ref.Hash())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.K2.Metadata.ID != "my-template" {
		t.Fatalf("expected 'my-template', got '%s'", result.K2.Metadata.ID)
	}
}

func TestTemplatingStore_ResolveTemplateInventory_NotFound(t *testing.T) {
	plan := NewActionPlan(nil)
	plan.Refs = []types.IK2TemplateRef{}

	ts := NewTemplatingStore(plan)
	_, err := ts.resolveTemplateInventory("nonexistent-hash")
	if err == nil {
		t.Fatal("expected error for missing template")
	}
}

func TestTemplatingStore_ApplyTemplate(t *testing.T) {
	tmpDir := t.TempDir()

	// Create template folder with files
	tplDir := filepath.Join(tmpDir, "templates", "kind1")
	os.MkdirAll(tplDir, 0755)

	// Create template yaml
	os.WriteFile(filepath.Join(tplDir, "k2.template.yaml"), []byte(`k2:
  metadata:
    id: test-tpl
    kind: template
  body:
    name: kind1
    parameters:
      name: kind1
`), 0644)

	// Create a README template file
	os.WriteFile(filepath.Join(tplDir, "README.md"), []byte("# {{ .name }}\n\n{{ .description }}"), 0644)

	// Create apply folder
	applyDir := filepath.Join(tmpDir, "services", "comp1")
	os.MkdirAll(applyDir, 0755)

	// Build plan
	inv := &Inventory{
		InventoryDir: tmpDir,
		InventoryKey: "k2.inventory.yaml",
	}
	plan := NewActionPlan(inv)

	tpl := &types.IK2Template{
		K2: types.IK2TemplateRoot{
			Metadata: types.IK2Metadata{
				ID:     "test-tpl",
				Kind:   "template",
				Path:   filepath.Join(tplDir, "k2.template.yaml"),
				Folder: tplDir,
			},
			Body: types.IK2TemplateBody{
				Name:       "kind1",
				Parameters: map[string]any{"name": "kind1", "description": "A template"},
			},
		},
	}

	ref := types.IK2TemplateRef{
		Source: types.K2TemplateRefSourceInventory,
		Params: map[string]string{"id": "test-tpl"},
	}
	hash := ref.Hash()
	plan.Templates[hash] = tpl

	apply := &types.IK2TemplateApply{
		K2: types.IK2TemplateApplyRoot{
			Metadata: types.IK2Metadata{
				ID:     "test-apply",
				Kind:   "template-apply",
				Path:   filepath.Join(applyDir, "k2.apply.yaml"),
				Folder: applyDir,
			},
			Body: types.IK2TemplateApplyBody{
				Template: ref,
				Vars:     map[string]any{"name": "mycomp", "description": "My component"},
			},
		},
	}
	plan.AddEntity(types.IK2Metadata{ID: "test-apply"}, apply)

	ts := NewTemplatingStore(plan)
	ok, err := ts.ApplyTemplate("test-apply", hash, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected ApplyTemplate to return true")
	}

	// Verify README was rendered
	readmePath := filepath.Join(applyDir, "README.md")
	content, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("failed to read README: %v", err)
	}
	if !strings.Contains(string(content), "mycomp") {
		t.Fatalf("expected rendered content to contain 'mycomp', got: %s", string(content))
	}

	// Verify .gitignore was created
	gitignorePath := filepath.Join(applyDir, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		t.Fatal("expected .gitignore to be created")
	}
	gitignoreContent, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}
	if !strings.Contains(string(gitignoreContent), "README.md") {
		t.Fatalf("expected .gitignore to contain 'README.md', got: %s", string(gitignoreContent))
	}
}

func TestTemplatingStore_ApplyTemplate_TemplateNotFound(t *testing.T) {
	plan := NewActionPlan(nil)
	ts := NewTemplatingStore(plan)
	_, err := ts.ApplyTemplate("test", "nonexistent-hash", true)
	if err == nil {
		t.Fatal("expected error for missing template")
	}
}

func TestTemplatingStore_DestroyTemplate(t *testing.T) {
	tmpDir := t.TempDir()

	applyDir := filepath.Join(tmpDir, "services", "comp1")
	os.MkdirAll(applyDir, 0755)

	// Create files that would be destroyed
	os.WriteFile(filepath.Join(applyDir, "README.md"), []byte("# test"), 0644)
	os.WriteFile(filepath.Join(applyDir, ".gitignore"), []byte("README.md\n!k2.apply.yaml"), 0644)
	os.WriteFile(filepath.Join(applyDir, "k2.apply.yaml"), []byte("apply file"), 0644)

	plan := NewActionPlan(nil)
	apply := &types.IK2TemplateApply{
		K2: types.IK2TemplateApplyRoot{
			Metadata: types.IK2Metadata{
				ID:     "test-apply",
				Folder: applyDir,
			},
			Body: types.IK2TemplateApplyBody{
				Template: types.IK2TemplateRef{
					Source: types.K2TemplateRefSourceInventory,
					Params: map[string]string{"id": "some-tpl"},
				},
				Vars: map[string]any{},
			},
		},
	}
	plan.AddEntity(types.IK2Metadata{ID: "test-apply"}, apply)

	ts := NewTemplatingStore(plan)
	err := ts.DestroyTemplate("test-apply")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// README.md should be removed
	if _, err := os.Stat(filepath.Join(applyDir, "README.md")); !os.IsNotExist(err) {
		t.Fatal("expected README.md to be removed")
	}
	// .gitignore should be removed
	if _, err := os.Stat(filepath.Join(applyDir, ".gitignore")); !os.IsNotExist(err) {
		t.Fatal("expected .gitignore to be removed")
	}
}

func TestTemplatingStore_DestroyTemplate_NoGitignore(t *testing.T) {
	tmpDir := t.TempDir()
	applyDir := filepath.Join(tmpDir, "services", "comp1")
	os.MkdirAll(applyDir, 0755)

	plan := NewActionPlan(nil)
	apply := &types.IK2TemplateApply{
		K2: types.IK2TemplateApplyRoot{
			Metadata: types.IK2Metadata{
				ID:     "test-apply",
				Folder: applyDir,
			},
			Body: types.IK2TemplateApplyBody{
				Template: types.IK2TemplateRef{
					Source: types.K2TemplateRefSourceInventory,
					Params: map[string]string{"id": "tpl"},
				},
				Vars: map[string]any{},
			},
		},
	}
	plan.AddEntity(types.IK2Metadata{ID: "test-apply"}, apply)

	ts := NewTemplatingStore(plan)
	err := ts.DestroyTemplate("test-apply")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTemplatingStore_DestroyTemplate_NonexistentFolder(t *testing.T) {
	plan := NewActionPlan(nil)
	apply := &types.IK2TemplateApply{
		K2: types.IK2TemplateApplyRoot{
			Metadata: types.IK2Metadata{
				ID:     "test-apply",
				Folder: "/nonexistent/path/that/does/not/exist",
			},
			Body: types.IK2TemplateApplyBody{
				Template: types.IK2TemplateRef{
					Source: types.K2TemplateRefSourceInventory,
					Params: map[string]string{"id": "tpl"},
				},
				Vars: map[string]any{},
			},
		},
	}
	plan.AddEntity(types.IK2Metadata{ID: "test-apply"}, apply)

	ts := NewTemplatingStore(plan)
	err := ts.DestroyTemplate("test-apply")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTemplatingStore_DestroyTemplateRef(t *testing.T) {
	tmpDir := t.TempDir()
	refDir := filepath.Join(tmpDir, "ref-hash")
	os.MkdirAll(refDir, 0755)
	os.WriteFile(filepath.Join(refDir, "file.txt"), []byte("data"), 0644)

	ts := NewTemplatingStore(NewActionPlan(nil))
	err := ts.destroyTemplateRef(refDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(refDir); !os.IsNotExist(err) {
		t.Fatal("expected ref dir to be removed")
	}
}

func TestTemplatingStore_DestroyTemplateRef_NonexistentFolder(t *testing.T) {
	ts := NewTemplatingStore(NewActionPlan(nil))
	err := ts.destroyTemplateRef("/nonexistent/dir")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTemplatingStore_CleanupEmptyDirs(t *testing.T) {
	tmpDir := t.TempDir()

	// Create nested empty dirs
	emptyDir := filepath.Join(tmpDir, "a", "b", "c")
	os.MkdirAll(emptyDir, 0755)

	ts := NewTemplatingStore(NewActionPlan(nil))
	err := ts.cleanupEmptyDirs(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The empty nested dirs should be removed
	if _, err := os.Stat(filepath.Join(tmpDir, "a")); !os.IsNotExist(err) {
		t.Fatal("expected empty dirs to be cleaned up")
	}
}

func TestTemplatingStore_CleanupEmptyDirs_WithFiles(t *testing.T) {
	tmpDir := t.TempDir()

	dir := filepath.Join(tmpDir, "a", "b")
	os.MkdirAll(dir, 0755)
	os.WriteFile(filepath.Join(dir, "keep.txt"), []byte("data"), 0644)

	ts := NewTemplatingStore(NewActionPlan(nil))
	err := ts.cleanupEmptyDirs(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Dir with file should still exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatal("expected dir with files to remain")
	}
}
