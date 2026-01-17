// Application State
const state = {
  encryptInputPath: null,
  decryptInputPath: null,
  infoInputPath: null,
  encryptKeyFile: null,
  decryptKeyFile: null,
};

// Check if electronAPI is available
if (typeof window.electronAPI === "undefined") {
  console.error(
    "electronAPI is not available. Make sure preload.js is loaded correctly.",
  );
}

// Tab Navigation
document.querySelectorAll(".nav-item").forEach((btn) => {
  btn.addEventListener("click", () => {
    const tab = btn.dataset.tab;
    switchTab(tab);
  });
});

function switchTab(tab) {
  // Update navigation
  document
    .querySelectorAll(".nav-item")
    .forEach((item) => item.classList.remove("active"));
  document.querySelector(`[data-tab="${tab}"]`).classList.add("active");

  // Update panels
  document
    .querySelectorAll(".tab-panel")
    .forEach((panel) => panel.classList.remove("active"));
  document.getElementById(`${tab}-panel`).classList.add("active");

  // Update header
  const titles = {
    encrypt: [
      "Encrypt Files & Folders",
      "Protect your data with military-grade encryption",
    ],
    decrypt: ["Decrypt Files & Folders", "Restore your encrypted data"],
    keygen: ["Generate Encryption Key", "Create cryptographically secure keys"],
    info: ["Container Information", "View metadata without decryption"],
    history: ["Operation History", "Review your recent encryption activities"],
  };

  document.getElementById("page-title").textContent = titles[tab][0];
  document.getElementById("page-subtitle").textContent = titles[tab][1];

  // Re-initialize Lucide icons
  lucide.createIcons();
}

// ===== ENCRYPT TAB =====

// Drag and drop for encrypt input
const encryptInputSelector = document.getElementById("encrypt-input-selector");
if (encryptInputSelector) {
  encryptInputSelector.addEventListener("dragover", (e) => {
    e.preventDefault();
    e.stopPropagation();
    encryptInputSelector.classList.add("drag-over");
  });

  encryptInputSelector.addEventListener("dragleave", (e) => {
    e.preventDefault();
    e.stopPropagation();
    encryptInputSelector.classList.remove("drag-over");
  });

  encryptInputSelector.addEventListener("drop", (e) => {
    e.preventDefault();
    e.stopPropagation();
    encryptInputSelector.classList.remove("drag-over");

    const files = e.dataTransfer.files;
    if (files.length > 0) {
      const path = files[0].path;
      state.encryptInputPath = path;
      updateFileDisplay("encrypt-input-display", path);
      getSuggestionsForPath(path, "encrypt-suggestions", "encrypt-output");
    }
  });
}

// Select folder for encryption
document
  .getElementById("encrypt-select-folder")
  .addEventListener("click", async () => {
    const path = await window.electronAPI.selectFolder();
    if (path) {
      state.encryptInputPath = path;
      updateFileDisplay("encrypt-input-display", path);
      getSuggestionsForPath(path, "encrypt-suggestions", "encrypt-output");
    }
  });

// Select file for encryption
document
  .getElementById("encrypt-select-file")
  .addEventListener("click", async () => {
    const path = await window.electronAPI.selectFile();
    if (path) {
      state.encryptInputPath = path;
      updateFileDisplay("encrypt-input-display", path);
      getSuggestionsForPath(path, "encrypt-suggestions", "encrypt-output");
    }
  });

// Save as dialog
document
  .getElementById("encrypt-save-as")
  .addEventListener("click", async () => {
    const path = await window.electronAPI.saveFile("output.ecrypt", [
      { name: "Ecrypt Files", extensions: ["ecrypt"] },
    ]);
    if (path) {
      document.getElementById("encrypt-output").value = path;
    }
  });

// Security method toggle
document.querySelectorAll('input[name="encrypt-method"]').forEach((radio) => {
  radio.addEventListener("change", (e) => {
    const usePassword = e.target.value === "password";
    document.getElementById("encrypt-password-section").style.display =
      usePassword ? "block" : "none";
    document.getElementById("encrypt-key-section").style.display = usePassword
      ? "none"
      : "block";
  });
});

