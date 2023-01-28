#!/usr/bin/env node
import path from "path";
import fs from "fs";
import { Inventory } from "./inventory/Inventory";

const printUsage = (params?: string[]): void => {
  console.warn("USAGE: npx k2 <k2 inventory file path> <action>");
  console.warn("received parameters", params);
};

const printVersion = (): void => {
  const packageJson = path.join(__dirname, "..", "package.json");
  const obj = JSON.parse(fs.readFileSync(packageJson, "utf-8"));

  console.warn("K2 by github:@tuxounet", obj.version);
};

const checkParams = (): {
  inventoryArg: string;
  commandArg: string;
  runFolder: string;
} => {
  const params = process.argv.filter(
    (item, index) => index >= process.argv.length - 2
  );

  try {
    const inventoryArg = params[0];
    if (inventoryArg.trim() === "" || inventoryArg.endsWith("node")) {
      throw new Error(
        "vous devez preciser un chemin vers un fichier d'inventaire k2"
      );
    }
    if (!fs.existsSync(inventoryArg)) {
      throw new Error("fichier d'inventaire introuvable");
    }
    const stat = fs.statSync(inventoryArg);
    if (!stat.isFile()) {
      throw new Error("chemin d'entrée doit etre un fichier d'inventaire k2");
    }

    const commandArg = params[1];
    if (
      commandArg !== undefined &&
      commandArg !== "" &&
      typeof commandArg === "string" &&
      commandArg !== __filename
    ) {
      return {
        runFolder: process.cwd(),
        inventoryArg,
        commandArg: commandArg.toLowerCase().trim(),
      };
    } else {
      throw new Error("format de la commande incorrect");
    }
  } catch (error) {
    printUsage(params);
    throw error;
  }
};
printVersion();
const params = checkParams();
console.info("parametres", params);
const run = async (): Promise<void> => {
  const fullPath = path.resolve(params.inventoryArg);
  const inventoryFolder = path.dirname(fullPath);
  const inventoryFilename = path.basename(fullPath);
  const inventory = new Inventory(inventoryFilename, inventoryFolder);

  await inventory.loadCommands();

  if (
    !Array.from(inventory.allowedCommands.keys()).includes(params.commandArg)
  ) {
    throw new Error("commande incorrecte");
  }

  await inventory.load();

  const handler = inventory.allowedCommands.get(params.commandArg);

  if (handler == null) throw new Error("commande introuvable");

  await handler(inventory);
};

run().catch((e) => {
  console.error(e);
  process.exit(1);
});
