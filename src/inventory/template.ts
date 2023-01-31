import { IK2Template } from "../types/IK2Template";
import { IK2TemplateRef } from "../types/templates/IK2TemplateRef";
import { IK2TemplateRefGitParams } from "../types/templates/IK2TemplateRefGitParams";
import { IK2TemplateRefInventoryParams } from "../types/templates/IK2TemplateRefInventoryParams";
import { Inventory } from "./Inventory";
import fs from "fs";
import path from "path";
import fg from "fast-glob";
import ejs from "ejs";
import { md5 } from "../helpers/hash";
import { exec } from "../helpers/exec";
import { IK2Apply } from "../types/IK2Apply";
import { IK2Inventory } from "../types/IK2Inventory";
import { executeScript } from "./scripts";
export function resolveTemplate(
  inventory: Inventory,
  ref: IK2TemplateRef
): IK2Template {
  switch (ref.source) {
    case "inventory": {
      const param = ref.params as IK2TemplateRefInventoryParams;
      const template = resolveInventoryTemplate(inventory, param);
      if (template == null) {
        throw new Error(`template introuvable dans l'inventaire ${param.id}`);
      }
      return template;
    }
    case "git": {
      const param = ref.params as IK2TemplateRefGitParams;
      const template = resolveGitTemplate(inventory, param);
      if (template == null) {
        throw new Error(`template non résolu ${param.repository}`);
      }
      return template;
    }
    default:
      throw new Error(`source de template non trouvé ${String(ref.source)}`);
  }
}

export async function applyTemplate(
  destinationFolder: string,
  request: IK2Apply,
  inventory: IK2Inventory,
  template: IK2Template,
  produceGitIgnore: boolean = true
): Promise<boolean> {
  console.info("apply template", destinationFolder);

  await executeScript(template, "bootstrap", destinationFolder);
  await executeScript(request, "bootstrap", destinationFolder);

  await executeScript(template, "pre", destinationFolder);
  await executeScript(request, "pre", destinationFolder);

  const allTemplateFiles = await fg(["**/*", "**/.gitignore"], {
    ignore: ["k2.template.yaml"],
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
  if (produceGitIgnore) {
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

      const appendContent = ignoreContent.filter(
        (item) => !lines.includes(item)
      );
      if (appendContent.length > 0) {
        fs.appendFileSync(ignorePath, "\n" + appendContent.join("\n"), {
          encoding: "utf-8",
        });
      }
    }
  }

  await executeScript(template, "post", destinationFolder);
  await executeScript(request, "post", destinationFolder);

  return notExistingSubApplies.length > 0;
}

function resolveInventoryTemplate(
  inventory: Inventory,
  refParams: IK2TemplateRefInventoryParams
): IK2Template {
  console.info("resolveInventoryTemplate", refParams);
  const template = inventory.templates.get(refParams.id);
  if (template === undefined) {
    throw new Error(
      "Impossible de trouver le template ayant l'id " + refParams.id
    );
  }
  return template;
}

function resolveGitTemplate(
  inventory: Inventory,
  refParams: IK2TemplateRefGitParams
): IK2Template {
  console.info("resolveGitTemplate", refParams);

  const id = md5(
    JSON.stringify({
      repository: refParams.repository,
      branch: refParams.branch,
    })
  );

  const templateRefPath = path.join(
    inventory.inventory.k2.metadata.folder,
    "refs",
    id
  );
  if (!fs.existsSync(templateRefPath)) {
    exec(
      `git clone  ${
        refParams.branch !== undefined
          ? `--branch ${refParams.branch} --single-branch`
          : ""
      } ${refParams.repository} ${templateRefPath}`,
      inventory.inventory.k2.metadata.folder
    );
  } else {
    exec(`git pull`, templateRefPath);
  }

  const targetPath = templateRefPath + "/" + refParams.path;

  if (!fs.existsSync(targetPath)) {
    throw new Error(
      "Impossible de trouver le fichier de template " + targetPath
    );
  }

  const templateFile = inventory.loadK2File<IK2Template>(targetPath);
  return templateFile;
}
