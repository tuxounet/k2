const json = require("@rollup/plugin-json");
const typescript = require("@rollup/plugin-typescript");
const shebang = require("rollup-plugin-preserve-shebang");
const noderesolve = require("@rollup/plugin-node-resolve").default;
const commonjs = require("@rollup/plugin-commonjs").default;
const pkg = require("./package.json");
module.exports.default = [
  {
    input: "src/index.ts",
    external: Object.keys(pkg.dependencies),
    output: {
      file: "dist/index.js",
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
];
