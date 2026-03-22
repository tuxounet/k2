package libs

import "fmt"

func WriteStackBanner(verb, stack string, layerCount int) {
	label := ""
	if layerCount > 0 {
		label = fmt.Sprintf(" (%d layers)", layerCount)
	}
	fmt.Printf("\n%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", boldStyle, resetColor)
	fmt.Printf("%s  K2 ▶ %s %s%s%s%s%s\n", boldStyle, verb, cyanColor, stack, resetColor, boldStyle, label+resetColor)
	fmt.Printf("%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n\n", boldStyle, resetColor)
}

func WriteStackStepStart(current, total int, verb, ref string) {
	fmt.Printf("%s[%d/%d]%s %s%s%s de %s%s%s%s\n", boldStyle, current, total, resetColor, cyanColor, verb, resetColor, boldStyle, cyanColor, ref, resetColor)
}

func WriteStackStepOk(plan, message string) {
	fmt.Printf("       %s✓%s %s %s\n\n", greenColor, resetColor, plan, message)
}

func WriteStackStepFail(plan, message string) {
	fmt.Printf("       %s✗%s %s — %s%s%s\n\n", redColor, resetColor, plan, redColor, message, resetColor)
}

func WriteStackStepSkip(plan, reason string) {
	fmt.Printf("       %s⊘%s %s — %s%s%s\n\n", yellowColor, resetColor, plan, yellowColor, reason, resetColor)
}

func WriteStackSummary(verbPast, stack string, successCount, total int, failures []string) {
	fmt.Printf("%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", boldStyle, resetColor)
	if successCount == total {
		fmt.Printf("  %s✓%s Stack %s%s%s %s %s(%d/%d layers)%s\n", greenColor, resetColor, cyanColor, stack, resetColor, verbPast, greenColor, successCount, total, resetColor)
	} else {
		fmt.Printf("  %s⚠%s Stack %s%s%s partiellement %s %s(%d/%d layers)%s\n", yellowColor, resetColor, cyanColor, stack, resetColor, verbPast, yellowColor, successCount, total, resetColor)
		for _, f := range failures {
			fmt.Printf("  %s✗%s %s\n", redColor, resetColor, f)
		}
	}
	fmt.Printf("%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n\n", boldStyle, resetColor)
}

func WriteStackLinksTable(links []StackLink) {
	if len(links) == 0 {
		return
	}
	fmt.Printf("  %sAccès aux services :%s\n", boldStyle, resetColor)
	fmt.Printf("  ──────────────────────────── ──────────────────────────────\n")
	for _, l := range links {
		fmt.Printf("  %-30s %s\n", l.Label, l.URL)
	}
	fmt.Println()
}

func WriteStackStatusTable(statuses []StackStatus) {
	fmt.Println()
	fmt.Printf("  %sLAYER                                    STATUS     URL%s\n", boldStyle, resetColor)
	fmt.Printf("  ──────────────────────────────────────── ────────── ──────────────────────────\n")
	for _, s := range statuses {
		statusDisplay := ""
		switch s.Status {
		case "UP":
			statusDisplay = fmt.Sprintf("%s✓ UP%s", greenColor, resetColor)
		case "DOWN":
			statusDisplay = fmt.Sprintf("%s✗ DOWN%s", redColor, resetColor)
		case "DEGRADED":
			statusDisplay = fmt.Sprintf("%s⚠ DEGRADED%s", yellowColor, resetColor)
		case "SHELL":
			statusDisplay = fmt.Sprintf("%s◆ SHELL%s", cyanColor, resetColor)
		default:
			statusDisplay = fmt.Sprintf("%s? %s%s", grayColor, s.Status, resetColor)
		}
		urlDisplay := ""
		if s.URL != "" {
			urlDisplay = fmt.Sprintf("%s%s%s", grayColor, s.URL, resetColor)
		}
		fmt.Printf("  %-40s %s  %s\n", s.Ref, statusDisplay, urlDisplay)
	}
	fmt.Println()
}

type StackLink struct {
	Label string
	URL   string
}

type StackStatus struct {
	Ref    string
	Status string
	URL    string
}
