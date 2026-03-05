package libs

import (
	"testing"
)

func TestRefsDir(t *testing.T) {
	if RefsDir != ".refs" {
		t.Fatalf("expected '.refs', got '%s'", RefsDir)
	}
}
