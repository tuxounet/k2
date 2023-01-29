import { Inventory } from "../inventory/Inventory";
import path from "path";
export default async function list(inventory: Inventory): Promise<void> {
  console.info("list");

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
          item.entry.k2.metadata.version || "latest",
          ")"
        );
    });
}
