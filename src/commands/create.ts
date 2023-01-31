import { IK2Template } from "../types/IK2Template";
import path from "path";
import { IK2Inventory } from "../types/IK2Inventory";
import { IK2Apply } from "../types/IK2Apply";
import { applyTemplate } from "../inventory/template";
import { Command } from "commander";
import { loadK2File } from "../inventory/files";

export default function create(): Command {
  const program = new Command("create");

  program.description("create init k2 inventory folder in current directory");

  program.action(async () => {
    const templateCreateTemplate = path.resolve(
      __dirname,
      "..",
      "create",
      "k2.template.yaml"
    );
    console.info("create", {
      cwd: process.cwd(),
      template: templateCreateTemplate,
    });

    const templateK2 = loadK2File<IK2Template>(templateCreateTemplate);
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
            source: "inventory",

            params: {
              id: "k2.cli.create.template",
            },
          },
        },
      },
    };

    let needReapply = await applyTemplate(
      process.cwd(),
      applyK2,
      inventory,
      templateK2,
      false
    );
    while (needReapply) {
      needReapply = await applyTemplate(
        process.cwd(),
        applyK2,
        inventory,
        templateK2
      );
    }
    console.info("created");
  });
  return program;
}
