import { initNavigation, switchView, viewModelFiles, handleSetPreview, handleScanLibrary, loadPreview } from "./modules/ui.js";

// Expose functions to global scope for onclick handlers
window.viewModelFiles = viewModelFiles;
window.loadPreview = loadPreview;
window.handleSetPreview = handleSetPreview;
window.handleScanLibrary = handleScanLibrary;

// Initialize app
document.addEventListener("DOMContentLoaded", () => {
    initNavigation();
    switchView("dashboard");
});
