import fs from "fs";
import path from "path";
import jsYaml from "js-yaml";
import { IK2 } from "../types/IK2";
import { IK2Inventory } from "../types/IK2Inventory";
import { IK2Template } from "../types/IK2Template";

import { inventoryKind } from "./kinds";
import { loadManyK2Files } from "./files";

export async function getInventory(
  sourceInventoryFile: string
): Promise<Inventory> {
  const inventoryFile = path.resolve(sourceInventoryFile);
  if (!fs.existsSync(inventoryFile)) {
    throw new Error(
      "current folder doesn't contains k2.inventor.yaml file at " +
        inventoryFile
    );
  }

  const inventoryFolder = path.dirname(inventoryFile);
  const inventoryFilename = path.basename(inventoryFile);
  const inventory = new Inventory(inventoryFilename, inventoryFolder);

  await inventory.load();

  return inventory;
}

export class Inventory {
  constructor(inventoryFilename: string, inventoryFolder: string) {
    this.inventoryFilePath = path.resolve(inventoryFolder, inventoryFilename);
    this.inventory = jsYaml.load(
      fs.readFileSync(this.inventoryFilePath, {
        encoding: "utf-8",
      })
    ) as IK2Inventory;
    this.inventory.k2.metadata.folder = path.dirname(this.inventoryFilePath);
    this.inventory.k2.metadata.path = this.inventoryFilePath;
    this.sources = new Map();
    this.templates = new Map();

    this.env = process.env.K2_ENV ?? "local";
  }

  readonly env: string;
  readonly inventoryFilePath: string;
  inventory: IK2Inventory;
  sources: Map<string, IK2>;
  templates: Map<string, IK2Template>;

  async load(): Promise<void> {
    const sourcesGlob = [
      "k2.inventory.yaml",
      ...this.inventory.k2.body.folders.applies,
    ];
    this.sources = await this.mapK2Files(sourcesGlob);
    const inventory = Array.from(this.sources.values()).find(
      (item) => item.k2.metadata.kind === inventoryKind
    );
    if (inventory == null) {
      throw new Error("fichier d'inventaire k2 introuvable");
    }
    this.inventory = inventory as IK2Inventory;

    const templatesGlob = [...this.inventory.k2.body.folders.templates];
    this.templates = await this.mapK2Files<IK2Template>(templatesGlob);
  }

  private async mapK2Files<T extends IK2>(
    searchGlob: string[]
  ): Promise<Map<string, T>> {
    const sources = new Map<string, T>();

    const sourcesEntries = await loadManyK2Files(
      searchGlob,
      this.inventory.k2.metadata.folder,
      this.inventory.k2.body.folders.ignore
    );
    sourcesEntries.forEach((item) => {
      sources.set(item.k2.metadata.id, item as T);
    });

    return sources;
  }
}
