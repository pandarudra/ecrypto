const { contextBridge, ipcRenderer } = require("electron");

// Expose secure API to renderer process
contextBridge.exposeInMainWorld("electronAPI", {
  // Dialog APIs
  selectFolder: () => ipcRenderer.invoke("dialog:selectFolder"),
  selectFile: (filters) => ipcRenderer.invoke("dialog:selectFile", filters),
  saveFile: (defaultName, filters) =>
    ipcRenderer.invoke("dialog:saveFile", defaultName, filters),

  // Crypto APIs
  encrypt: (data) => ipcRenderer.invoke("api:encrypt", data),
  decrypt: (data) => ipcRenderer.invoke("api:decrypt", data),
  generateKey: (data) => ipcRenderer.invoke("api:keygen", data),
  getInfo: (data) => ipcRenderer.invoke("api:info", data),

  // History APIs
  getHistory: () => ipcRenderer.invoke("api:history"),
  undoOperation: (operationId) => ipcRenderer.invoke("api:undo", operationId),

  // AI APIs
  suggestPath: (inputPath) => ipcRenderer.invoke("api:suggest-path", inputPath),
  checkPassword: (password) =>
    ipcRenderer.invoke("api:check-password", password),

  // Progress tracking
  subscribeProgress: () => ipcRenderer.invoke("api:subscribe-progress"),
  onProgress: (callback) => {
    ipcRenderer.on("progress-update", (event, data) => callback(data));
  },
});
