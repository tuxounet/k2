import path from "path";
import { Command } from "commander";
import { getInventory } from "../inventory/inventory";
export default function list(): Command {
  const program = new Command("list");
  program.description("list all elements in current inventory folder");
  program.option(
    "-i, --inventory <value>",
    "inventory file",
    path.join(process.cwd(), "k2.inventory.yaml")
  );

  program.action(async () => {
    console.info("list", program.opts());

    const inventoryPath = program.getOptionValue("inventory");
    const inventory = await getInventory(inventoryPath);
    const tree = inventory.toTree();

    console.info(JSON.stringify(tree, null, 2));
  });
  return program;
}