// Password visibility toggle
document
  .getElementById("encrypt-toggle-password")
  .addEventListener("click", (e) => {
    const input = document.getElementById("encrypt-password");
    const button = e.currentTarget;
    const icon = button.querySelector("i");

    if (input.type === "password") {
      input.type = "text";
      icon.setAttribute("data-lucide", "eye-off");
    } else {
      input.type = "password";
      icon.setAttribute("data-lucide", "eye");
    }

    // Re-initialize the icon
    lucide.createIcons();
  });

// Password strength check
document
  .getElementById("encrypt-password")
  .addEventListener("input", async (e) => {
    const password = e.target.value;
    const strengthElement = document.getElementById(
      "encrypt-password-strength",
    );

    if (password.length > 0) {
      const result = await window.electronAPI.checkPassword(password);
      if (result.success && result.data) {
        displayPasswordStrength(result.data.data, "encrypt-password-strength");
      } else {
        strengthElement.classList.add("visible");
        strengthElement.innerHTML =
          '<p style="color: var(--danger);">Unable to check password strength</p>';
      }
    } else {
      strengthElement.classList.remove("visible");
    }
  });

// Select key file
document
  .getElementById("encrypt-select-key")
  .addEventListener("click", async () => {
    const path = await window.electronAPI.selectFile([
      { name: "Key Files", extensions: ["key", "txt"] },
      { name: "All Files", extensions: ["*"] },
    ]);
    if (path) {
      state.encryptKeyFile = path;
      document.getElementById("encrypt-key-file").value = path;
    }
  });

// Encrypt button
document
  .getElementById("encrypt-button")
  .addEventListener("click", async () => {
    if (!state.encryptInputPath) {
      showToast("error", "Please select a file or folder to encrypt");
      return;
    }

    const outputPath = document.getElementById("encrypt-output").value;
    if (!outputPath) {
      showToast("error", "Please specify an output path");
      return;
    }

    const useKey =
      document.querySelector('input[name="encrypt-method"]:checked').value ===
      "key";
    const password = document.getElementById("encrypt-password").value;
    const keyFile = state.encryptKeyFile;

    if (!useKey && !password) {
      showToast("error", "Please enter a passphrase");
      return;
    }

    if (useKey && !keyFile) {
      showToast("error", "Please select a key file");
      return;
    }

    showProgress("Encrypting...");

    const result = await window.electronAPI.encrypt({
      inputPath: state.encryptInputPath,
      outputPath,
      password: useKey ? undefined : password,
      keyFile: useKey ? keyFile : undefined,
      useKey,
    });

    console.log(result);

    hideProgress();

    if (result.success) {
      showToast("success", "Encryption completed successfully!");
      clearEncryptForm();
    } else {
      showToast("error", `Encryption failed: ${result.error}`);
    }
  });

// ===== DECRYPT TAB =====

// Drag and drop for decrypt input
const decryptFileSelector = document.querySelector(
  "#decrypt-panel .file-selector",
);
if (decryptFileSelector) {
  decryptFileSelector.addEventListener("dragover", (e) => {
    e.preventDefault();
    e.stopPropagation();
    decryptFileSelector.classList.add("drag-over");
  });

  decryptFileSelector.addEventListener("dragleave", (e) => {
    e.preventDefault();
    e.stopPropagation();
    decryptFileSelector.classList.remove("drag-over");
  });

  decryptFileSelector.addEventListener("drop", (e) => {
    e.preventDefault();
    e.stopPropagation();
    decryptFileSelector.classList.remove("drag-over");

    const files = e.dataTransfer.files;
    if (files.length > 0) {
      const path = files[0].path;
      state.decryptInputPath = path;
      updateFileDisplay("decrypt-input-display", path);
    }
  });
}

