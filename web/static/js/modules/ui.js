import { fetchModels, fetchModel, fetchModelFiles, fetchLibraries, fetchCollections, fetchTags, setModelPreview, scanLibrary, getFileDownloadUrl } from "./api.js";
import { loadDetailPreview, loadCardPreview, loadImagePreview, loadCardImagePreview } from "./model-viewer.js";
import { is3DFile, isImageFile } from "./three-utils.js";
import { rendererPool } from "./renderer-pool.js";

export function initNavigation() {
    document.querySelectorAll(".nav-item").forEach(item => {
        item.addEventListener("click", (e) => {
            e.preventDefault();
            const view = item.dataset.view;
            switchView(view);
        });
    });
}

export function switchView(viewName) {
    
    document.querySelectorAll(".view").forEach(v => v.classList.remove("active"));
    document.querySelectorAll(".nav-item").forEach(n => n.classList.remove("active"));
    
    document.getElementById(viewName).classList.add("active");
    const navItem = document.querySelector(`[data-view="${viewName}"]`);
    if (navItem) navItem.classList.add("active");
    
    loadViewData(viewName);
}

async function loadViewData(view) {
    switch(view) {
        case "dashboard":
            await loadDashboard();
            break;
        case "models":
            await loadModels();
            break;
        case "libraries":
            await loadLibraries();
            break;
        case "collections":
            await loadCollections();
            break;
    }
}

async function loadDashboard() {
    try {
        const [models, libs, colls, tags] = await Promise.all([
            fetchModels(),
            fetchLibraries(),
            fetchCollections(),
            fetchTags()
        ]);

        document.getElementById("statModels").textContent = models?.length || 0;
        document.getElementById("statLibraries").textContent = libs?.length || 0;
        document.getElementById("statCollections").textContent = colls?.length || 0;
        document.getElementById("statTags").textContent = tags?.length || 0;

        renderModels(models?.slice(0, 6) || [], "recentModels");
    } catch (error) {
        console.error("Failed to load dashboard:", error);
    }
}

async function loadModels() {
    try {
        const models = await fetchModels();
        renderModels(models || [], "modelsList");
    } catch (error) {
        console.error("Failed to load models:", error);
    }
}

async function loadLibraries() {
    try {
        const libs = await fetchLibraries();
        renderLibraries(libs || []);
    } catch (error) {
        console.error("Failed to load libraries:", error);
    }
}

async function loadCollections() {
    try {
        const colls = await fetchCollections();
        renderCollections(colls || []);
    } catch (error) {
        console.error("Failed to load collections:", error);
    }
}

async function renderModels(models, containerId) {
    const container = document.getElementById(containerId);
    if (!models || models.length === 0) {
        container.innerHTML = "<p style=\"color: #666;\">No models found</p>";
        return;
    }
    
    const libs = await fetchLibraries();
    const libMap = {};
    libs.forEach(lib => libMap[lib.id] = lib.name);
    
    const html = models.map((m, idx) => `
        <div class="model-card" onclick="window.viewModelFiles(${m.id})">
            <div class="model-preview" id="${containerId}-preview-${idx}" data-model-id="${m.id}" style="position: relative; width: 100%; height: 300px; background: #0f0f23; border-radius: 8px 8px 0 0; display: flex; align-items: center; justify-content: center; color: #666;">
                <div style="font-size: 48px;">📦</div>
            </div>
            <div class="model-info">
                <div class="model-title">${m.name}</div>
                <div class="model-desc">${m.description || "No description"}</div>
                <div class="model-meta">
                    <span>📁 ${libMap[m.library_id] || "Library " + m.library_id}</span>
                    <span>📅 ${new Date(m.created_at).toLocaleDateString()}</span>
                </div>
            </div>
        </div>
    `).join("");
    
    container.innerHTML = html;
    
    setTimeout(() => {
        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting && !entry.target.dataset.loaded) {
                    entry.target.dataset.loaded = "true";
                    const modelId = entry.target.dataset.modelId;
                    loadModelPreview(modelId, entry.target);
                }
            });
        }, { rootMargin: "50px" });
        
        models.forEach((m, idx) => {
            const previewEl = document.getElementById(containerId + "-preview-" + idx);
            if (previewEl) observer.observe(previewEl);
        });
    }, 200);
}

