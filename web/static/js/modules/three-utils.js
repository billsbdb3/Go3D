export async function waitForThree() {
    while (typeof THREE === "undefined") {
        await new Promise(resolve => setTimeout(resolve, 100));
    }
}

export function createScene() {
    const scene = new THREE.Scene();
    scene.background = new THREE.Color(0x0f0f23);
    return scene;
}

export function createLights() {
    const ambient = new THREE.AmbientLight(0xffffff, 0.6);
    const directional = new THREE.DirectionalLight(0xffffff, 0.8);
    directional.position.set(1, 1, 1);
    return [ambient, directional];
}

export function centerAndScaleObject(object, targetSize) {
    const box = new THREE.Box3().setFromObject(object);
    const center = new THREE.Vector3();
    box.getCenter(center);
    
    object.traverse((child) => {
        if (child.isMesh && child.geometry) {
            child.geometry.translate(-center.x, -center.y, -center.z);
        }
    });
    
    box.setFromObject(object);
    const size = new THREE.Vector3();
    box.getSize(size);
    const maxDim = Math.max(size.x, size.y, size.z);
    const scale = targetSize / maxDim;
    
    object.scale.setScalar(scale);
    return { size, scale };
}

export function centerAndScaleGeometry(geometry, targetSize) {
    geometry.computeBoundingBox();
    const center = new THREE.Vector3();
    geometry.boundingBox.getCenter(center);
    geometry.translate(-center.x, -center.y, -center.z);
    geometry.computeBoundingBox();
    
    const size = new THREE.Vector3();
    geometry.boundingBox.getSize(size);
    const maxDim = Math.max(size.x, size.y, size.z);
    const scale = targetSize / maxDim;
    
    return { size, scale };
}

export function getFileExtension(filename) {
    return filename.toLowerCase().split(".").pop();
}

export function is3DFile(filename) {
    const ext = getFileExtension(filename);
    return ["stl", "obj", "3mf"].includes(ext);
}

export function isImageFile(filename) {
    const ext = getFileExtension(filename);
    return ["png", "jpg", "jpeg"].includes(ext);
}