document
  .getElementById("decrypt-select-file")
  .addEventListener("click", async () => {
    const path = await window.electronAPI.selectFile([
      { name: "Ecrypt Files", extensions: ["ecrypt"] },
      { name: "All Files", extensions: ["*"] },
    ]);
    if (path) {
      state.decryptInputPath = path;
      updateFileDisplay("decrypt-input-display", path);
    }
  });

document
  .getElementById("decrypt-select-folder")
  .addEventListener("click", async () => {
    const path = await window.electronAPI.selectFolder();
    if (path) {
      document.getElementById("decrypt-output").value = path;
    }
  });

document.querySelectorAll('input[name="decrypt-method"]').forEach((radio) => {
  radio.addEventListener("change", (e) => {
    const usePassword = e.target.value === "password";
    document.getElementById("decrypt-password-section").style.display =
      usePassword ? "block" : "none";
    document.getElementById("decrypt-key-section").style.display = usePassword
      ? "none"
      : "block";
  });
});

document
  .getElementById("decrypt-toggle-password")
  .addEventListener("click", (e) => {
    const input = document.getElementById("decrypt-password");
    const button = e.currentTarget;
    const icon = button.querySelector("i");

    if (input.type === "password") {
      input.type = "text";
      icon.setAttribute("data-lucide", "eye-off");
    } else {
      input.type = "password";
      icon.setAttribute("data-lucide", "eye");
    }

    // Re-initialize the icon
    lucide.createIcons();
  });

document
  .getElementById("decrypt-select-key")
  .addEventListener("click", async () => {
    const path = await window.electronAPI.selectFile([
      { name: "Key Files", extensions: ["key", "txt"] },
      { name: "All Files", extensions: ["*"] },
    ]);
    if (path) {
      state.decryptKeyFile = path;
      document.getElementById("decrypt-key-file").value = path;
    }
  });

document
  .getElementById("decrypt-button")
  .addEventListener("click", async () => {
    if (!state.decryptInputPath) {
      showToast("error", "Please select a container file to decrypt");
      return;
    }

    const outputPath = document.getElementById("decrypt-output").value;
    if (!outputPath) {
      showToast("error", "Please specify an output path");
      return;
    }

    const useKey =
      document.querySelector('input[name="decrypt-method"]:checked').value ===
      "key";
    const password = document.getElementById("decrypt-password").value;
    const keyFile = state.decryptKeyFile;

    if (!useKey && !password) {
      showToast("error", "Please enter a passphrase");
      return;
    }

    if (useKey && !keyFile) {
      showToast("error", "Please select a key file");
      return;
    }

    showProgress("Decrypting...");

    const result = await window.electronAPI.decrypt({
      inputPath: state.decryptInputPath,
      outputPath,
      password: useKey ? undefined : password,
      keyFile: useKey ? keyFile : undefined,
      useKey,
    });

    hideProgress();

    if (result.success) {
      showToast("success", "Decryption completed successfully!");
      clearDecryptForm();
    } else {
      showToast("error", `Decryption failed: ${result.error}`);
    }
  });

// ===== KEYGEN TAB =====

document
  .getElementById("keygen-save-as")
  .addEventListener("click", async () => {
    const path = await window.electronAPI.saveFile("encryption.key", [
      { name: "Key Files", extensions: ["key"] },
      { name: "Text Files", extensions: ["txt"] },
    ]);
    if (path) {
      document.getElementById("keygen-output").value = path;
    }
  });

document.getElementById("keygen-button").addEventListener("click", async () => {
  const outputPath = document.getElementById("keygen-output").value;
  if (!outputPath) {
    showToast("error", "Please specify an output path for the key");
    return;
  }

  showProgress("Generating key...");

  const result = await window.electronAPI.generateKey({ outputPath });

  hideProgress();

  if (result.success) {
    const keyElement = document.getElementById("keygen-key-value");
    console.log("Full API response:", result);
    console.log("result.data:", result.data);

    // The API response structure is: result.data.data.key
    const key = result.data.data?.key || result.data.key;
    console.log("Extracted key:", key);

    keyElement.textContent = key;
    console.log(
      "Key element textContent after setting:",
      keyElement.textContent,
    );
    document.getElementById("keygen-result").style.display = "block";
    showToast("success", "Key generated successfully!");
  } else {
    showToast("error", `Key generation failed: ${result.error}`);
  }
});

