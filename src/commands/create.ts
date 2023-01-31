import path from "path";
import { IK2Inventory } from "../types/IK2Inventory";
import { IK2Apply } from "../types/IK2Apply";
import { applyTemplate, resolveTemplate } from "../inventory/template";
import { Command } from "commander";
import { loadK2File } from "../inventory/files";

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

    const applyFilePath = path.join(__dirname, "..", "create", "k2.apply.yaml");

    const applyK2 = loadK2File<IK2Apply>(applyFilePath);

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
