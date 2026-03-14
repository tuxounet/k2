package stores

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/tuxounet/k2/libs"
	"gopkg.in/yaml.v3"
)

type layerRecipeType string

const (
	recipeShell   layerRecipeType = "shell"
	recipeUnknown layerRecipeType = "unknown"
)

func layerResolvePath(rootDir, layerPath, plan string) string {
	return filepath.Join(rootDir, layerPath, plan)
}

func layerDetectType(planDir string) layerRecipeType {
	if _, err := os.Stat(filepath.Join(planDir, "verbs", "up.sh")); err == nil {
		return recipeShell
	}
	return recipeUnknown
}

func layerRef(layerPath, plan string) string {
	return filepath.Base(layerPath) + "/" + plan
}

func layerStart(planDir string) error {
	cmd := exec.Command("bash", "verbs/up.sh")
	cmd.Dir = planDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func layerStop(planDir string) error {
	downScript := filepath.Join(planDir, "verbs", "down.sh")
	if _, err := os.Stat(downScript); os.IsNotExist(err) {
		return nil
	}
	cmd := exec.Command("bash", "verbs/down.sh")
	cmd.Dir = planDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func layerStatus(planDir string) string {
	statusScript := filepath.Join(planDir, "verbs", "status.sh")
	if _, err := os.Stat(statusScript); os.IsNotExist(err) {
		rt := layerDetectType(planDir)
		if rt == recipeShell {
			return "PRESENT"
		}
		return "UNKNOWN"
	}
	cmd := exec.Command("bash", "verbs/status.sh")
	cmd.Dir = planDir
	output, err := cmd.Output()
	if err != nil {
		return "DOWN"
	}
	status := strings.TrimSpace(string(output))
	if status == "" {
		return "UNKNOWN"
	}
	return strings.ToUpper(status)
}

func layerHealthcheck(planDir string) string {
	if layerRunHook(planDir, "healthcheck") != nil {
		return "FAIL"
	}
	return "OK"
}

func layerRunHook(planDir, hookName string) error {
	applyFile := filepath.Join(planDir, "k2.apply.yaml")
	if _, err := os.Stat(applyFile); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(applyFile)
	if err != nil {
		return nil
	}

	var doc struct {
		K2 struct {
			Body struct {
				Hooks map[string]string `yaml:"hooks"`
			} `yaml:"body"`
		} `yaml:"k2"`
	}
	if yaml.Unmarshal(data, &doc) != nil {
		return nil
	}

	hookContent, ok := doc.K2.Body.Hooks[hookName]
	if !ok || hookContent == "" {
		return nil
	}

	cmd := exec.Command("bash", "-c", hookContent)
	cmd.Dir = planDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func layerGetLinks(planDir string) []libs.StackLink {
	linksFile := filepath.Join(planDir, "links.env")
	f, err := os.Open(linksFile)
	if err != nil {
		return nil
	}
	defer f.Close()

	var links []libs.StackLink
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			links = append(links, libs.StackLink{
				Label: parts[0],
				URL:   expandEnvWithDefaults(parts[1]),
			})
		}
	}
	return links
}

func layerGetURL(planDir string) string {
	links := layerGetLinks(planDir)
	if len(links) > 0 {
		return links[0].URL
	}
	return ""
}

func layerLogs(planDir string) error {
	logsScript := filepath.Join(planDir, "verbs", "logs.sh")
	if _, err := os.Stat(logsScript); os.IsNotExist(err) {
		return fmt.Errorf("no logs verb available for '%s'", filepath.Base(planDir))
	}
	cmd := exec.Command("bash", "verbs/logs.sh")
	cmd.Dir = planDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func layerBuild(planDir string) error {
	buildScript := filepath.Join(planDir, "verbs", "build.sh")
	if _, err := os.Stat(buildScript); os.IsNotExist(err) {
		return nil
	}
	cmd := exec.Command("bash", "verbs/build.sh")
	cmd.Dir = planDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func layerShell(planDir string) error {
	shellScript := filepath.Join(planDir, "verbs", "shell.sh")
	if _, err := os.Stat(shellScript); os.IsNotExist(err) {
		return fmt.Errorf("no shell verb available for '%s'", filepath.Base(planDir))
	}
	cmd := exec.Command("bash", "verbs/shell.sh")
	cmd.Dir = planDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func layerListVerbs(planDir string) []string {
	verbsDir := filepath.Join(planDir, "verbs")
	entries, err := os.ReadDir(verbsDir)
	if err != nil {
		return nil
	}
	var verbs []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sh") {
			verbs = append(verbs, strings.TrimSuffix(e.Name(), ".sh"))
		}
	}
	sort.Strings(verbs)
	return verbs
}

func layerRunVerb(planDir, verb string, args []string) error {
	verbScript := filepath.Join(planDir, "verbs", verb+".sh")
	if _, err := os.Stat(verbScript); os.IsNotExist(err) {
		available := layerListVerbs(planDir)
		if len(available) > 0 {
			return fmt.Errorf("verb '%s' not found. Available: %s", verb, strings.Join(available, ", "))
		}
		return fmt.Errorf("verb '%s' not found (no verbs available)", verb)
	}

	cmdArgs := []string{verbScript}
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command("bash", cmdArgs...)
	cmd.Dir = planDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

var envDefaultRegex = regexp.MustCompile(`\$\{([^}:]+):-([^}]*)\}`)

func expandEnvWithDefaults(s string) string {
	s = envDefaultRegex.ReplaceAllStringFunc(s, func(match string) string {
		groups := envDefaultRegex.FindStringSubmatch(match)
		val := os.Getenv(groups[1])
		if val == "" {
			return groups[2]
		}
		return val
	})
	return os.ExpandEnv(s)
}

func loadDotEnv(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			os.Setenv(parts[0], parts[1])
		}
	}
	return nil
}

func exportEnvMap(env map[string]string) {
	for k, v := range env {
		os.Setenv(k, v)
	}
}

func loadDefaultsEnv(planDir string) {
	defaultsFile := filepath.Join(planDir, "defaults.env")
	loadDotEnv(defaultsFile)
}
