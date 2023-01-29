import childProc from "child_process";

export function exec(cmd: string, cwd: string): string {
  console.debug("exec", cmd, "inside", cwd);
  const result = childProc.execSync(cmd, {
    encoding: "utf-8",
    cwd,
  });
  console.debug(result);
  return result;
}
