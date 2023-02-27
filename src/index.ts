#!/usr/bin/env node
import path from "path";
import fs from "fs";
import { Command } from "commander";
import apply from "./commands/apply";
import clean from "./commands/clean";
import list from "./commands/list";
import create from "./commands/create";

const packageJson = path.join(__dirname, "..", "package.json");
const obj = JSON.parse(fs.readFileSync(packageJson, "utf-8"));
console.warn("K2 Build System CLI by github:@tuxounet", obj.version);

const program = new Command();
program.version(obj.version).description("K2 Build System CLI");

program.addCommand(create());
program.addCommand(list());
program.addCommand(apply());
program.addCommand(clean());
program.parse(process.argv);
program.showHelpAfterError();