document.getElementById("keygen-copy").addEventListener("click", async () => {
  const keyElement = document.getElementById("keygen-key-value");
  console.log("Copy button clicked");
  console.log("Key element:", keyElement);
  console.log("Key element exists:", !!keyElement);
  console.log(
    "Key element textContent:",
    keyElement ? keyElement.textContent : "element not found",
  );

  const key = keyElement ? keyElement.textContent : "";
  console.log("Attempting to copy key, length:", key ? key.length : 0);

  if (!key || key.trim() === "") {
    showToast("error", "No key to copy");
    return;
  }

  try {
    // Try using the Clipboard API
    if (navigator.clipboard && navigator.clipboard.writeText) {
      await navigator.clipboard.writeText(key);
      showToast("info", "Key copied to clipboard");
    } else {
      // Fallback method for older browsers
      const textarea = document.createElement("textarea");
      textarea.value = key;
      textarea.style.position = "fixed";
      textarea.style.opacity = "0";
      document.body.appendChild(textarea);
      textarea.select();
      document.execCommand("copy");
      document.body.removeChild(textarea);
      showToast("info", "Key copied to clipboard");
    }
  } catch (err) {
    console.error("Failed to copy to clipboard:", err);
    showToast("error", "Failed to copy to clipboard: " + err.message);
  }
});

// ===== INFO TAB =====

document
  .getElementById("info-select-file")
  .addEventListener("click", async () => {
    const path = await window.electronAPI.selectFile([
      { name: "Ecrypt Files", extensions: ["ecrypt"] },
      { name: "All Files", extensions: ["*"] },
    ]);
    if (path) {
      state.infoInputPath = path;
      updateFileDisplay("info-input-display", path);
    }
  });

document.getElementById("info-button").addEventListener("click", async () => {
  if (!state.infoInputPath) {
    showToast("error", "Please select a container file");
    return;
  }

  showProgress("Reading container info...");

  const result = await window.electronAPI.getInfo({
    filePath: state.infoInputPath,
  });

  hideProgress();
  // console.log(result);

  if (result.success) {
    displayContainerInfo(result.data.data);
  } else {
    showToast("error", `Failed to read info: ${result.error}`);
  }
});

// ===== HISTORY TAB =====

document
  .getElementById("refresh-history")
  .addEventListener("click", loadHistory);

async function loadHistory() {
  const result = await window.electronAPI.getHistory();

  console.log("History API response:", result);
  console.log("result.data:", result.data);

  const historyList = document.getElementById("history-list");

  // Handle nested response structure: result.data.data
  let historyData = [];
  if (result.success && result.data) {
    // Check for nested data structure (result.data.data)
    if (result.data.data) {
      console.log("result.data.data:", result.data.data);
      console.log("Is array?", Array.isArray(result.data.data));
      console.log("result.data.data.operations:", result.data.data.operations);
      historyData = Array.isArray(result.data.data)
        ? result.data.data
        : result.data.data.operations || [];
    } else {
      historyData = Array.isArray(result.data)
        ? result.data
        : result.data.operations || [];
    }
  }

  console.log("Processed history data:", historyData);
  console.log("History data length:", historyData.length);

  if (historyData.length === 0) {
    historyList.innerHTML =
      '<p class="placeholder-text">No operations in history</p>';
    return;
  }

  historyList.innerHTML = "";

  historyData.forEach((item) => {
    const historyItem = document.createElement("div");
    historyItem.className = "history-item";

    const historyDetails = document.createElement("div");
    historyDetails.className = "history-details";

    const historyOperation = document.createElement("div");
    historyOperation.className = "history-operation";
    historyOperation.textContent = item.type || "Unknown Operation";

    const historyPath = document.createElement("div");
    historyPath.className = "history-path";
    historyPath.textContent = `${item.input_path || "N/A"} → ${item.output_path || "N/A"}`;

    const historyTime = document.createElement("div");
    historyTime.className = "history-time";
    historyTime.textContent = item.timestamp || "N/A";

    historyDetails.appendChild(historyOperation);
    historyDetails.appendChild(historyPath);
    historyDetails.appendChild(historyTime);

    const historyActions = document.createElement("div");
    historyActions.className = "history-actions";

    if (item.type === "encrypt") {
      const undoBtn = document.createElement("button");
      undoBtn.className = "btn btn-small btn-secondary";
      undoBtn.onclick = () => undoOperation(item.id);

      const icon = document.createElement("i");
      icon.setAttribute("data-lucide", "rotate-ccw");

      undoBtn.appendChild(icon);
      undoBtn.appendChild(document.createTextNode(" Undo"));
      historyActions.appendChild(undoBtn);
    }

    historyItem.appendChild(historyDetails);
    historyItem.appendChild(historyActions);
    historyList.appendChild(historyItem);
  });

  // Re-initialize Lucide icons for the history list
  lucide.createIcons();
}

