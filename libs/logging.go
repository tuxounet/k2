package libs

import "fmt"

func WriteOutput(line string) {
	fmt.Printf("OUT: %s\n", line)
}

func WriteOutputf(format string, a ...any) {
	fmt.Printf(fmt.Sprintf("OUT: %s", format), a...)
}

func WriteError(err error) error {
	fmt.Printf("ERROR: %s\n", err.Error())
	return err
}
func WriteErrorString(err string) error {
	errObj := fmt.Errorf("%s", err)
	fmt.Printf("ERROR: %s\n", errObj)
	return errObj
}

func WriteErrorf(format string, a ...any) error {
	err := fmt.Errorf(format, a...)
	fmt.Printf("ERROR: %s\n", err.Error())
	return err
}
