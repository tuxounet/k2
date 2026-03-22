package types

import (
	"testing"
)

func TestIK2TemplateRefHash_Deterministic(t *testing.T) {
	ref := IK2TemplateRef{
		Source: K2TemplateRefSourceInventory,
		Params: map[string]string{"id": "test-template"},
	}
	h1 := ref.Hash()
	h2 := ref.Hash()
	if h1 != h2 {
		t.Fatalf("hash should be deterministic, got %s and %s", h1, h2)
	}
}

func TestIK2TemplateRefHash_DifferentForDifferentSources(t *testing.T) {
	ref1 := IK2TemplateRef{
		Source: K2TemplateRefSourceInventory,
		Params: map[string]string{"id": "test-template"},
	}
	ref2 := IK2TemplateRef{
		Source: K2TemplateRefSourceGit,
		Params: map[string]string{"id": "test-template"},
	}
	if ref1.Hash() == ref2.Hash() {
		t.Fatal("hashes should differ for different sources")
	}
}

func TestIK2TemplateRefHash_DifferentForDifferentParams(t *testing.T) {
	ref1 := IK2TemplateRef{
		Source: K2TemplateRefSourceInventory,
		Params: map[string]string{"id": "template-a"},
	}
	ref2 := IK2TemplateRef{
		Source: K2TemplateRefSourceInventory,
		Params: map[string]string{"id": "template-b"},
	}
	if ref1.Hash() == ref2.Hash() {
		t.Fatal("hashes should differ for different params")
	}
}

func TestIK2TemplateRefHash_NotEmpty(t *testing.T) {
	ref := IK2TemplateRef{
		Source: K2TemplateRefSourceGit,
		Params: map[string]string{
			"repository": "https://github.com/tuxounet/k2.git",
			"branch":     "main",
			"path":       "samples/templates/fromGit1/k2.template.yaml",
		},
	}
	h := ref.Hash()
	if h == "" {
		t.Fatal("hash should not be empty")
	}
	if len(h) != 64 {
		t.Fatalf("expected sha256 hex length of 64, got %d", len(h))
	}
}

func TestK2TemplateRefSourceConstants(t *testing.T) {
	if K2TemplateRefSourceInventory != "inventory" {
		t.Fatalf("expected 'inventory', got '%s'", K2TemplateRefSourceInventory)
	}
	if K2TemplateRefSourceGit != "git" {
		t.Fatalf("expected 'git', got '%s'", K2TemplateRefSourceGit)
	}
}

