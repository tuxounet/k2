import fs from "fs";
import path from "path";
import jsYaml from "js-yaml";
import { IK2 } from "../types/IK2";
import { IK2Inventory } from "../types/IK2Inventory";
import { IK2Template } from "../types/IK2Template";

import { inventoryKind } from "./kinds";
import { loadManyK2Files } from "./files";
import { INode } from "../types/map/INode";
import { idName, idNs, idRoot, idVersion } from "../helpers/namespace";

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

  toTree(): INode {
    const allNodes = Array.from(this.sources.values()).map((item) => {
      const node: INode = {
        id: item.k2.metadata.id,
        name: idName(item.k2.metadata.id),
        namespace: idNs(item.k2.metadata.id),
        kind: item.k2.metadata.kind,
        version: idVersion(item.k2.metadata.version),
      };

      return node;
    });

    const allNamespaces = new Map<string, INode>();

    const upsertNamespace = (ns: string): INode => {
      const foundNs = allNamespaces.get(ns);
      if (foundNs === undefined) {
        const node: INode = {
          id: ns,
          name: idName(ns),
          namespace: idNs(ns),
          kind: "namespace",
          version: "latest",
          childs: [],
        };
        allNamespaces.set(ns, node);

        if (node.namespace.trim() !== "" && node.name.trim() !== "") {
          const nsNode = upsertNamespace(node.namespace);
          if (nsNode.childs == null) nsNode.childs = [];
          nsNode.childs.push(node);
        }

        return node;
      }
      return foundNs;
    };

    allNodes.forEach((item) => {
      const nsNode = upsertNamespace(item.namespace);
      if (nsNode.childs == null) nsNode.childs = [];
      nsNode.childs.push(item);
    });

    const root = idRoot(this.inventory.k2.metadata.id);
    const rootNode =
      allNamespaces.get(root) ?? Array.from(allNamespaces.values())[0];

    return rootNode;
  }
}
