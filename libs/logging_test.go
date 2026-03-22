package libs

import (
	"fmt"
	"testing"
)

func TestWriteOutput(t *testing.T) {
	// Just verify it doesn't panic
	WriteOutput("test message")
}

func TestWriteOutputf(t *testing.T) {
	WriteOutputf("test %s %d\n", "message", 42)
}

func TestWriteError(t *testing.T) {
	err := WriteError(fmt.Errorf("test error"))
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "test error" {
		t.Fatalf("expected 'test error', got '%s'", err.Error())
	}
}

func TestWriteErrorString(t *testing.T) {
	err := WriteErrorString("test error")
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "test error" {
		t.Fatalf("expected 'test error', got '%s'", err.Error())
	}
}

func TestWriteErrorf(t *testing.T) {
	err := WriteErrorf("error: %s %d", "test", 42)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "error: test 42" {
		t.Fatalf("expected 'error: test 42', got '%s'", err.Error())
	}
}
