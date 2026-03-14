package cmds

import (
	"context"
	"fmt"

	"github.com/tuxounet/k2/libs"
	"github.com/tuxounet/k2/stores"
	"github.com/urfave/cli/v3"
)

var stacksInventoryFile string

var StacksCmd = &cli.Command{
	Name:    "stacks",
	Aliases: []string{"ls"},
	Usage:   "list available stacks from the inventory",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "inventory",
			Value:       "",
			Usage:       "inventory file path",
			Destination: &stacksInventoryFile,
		},
	},
	Action: func(_ context.Context, _ *cli.Command) error {
		return doListStacks()
	},
}

func doListStacks() error {
	if stacksInventoryFile == "" {
		stacksInventoryFile = "./k2.inventory.yaml"
	}

	inventory, err := stores.NewInventory(stacksInventoryFile)
	if err != nil {
		return err
	}

	stacks, err := inventory.ListStacks()
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf("  %sStacks disponibles :%s\n", libs.BoldStyle(), libs.ResetCol())
	fmt.Println("  ────────────────────────────────────────")

	if len(stacks) == 0 {
		fmt.Printf("  %sAucune stack trouvée%s\n", libs.GrayColor(), libs.ResetCol())
	} else {
		for _, s := range stacks {
			if s.Description != "" {
				fmt.Printf("  %s%s%s  %s— %s (%d layers)%s\n",
					libs.CyanColor(), s.Name, libs.ResetCol(),
					libs.GrayColor(), s.Description, s.LayerCount, libs.ResetCol())
			} else {
				fmt.Printf("  %s%s%s  %s(%d layers)%s\n",
					libs.CyanColor(), s.Name, libs.ResetCol(),
					libs.GrayColor(), s.LayerCount, libs.ResetCol())
			}
		}
	}
	fmt.Println()
	return nil
}