async function undoOperation(operationId) {
  showProgress("Undoing operation...");

  const result = await window.electronAPI.undoOperation(operationId);

  hideProgress();

  if (result.success) {
    showToast("success", "Operation undone successfully");
    loadHistory();
  } else {
    showToast("error", `Undo failed: ${result.error}`);
  }
}

// ===== UTILITY FUNCTIONS =====

function updateFileDisplay(elementId, path) {
  const display = document.getElementById(elementId);
  display.innerHTML = `<span class="selected-path">${path}</span>`;
}

async function getSuggestionsForPath(
  inputPath,
  suggestionsElementId,
  outputInputId,
) {
  const result = await window.electronAPI.suggestPath(inputPath);

  if (!result.success || !result.data) return;

  // Handle both array and non-array responses
  let suggestions = Array.isArray(result.data) ? result.data : [];
  if (suggestions.length === 0) return;

  const suggestionsElement = document.getElementById(suggestionsElementId);
  suggestionsElement.innerHTML = suggestions
    .map(
      (suggestion) =>
        `<span class="suggestion-chip" onclick="applySuggestion('${outputInputId}', '${suggestion.path || suggestion}')">${suggestion.path || suggestion}</span>`,
    )
    .join("");
}

function applySuggestion(inputId, path) {
  document.getElementById(inputId).value = path;
}

function displayPasswordStrength(data, elementId) {
  const strengthElement = document.getElementById(elementId);

  // Validate data
  if (!data || !data.strength) {
    strengthElement.classList.add("visible");
    strengthElement.innerHTML =
      '<p style="color: var(--danger);">Unable to check password strength</p>';
    return;
  }

  strengthElement.classList.add("visible");
  const strength = data.strength.toLowerCase().replace(/ /g, "-");

  strengthElement.innerHTML = `
        <div class="strength-bar">
            <div class="strength-fill ${strength}"></div>
        </div>
        <div style="font-size: 0.875rem;">
            <strong>Strength:</strong> ${data.strength}
        </div>
        ${
          data.suggestions && data.suggestions.length > 0
            ? `
            <div style="margin-top: 0.5rem; font-size: 0.875rem;">
                <strong>Tips:</strong>
                <ul style="margin-left: 1.25rem; margin-top: 0.25rem;">
                    ${data.suggestions.map((s) => `<li>${s}</li>`).join("")}
                </ul>
            </div>
        `
            : ""
        }
    `;
}

