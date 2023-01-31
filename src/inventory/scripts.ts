import { exec } from "../helpers/exec";
import { IK2Apply } from "../types/IK2Apply";
import { IK2Template } from "../types/IK2Template";

export async function executeScript(
  askedBy: IK2Template | IK2Apply,
  script: string,
  cwd: string
): Promise<void> {
  if (askedBy.k2.body.scripts === undefined) return;
  const scriptsBloc = askedBy.k2.body.scripts as Record<string, string[]>;

  const scriptValue = scriptsBloc[script];
  if (scriptValue === undefined) return;

  console.info("exec script", script, "from", askedBy.k2.metadata.id);
  for (const line of scriptValue) {
    exec(line, cwd);
  }
}
