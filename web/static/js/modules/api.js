import { API_BASE } from "./config.js";

export async function fetchModels() {
    const response = await fetch(`${API_BASE}/models`);
    if (!response.ok) throw new Error("Failed to fetch models");
    return response.json();
}

export async function fetchModel(id) {
    const response = await fetch(`${API_BASE}/models/${id}`);
    if (!response.ok) throw new Error("Failed to fetch model");
    return response.json();
}

export async function fetchModelFiles(id) {
    const response = await fetch(`${API_BASE}/models/${id}/files`);
    if (!response.ok) throw new Error("Failed to fetch model files");
    return response.json();
}

export async function fetchFile(id) {
    const response = await fetch(`${API_BASE}/files/${id}`);
    if (!response.ok) throw new Error("Failed to fetch file");
    return response.json();
}

export async function fetchLibraries() {
    const response = await fetch(`${API_BASE}/libraries`);
    if (!response.ok) throw new Error("Failed to fetch libraries");
    return response.json();
}

export async function fetchCollections() {
    const response = await fetch(`${API_BASE}/collections`);
    if (!response.ok) throw new Error("Failed to fetch collections");
    return response.json();
}

export async function fetchTags() {
    const response = await fetch(`${API_BASE}/tags`);
    if (!response.ok) throw new Error("Failed to fetch tags");
    return response.json();
}

export async function setModelPreview(modelId, fileId) {
    const response = await fetch(`${API_BASE}/models/${modelId}/preview`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ file_id: fileId })
    });
    if (!response.ok) throw new Error("Failed to set preview");
}

export async function scanLibrary(id) {
    const response = await fetch(`${API_BASE}/libraries/${id}/scan`, {
        method: "POST"
    });
    if (!response.ok) throw new Error("Failed to scan library");
    return response.json();
}

export function getFileDownloadUrl(fileId) {
    return `${API_BASE}/files/${fileId}/download`;
}
