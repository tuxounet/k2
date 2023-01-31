import { IK2 } from "../types/IK2";
import fs from "fs";
import path from "path";
import jsYaml from "js-yaml";

export function loadK2File<T extends IK2>(filePath: string): T {
  const item = jsYaml.load(
    fs.readFileSync(filePath, {
      encoding: "utf-8",
    })
  ) as T;
  item.k2.metadata.path = filePath;
  item.k2.metadata.folder = path.dirname(filePath);
  return item;
}
