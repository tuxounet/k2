package libs

import (
	"testing"
)

func TestRenderTemplate_SimpleString(t *testing.T) {
	result, err := RenderTemplate("hello world", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "hello world" {
		t.Fatalf("expected 'hello world', got '%s'", string(result))
	}
}

func TestRenderTemplate_WithVariable(t *testing.T) {
	data := map[string]any{"name": "k2"}
	result, err := RenderTemplate("hello {{ .name }}", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "hello k2" {
		t.Fatalf("expected 'hello k2', got '%s'", string(result))
	}
}

func TestRenderTemplate_WithMultipleVariables(t *testing.T) {
	data := map[string]any{"name": "k2", "description": "template engine"}
	result, err := RenderTemplate("{{ .name }} is a {{ .description }}", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "k2 is a template engine"
	if string(result) != expected {
		t.Fatalf("expected '%s', got '%s'", expected, string(result))
	}
}

func TestRenderTemplate_WithSprigFunctions(t *testing.T) {
	data := map[string]any{"name": "hello"}
	result, err := RenderTemplate("{{ .name | upper }}", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "HELLO" {
		t.Fatalf("expected 'HELLO', got '%s'", string(result))
	}
}

func TestRenderTemplate_WithRange(t *testing.T) {
	data := map[string]any{"items": []string{"a", "b", "c"}}
	result, err := RenderTemplate("{{ range .items }}{{ . }}{{ end }}", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "abc" {
		t.Fatalf("expected 'abc', got '%s'", string(result))
	}
}

func TestRenderTemplate_InvalidTemplate(t *testing.T) {
	_, err := RenderTemplate("{{ .invalid }", nil)
	if err == nil {
		t.Fatal("expected error for invalid template syntax")
	}
}

func TestRenderTemplate_MissingKey(t *testing.T) {
	data := map[string]any{}
	_, err := RenderTemplate("{{ .missing }}", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderTemplate_EmptyString(t *testing.T) {
	result, err := RenderTemplate("", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "" {
		t.Fatalf("expected empty string, got '%s'", string(result))
	}
}

func TestRenderTemplate_NestedMap(t *testing.T) {
	data := map[string]any{
		"obj": map[string]any{
			"a": 1,
			"b": "stc",
		},
	}
	result, err := RenderTemplate("{{ .obj.b }}", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "stc" {
		t.Fatalf("expected 'stc', got '%s'", string(result))
	}
}

func TestMergeMaps_Empty(t *testing.T) {
	result := MergeMaps()
	if len(result) != 0 {
		t.Fatalf("expected empty map, got %v", result)
	}
}

func TestMergeMaps_SingleMap(t *testing.T) {
	m := map[string]any{"a": 1, "b": "two"}
	result := MergeMaps(m)
	if len(result) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(result))
	}
	if result["a"] != 1 {
		t.Fatalf("expected a=1, got a=%v", result["a"])
	}
	if result["b"] != "two" {
		t.Fatalf("expected b='two', got b=%v", result["b"])
	}
}

func TestMergeMaps_MultipleMaps(t *testing.T) {
	m1 := map[string]any{"a": 1, "b": "two"}
	m2 := map[string]any{"b": "overridden", "c": 3}
	result := MergeMaps(m1, m2)
	if len(result) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(result))
	}
	if result["a"] != 1 {
		t.Fatalf("expected a=1, got a=%v", result["a"])
	}
	if result["b"] != "overridden" {
		t.Fatalf("expected b='overridden', got b=%v", result["b"])
	}
	if result["c"] != 3 {
		t.Fatalf("expected c=3, got c=%v", result["c"])
	}
}

func TestMergeMaps_ThreeMaps(t *testing.T) {
	m1 := map[string]any{"x": 1}
	m2 := map[string]any{"y": 2}
	m3 := map[string]any{"x": 99, "z": 3}
	result := MergeMaps(m1, m2, m3)
	if result["x"] != 99 {
		t.Fatalf("expected x=99 (last write wins), got x=%v", result["x"])
	}
	if result["y"] != 2 {
		t.Fatalf("expected y=2, got y=%v", result["y"])
	}
	if result["z"] != 3 {
		t.Fatalf("expected z=3, got z=%v", result["z"])
	}
}

func TestMergeMaps_NilValues(t *testing.T) {
	m1 := map[string]any{"a": nil}
	m2 := map[string]any{"b": "value"}
	result := MergeMaps(m1, m2)
	if result["a"] != nil {
		t.Fatalf("expected a=nil, got a=%v", result["a"])
	}
	if result["b"] != "value" {
		t.Fatalf("expected b='value', got b=%v", result["b"])
	}
}
