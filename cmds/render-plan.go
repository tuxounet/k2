package cmds

import (
	"context"

	"github.com/tuxounet/k2/libs"
	"github.com/tuxounet/k2/stores"

	"github.com/urfave/cli/v3"
)

var RenderPlanCmd = &cli.Command{
	Name:  "render-plan",
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
		return doRenderPlan()
	},
}

func doRenderPlan() error {
	libs.WriteTitle("Plan %s", initialInventoryFile)

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
	libs.WriteDetail("%d actions planned", len(plan.Tasks))
	for _, r := range plan.Tasks {
		libs.WriteStep(libs.IconPlan, "%v", r)
	}

	return nil
}
