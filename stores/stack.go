package stores

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/tuxounet/k2/libs"
	"github.com/tuxounet/k2/types"

	"gopkg.in/yaml.v3"
)

type StackStore struct {
	RootDir    string
	StacksDir  string
	Name       string
	Definition *types.IK2Stack
	Debug      bool
}

func NewStackStore(rootDir string, stackName string, debug bool) (*StackStore, error) {
	absRoot, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, err
	}

	stacksDir := filepath.Join(absRoot, "stacks")
	stackFile := filepath.Join(stacksDir, stackName+".yaml")

	if _, err := os.Stat(stackFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("stack '%s' not found: %s", stackName, stackFile)
	}

	data, err := os.ReadFile(stackFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read stack file: %w", err)
	}

	var def types.IK2Stack
	if err := yaml.Unmarshal(data, &def); err != nil {
		return nil, fmt.Errorf("cannot parse stack '%s': %w", stackName, err)
	}

	return &StackStore{
		RootDir:    absRoot,
		StacksDir:  stacksDir,
		Name:       stackName,
		Definition: &def,
		Debug:      debug,
	}, nil
}

func (s *StackStore) logDebug(format string, a ...any) {
	if s.Debug {
		libs.WriteSubStep(format, a...)
	}
}

func (s *StackStore) loadEnv() {
	loadDotEnv(filepath.Join(s.RootDir, ".env"))
	exportEnvMap(s.Definition.Stack.Env)
}

func (s *StackStore) exportLayerEnv(idx int) {
	layer := s.Definition.Stack.Layers[idx]
	exportEnvMap(layer.Env)
	planDir := layerResolvePath(s.RootDir, layer.Layer, layer.Plan)
	loadDefaultsEnv(planDir)
}

func (s *StackStore) doRender() {
	inventoryFile := filepath.Join(s.RootDir, "k2.inventory.yaml")
	if _, err := os.Stat(inventoryFile); os.IsNotExist(err) {
		s.logDebug("No k2.inventory.yaml found, skipping rendering")
		return
	}

	libs.WriteSubStep("Rendering templates...")
	inventory, err := NewInventory(inventoryFile)
	if err != nil {
		libs.WriteSubStep("Rendering failed (using existing files): %v", err)
		return
	}
	plan, err := inventory.Plan()
	if err != nil {
		libs.WriteSubStep("Rendering plan failed: %v", err)
		return
	}
	if err := inventory.Apply(plan); err != nil {
		libs.WriteSubStep("Rendering apply failed: %v", err)
	}
}

func (s *StackStore) Up() error {
	layers := s.Definition.Stack.Layers
	layerCount := len(layers)

	s.loadEnv()
	libs.WriteStackBanner("up", s.Name, layerCount)
	s.doRender()

	successCount := 0
	var failures []string
	var allLinks []libs.StackLink

	for i, layer := range layers {
		ref := layerRef(layer.Layer, layer.Plan)
		planDir := layerResolvePath(s.RootDir, layer.Layer, layer.Plan)

		libs.WriteStackStepStart(i+1, layerCount, "▶ Démarrage", ref)
		s.logDebug("path=%s", planDir)

		if _, err := os.Stat(planDir); os.IsNotExist(err) {
			libs.WriteStackStepSkip(layer.Plan, "dossier introuvable")
			continue
		}

		s.exportLayerEnv(i)

		rt := layerDetectType(planDir)
		if rt == recipeUnknown {
			libs.WriteStackStepSkip(layer.Plan, "aucune recette détectée")
			continue
		}

		libs.WriteSubStep("verbs/up.sh")

		layerRunHook(planDir, "pre_start")

		if err := layerStart(planDir); err != nil {
			libs.WriteStackStepFail(layer.Plan, "échec (voir logs)")
			failures = append(failures, ref+" — échec")
			continue
		}

		layerRunHook(planDir, "post_start")
		libs.WriteStackStepOk(layer.Plan, "démarré")
		successCount++

		allLinks = append(allLinks, layerGetLinks(planDir)...)
	}

	libs.WriteStackSummary("démarrée", s.Name, successCount, layerCount, failures)
	libs.WriteStackLinksTable(allLinks)
	return nil
}

