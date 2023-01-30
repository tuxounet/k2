import { IK2Template } from "../types/IK2Template";
import path from "path";
import fs from "fs";
import fg from "fast-glob";
import ejs from "ejs";
import { IK2Inventory } from "../types/IK2Inventory";
import { IK2Apply } from "../types/IK2Apply";
import { templateApplyKind } from "../inventory/kinds";
import { resolveTemplate } from "../inventory/template";
import { Command } from "commander";
import { getInventory } from "../inventory/getInventory";

export default function apply(): Command {
  const program = new Command("apply");

  program.description("apply all elements in current inventory folder");
  program.option(
    "-i, --inventory <value>",
    "inventory file",
    path.join(process.cwd(), "k2.inventory.yaml")
  );
  program.action(async (inventoryPath: string) => {
    const run = async () => {
      let reapply = await doApply(program.getOptionValue("inventory"));
      while (reapply === true) {
        reapply = await doApply(program.getOptionValue("inventory"));
      }
    };
    run().catch((e) => {
      console.error(e);
      process.exit(1);
    });
  });
  return program;
}

async function doApply(inventoryPath: string): Promise<boolean> {
  console.debug("do apply");
  const inventory = await getInventory(inventoryPath);

  const allRequests = Array.from(inventory.sources.values())
    .filter((item) => item.k2.metadata.kind === templateApplyKind)
    .map((item) => item as IK2Apply)
    .map((item) => {
      return {
        request: item,
        path: item.k2.metadata.path,
        folder: path.dirname(item.k2.metadata.path),
        template: resolveTemplate(inventory, item.k2.body.template),
      };
    })
    .filter((item) => item.template !== undefined)
    .filter((item) => item.path !== undefined)
    .map(
      async (item) =>
        await applyTemplate(
          item.folder,
          item.request,
          inventory.inventory,
          item.template
        )
    );

  const results = await Promise.all(allRequests);
  if (results.filter((item) => item === true).length > 0) {
    console.warn("need reapply");
    return true;
  }
  return false;
}

async function applyTemplate(
  destinationFolder: string,
  request: IK2Apply,
  inventory: IK2Inventory,
  template: IK2Template
): Promise<boolean> {
  console.info("apply template", destinationFolder);

  const allTemplateFiles = await fg(["**/*", "**/.gitignore"], {
    markDirectories: true,
    onlyFiles: false,
    cwd: template.k2.metadata.folder,
  });

  const allCopies = allTemplateFiles
    .map((item) => {
      return {
        item,
        filename: path.basename(item),
        sourcePath: path.join(template.k2.metadata.folder, item),
        isDirectory: item.endsWith("/"),
        destinationPath: path.join(destinationFolder, item),
      };
    })
    .filter((item) => item.sourcePath !== template.k2.metadata.path);

  await Promise.all(
    allCopies
      .filter((item) => item.isDirectory)
      .filter((item) => !fs.existsSync(item.destinationPath))
      .map(
        async (item) =>
          await fs.promises.mkdir(item.destinationPath, { recursive: true })
      )
  );

  const notExistingSubApplies = allCopies
    .filter((item) => !item.isDirectory)
    .filter((item) => item.filename.trim().toLowerCase() === "k2.apply.yaml")
    .filter((item) => !fs.existsSync(item.destinationPath));

  await Promise.all(
    allCopies
      .filter((item) => !item.isDirectory)
      .map(async (item) => {
        return await (async () => {
          try {
            const input = await fs.promises.readFile(item.sourcePath, {
              encoding: "utf-8",
            });
            const output = await ejs.render(
              input,
              {
                ...template.k2.body.parameters,
                ...inventory.k2.body.vars,
                ...request.k2.body.vars,
              },
              { async: true }
            );

            await fs.promises.writeFile(item.destinationPath, output, {
              encoding: "utf-8",
            });
          } catch (e) {
            throw new Error(item.sourcePath + "!" + String(e));
          }
        })();
      })
  );

  const ignoreContent = [];
  ignoreContent.push("!" + path.basename(request.k2.metadata.path));
  ignoreContent.push(
    ...allTemplateFiles
      .filter((item) => !item.endsWith("/"))
      .filter((item) => item !== ".gitignore")
  );

  const ignorePath = path.join(destinationFolder, ".gitignore");
  if (!fs.existsSync(ignorePath)) {
    fs.writeFileSync(ignorePath, ignoreContent.join("\n"), {
      encoding: "utf-8",
    });
  } else {
    const body = fs.readFileSync(ignorePath, { encoding: "utf-8" });
    const lines = body.split("\n");

    const appendContent = ignoreContent.filter((item) => !lines.includes(item));
    if (appendContent.length > 0) {
      fs.appendFileSync(ignorePath, "\n" + appendContent.join("\n"), {
        encoding: "utf-8",
      });
    }
  }

  return notExistingSubApplies.length > 0;
}
