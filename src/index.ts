#!/usr/bin/env node
import path from "path";
import fs from "fs";
import { Inventory } from "./inventory/Inventory";

const printUsage = () => {
  console.warn("USAGE: npx k2 <k2 inventory file path> <action>");
};

const printVersion = () => {
  console.warn("K2 by github:@tuxounet");
};

const checkParams = () => {
  try {
    const inventory_arg = process.argv[process.argv.length - 2];
    if (!inventory_arg || inventory_arg.endsWith("node")) {
      throw "vous devez preciser un chemin vers un fichier d'inventaire k2";
    }
    if (!fs.existsSync(inventory_arg)) {
      throw "fichier d'inventaire introuvable";
    }
    const stat = fs.statSync(inventory_arg);
    if (!stat.isFile()) {
      throw "chemin d'entrée doit etre un fichier d'inventaire k2";
    }

    const command_arg = process.argv[process.argv.length - 1];
    if (
      command_arg !== undefined &&
      command_arg !== "" &&
      typeof command_arg === "string" &&
      command_arg !== __filename
    ) {
      return {
        inventory_arg,
        command_arg: command_arg.toLowerCase().trim(),
      };
    } else {
      throw "format de la commande incorrect";
    }
  } catch (error) {
    printUsage();
    throw error;
  }
};
printVersion();
const params = checkParams();
console.info("parametres", params);
const run = async () => {
  const full_path = path.resolve(params.inventory_arg);
  const inventory_folder = path.dirname(full_path);
  const inventory_filename = path.basename(full_path);
  const inventory = new Inventory(inventory_filename, inventory_folder);

  await inventory.loadCommands();

  if (
    !Array.from(inventory.allowedCommands.keys()).includes(params.command_arg)
  ) {
    throw "commande incorrecte";
  }

  await inventory.load();

  const handler = inventory.allowedCommands.get(params.command_arg);
  {
    if (!handler) throw "commande introuvable";
  }

  await handler(inventory);
};

run().catch((e) => {
  console.error(e);
  process.exit(1);
});
