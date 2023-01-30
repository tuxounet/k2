import path from "path";
import { Command } from "commander";
import { getInventory } from "../inventory/getInventory";
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

    Array.from(inventory.sources.values())
      .map((item) => item)
      .map((item) => {
        return {
          request: item,
          path: item.k2.metadata.path,
          folder: path.dirname(item.k2.metadata.path),
          entry: inventory.sources.get(item.k2.metadata.id),
        };
      })
      .forEach((item) => {
        item.entry != null &&
          console.info(
            item.entry.k2.metadata.id,
            "(",
            item.entry.k2.metadata.kind,
            "/",
            item.entry.k2.metadata.version ?? "latest",
            ")"
          );
      });
  });
  return program;
}
