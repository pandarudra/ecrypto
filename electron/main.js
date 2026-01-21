const { app, BrowserWindow, ipcMain, dialog, Menu } = require("electron");
const path = require("path");
const { spawn, spawnSync } = require("child_process");
const axios = require("axios");

let mainWindow;
let goServer;
const API_PORT = 8765;
const API_URL = `http://127.0.0.1:${API_PORT}`; // Use IPv4 explicitly

// Handle --version flag before app starts
if (process.argv.includes("--version") || process.argv.includes("-v")) {
  const ext = process.platform === "win32" ? ".exe" : "";
  const binaryPath = path.join(__dirname, "..", "bin", `ecrypto${ext}`);
  const result = spawnSync(binaryPath, ["--version"], { encoding: "utf8" });
  console.log(result.stdout || result.stderr);
  process.exit(0);
}

// Determine Go binary path
function getGoBinaryPath() {
  const ext = process.platform === "win32" ? ".exe" : "";
  if (app.isPackaged) {
    // Production: binary in resources
    return path.join(process.resourcesPath, "bin", `ecrypto${ext}`);
  }
  // Development: binary in project directory
  return path.join(__dirname, "..", "bin", `ecrypto${ext}`);
}

// Start Go API server
function startGoServer() {
  return new Promise((resolve, reject) => {
    const binaryPath = getGoBinaryPath();
    console.log("Starting Go server:", binaryPath);

    goServer = spawn(binaryPath, ["--serve", `--port=${API_PORT}`]);

    goServer.stdout.on("data", (data) => {
      console.log(`[Go Server]: ${data}`);
      if (data.toString().includes("Server started")) {
        resolve();
      }
    });

    goServer.stderr.on("data", (data) => {
      console.log(data);
      const message = data.toString();
      console.log(`[Go Server STDERR]: ${message}`);
      // Go's log package writes to stderr by default, even for info messages
      if (
        message.includes("Server started") ||
        message.includes("http://localhost")
      ) {
        console.log(`[Go Server]: ${message}`);
        if (message.includes("Server started")) {
          resolve();
        }
      } else {
        console.error(`[Go Server Error]: ${message}`);
      }
    });

    goServer.on("error", (err) => {
      console.error("Failed to start Go server:", err);
      reject(err);
    });

    goServer.on("close", (code) => {
      console.log(`Go server exited with code ${code}`);
    });

    // Timeout fallback
    setTimeout(() => {
      resolve(); // Assume started after 2 seconds
    }, 2000);
  });
}

// Create main window
function createWindow() {
  mainWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    minWidth: 800,
    minHeight: 600,
    webPreferences: {
      preload: path.join(__dirname, "preload.js"),
      contextIsolation: true,
      nodeIntegration: false,
    },
    icon: path.join(__dirname, "assets", "icon.png"),
    title: "Ecrypto - Military-Grade Encryption",
  });

  mainWindow.loadFile(path.join(__dirname, "renderer", "index.html"));

  // Open DevTools in development
  if (!app.isPackaged) {
    mainWindow.webContents.openDevTools();
  }

  mainWindow.on("closed", () => {
    mainWindow = null;
  });
}

// App lifecycle
app.whenReady().then(async () => {
  // Remove the menu bar
  Menu.setApplicationMenu(null);
  try {
    await startGoServer();
    createWindow();
  } catch (err) {
    console.error("Failed to start application:", err);
    app.quit();
  }

  app.on("activate", () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow();
    }
  });
});

app.on("window-all-closed", () => {
  if (process.platform !== "darwin") {
    app.quit();
  }
});

app.on("quit", () => {
  if (goServer) {
    goServer.kill();
  }
});

// IPC Handlers

// Select folder
ipcMain.handle("dialog:selectFolder", async () => {
  const result = await dialog.showOpenDialog(mainWindow, {
    properties: ["openDirectory"],
  });
  return result.canceled ? null : result.filePaths[0];
});

// Select file
ipcMain.handle("dialog:selectFile", async (event, filters) => {
  const result = await dialog.showOpenDialog(mainWindow, {
    properties: ["openFile"],
    filters: filters || [
      { name: "Ecrypt Files", extensions: ["ecrypt"] },
      { name: "All Files", extensions: ["*"] },
    ],
  });
  return result.canceled ? null : result.filePaths[0];
});

// Save file dialog
ipcMain.handle("dialog:saveFile", async (event, defaultName, filters) => {
  const result = await dialog.showSaveDialog(mainWindow, {
    defaultPath: defaultName,
    filters: filters || [
      { name: "Ecrypt Files", extensions: ["ecrypt"] },
      { name: "All Files", extensions: ["*"] },
    ],
  });
  return result.canceled ? null : result.filePath;
});

// API Proxy Handlers
ipcMain.handle("api:encrypt", async (event, data) => {
  try {
    const response = await axios.post(`${API_URL}/encrypt`, data);
    return { success: true, data: response.data };
  } catch (error) {
    return {
      success: false,
      error: error.response?.data?.error || error.message,
    };
  }
});

ipcMain.handle("api:decrypt", async (event, data) => {
  try {
    const response = await axios.post(`${API_URL}/decrypt`, data);
    return { success: true, data: response.data };
  } catch (error) {
    return {
      success: false,
      error: error.response?.data?.error || error.message,
    };
  }
});

ipcMain.handle("api:keygen", async (event, data) => {
  try {
    const response = await axios.post(`${API_URL}/keygen`, data);
    return { success: true, data: response.data };
  } catch (error) {
    return {
      success: false,
      error: error.response?.data?.error || error.message,
    };
  }
});

ipcMain.handle("api:info", async (event, data) => {
  try {
    const response = await axios.post(`${API_URL}/info`, data);
    return { success: true, data: response.data };
  } catch (error) {
    return {
      success: false,
      error: error.response?.data?.error || error.message,
    };
  }
});

ipcMain.handle("api:history", async () => {
  try {
    const response = await axios.get(`${API_URL}/history`);
    return { success: true, data: response.data };
  } catch (error) {
    return {
      success: false,
      error: error.response?.data?.error || error.message,
    };
  }
});

ipcMain.handle("api:undo", async (event, operationId) => {
  try {
    const response = await axios.post(`${API_URL}/undo`, { operationId });
    return { success: true, data: response.data };
  } catch (error) {
    return {
      success: false,
      error: error.response?.data?.error || error.message,
    };
  }
});

ipcMain.handle("api:suggest-path", async (event, inputPath) => {
  try {
    const response = await axios.post(`${API_URL}/suggest-path`, {
      path: inputPath,
    });
    return { success: true, data: response.data };
  } catch (error) {
    return {
      success: false,
      error: error.response?.data?.error || error.message,
    };
  }
});

ipcMain.handle("api:check-password", async (event, password) => {
  try {
    const response = await axios.post(`${API_URL}/check-password`, {
      password,
    });
    return { success: true, data: response.data };
  } catch (error) {
    return {
      success: false,
      error: error.response?.data?.error || error.message,
    };
  }
});

// Progress updates via SSE
ipcMain.handle("api:subscribe-progress", (event) => {
  const eventSource = new EventSource(`${API_URL}/progress`);

  eventSource.onmessage = (e) => {
    const data = JSON.parse(e.data);
    mainWindow.webContents.send("progress-update", data);
  };

  eventSource.onerror = (err) => {
    console.error("SSE Error:", err);
    eventSource.close();
  };

  return { success: true };
});
