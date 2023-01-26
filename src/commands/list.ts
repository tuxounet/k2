import { Inventory } from "../inventory/Inventory";
import path from "path";
import { IK2 } from "../types/IK2";
export default async function list(inventory: Inventory) {
  console.info("list");

  Array.from(inventory.sources.values())

    .map((item) => item as IK2)
    .map((item) => {
      return {
        request: item,
        path: item.k2.metadata.path,
        folder: path.dirname(item.k2.metadata.path),
        entry: inventory.sources.get(item.k2.metadata.id),
      };
    })
    .forEach((item) => {
      item.entry &&
        console.info(
          item.entry.k2.metadata.id,
          "(",
          item.entry.k2.metadata.kind,
          "/",
          item.entry.k2.body.template,
          ")"
        );
    });
}
