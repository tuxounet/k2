import { IK2Template } from "../types/IK2Template";
import { Inventory } from "../inventory/Inventory";
import path from "path";
import fs from "fs";
import fg from "fast-glob";
import childProc from "child_process";
import { IK2Apply } from "../types/IK2Apply";
import { templateApplyKind } from "../inventory/kinds";
import { resolveTemplate } from "../inventory/template";
export default async function clean(inventory: Inventory): Promise<void> {
  console.info("clean");

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
    .filter((item) => item.template !== undefined && item.path !== undefined)
    .map(
      (item) =>
        item.template != null && cleanTemplate(item.template, item.folder)
    );

  await Promise.all(allRequests);
}

async function cleanTemplate(
  template: IK2Template,
  destinationFolder: string
): Promise<void> {
  console.info("cleaning folder", template.k2.body.name, destinationFolder);
  const allTemplatedFiles = await fg(["**/*", "**/.*"], {
    markDirectories: true,
    onlyFiles: false,
    dot: true,
    cwd: destinationFolder,
  });

  if (!allTemplatedFiles.includes(".gitignore")) {
    return;
  }
  const gitIgnorePath = path.join(destinationFolder, ".gitignore");

  const gitIgnoreContent = await fs.promises.readFile(gitIgnorePath, "utf-8");
  gitIgnoreContent.split("\n").forEach((item) => {
    const targetContent = path.join(destinationFolder, item);
    childProc.execSync("rm -rf " + targetContent);
  });
  childProc.execSync("rm -rf .gitignore", { cwd: destinationFolder });
  childProc.execSync("find . -empty -type d -delete", {
    cwd: destinationFolder,
  });
}
