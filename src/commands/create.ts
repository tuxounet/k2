import path from "path";
import { IK2Inventory } from "../types/IK2Inventory";
import { IK2Apply } from "../types/IK2Apply";
import { applyTemplate, resolveTemplate } from "../inventory/template";
import { Command } from "commander";

export default function create(): Command {
  const program = new Command("create");

  program.description("create init k2 inventory folder in current directory");

  program.action(async () => {
    console.info("create", {
      cwd: process.cwd(),
    });

    const inventory: IK2Inventory = {
      k2: {
        metadata: {
          id: "k2.cli.create.inventory",
          kind: "inventory",
          folder: process.cwd(),
          path: path.resolve(process.cwd(), "k2.inventory.yaml"),
        },
        body: {
          folders: {
            applies: [],
            templates: [],
            ignore: [],
          },
        },
      },
    };

    const applyK2: IK2Apply = {
      k2: {
        metadata: {
          id: "k2.cli.create.init",
          kind: "template-apply",
          folder: process.cwd(),
          path: path.resolve(process.cwd(), "k2.apply.yaml"),
        },
        body: {
          scripts: {},
          template: {
            source: "git",
            params: {
              repository:
                "https://github.com/tuxounet/k2-blocks-structures.git",
              branch: "dev",
              path: "root/create-root/k2.template.yaml",
            },
          },
        },
      },
    };
    const template = resolveTemplate(process.cwd(), applyK2.k2.body.template);

    let needReapply = await applyTemplate(
      process.cwd(),
      applyK2,
      inventory,
      Promise.resolve(template),
      false
    );
    while (needReapply) {
      needReapply = await applyTemplate(
        process.cwd(),
        applyK2,
        inventory,
        Promise.resolve(template)
      );
    }
    console.info("created");
  });
  return program;
}