func TestIK2TemplateApply_ExecuteScriptsNoOp(t *testing.T) {
	apply := &IK2TemplateApply{
		K2: IK2TemplateApplyRoot{
			Metadata: IK2Metadata{
				ID:     "test",
				Kind:   "template-apply",
				Folder: t.TempDir(),
			},
			Body: IK2TemplateApplyBody{
				Vars: map[string]any{"name": "test"},
			},
		},
	}
	// Empty scripts should be no-op
	if err := apply.ExecutePre(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := apply.ExecutePost(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := apply.ExecuteNuke(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := apply.ExecuteBootstrap(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIK2TemplateApply_ExecuteBootstrapWithScript(t *testing.T) {
	tmpDir := t.TempDir()
	apply := &IK2TemplateApply{
		K2: IK2TemplateApplyRoot{
			Metadata: IK2Metadata{
				ID:     "test-apply",
				Kind:   "template-apply",
				Folder: tmpDir,
			},
			Body: IK2TemplateApplyBody{
				Vars: map[string]any{"name": "mytest"},
			},
		},
	}
	apply.K2.Body.Scripts.Bootstrap = []string{"echo hello"}

	err := apply.ExecuteBootstrap()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIK2TemplateApply_ExecutePreWithScript(t *testing.T) {
	tmpDir := t.TempDir()
	apply := &IK2TemplateApply{
		K2: IK2TemplateApplyRoot{
			Metadata: IK2Metadata{
				ID:     "test-apply",
				Kind:   "template-apply",
				Folder: tmpDir,
			},
			Body: IK2TemplateApplyBody{
				Vars: map[string]any{"name": "mytest"},
			},
		},
	}
	apply.K2.Body.Scripts.Pre = []string{"echo pre"}

	err := apply.ExecutePre()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIK2TemplateApply_ExecutePostWithScript(t *testing.T) {
	tmpDir := t.TempDir()
	apply := &IK2TemplateApply{
		K2: IK2TemplateApplyRoot{
			Metadata: IK2Metadata{
				ID:     "test-apply",
				Kind:   "template-apply",
				Folder: tmpDir,
			},
			Body: IK2TemplateApplyBody{
				Vars: map[string]any{"name": "mytest"},
			},
		},
	}
	apply.K2.Body.Scripts.Post = []string{"echo post"}

	err := apply.ExecutePost()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIK2TemplateApply_ExecuteNukeWithScript(t *testing.T) {
	tmpDir := t.TempDir()
	apply := &IK2TemplateApply{
		K2: IK2TemplateApplyRoot{
			Metadata: IK2Metadata{
				ID:     "test-apply",
				Kind:   "template-apply",
				Folder: tmpDir,
			},
			Body: IK2TemplateApplyBody{
				Vars: map[string]any{"name": "mytest"},
			},
		},
	}
	apply.K2.Body.Scripts.Nuke = []string{"echo nuke"}

	err := apply.ExecuteNuke()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIK2TemplateApply_ExecuteScriptWithTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	apply := &IK2TemplateApply{
		K2: IK2TemplateApplyRoot{
			Metadata: IK2Metadata{
				ID:     "test-apply",
				Kind:   "template-apply",
				Folder: tmpDir,
			},
			Body: IK2TemplateApplyBody{
				Vars: map[string]any{"name": "world"},
			},
		},
	}
	apply.K2.Body.Scripts.Pre = []string{"echo {{ .name }}"}

	err := apply.ExecutePre()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIK2Template_ExecuteScriptsNoOp(t *testing.T) {
	tpl := &IK2Template{
		K2: IK2TemplateRoot{
			Metadata: IK2Metadata{
				ID:   "test-tpl",
				Kind: "template",
			},
			Body: IK2TemplateBody{
				Name:       "test",
				Parameters: map[string]any{"name": "test"},
			},
		},
	}
	apply := &IK2TemplateApply{
		K2: IK2TemplateApplyRoot{
			Metadata: IK2Metadata{
				ID:     "test-apply",
				Folder: t.TempDir(),
			},
			Body: IK2TemplateApplyBody{
				Vars: map[string]any{},
			},
		},
	}

	if err := tpl.ExecutePre(apply); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := tpl.ExecutePost(apply); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := tpl.ExecuteNuke(apply); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := tpl.ExecuteBootstrap(apply); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIK2Template_ExecuteBootstrapWithScript(t *testing.T) {
	tmpDir := t.TempDir()
	tpl := &IK2Template{
		K2: IK2TemplateRoot{
			Metadata: IK2Metadata{
				ID:   "test-tpl",
				Kind: "template",
			},
			Body: IK2TemplateBody{
				Name:       "test",
				Parameters: map[string]any{"name": "test"},
			},
		},
	}
	tpl.K2.Body.Scripts.Bootstrap = []string{"echo template boot of {{ .name }}"}

	apply := &IK2TemplateApply{
		K2: IK2TemplateApplyRoot{
			Metadata: IK2Metadata{
				ID:     "test-apply",
				Folder: tmpDir,
			},
			Body: IK2TemplateApplyBody{
				Vars: map[string]any{"name": "mycomp"},
			},
		},
	}

	err := tpl.ExecuteBootstrap(apply)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIK2Template_ExecutePreWithScript(t *testing.T) {
	tmpDir := t.TempDir()
	tpl := &IK2Template{
		K2: IK2TemplateRoot{
			Metadata: IK2Metadata{ID: "tpl1"},
			Body: IK2TemplateBody{
				Name:       "test",
				Parameters: map[string]any{"name": "test"},
			},
		},
	}
	tpl.K2.Body.Scripts.Pre = []string{"echo pre"}

	apply := &IK2TemplateApply{
		K2: IK2TemplateApplyRoot{
			Metadata: IK2Metadata{ID: "apply1", Folder: tmpDir},
			Body:     IK2TemplateApplyBody{Vars: map[string]any{}},
		},
	}

	if err := tpl.ExecutePre(apply); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIK2Template_ExecutePostWithScript(t *testing.T) {
	tmpDir := t.TempDir()
	tpl := &IK2Template{
		K2: IK2TemplateRoot{
			Metadata: IK2Metadata{ID: "tpl1"},
			Body: IK2TemplateBody{
				Name:       "test",
				Parameters: map[string]any{"name": "test"},
			},
		},
	}
	tpl.K2.Body.Scripts.Post = []string{"echo post of {{ .name }}"}

	apply := &IK2TemplateApply{
		K2: IK2TemplateApplyRoot{
			Metadata: IK2Metadata{ID: "apply1", Folder: tmpDir},
			Body:     IK2TemplateApplyBody{Vars: map[string]any{"name": "comp"}},
		},
	}

	if err := tpl.ExecutePost(apply); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIK2Metadata_Fields(t *testing.T) {
	m := IK2Metadata{
		ID:      "my-id",
		Kind:    "template",
		Version: "1.0.0",
		Path:    "/some/path",
		Folder:  "/some",
	}
	if m.ID != "my-id" {
		t.Fatalf("expected 'my-id', got '%s'", m.ID)
	}
	if m.Kind != "template" {
		t.Fatalf("expected 'template', got '%s'", m.Kind)
	}
	if m.Version != "1.0.0" {
		t.Fatalf("expected '1.0.0', got '%s'", m.Version)
	}
	if m.Path != "/some/path" {
		t.Fatalf("expected '/some/path', got '%s'", m.Path)
	}
	if m.Folder != "/some" {
		t.Fatalf("expected '/some', got '%s'", m.Folder)
	}
}