async function loadModelPreview(modelId, container) {
    try {
        const model = await fetchModel(modelId);
        const files = await fetchModelFiles(modelId);
        
        if (!files || files.length === 0) return;
        
        let previewFile = null;
        if (model.preview_file_id) {
            previewFile = files.find(f => f.id === model.preview_file_id);
        }
        if (!previewFile) {
            previewFile = files.find(f => is3DFile(f.filename) || isImageFile(f.filename));
        }
        
        if (previewFile) {
            if (isImageFile(previewFile.filename)) {
                loadCardImagePreview(previewFile.id, container);
            } else {
                loadCardPreview(previewFile.id, container);
            }
        }
    } catch (error) {
        console.error("Error loading model preview:", error);
    }
}

export async function viewModelFiles(modelId) {
    try {
        const files = await fetchModelFiles(modelId);
        const model = await fetchModel(modelId);
        
        document.getElementById("model-detail-name").textContent = model.name;
        
        const content = document.getElementById("model-detail-content");
        
        if (files.length === 0) {
            content.innerHTML = "<p style=\"color: #666;\">No files found</p>";
        } else {
            const html = files.map((f, idx) => {
                const is3D = is3DFile(f.filename);
                const isImage = isImageFile(f.filename);
                const downloadUrl = getFileDownloadUrl(f.id);
                const isPreview = model.preview_file_id === f.id;
                
                const slicerLinks = is3D ? `
                    <div style="margin-top: 10px;">
                        <strong style="color: #999; font-size: 0.9em;">Open in:</strong>
                        <div style="display: flex; gap: 8px; margin-top: 8px; flex-wrap: wrap;">
                            <a href="prusaslicer://open?file=${encodeURIComponent(downloadUrl)}" style="background: #ff6b35; color: white; padding: 6px 12px; border-radius: 4px; text-decoration: none; font-size: 0.85em;">PrusaSlicer</a>
                            <a href="bambu-studio://open?file=${encodeURIComponent(downloadUrl)}" style="background: #00ae42; color: white; padding: 6px 12px; border-radius: 4px; text-decoration: none; font-size: 0.85em;">Bambu Studio</a>
                            <a href="orcaslicer://open?file=${encodeURIComponent(downloadUrl)}" style="background: #4a90e2; color: white; padding: 6px 12px; border-radius: 4px; text-decoration: none; font-size: 0.85em;">OrcaSlicer</a>
                            <a href="cura://open?file=${encodeURIComponent(downloadUrl)}" style="background: #0066b3; color: white; padding: 6px 12px; border-radius: 4px; text-decoration: none; font-size: 0.85em;">Cura</a>
                        </div>
                    </div>
                ` : "";
                
                return `
                    <div style="background: #16213e; border-radius: 8px; padding: 20px; margin-bottom: 20px;">
                        <div style="display: flex; gap: 20px; align-items: flex-start;">
                            ${(is3D || isImage) ? `
                                <div id="detail-preview-${idx}" data-file-id="${f.id}" data-is3d="${is3D}" data-size="${f.size}" style="width: 300px; height: 300px; background: #0f0f23; border-radius: 8px; flex-shrink: 0; display: flex; align-items: center; justify-content: center; position: relative;">
                                    ${is3D ? (f.size > 10000000 ? `<button onclick="window.loadPreview(${f.id}, ${idx}, true)" style="background: #667eea; color: white; padding: 10px 20px; border-radius: 6px; border: none; cursor: pointer;">Load 3D Preview<br><small>(${(f.size / 1024 / 1024).toFixed(1)} MB)</small></button>` : `<div style="color: #666;">Scroll to load...</div>`) : ""}
                                </div>
                            ` : ""}
                            <div style="flex: 1;">
                                <h3 style="margin: 0 0 10px 0;">${f.filename} ${isPreview ? "<span style=\"color: #4ade80; font-size: 0.8em;\">★ Preview</span>" : ""}</h3>
                                <div style="color: #666; margin-bottom: 15px;">${(f.size / 1024).toFixed(1)} KB</div>
                                <div style="display: flex; gap: 10px; flex-wrap: wrap;">
                                    <a href="${downloadUrl}" download style="background: #667eea; color: white; padding: 10px 20px; border-radius: 6px; text-decoration: none; display: inline-block;">Download</a>
                                    ${(is3D || isImage) && !isPreview ? `<button onclick="window.handleSetPreview(${modelId}, ${f.id})" style="background: #10b981; color: white; padding: 10px 20px; border-radius: 6px; border: none; cursor: pointer;">Set as Preview</button>` : ""}
                                </div>
                                ${slicerLinks}
                            </div>
                        </div>
                    </div>
                `;
            }).join("");
            
            content.innerHTML = html;
            
            // Auto-load images immediately
            files.forEach((f, idx) => {
                const container = document.getElementById("detail-preview-" + idx);
                if (!container) return;
                
                if (isImageFile(f.filename)) {
                    loadImagePreview(f.id, container);
                }
            });
            
            // Lazy load small 3D files when scrolled into view
            const observer = new IntersectionObserver((entries) => {
                entries.forEach(entry => {
                    if (entry.isIntersecting && !entry.target.dataset.loaded) {
                        entry.target.dataset.loaded = "true";
                        const fileId = parseInt(entry.target.dataset.fileId);
                        const is3D = entry.target.dataset.is3d === "true";
                        const size = parseInt(entry.target.dataset.size);
                        
                        // Only auto-load if 3D and under 10MB
                        if (is3D && size <= 10000000) {
                            loadDetailPreview(fileId, entry.target);
                        }
                        observer.unobserve(entry.target);
                    }
                });
            }, { rootMargin: "200px" });
            
            // Only observe small 3D files
            files.forEach((f, idx) => {
                if (is3DFile(f.filename) && f.size <= 10000000) {
                    const container = document.getElementById("detail-preview-" + idx);
                    if (container) {
                        observer.observe(container);
                    }
                }
            });
        }
        
        switchView("model-detail");
    } catch (error) {
        console.error("Failed to view model files:", error);
    }
}

