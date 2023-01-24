import { IK2Template } from "../../types/IK2Template";
import { Inventory } from "../Inventory";
import path from "path";
import fs from "fs";
import fg from "fast-glob";
import ejs from "ejs";
import { IK2Inventory } from "../../types/IK2Inventory";
import { IK2Apply } from "../../types/IK2Apply";
export default async function apply(inventory: Inventory) {
  console.info("apply");

  const all_templates_requests = Array.from(inventory.sources.values())
    .filter((item) => item.k2.metadata.kind === "template-apply")    
    .map((item) => item as IK2Apply)
    .map((item) => {
      return {
        request: item,
        path: item.k2.metadata.path,
        folder: path.dirname(item.k2.metadata.path),

        template: inventory.templates.get(String(item.k2.body.template)) as
          | IK2Template
          | undefined,
      };
    })
    .map((item) => {
      console.info(item)
      return item
    })
    .filter((item) => item.template !== undefined && item.path !== undefined)
    .map(
      (item) =>
        item.template &&
        applyTemplate(
          item.template,
          item.folder,
          item.request,
          inventory.inventory
        )
    );

  await Promise.all(all_templates_requests);
}

async function applyTemplate(
  template: IK2Template,
  destination_folder: string,
  request: IK2Apply,
  inventory: IK2Inventory
) {
  console.info("apply template", template.k2.body.name, destination_folder);
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
        destinationPath: path.join(destination_folder, item),
      };
    })
    .filter((item) => item.sourcePath !== template.k2.metadata.path);
  await Promise.all(
    allCopies
      .filter((item) => item.isDirectory)
      .filter((item) => !fs.existsSync(item.destinationPath))
      .map((item) =>
        fs.promises.mkdir(item.destinationPath, { recursive: true })
      )
  );

  await Promise.all(
    allCopies
      .filter((item) => !item.isDirectory)
      .map((item) => {
        return (async () => {
          try {
            const input = await fs.promises.readFile(item.sourcePath, {
              encoding: "utf-8",
            });
            const output = await ejs.render(
              input,
              {
                ...inventory.k2.body.vars,
                ...request.k2.body.vars,
              },
              { async: true }
            );

            await fs.promises.writeFile(item.destinationPath, output, {
              encoding: "utf-8",
            });
          } catch (e) {
            throw item.sourcePath + "!" + e;
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

  const ignorePath = path.join(destination_folder, ".gitignore");
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