func (s *StackStore) Down() error {
	layers := s.Definition.Stack.Layers
	layerCount := len(layers)

	s.loadEnv()
	libs.WriteStackBanner("down", s.Name, layerCount)

	successCount := 0
	var failures []string

	for i := layerCount - 1; i >= 0; i-- {
		layer := layers[i]
		ref := layerRef(layer.Layer, layer.Plan)
		planDir := layerResolvePath(s.RootDir, layer.Layer, layer.Plan)
		displayIdx := layerCount - i

		libs.WriteStackStepStart(displayIdx, layerCount, "■ Arrêt", ref)

		if _, err := os.Stat(planDir); os.IsNotExist(err) {
			libs.WriteStackStepSkip(layer.Plan, "dossier introuvable")
			continue
		}

		s.exportLayerEnv(i)

		rt := layerDetectType(planDir)
		if rt == recipeUnknown {
			libs.WriteStackStepSkip(layer.Plan, "aucune recette détectée")
			continue
		}

		libs.WriteSubStep("verbs/down.sh")

		layerRunHook(planDir, "pre_stop")

		if err := layerStop(planDir); err != nil {
			libs.WriteStackStepFail(layer.Plan, "échec")
			failures = append(failures, ref+" — échec")
			continue
		}

		layerRunHook(planDir, "post_stop")
		libs.WriteStackStepOk(layer.Plan, "arrêté")
		successCount++
	}

	libs.WriteStackSummary("arrêtée", s.Name, successCount, layerCount, failures)
	return nil
}

func (s *StackStore) Rm() error {
	return s.execVerbOnLayers("rm", "", nil)
}

// execVerbOnLayers runs a named verb on every layer of the stack (optionally
// filtered to a single layer). Layers that do not have the verb script are
// silently skipped. Hooks named pre_<verb> and post_<verb> are called when
// defined in k2.apply.yaml.
func (s *StackStore) execVerbOnLayers(verb, targetLayer string, args []string) error {
	layers := s.Definition.Stack.Layers
	layerCount := len(layers)

	s.loadEnv()
	libs.WriteStackBanner(verb, s.Name, layerCount)
	s.doRender()

	successCount := 0
	skipCount := 0
	var failures []string

	for i, layer := range layers {
		ref := layerRef(layer.Layer, layer.Plan)
		planDir := layerResolvePath(s.RootDir, layer.Layer, layer.Plan)

		if targetLayer != "" {
			if !strings.Contains(ref, targetLayer) && layer.Plan != targetLayer &&
				!strings.Contains(layer.Layer+"/"+layer.Plan, targetLayer) {
				continue
			}
		}

		libs.WriteStackStepStart(i+1, layerCount, "▶ "+verb, ref)
		s.logDebug("path=%s", planDir)

		if _, err := os.Stat(planDir); os.IsNotExist(err) {
			libs.WriteStackStepSkip(layer.Plan, "dossier introuvable")
			continue
		}

		if !layerHasVerb(planDir, verb) {
			libs.WriteStackStepSkip(layer.Plan, "pas de verbs/"+verb+".sh")
			skipCount++
			continue
		}

		s.exportLayerEnv(i)

		libs.WriteSubStep("verbs/%s.sh", verb)

		layerRunHook(planDir, "pre_"+verb)

		if err := layerRunVerb(planDir, verb, args); err != nil {
			libs.WriteStackStepFail(layer.Plan, "échec (voir logs)")
			failures = append(failures, ref+" — échec")
			continue
		}

		layerRunHook(planDir, "post_"+verb)
		libs.WriteStackStepOk(layer.Plan, verb+" terminé")
		successCount++
	}

	if targetLayer != "" && successCount == 0 && len(failures) == 0 && skipCount == 0 {
		return fmt.Errorf("layer '%s' not found in stack '%s'", targetLayer, s.Name)
	}

	libs.WriteStackSummary(verb+" terminé", s.Name, successCount, layerCount-skipCount, failures)
	return nil
}

func (s *StackStore) Build(targetLayer string) error {
	return s.execVerbOnLayers("build", targetLayer, nil)
}

func (s *StackStore) Exec(verb string, args []string) error {
	return s.execVerbOnLayers(verb, "", args)
}

func (s *StackStore) Restart() error {
	if err := s.Down(); err != nil {
		return err
	}
	return s.Up()
}

func (s *StackStore) Status() error {
	layers := s.Definition.Stack.Layers
	layerCount := len(layers)

	s.loadEnv()
	libs.WriteStackBanner("status", s.Name, layerCount)

	var statuses []libs.StackStatus
	for i, layer := range layers {
		ref := layerRef(layer.Layer, layer.Plan)
		planDir := layerResolvePath(s.RootDir, layer.Layer, layer.Plan)
		s.exportLayerEnv(i)

		status := layerStatus(planDir)
		url := layerGetURL(planDir)

		statuses = append(statuses, libs.StackStatus{
			Ref:    ref,
			Status: status,
			URL:    url,
		})
	}

	libs.WriteStackStatusTable(statuses)
	return nil
}

