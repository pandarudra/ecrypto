#!/usr/bin/env node

const { spawn } = require("child_process");
const path = require("path");
const os = require("os");

const platform = os.platform();
const binary = platform === "win32" ? "ecrypto.exe" : "ecrypto";

const binPath = path.join(__dirname, binary);

const child = spawn(binPath, process.argv.slice(2), { stdio: "inherit" });

child.on("exit", (code) => {
  process.exit(code);
});
