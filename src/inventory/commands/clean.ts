import { IK2Template } from "../../types/IK2Template";
import { Inventory } from "../Inventory";
import path from "path";
import fs from "fs";
import fg from "fast-glob";
import childProc from "child_process";
import { IK2Inventory } from "../../types/IK2Inventory";
import { IK2Apply } from "../../types/IK2Apply";
export default async function clean(inventory: Inventory) {
  console.info("clean");

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
    .filter((item) => item.template !== undefined && item.path !== undefined)
    .map((item) => item.template && cleanTemplate(item.template, item.folder));

  await Promise.all(all_templates_requests);
}

async function cleanTemplate(
  template: IK2Template,
  destination_folder: string
) {
  console.info("cleaning folder", template.k2.body.name, destination_folder);
  const allTemplatedFiles = await fg(["**/*", "**/.*"], {
    markDirectories: true,
    onlyFiles: false,
    dot: true,
    cwd: destination_folder,
  });

  if (!allTemplatedFiles.includes(".gitignore")) {
    return;
  }
  const gitIgnorePath = path.join(destination_folder, ".gitignore");

  const gitIgnoreContent = await fs.promises.readFile(gitIgnorePath, "utf-8");
  gitIgnoreContent.split("\n").map((item) => {
    const targetContent = path.join(destination_folder, item);
    childProc.execSync("rm -rf " + targetContent);
  });
  childProc.execSync("rm -rf .gitignore", { cwd: destination_folder });
  childProc.execSync("find . -empty -type d -delete", {
    cwd: destination_folder,
  });
}
