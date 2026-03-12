#!/usr/bin/env node

const { main } = require("../src/index");

main(process.argv.slice(2)).catch((error) => {
  process.stderr.write(`${error.message}\n`);
  process.exitCode = 1;
});
