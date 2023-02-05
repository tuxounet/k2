import { getInventory, Inventory } from "../inventory/Inventory";
import path from "path";
import fs from "fs";
import fg from "fast-glob";
import { IK2Apply } from "../types/IK2Apply";
import { templateApplyKind } from "../inventory/kinds";
import { exec } from "../helpers/exec";
import { Command } from "commander";

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
      .map(async (item) => await cleanTemplate(item.folder));

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
  const ops = gitIgnoreContent.split("\n").map(async (item) => {
    if (path.basename(item).startsWith("!")) return;
    const targetContent = path.join(destinationFolder, item);
    await exec("rm -rf " + targetContent, destinationFolder);
  });
  await Promise.all(ops);
  await exec("rm -rf .gitignore", destinationFolder);
  await exec("find . -empty -type d -delete", destinationFolder);
}

async function cleanupRefs(inventory: Inventory): Promise<void> {
  const templateRefs = path.join(
    inventory.inventory.k2.metadata.folder,
    "refs"
  );
  if (fs.existsSync(templateRefs)) {
    await exec(
      `rm -rf ${templateRefs}`,
      inventory.inventory.k2.metadata.folder
    );
  }
}
