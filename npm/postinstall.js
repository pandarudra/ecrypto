const https = require("https");
const fs = require("fs");
const path = require("path");
const os = require("os");

const version = process.env.npm_package_version
  ? "v" + process.env.npm_package_version
  : "v1.0.3";

const platform = os.platform();
const arch = os.arch();

const binDir = path.join(__dirname, "bin");
if (!fs.existsSync(binDir)) fs.mkdirSync(binDir, { recursive: true });

function getBinaryName() {
  if (platform === "win32") {
    return arch === "arm64"
      ? "ecrypto-windows-arm64.exe"
      : "ecrypto-windows-amd64.exe";
  }
  if (platform === "darwin") {
    return arch === "arm64"
      ? "ecrypto-darwin-arm64"
      : "ecrypto-darwin-amd64";
  }
  if (platform === "linux") {
    return arch === "arm64"
      ? "ecrypto-linux-arm64"
      : "ecrypto-linux-amd64";
  }

  throw new Error(`Unsupported platform: ${platform} ${arch}`);
}

const filename = getBinaryName();
const output = path.join(binDir, platform === "win32" ? "ecrypto.exe" : "ecrypto");

const downloadUrl = `https://github.com/pandarudra/ecrypto/releases/download/${version}/${filename}`;

console.log(`Downloading ECRYPTO binary for ${platform}-${arch}:`);
console.log(downloadUrl);

https.get(downloadUrl, (res) => {
  if (res.statusCode !== 200) {
    console.error("Failed to download:", res.statusCode);
    process.exit(1);
  }

  const file = fs.createWriteStream(output);
  res.pipe(file);

  file.on("finish", () => {
    file.close();
    if (platform !== "win32") {
      fs.chmodSync(output, 0o755);
    }
    console.log("✔ ECRYPTO installed successfully");
  });
});
