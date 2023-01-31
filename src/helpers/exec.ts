import childProc from "child_process";

export async function exec(cmd: string, cwd: string): Promise<number> {
  console.debug("exec async", cmd, "inside", cwd);
  return await new Promise<number>((resolve, reject) => {
    let executable = cmd;
    let args: string[] = [];
    const cmdSegments = cmd.split(" ");
    if (cmdSegments.length > 0) {
      executable = cmdSegments[0];
      cmdSegments.splice(0, 1);
      args = cmdSegments;
    }
    console.info(">", executable, args);
    const ps = childProc.execFile(executable, args, {
      encoding: "utf-8",
      cwd,
      shell: true,
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
