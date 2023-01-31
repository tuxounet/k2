import path from "path";
import { Command } from "commander";
import { INode } from "../types/map/INode";
import { idName, idNs, idRoot, idVersion } from "../helpers/namespace";
import { getInventory } from "../inventory/Inventory";
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

    const allNodes = Array.from(inventory.sources.values()).map((item) => {
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

    const root = idRoot(inventory.inventory.k2.metadata.id);
    const rootNode = allNamespaces.get(root);
    console.info(JSON.stringify(rootNode, null, 2));
  });
  return program;
}
