package main

import (
	"context"
	_ "embed"

	"log"
	"os"
	"strings"

	"github.com/tuxounet/k2/cmds"
	"github.com/tuxounet/k2/libs"
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
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			libs.WriteBanner(strings.TrimSpace(version))
			return ctx, nil
		},
		Commands: []*cli.Command{
			cmds.RenderPlanCmd,
			cmds.RenderCmd,
			cmds.UnrenderCmd,
			cmds.StackCmd,
		},
	}

	if err := rootCmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}
