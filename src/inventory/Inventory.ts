import fs from "fs";
import path from "path";
import fg from "fast-glob";
import jsYaml from "js-yaml";
import { IK2 } from "../types/IK2";
import { IK2Inventory } from "../types/IK2Inventory";
import { IK2Template } from "../types/IK2Template";
import applyCommand from "../commands/apply";
import cleanCommand from "../commands/clean";

export class Inventory {
  constructor(
    readonly inventoryFilename: string,
    readonly inventoryFolder: string
  ) {
    this.inventory_path = path.join(inventoryFolder, inventoryFilename);
    this.inventory = jsYaml.load(
      fs.readFileSync(path.join(this.inventoryFolder, this.inventoryFilename), {
        encoding: "utf-8",
      })
    ) as IK2Inventory;
    this.sources = new Map();
    this.templates = new Map();
    this.allowedCommands = new Map();

    this.env = process.env.K2_ENV ?? "local";
  }

  readonly env: string;
  readonly inventory_path: string;
  inventory: IK2Inventory;
  sources: Map<string, IK2>;
  templates: Map<string, IK2Template>;
  allowedCommands: Map<string, Function>;
  async load(): Promise<void> {
    const sourcesGlob = [
      this.inventoryFilename,
      ...this.inventory.k2.body.folders.sources,
    ];
    this.sources = await this.loadK2Files(sourcesGlob);
    const inventory = Array.from(this.sources.values()).find(
      (item) => item.k2.metadata.kind === "k2.inventory"
    );
    if (inventory == null) {
      throw new Error("fichier d'inventaire k2 introuvable");
    }
    this.inventory = inventory as IK2Inventory;

    const templatesGlob = [...this.inventory.k2.body.folders.templates];
    this.templates = await this.loadK2Files<IK2Template>(templatesGlob);
  }

  async loadCommands(): Promise<void> {
    this.allowedCommands = new Map();
    this.allowedCommands.set("apply", applyCommand);
    this.allowedCommands.set("clean", cleanCommand);
    this.allowedCommands.set("list", cleanCommand);
  }

  private async loadK2Files<T extends IK2>(
    searchGlob: string[]
  ): Promise<Map<string, T>> {
    const sources = new Map<string, T>();

    const sourcesEntries = await fg(searchGlob, {
      onlyFiles: true,
      cwd: this.inventoryFolder,
      ignore: this.inventory.k2.body.folders.ignore,
    });

    sourcesEntries
      .map((item) => {
        return {
          path: path.join(this.inventoryFolder, item),
          body: jsYaml.load(
            fs.readFileSync(path.join(this.inventoryFolder, item), {
              encoding: "utf-8",
            })
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
