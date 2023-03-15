const json = require("@rollup/plugin-json");
const typescript = require("@rollup/plugin-typescript");
const shebang = require("rollup-plugin-preserve-shebang");
const noderesolve = require("@rollup/plugin-node-resolve").default;
const commonjs = require("@rollup/plugin-commonjs").default;
const dts = require("rollup-plugin-dts").default;
const pkg = require("./package.json");
module.exports.default = [
  {
    input: "src/bin.ts",
    external: pkg.dependencies ? Object.keys(pkg.dependencies) : [],
    output: {
      file: "dist/bin.js",
      format: "cjs",
      sourcemap: "inline",
    },
    plugins: [
      shebang({ shebang: "#!/usr/bin/env node" }),
      typescript(),
      commonjs(),
      noderesolve(),
      json(),
    ],
  },
  {
    input: "src/lib.ts",
    external: pkg.dependencies ? Object.keys(pkg.dependencies) : [],
    output: {
      file: `dist/index.d.ts`,
      sourcemap: "inline",
      format: "es",
    },
    plugins: [dts()],
  },
  {
    input: "src/lib.ts",
    external: pkg.dependencies ? Object.keys(pkg.dependencies) : [],
    output: {
      file: `dist/index.js`,
      sourcemap: "inline",
      format: "es",
    },
    plugins: [typescript(), commonjs(), noderesolve(), json()],
  },
];
