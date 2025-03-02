package libs

import (
	"os"
	"os/exec"
)

func ExecCommand(script string, cwd string, tplData any) error {

	_, err := RenderTemplate(script, tplData)
	if err != nil {
		return err
	}

	cmd := exec.Command("/bin/sh", "-c", script)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = cwd

	err = cmd.Run()

	if err != nil {
		return err
	}

	return nil
}
