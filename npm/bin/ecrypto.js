#!/usr/bin/env node

const path = require("path");
const { spawn } = require("child_process");
const os = require("os");

const exe = os.platform() === "win32" ? "ecrypto.exe" : "ecrypto";
const binaryPath = path.join(__dirname, exe);

const args = process.argv.slice(2);

const child = spawn(binaryPath, args, { stdio: "inherit" });

child.on("exit", (code) => {
  process.exit(code);
});
