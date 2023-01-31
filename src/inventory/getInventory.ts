import path from "path";
import fs from "fs";
import { Inventory } from "./Inventory";
export async function getInventory(inventoryPath: string): Promise<Inventory> {
  if (!fs.existsSync(inventoryPath)) {
    throw new Error(
      "current folder doesn't contains k2.inventor.yaml file at " +
        inventoryPath
    );
  }

  const inventoryFolder = path.dirname(inventoryPath);
  const inventoryFilename = path.basename(inventoryPath);
  const inventory = new Inventory(inventoryFilename, inventoryFolder);

  await inventory.load();

  return inventory;
}