export async function handleSetPreview(modelId, fileId) {
    try {
        await setModelPreview(modelId, fileId);
        viewModelFiles(modelId);
    } catch (error) {
        console.error("Failed to set preview:", error);
    }
}

function renderLibraries(libs) {
    const html = libs.map(lib => `
        <div class="library-card">
            <h3>${lib.name}</h3>
            <p>${lib.path}</p>
            <div class="library-actions">
                <button onclick="window.handleScanLibrary(${lib.id})">Scan</button>
            </div>
        </div>
    `).join("");
    document.getElementById("librariesList").innerHTML = html || "<p style=\"color: #666;\">No libraries found</p>";
}

function renderCollections(colls) {
    const html = colls.map(c => `
        <div class="collection-card">
            <h3>${c.name}</h3>
            <p>${c.description || "No description"}</p>
        </div>
    `).join("");
    document.getElementById("collectionsList").innerHTML = html || "<p style=\"color: #666;\">No collections found</p>";
}

export async function handleScanLibrary(id) {
    try {
        await scanLibrary(id);
        alert("Library scan started");
    } catch (error) {
        console.error("Failed to scan library:", error);
        alert("Failed to scan library");
    }
}

export function loadPreview(fileId, idx, is3D) {
    const container = document.getElementById("detail-preview-" + idx);
    if (!container) return;
    
    if (is3D) {
        loadDetailPreview(fileId, container);
    } else {
        loadImagePreview(fileId, container);
    }
}
