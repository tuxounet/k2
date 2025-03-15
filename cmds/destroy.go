package cmds

import (
	"context"

	"github.com/tuxounet/k2/libs"
	"github.com/tuxounet/k2/stores"

	"github.com/urfave/cli/v3"
)

var DestroyCmd = &cli.Command{
	Name:  "destroy",
	Usage: "destroy the generated files",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "inventory",
			Value:       "",
			Usage:       "inventory to use",
			Destination: &initialInventoryFile,
		},
	},
	Action: func(context.Context, *cli.Command) error {
		return doDestroy()
	},
}

func doDestroy() error {
	libs.WriteOutputf("Destroy inventory %s\n", initialInventoryFile)

	if initialInventoryFile == "" {
		initialInventoryFile = "./k2.inventory.yaml"
	}

	inventory, err := stores.NewInventory(initialInventoryFile)
	if err != nil {
		return err
	}

	plan, err := inventory.Plan()
	if err != nil {
		return err
	}

	err = inventory.Destroy(plan)
	if err != nil {
		return err
	}

	return nil
}
