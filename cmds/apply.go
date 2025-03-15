package cmds

import (
	"context"

	"github.com/tuxounet/k2/libs"
	"github.com/tuxounet/k2/stores"
	"github.com/urfave/cli/v3"
)

var ApplyCmd = &cli.Command{
	Name:  "apply",
	Usage: "apply all elements in current inventory folder",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "inventory",
			Value:       "",
			Usage:       "inventory to use",
			Destination: &initialInventoryFile,
		},
	},
	Action: func(context.Context, *cli.Command) error {
		return doApply()
	},
}

func doApply() error {
	libs.WriteOutputf("Applying inventory %s\n", initialInventoryFile)

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

	err = inventory.Apply(plan)
	if err != nil {
		return err
	}

	return nil
}
