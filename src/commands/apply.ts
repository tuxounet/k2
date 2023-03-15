import path from "path";
import { IK2Apply } from "../types/IK2Apply";
import { templateApplyKind } from "../inventory/kinds";
import { applyTemplate, resolveTemplate } from "../inventory/template";
import { Command } from "commander";
import { getInventory } from "../inventory/inventory";

export default function apply(): Command {
  const program = new Command("apply");

  program.description("apply all elements in current inventory folder");
  program.option(
    "-i, --inventory <value>",
    "inventory file",
    path.join(process.cwd(), "k2.inventory.yaml")
  );
  program.action(async () => {
    const run = async (): Promise<void> => {
      let reapply = await doApply(program.getOptionValue("inventory"));
      while (reapply) {
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

  const allRequestsQuery = Array.from(inventory.sources.values())
    .filter((item) => item.k2.metadata.kind === templateApplyKind)
    .map((item) => item as IK2Apply)
    .map((item) => {
      return {
        request: item,
        path: item.k2.metadata.path,
        folder: path.dirname(item.k2.metadata.path),
        templateRef: item.k2.body.template,
      };
    });

  const allRequests = [];
  for (const request of allRequestsQuery) {
    const template = await resolveTemplate(
      inventory.inventory.k2.metadata.folder,
      request.templateRef
    );
    allRequests.push({
      ...request,
      template,
    });
  }

  let needReapply = false;
  for (const item of allRequests
    .filter((item) => item.template !== undefined)
    .filter((item) => item.path !== undefined)) {
    const need = await applyTemplate(
      item.folder,
      item.request,
      inventory.inventory,
      item.template
    );
    if (need) needReapply = true;
  }

  return needReapply;
}
