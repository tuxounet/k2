#!/usr/bin/env node

const args1 = process.argv[process.argv.length - 2];
const args2 = process.argv[process.argv.length - 1];
const childProc = require("child_process");
const proc = childProc.exec(`npx ts-node ./src/index.ts ${args1} ${args2}`, {
  encoding: "utf-8",
  cwd: __dirname
});

proc.stdin.pipe(process.stdin);
proc.stdout.pipe(process.stdout);
proc.stderr.pipe(process.stderr);
proc.on("exit", (code, signal) => {
  if (code === 0) {
    process.exit(0);
  } else {
    console.error(code, signal);
    process.exit(code || 254);
  }
});
