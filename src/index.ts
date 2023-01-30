#!/usr/bin/env node
import path from "path";
import fs from "fs";
import { Command } from "commander";
import apply from "./commands/apply";
import clean from "./commands/clean";
import list from "./commands/list";

const packageJson = path.join(__dirname, "..", "package.json");
const obj = JSON.parse(fs.readFileSync(packageJson, "utf-8"));
console.warn("K2 by github:@tuxounet");

const program = new Command();
program
  .version(obj.version)
  .description("An example CLI for managing a directory");

program.addCommand(list());
program.addCommand(apply());
program.addCommand(clean());
program.parse(process.argv);
program.showHelpAfterError();
