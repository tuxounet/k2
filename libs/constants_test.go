package libs

import (
	"testing"
)

func TestRefsDir(t *testing.T) {
	if RefsDir != ".refs" {
		t.Fatalf("expected '.refs', got '%s'", RefsDir)
	}
}

func TestDefaultInventoryFile(t *testing.T) {
	if DefaultInventoryFile != "k2.inventory.yaml" {
		t.Fatalf("expected 'k2.inventory.yaml', got '%s'", DefaultInventoryFile)
	}
}
