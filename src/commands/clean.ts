import { IK2Template } from "../types/IK2Template";
import { Inventory } from "../inventory/Inventory";
import path from "path";
import fs from "fs";
import fg from "fast-glob";
import childProc from "child_process";
import { IK2Apply } from "../types/IK2Apply";
import { templateApplyKind } from "../inventory/kinds";
import { resolveTemplate } from "../inventory/template";
import { exec } from "../helpers/exec";
import { Command } from "commander";
import { getInventory } from "../inventory/getInventory";
export default function clean(): Command {
  const program = new Command("clean");
  program.description("clean all elements in current inventory folder");
  program.option(
    "-i, --inventory <value>",
    "inventory file",
    path.join(process.cwd(), "k2.inventory.yaml")
  );
  program.action(async () => {
    console.info("clean", program.opts());
    const inventoryPath = program.getOptionValue("inventory");
    const inventory = await getInventory(inventoryPath);

    const allRequests = Array.from(inventory.sources.values())
      .filter((item) => item.k2.metadata.kind === templateApplyKind)
      .map((item) => item as IK2Apply)
      .map((item) => {
        return {
          request: item,
          path: item.k2.metadata.path,
          folder: path.dirname(item.k2.metadata.path),
        };
      })
      .filter((item) => item.path !== undefined)
      .map((item) => cleanTemplate(item.folder));

    await Promise.all(allRequests);
    await cleanupRefs(inventory);
  });
  return program;
}

async function cleanTemplate(destinationFolder: string): Promise<void> {
  console.info("cleaning folder", destinationFolder);
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

async function cleanupRefs(inventory: Inventory): Promise<void> {
  const templateRefs = path.join(
    inventory.inventory.k2.metadata.folder,
    "refs"
  );
  if (fs.existsSync(templateRefs)) {
    exec(`rm -rf ${templateRefs}`, inventory.inventory.k2.metadata.folder);
  }
}
