package cmds

import (
	"context"

	"github.com/tuxounet/k2/libs"
	"github.com/tuxounet/k2/stores"

	"github.com/urfave/cli/v3"
)

var PlanCmd = &cli.Command{
	Name:  "plan",
	Usage: "plan all elements in current inventory folder",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "inventory",
			Value:       "",
			Usage:       "inventory to use",
			Destination: &initialInventoryFile,
		},
	},
	Action: func(context.Context, *cli.Command) error {
		return doPlan()
	},
}

func doPlan() error {
	libs.WriteOutputf("Planning inventory %s\n", initialInventoryFile)

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
	libs.WriteOutputf("PLAN RESULT: %d\n", len(plan.Tasks))
	for _, r := range plan.Tasks {
		libs.WriteOutputf("WILL DO ACTION: %v\n", r)
	}

	return nil
}
