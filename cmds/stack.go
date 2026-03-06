package cmds

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/tuxounet/k2/stores"
	"github.com/urfave/cli/v3"
)

var stackInventoryFile string
var stackDebug bool

func getStackRootDir() string {
	if stackInventoryFile != "" {
		abs, err := filepath.Abs(stackInventoryFile)
		if err != nil {
			return "."
		}
		return filepath.Dir(abs)
	}
	return "."
}

func newStack(cmd *cli.Command) (*stores.StackStore, error) {
	stackName := cmd.Args().First()
	if stackName == "" {
		return nil, fmt.Errorf("stack name required")
	}
	return stores.NewStackStore(getStackRootDir(), stackName, stackDebug)
}

var stackFlags = []cli.Flag{
	&cli.StringFlag{
		Name:        "inventory",
		Value:       "",
		Usage:       "inventory file path (used to determine project root)",
		Destination: &stackInventoryFile,
	},
	&cli.BoolFlag{
		Name:        "debug",
		Usage:       "enable debug mode",
		Destination: &stackDebug,
	},
}

var StackCmd = &cli.Command{
	Name:    "stack",
	Aliases: []string{"s"},
	Usage:   "manage stacks of services (up, down, status, logs, ...)",
	Flags:   stackFlags,
	Commands: []*cli.Command{
		stackUpCmd,
		stackDownCmd,
		stackRestartCmd,
		stackStatusCmd,
		stackLogsCmd,
		stackHealthcheckCmd,
		stackShellCmd,
		stackUrlsCmd,
		stackRunCmd,
		stackListCmd,
		stackLayersCmd,
	},
}

var stackUpCmd = &cli.Command{
	Name:  "up",
	Usage: "render + start all layers in a stack",
	Flags: stackFlags,
	Action: func(_ context.Context, cmd *cli.Command) error {
		s, err := newStack(cmd)
		if err != nil {
			return err
		}
		return s.Up()
	},
}

var stackDownCmd = &cli.Command{
	Name:  "down",
	Usage: "stop all layers in a stack (reverse order)",
	Flags: stackFlags,
	Action: func(_ context.Context, cmd *cli.Command) error {
		s, err := newStack(cmd)
		if err != nil {
			return err
		}
		return s.Down()
	},
}

var stackRestartCmd = &cli.Command{
	Name:  "restart",
	Usage: "restart a stack (down then up)",
	Flags: stackFlags,
	Action: func(_ context.Context, cmd *cli.Command) error {
		s, err := newStack(cmd)
		if err != nil {
			return err
		}
		return s.Restart()
	},
}

var stackStatusCmd = &cli.Command{
	Name:  "status",
	Usage: "show status of each layer in a stack",
	Flags: stackFlags,
	Action: func(_ context.Context, cmd *cli.Command) error {
		s, err := newStack(cmd)
		if err != nil {
			return err
		}
		return s.Status()
	},
}

var stackLogsCmd = &cli.Command{
	Name:  "logs",
	Usage: "show logs of a stack [optional: specific layer]",
	Flags: stackFlags,
	Action: func(_ context.Context, cmd *cli.Command) error {
		stackName := cmd.Args().First()
		if stackName == "" {
			return fmt.Errorf("stack name required")
		}
		s, err := stores.NewStackStore(getStackRootDir(), stackName, stackDebug)
		if err != nil {
			return err
		}
		targetLayer := cmd.Args().Get(1)
		return s.Logs(targetLayer)
	},
}

var stackHealthcheckCmd = &cli.Command{
	Name:  "healthcheck",
	Usage: "check health of each layer in a stack",
	Flags: stackFlags,
	Action: func(_ context.Context, cmd *cli.Command) error {
		s, err := newStack(cmd)
		if err != nil {
			return err
		}
		return s.Healthcheck()
	},
}

var stackShellCmd = &cli.Command{
	Name:  "shell",
	Usage: "open a shell in a stack layer",
	Flags: stackFlags,
	Action: func(_ context.Context, cmd *cli.Command) error {
		stackName := cmd.Args().First()
		if stackName == "" {
			return fmt.Errorf("stack name required")
		}
		s, err := stores.NewStackStore(getStackRootDir(), stackName, stackDebug)
		if err != nil {
			return err
		}
		targetLayer := cmd.Args().Get(1)
		return s.Shell(targetLayer)
	},
}

var stackUrlsCmd = &cli.Command{
	Name:  "urls",
	Usage: "show access URLs for a stack",
	Flags: stackFlags,
	Action: func(_ context.Context, cmd *cli.Command) error {
		s, err := newStack(cmd)
		if err != nil {
			return err
		}
		return s.Urls()
	},
}

var stackRunCmd = &cli.Command{
	Name:  "run",
	Usage: "execute a verb on a specific layer: run <stack> <layer> [verb] [args...]",
	Flags: stackFlags,
	Action: func(_ context.Context, cmd *cli.Command) error {
		stackName := cmd.Args().First()
		if stackName == "" {
			return fmt.Errorf("usage: k2 stack run <stack> <layer> [verb] [args...]")
		}
		targetLayer := cmd.Args().Get(1)
		if targetLayer == "" {
			return fmt.Errorf("usage: k2 stack run <stack> <layer> [verb] [args...]")
		}
		s, err := stores.NewStackStore(getStackRootDir(), stackName, stackDebug)
		if err != nil {
			return err
		}
		verb := cmd.Args().Get(2)
		var args []string
		if cmd.Args().Len() > 3 {
			args = cmd.Args().Slice()[3:]
		}
		return s.Run(targetLayer, verb, args)
	},
}

var stackListCmd = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Usage:   "list available stacks",
	Flags:   stackFlags,
	Action: func(_ context.Context, _ *cli.Command) error {
		return stores.ListStacks(getStackRootDir())
	},
}

var stackLayersCmd = &cli.Command{
	Name:  "layers",
	Usage: "list available layers",
	Flags: stackFlags,
	Action: func(_ context.Context, _ *cli.Command) error {
		return stores.ListLayers(getStackRootDir())
	},
}
