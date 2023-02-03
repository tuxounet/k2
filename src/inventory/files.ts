import { IK2 } from "../types/IK2";
import fs from "fs";
import path from "path";
import jsYaml from "js-yaml";
import fg from "fast-glob";
export function loadK2File<T extends IK2>(filePath: string): T {
  try {
    const item = jsYaml.load(
      fs.readFileSync(filePath, {
        encoding: "utf-8",
      })
    ) as T;
    item.k2.metadata.path = filePath;
    item.k2.metadata.folder = path.dirname(filePath);
    return item;
  } catch (e) {
    console.error("file format error", filePath);
    throw e;
  }
}

export async function loadManyK2Files<T extends IK2>(
  globs: string[],
  folder: string,
  ignore?: string[]
): Promise<T[]> {
  const files = await fg(globs, {
    onlyFiles: true,
    cwd: folder,
    ignore,
  });
  return files
    .map((item) => path.join(folder, item))
    .map((item) => {
      return {
        path: item,
        body: loadK2File<T>(item),
      };
    })
    .filter((item) => item.body.k2.metadata.kind)
    .map((item) => item.body);
}
