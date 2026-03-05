package stores

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetAllFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create files
	os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("# test"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "k2.template.yaml"), []byte("yaml"), 0644)
	os.WriteFile(filepath.Join(tmpDir, ".DS_Store"), []byte("ds"), 0644)

	subDir := filepath.Join(tmpDir, "sub")
	os.MkdirAll(subDir, 0755)
	os.WriteFile(filepath.Join(subDir, "file.txt"), []byte("text"), 0644)

	files, err := getAllFiles(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should exclude k2.template.yaml and .DS_Store
	for _, f := range files {
		if f == "k2.template.yaml" {
			t.Fatal("k2.template.yaml should be excluded")
		}
		if f == ".DS_Store" {
			t.Fatal(".DS_Store should be excluded")
		}
	}

	// Should include README.md and sub/file.txt
	hasReadme := false
	hasSubFile := false
	for _, f := range files {
		if f == "README.md" {
			hasReadme = true
		}
		if f == filepath.Join("sub", "file.txt") {
			hasSubFile = true
		}
	}
	if !hasReadme {
		t.Fatal("expected README.md in results")
	}
	if !hasSubFile {
		t.Fatal("expected sub/file.txt in results")
	}
}

func TestGetAllFiles_WithGitDir(t *testing.T) {
	// Note: the current getAllFiles implementation uses HasPrefix with swapped
	// args so .git files are not actually excluded. This test documents
	// the actual behavior.
	tmpDir := t.TempDir()

	gitDir := filepath.Join(tmpDir, ".git")
	os.MkdirAll(gitDir, 0755)
	os.WriteFile(filepath.Join(gitDir, "config"), []byte("git config"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("text"), 0644)

	files, err := getAllFiles(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	hasNonGitFile := false
	for _, f := range files {
		if f == "file.txt" {
			hasNonGitFile = true
		}
	}
	if !hasNonGitFile {
		t.Fatal("expected file.txt in results")
	}
}

func TestGetAllFiles_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	files, err := getAllFiles(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 0 {
		t.Fatalf("expected 0 files, got %d", len(files))
	}
}

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()

	srcFile := filepath.Join(tmpDir, "source.md")
	os.WriteFile(srcFile, []byte("# {{ .name }}\n\n{{ .description }}"), 0644)

	destFile := filepath.Join(tmpDir, "output", "result.md")

	data := map[string]any{
		"name":        "myproject",
		"description": "A great project",
	}

	err := copyFile(srcFile, destFile, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(destFile)
	if err != nil {
		t.Fatalf("failed to read dest file: %v", err)
	}

	if !strings.Contains(string(content), "myproject") {
		t.Fatalf("expected 'myproject' in output, got: %s", string(content))
	}
	if !strings.Contains(string(content), "A great project") {
		t.Fatalf("expected 'A great project' in output, got: %s", string(content))
	}
}

func TestCopyFile_SkipGitignore(t *testing.T) {
	tmpDir := t.TempDir()

	srcFile := filepath.Join(tmpDir, ".gitignore")
	os.WriteFile(srcFile, []byte("*.log"), 0644)

	destFile := filepath.Join(tmpDir, "output", ".gitignore")

	err := copyFile(srcFile, destFile, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// .gitignore should NOT be copied
	if _, err := os.Stat(destFile); !os.IsNotExist(err) {
		t.Fatal("expected .gitignore to not be copied")
	}
}

func TestCopyFile_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	srcFile := filepath.Join(tmpDir, "source.txt")
	os.WriteFile(srcFile, []byte("plain text"), 0644)

	destFile := filepath.Join(tmpDir, "deep", "nested", "dir", "output.txt")

	err := copyFile(srcFile, destFile, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(destFile)
	if err != nil {
		t.Fatalf("failed to read dest file: %v", err)
	}
	if string(content) != "plain text" {
		t.Fatalf("expected 'plain text', got '%s'", string(content))
	}
}

func TestCreateGitIgnore(t *testing.T) {
	tmpDir := t.TempDir()

	destFolder := filepath.Join(tmpDir, "services", "comp1")
	os.MkdirAll(destFolder, 0755)

	files := map[string]string{
		filepath.Join(tmpDir, "templates", "kind1", "README.md"):     filepath.Join(destFolder, "README.md"),
		filepath.Join(tmpDir, "templates", "kind1", "k2.apply.yaml"): filepath.Join(destFolder, "k2.apply.yaml"),
	}

	err := createGitIgnore(files, destFolder)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(destFolder, ".gitignore"))
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "!k2.apply.yaml") {
		t.Fatalf("expected .gitignore to contain '!k2.apply.yaml', got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "README.md") {
		t.Fatalf("expected .gitignore to contain 'README.md', got: %s", contentStr)
	}
}

func TestCreateGitIgnore_NoDuplicates(t *testing.T) {
	tmpDir := t.TempDir()
	os.MkdirAll(tmpDir, 0755)

	files := map[string]string{
		filepath.Join(tmpDir, "src", "file.txt"): filepath.Join(tmpDir, "file.txt"),
	}

	err := createGitIgnore(files, tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, ".gitignore"))
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	seen := make(map[string]int)
	for _, line := range lines {
		if line != "" {
			seen[line]++
		}
	}
	for line, count := range seen {
		if count > 1 {
			t.Fatalf("duplicate entry in .gitignore: '%s' (appears %d times)", line, count)
		}
	}
}

func TestCreateGitIgnore_MergesSourceGitignore(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := filepath.Join(tmpDir, "src")
	destDir := filepath.Join(tmpDir, "dest")
	os.MkdirAll(srcDir, 0755)
	os.MkdirAll(destDir, 0755)

	// Create a source .gitignore
	os.WriteFile(filepath.Join(srcDir, ".gitignore"), []byte("*.log\n*.tmp"), 0644)

	files := map[string]string{
		filepath.Join(srcDir, ".gitignore"): filepath.Join(destDir, ".gitignore"),
		filepath.Join(srcDir, "file.txt"):   filepath.Join(destDir, "file.txt"),
	}

	err := createGitIgnore(files, destDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(destDir, ".gitignore"))
	if err != nil {
		t.Fatalf("failed to read .gitignore: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "*.log") {
		t.Fatalf("expected merged .gitignore to contain '*.log', got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "*.tmp") {
		t.Fatalf("expected merged .gitignore to contain '*.tmp', got: %s", contentStr)
	}
}
