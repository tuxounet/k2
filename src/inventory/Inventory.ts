import fs from "fs";
import path from "path";
import fg from "fast-glob";
import jsYaml from "js-yaml";
import { IK2 } from "../types/IK2";
import { IK2Inventory } from "../types/IK2Inventory";
import { IK2Template } from "../types/IK2Template";

import { inventoryKind } from "./kinds";

export class Inventory {
  constructor(inventoryFilename: string, inventoryFolder: string) {
    this.inventory_path = path.resolve(inventoryFolder, inventoryFilename);
    this.inventory = jsYaml.load(
      fs.readFileSync(this.inventory_path, {
        encoding: "utf-8",
      })
    ) as IK2Inventory;
    this.inventory.k2.metadata.folder = path.dirname(this.inventory_path);
    this.inventory.k2.metadata.path = this.inventory_path;
    this.sources = new Map();
    this.templates = new Map();

    this.env = process.env.K2_ENV ?? "local";
  }

  readonly env: string;
  readonly inventory_path: string;
  inventory: IK2Inventory;
  sources: Map<string, IK2>;
  templates: Map<string, IK2Template>;

  async load(): Promise<void> {
    const sourcesGlob = [
      "k2.inventory.yaml",
      ...this.inventory.k2.body.folders.applies,
    ];
    this.sources = await this.loadK2Files(sourcesGlob);
    const inventory = Array.from(this.sources.values()).find(
      (item) => item.k2.metadata.kind === inventoryKind
    );
    if (inventory == null) {
      throw new Error("fichier d'inventaire k2 introuvable");
    }
    this.inventory = inventory as IK2Inventory;

    const templatesGlob = [...this.inventory.k2.body.folders.templates];
    this.templates = await this.loadK2Files<IK2Template>(templatesGlob);
  }

  loadK2File<T extends IK2>(filePath: string): T {
    const item = jsYaml.load(
      fs.readFileSync(filePath, {
        encoding: "utf-8",
      })
    ) as T;
    item.k2.metadata.path = filePath;
    item.k2.metadata.folder = path.dirname(filePath);
    return item;
  }

  private async loadK2Files<T extends IK2>(
    searchGlob: string[]
  ): Promise<Map<string, T>> {
    const sources = new Map<string, T>();

    const sourcesEntries = await fg(searchGlob, {
      onlyFiles: true,
      cwd: this.inventory.k2.metadata.folder,
      ignore: this.inventory.k2.body.folders.ignore,
    });

    sourcesEntries
      .map((item) => {
        return {
          path: path.join(this.inventory.k2.metadata.folder, item),
          body: jsYaml.load(
            fs.readFileSync(
              path.join(this.inventory.k2.metadata.folder, item),
              {
                encoding: "utf-8",
              }
            )
          ) as IK2,
        };
      })
      .map((item) => {
        item.body.k2.metadata.path = item.path;
        item.body.k2.metadata.folder = path.dirname(item.path);
        return item;
      })
      .filter((item) => item.body.k2.metadata.kind)
      .forEach((item) => {
        sources.set(item.body.k2.metadata.id, item.body as T);
      });

    return sources;
  }
}
