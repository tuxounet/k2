import { IK2Template } from "../types/IK2Template";
import { Inventory } from "../inventory/Inventory";
import path from "path";
import fs from "fs";
import fg from "fast-glob";
import ejs from "ejs";
import { IK2Inventory } from "../types/IK2Inventory";
import { IK2Apply } from "../types/IK2Apply";
import { templateApplyKind } from "../inventory/kinds";
import { resolveTemplate } from "../inventory/template";
export default async function apply(inventory: Inventory): Promise<void> {
  console.info("apply");

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
    .map(async (item) => {
      return await applyTemplate(
        item.folder,
        item.request,
        inventory.inventory,
        item.template
      );
    });

  await Promise.all(allRequests);
}

async function applyTemplate(
  destinationFolder: string,
  request: IK2Apply,
  inventory: IK2Inventory,
  template: IK2Template
): Promise<void> {
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
}
