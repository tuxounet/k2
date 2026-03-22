package libs

import "fmt"

const (
	resetColor  = "\033[0m"
	redColor    = "\033[31m"
	greenColor  = "\033[32m"
	yellowColor = "\033[33m"
	cyanColor   = "\033[36m"
	grayColor   = "\033[90m"
	boldStyle   = "\033[1m"

	IconApply   = greenColor + "✓" + resetColor
	IconDestroy = redColor + "✗" + resetColor
	IconResolve = cyanColor + "◆" + resetColor
	IconPlan    = yellowColor + "▪" + resetColor
)

func CyanColor() string { return cyanColor }
func GreenCol() string  { return greenColor }
func GrayColor() string { return grayColor }
func BoldStyle() string { return boldStyle }
func ResetCol() string  { return resetColor }
func RedColor() string  { return redColor }
func YellowCol() string { return yellowColor }

func WriteBanner(version string) {
	fmt.Printf("%s%sK2%s %s%s%s by %sKrux%s %s(github.com/tuxounet/k2)%s\n", boldStyle, cyanColor, resetColor, yellowColor, version, resetColor, boldStyle, resetColor, grayColor, resetColor)
}

func WriteTitle(format string, a ...any) {
	fmt.Printf("%s%s▸ %s%s\n", boldStyle, cyanColor, fmt.Sprintf(format, a...), resetColor)
}

func WriteStep(icon string, format string, a ...any) {
	fmt.Printf("  %s %s\n", icon, fmt.Sprintf(format, a...))
}

func WriteDetail(format string, a ...any) {
	fmt.Printf("  %s%s%s\n", grayColor, fmt.Sprintf(format, a...), resetColor)
}

func WriteSubStep(format string, a ...any) {
	fmt.Printf("    %s↳ %s%s\n", grayColor, fmt.Sprintf(format, a...), resetColor)
}

func WriteOutput(line string) {
	fmt.Printf("  %s\n", line)
}

func WriteOutputf(format string, a ...any) {
	fmt.Printf(fmt.Sprintf("  %s", format), a...)
}

func WriteError(err error) error {
	fmt.Printf("  %s✘ %s%s\n", redColor, err.Error(), resetColor)
	return err
}

func WriteErrorString(err string) error {
	errObj := fmt.Errorf("%s", err)
	fmt.Printf("  %s✘ %s%s\n", redColor, errObj, resetColor)
	return errObj
}

func WriteErrorf(format string, a ...any) error {
	err := fmt.Errorf(format, a...)
	fmt.Printf("  %s✘ %s%s\n", redColor, err.Error(), resetColor)
	return err
}