function displayContainerInfo(data) {
  const infoResult = document.getElementById("info-result");
  infoResult.style.display = "block";

  infoResult.innerHTML = `
        <div class="info-item">
            <span class="info-label">Magic Header:</span>
            <span class="info-value">${data.magic || "N/A"}</span>
        </div>
        <div class="info-item">
            <span class="info-label">Version:</span>
            <span class="info-value">${data.version || "N/A"}</span>
        </div>
        <div class="info-item">
            <span class="info-label">KDF Type:</span>
            <span class="info-value">${data.kdfType || "N/A"}</span>
        </div>
        <div class="info-item">
            <span class="info-label">Argon2 Memory:</span>
            <span class="info-value">${data.argonMemory || "N/A"} MB</span>
        </div>
        <div class="info-item">
            <span class="info-label">Argon2 Time:</span>
            <span class="info-value">${data.argonTime || "N/A"} iterations</span>
        </div>
        <div class="info-item">
            <span class="info-label">Argon2 Parallelism:</span>
            <span class="info-value">${data.argonParallelism || "N/A"}</span>
        </div>
        <div class="info-item">
            <span class="info-label">Container Size:</span>
            <span class="info-value">${formatBytes(data.size) || "N/A"}</span>
        </div>
    `;
}

function showProgress(title) {
  document.getElementById("progress-title").textContent = title;
  document.getElementById("progress-text").textContent = "Processing...";
  document.getElementById("progress-fill").style.width = "50%";
  document.getElementById("progress-modal").style.display = "flex";
}

function hideProgress() {
  document.getElementById("progress-modal").style.display = "none";
}

function showToast(type, message) {
  const container = document.getElementById("toast-container");
  const toast = document.createElement("div");
  toast.className = `toast ${type}`;

  const icons = {
    success: "check-circle",
    error: "x-circle",
    info: "info",
  };

  // Create icon element
  const iconElement = document.createElement("i");
  iconElement.setAttribute("data-lucide", icons[type]);
  iconElement.className = "toast-icon";

  // Create message element
  const messageElement = document.createElement("div");
  messageElement.className = "toast-message";
  messageElement.textContent = message;

  toast.appendChild(iconElement);
  toast.appendChild(messageElement);
  container.appendChild(toast);

  // Initialize Lucide icons for the new toast
  lucide.createIcons();

  setTimeout(() => {
    toast.style.animation = "slideIn 0.3s reverse";
    setTimeout(() => toast.remove(), 300);
  }, 3000);
}

function formatBytes(bytes) {
  if (!bytes) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + " " + sizes[i];
}

function clearEncryptForm() {
  state.encryptInputPath = null;
  state.encryptKeyFile = null;

  const displayElement = document.getElementById("encrypt-input-display");
  displayElement.innerHTML = "";

  const placeholder = document.createElement("span");
  placeholder.className = "placeholder";

  const icon = document.createElement("i");
  icon.setAttribute("data-lucide", "folder");

  placeholder.appendChild(icon);
  placeholder.appendChild(
    document.createTextNode(" Click to select folder or file"),
  );
  displayElement.appendChild(placeholder);

  document
    .getElementById("encrypt-password-strength")
    .classList.remove("visible");

  // Re-initialize Lucide icons
  lucide.createIcons();
}

function clearDecryptForm() {
  state.decryptInputPath = null;
  state.decryptKeyFile = null;

  const displayElement = document.getElementById("decrypt-input-display");
  displayElement.innerHTML = "";

  const placeholder = document.createElement("span");
  placeholder.className = "placeholder";

  const icon = document.createElement("i");
  icon.setAttribute("data-lucide", "box");

  placeholder.appendChild(icon);
  placeholder.appendChild(
    document.createTextNode(" Click to select .ecrypt file"),
  );
  displayElement.appendChild(placeholder);

  // Re-initialize Lucide icons
  lucide.createIcons();
}

// Progress updates from backend
if (window.electronAPI && window.electronAPI.onProgress) {
  window.electronAPI.onProgress((data) => {
    if (data.percentage !== undefined) {
      document.getElementById("progress-fill").style.width =
        data.percentage + "%";
      document.getElementById("progress-text").textContent =
        `${data.current} / ${data.total} files processed`;
      document.getElementById("progress-file").textContent =
        data.filename || "";
    }
  });
}

// Initial load
loadHistory();
