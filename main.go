package main

import (
	"context"
	_ "embed"

	"log"
	"os"

	"github.com/tuxounet/k2/cmds"
	"github.com/urfave/cli/v3"
)

//go:embed version.txt
var version string

func main() {

	rootCmd := &cli.Command{
		Name:                  "k2",
		Version:               version,
		EnableShellCompletion: true,
		Description:           "k2 is a template engine",
		Commands: []*cli.Command{
			cmds.PlanCmd,
			cmds.ApplyCmd,
			cmds.DestroyCmd,
		},
	}

	if err := rootCmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}