func (s *StackStore) Logs(targetLayer string) error {
	layers := s.Definition.Stack.Layers
	layerCount := len(layers)

	s.loadEnv()

	if targetLayer != "" {
		for i, layer := range layers {
			ref := layerRef(layer.Layer, layer.Plan)
			planDir := layerResolvePath(s.RootDir, layer.Layer, layer.Plan)
			if strings.Contains(ref, targetLayer) || layer.Plan == targetLayer {
				s.exportLayerEnv(i)
				libs.WriteTitle("Logs de %s", ref)
				return layerLogs(planDir)
			}
		}
		return fmt.Errorf("layer '%s' not found in stack '%s'", targetLayer, s.Name)
	}

	libs.WriteStackBanner("logs", s.Name, layerCount)

	var cmds []*os.Process
	for i, layer := range layers {
		planDir := layerResolvePath(s.RootDir, layer.Layer, layer.Plan)
		s.exportLayerEnv(i)

		logsScript := filepath.Join(planDir, "verbs", "logs.sh")
		if _, err := os.Stat(logsScript); err == nil {
			cmd := logShellCmd(planDir)
			if err := cmd.Start(); err == nil {
				cmds = append(cmds, cmd.Process)
			}
		}
	}

	if len(cmds) == 0 {
		libs.WriteDetail("No logs available")
		return nil
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	for _, p := range cmds {
		p.Kill()
	}
	return nil
}

func logShellCmd(planDir string) *exec.Cmd {
	cmd := exec.Command("bash", "verbs/logs.sh")
	cmd.Dir = planDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func (s *StackStore) Healthcheck() error {
	layers := s.Definition.Stack.Layers
	layerCount := len(layers)

	s.loadEnv()
	libs.WriteStackBanner("healthcheck", s.Name, layerCount)

	var statuses []libs.StackStatus
	for i, layer := range layers {
		ref := layerRef(layer.Layer, layer.Plan)
		planDir := layerResolvePath(s.RootDir, layer.Layer, layer.Plan)

		libs.WriteStackStepStart(i+1, layerCount, "♥ Vérification", ref)
		s.exportLayerEnv(i)

		health := layerHealthcheck(planDir)
		switch health {
		case "OK":
			libs.WriteStackStepOk(layer.Plan, "OK")
		default:
			libs.WriteStackStepFail(layer.Plan, health)
		}

		statuses = append(statuses, libs.StackStatus{
			Ref:    ref,
			Status: health,
		})
	}

	libs.WriteStackStatusTable(statuses)
	return nil
}

func (s *StackStore) Shell(targetLayer string) error {
	layers := s.Definition.Stack.Layers

	s.loadEnv()

	if targetLayer != "" {
		for i, layer := range layers {
			ref := layerRef(layer.Layer, layer.Plan)
			planDir := layerResolvePath(s.RootDir, layer.Layer, layer.Plan)
			if strings.Contains(ref, targetLayer) || layer.Plan == targetLayer ||
				strings.Contains(layer.Layer+"/"+layer.Plan, targetLayer) {
				s.exportLayerEnv(i)
				libs.WriteTitle("Shell dans %s", ref)
				return layerShell(planDir)
			}
		}
		return fmt.Errorf("layer '%s' not found in stack '%s'", targetLayer, s.Name)
	}

	for i, layer := range layers {
		planDir := layerResolvePath(s.RootDir, layer.Layer, layer.Plan)
		rt := layerDetectType(planDir)
		if rt != recipeUnknown {
			ref := layerRef(layer.Layer, layer.Plan)
			s.exportLayerEnv(i)
			libs.WriteTitle("Shell dans %s (premier layer disponible)", ref)
			return layerShell(planDir)
		}
	}

	return fmt.Errorf("no layer with recipe found in stack '%s'", s.Name)
}

func (s *StackStore) Urls() error {
	layers := s.Definition.Stack.Layers
	layerCount := len(layers)

	s.loadEnv()
	libs.WriteStackBanner("urls", s.Name, layerCount)

	var allLinks []libs.StackLink
	for i, layer := range layers {
		planDir := layerResolvePath(s.RootDir, layer.Layer, layer.Plan)
		s.exportLayerEnv(i)
		allLinks = append(allLinks, layerGetLinks(planDir)...)
	}

	if len(allLinks) == 0 {
		libs.WriteDetail("Aucune URL d'accès définie dans la stack '%s'", s.Name)
		return nil
	}

	libs.WriteStackLinksTable(allLinks)
	return nil
}

func (s *StackStore) Run(targetLayer, verb string, args []string) error {
	layers := s.Definition.Stack.Layers

	s.loadEnv()

	for i, layer := range layers {
		ref := layerRef(layer.Layer, layer.Plan)
		planDir := layerResolvePath(s.RootDir, layer.Layer, layer.Plan)

		if strings.Contains(ref, targetLayer) || layer.Plan == targetLayer ||
			strings.Contains(layer.Layer+"/"+layer.Plan, targetLayer) {
			s.exportLayerEnv(i)

			if verb == "" {
				libs.WriteTitle("Verbes disponibles pour %s :", ref)
				verbs := layerListVerbs(planDir)
				if len(verbs) > 0 {
					for _, v := range verbs {
						fmt.Printf("    %s%s%s\n", libs.CyanColor(), v, libs.ResetCol())
					}
				} else {
					libs.WriteDetail("Aucun verbe disponible")
				}
				return nil
			}

			s.doRender()

			libs.WriteTitle("Exécution du verbe %s sur %s", verb, ref)
			if err := layerRunVerb(planDir, verb, args); err != nil {
				return fmt.Errorf("verb '%s' failed on %s: %w", verb, ref, err)
			}
			libs.WriteStep(libs.IconApply, "Verbe %s terminé avec succès", verb)
			return nil
		}
	}

	return fmt.Errorf("layer '%s' not found in stack '%s'", targetLayer, s.Name)
}

func ListStacks(rootDir string) error {
	absRoot, err := filepath.Abs(rootDir)
	if err != nil {
		return err
	}
	stacksDir := filepath.Join(absRoot, "stacks")
	entries, err := os.ReadDir(stacksDir)
	if err != nil {
		return fmt.Errorf("cannot read stacks directory: %w", err)
	}

	fmt.Println()
	fmt.Printf("  %sStacks disponibles :%s\n", libs.BoldStyle(), libs.ResetCol())
	fmt.Println("  ────────────────────────────────────────")

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		stackName := strings.TrimSuffix(e.Name(), ".yaml")
		stackFile := filepath.Join(stacksDir, e.Name())

		data, err := os.ReadFile(stackFile)
		description := ""
		if err == nil {
			var def types.IK2Stack
			if yaml.Unmarshal(data, &def) == nil && def.Stack.Description != "" {
				description = strings.TrimSpace(def.Stack.Description)
			}
		}

		if description != "" {
			fmt.Printf("  %s%s%s  %s— %s%s\n", libs.CyanColor(), stackName, libs.ResetCol(), libs.GrayColor(), description, libs.ResetCol())
		} else {
			fmt.Printf("  %s%s%s\n", libs.CyanColor(), stackName, libs.ResetCol())
		}
	}
	fmt.Println()
	return nil
}

func ListLayers(rootDir string) error {
	absRoot, err := filepath.Abs(rootDir)
	if err != nil {
		return err
	}
	layersDir := filepath.Join(absRoot, "layers")
	entries, err := os.ReadDir(layersDir)
	if err != nil {
		return fmt.Errorf("cannot read layers directory: %w", err)
	}

	fmt.Println()
	fmt.Printf("  %sLayers disponibles :%s\n", libs.BoldStyle(), libs.ResetCol())
	fmt.Println("  ────────────────────────────────────────")

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		layerName := e.Name()
		layerDir := filepath.Join(layersDir, layerName)

		planEntries, err := os.ReadDir(layerDir)
		if err != nil {
			continue
		}

		var plans []string
		for _, pe := range planEntries {
			if pe.IsDir() {
				plans = append(plans, pe.Name())
			}
		}

		if len(plans) > 0 {
			fmt.Printf("  %s%s/%s\n", libs.CyanColor(), layerName, libs.ResetCol())
			for _, p := range plans {
				planPath := filepath.Join(layerDir, p)
				indicator := fmt.Sprintf("%s○%s", libs.GrayColor(), libs.ResetCol())
				if _, err := os.Stat(filepath.Join(planPath, "verbs", "up.sh")); err == nil {
					indicator = fmt.Sprintf("%s◆%s", libs.CyanColor(), libs.ResetCol())
				}
				fmt.Printf("    %s %s\n", indicator, p)
			}
		}
	}
	fmt.Println()
	return nil
}
