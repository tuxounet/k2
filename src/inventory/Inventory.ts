import fs from "fs";
import path from "path";
import fg from "fast-glob";
import jsYaml from "js-yaml";
import { IK2 } from "../types/IK2";
import { IK2Inventory } from "../types/IK2Inventory";
import { IK2Template } from "../types/IK2Template";

export class Inventory {
  constructor(
    readonly inventory_filename: string,
    readonly inventory_folder: string
  ) {
    this.inventory_path = path.join(inventory_folder, inventory_filename);
    this.inventory = jsYaml.load(
      fs.readFileSync(
        path.join(this.inventory_folder, this.inventory_filename),
        {
          encoding: "utf-8",
        }
      )
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
  async load() {
    const sources_glob = [
      this.inventory_filename,
      ...this.inventory.k2.body.folders.sources,
    ];
    this.sources = await this.loadK2Files(sources_glob);
    const inventory = Array.from(this.sources.values()).find(
      (item) => item.k2.metadata.kind === "k2.inventory"
    );
    if (!inventory) {
      throw "fichier d'inventaire k2 introuvable";
    }
    this.inventory = inventory as IK2Inventory;

    const templates_glob = [...this.inventory.k2.body.folders.templates];
    this.templates = await this.loadK2Files<IK2Template>(templates_glob);
  }

  async loadCommands() {
    const search_glob = ["commands/*.js"];
    const entries = await fg(search_glob, {
      onlyFiles: true,
      cwd: __dirname,
    });
    this.allowedCommands = new Map();
    entries
      .map((item) => {
        return {
          source: path.join(__dirname, item),
          require: path.join(__dirname, item.replace(".js", "")),
          action: path.basename(item).replace(".js", ""),
        };
      })
      .map((item) => {
        return {
          ...item,
          handler: require(item.require).default,
        };
      })
      .filter((item) => item.handler !== undefined)
      .forEach((item) => {
        this.allowedCommands.set(
          item.action.toLowerCase().trim(),
          item.handler
        );
      });
  }

  private async loadK2Files<T extends IK2>(search_glob: string[]) {
    const sources = new Map<string, T>();
    

    const sources_entries = await fg(search_glob, {
      onlyFiles: true,
      cwd: this.inventory_folder,
      ignore: this.inventory.k2.body.folders.ignore,
    });

    sources_entries
      .map((item) => {
        return {
          path: path.join(this.inventory_folder, item),
          body: jsYaml.load(
            fs.readFileSync(path.join(this.inventory_folder, item), {
              encoding: "utf-8",
            })
          ) as IK2,
        };
      })
      .map((item) => {
        item.body.k2.metadata.path = item.path;
        item.body.k2.metadata.folder = path.dirname(item.path);
        console.dir(item)
        return item;
      })
      .filter(
        (item) =>
          item.body &&
          item.body.k2 &&
          item.body.k2.metadata &&
          item.body.k2.metadata.id &&
          item.body.k2.metadata.kind
      )
      .forEach((item) => {
        sources.set(item.body.k2.metadata.id, item.body as T);
      });

    return sources;
  }
}
