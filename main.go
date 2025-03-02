package main

import (
	"context"
	"k2/cmds"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {

	rootCmd := &cli.Command{
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
