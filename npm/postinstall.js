#!/usr/bin/env node

const os = require("os");
const path = require("path");
const https = require("https");
const fs = require("fs");

const version = "v1.0.6";

function getBinaryName() {
  const platform = os.platform();
  const arch = os.arch();

  const archMap = {
    x64: "amd64",
    arm64: "arm64",
  };

  const platformMap = {
    win32: "windows",
    darwin: "darwin",
    linux: "linux",
  };

  const goos = platformMap[platform];
  const goarch = archMap[arch];

  if (!goos || !goarch) {
    console.error(`Unsupported platform or architecture: ${platform} ${arch}`);
    process.exit(1);
    return;
  }

  return goos === "windows"
    ? `ecrypto-${goos}-${goarch}.exe`
    : `ecrypto-${goos}-${goarch}`;
}

function download(url, dest, cb) {
  https.get(url, (res) => {
    // Follow redirects
    if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
      return download(res.headers.location, dest, cb);
    }

    if (res.statusCode !== 200) {
      console.error("❌ Failed to download:", res.statusCode, url);
      process.exit(1);
    }

    const file = fs.createWriteStream(dest);
    res.pipe(file);

    file.on("finish", () => file.close(cb));
  });
}

const fileName = getBinaryName();
const downloadUrl = `https://github.com/pandarudra/ecrypto/releases/download/${version}/${fileName}`;

console.log(`⬇️  Downloading ECRYPTO binary: ${downloadUrl}`);

// bin folder should be inside npm/ folder
const binDir = path.join(__dirname, "bin");
const binaryName = os.platform() === "win32" ? "ecrypto.exe" : "ecrypto";
const outputPath = path.join(binDir, binaryName);

// Ensure bin/ exists
fs.mkdirSync(binDir, { recursive: true });

download(downloadUrl, outputPath, () => {
  if (os.platform() !== "win32") {
    fs.chmodSync(outputPath, 0o755);
  }
  console.log("✅ ECRYPTO installed successfully.");
});
