import childProc from "child_process";

export async function exec(cmd: string, cwd: string): Promise<number> {
  console.debug("exec async", cmd, "inside", cwd);
  return await new Promise<number>((resolve, reject) => {
    const ps = childProc.exec(String(cmd), {
      encoding: "utf-8",
      cwd,
    });
    if (ps.stdin != null) ps.stdin.pipe(process.stdin);
    if (ps.stdout != null) ps.stdout.pipe(process.stdout);
    if (ps.stderr != null) ps.stderr.pipe(process.stderr);

    ps.on("exit", (code) => {
      if (code === 0 || code === null) resolve(0);
      else reject(code);
    });
  });
}
