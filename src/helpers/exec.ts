import childProc from "child_process";

export async function exec(
  cmd: string,
  cwd: string,
  allowFailure = false
): Promise<number> {
  return await new Promise<number>((resolve, reject) => {
    let executable = cmd;
    let args: string[] = [];
    const cmdSegments = cmd.split(" ");
    if (cmdSegments.length > 0) {
      executable = cmdSegments[0];
      cmdSegments.splice(0, 1);
      args = cmdSegments;
    }
    console.debug(">", executable, args, cwd);
    const ps = childProc.execFile(executable, args, {
      encoding: "utf-8",
      cwd,
      shell: true,
    });

    ps.on("exit", (code) => {
      if (code === 0 || code === null) {
        resolve(0);
        return;
      }
      if (allowFailure) {
        resolve(code);
        return;
      }
      reject(code);
    });
  });
}
