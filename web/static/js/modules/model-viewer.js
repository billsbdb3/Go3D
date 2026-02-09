import { fetchFile } from "./api.js";
import { rendererPool } from "./renderer-pool.js";
import { getFileDownloadUrl } from "./api.js";
import { PREVIEW_SIZE, DETAIL_CAMERA_DISTANCE, CARD_CAMERA_DISTANCE, MODEL_SCALE } from "./config.js";
import { waitForThree, createScene, createLights, getFileExtension } from "./three-utils.js";

export async function loadDetailPreview(fileId, container) {
    rendererPool.cleanup(container);
    container.innerHTML = `
        <div class="loading-spinner" style="display:flex;flex-direction:column;align-items:center;justify-content:center;height:100%;color:#999;">
            <div style="border: 4px solid #333; border-top: 4px solid #667eea; border-radius: 50%; width: 40px; height: 40px; animation: spin 1s linear infinite;"></div>
            <div style="margin-top: 10px;">Loading...</div>
        </div>
    `;
    
    if (!document.getElementById("spinner-style")) {
        const style = document.createElement("style");
        style.id = "spinner-style";
        style.textContent = "@keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }";
        document.head.appendChild(style);
    }
    
    await waitForThree();
    
    const scene = createScene();
    const camera = new THREE.PerspectiveCamera(50, 1, 0.1, 1000);
    camera.position.set(DETAIL_CAMERA_DISTANCE, DETAIL_CAMERA_DISTANCE, DETAIL_CAMERA_DISTANCE);
    
    const controls = new OrbitControls(camera, container);
    controls.enableDamping = true;
    controls.dampingFactor = 0.05;
    
    const gridHelper = new THREE.GridHelper(100, 20, 0x444444, 0x222222);
    scene.add(gridHelper);
    
    const axesHelper = new THREE.AxesHelper(30);
    scene.add(axesHelper);
    
    const lights = createLights();
    lights.forEach(light => scene.add(light));
    
    try {
        const fileInfo = await fetchFile(fileId);
        const ext = getFileExtension(fileInfo.filename);
        const url = getFileDownloadUrl(fileId);
        
        if (ext === "3mf") {
            await load3MF(url, scene, controls, MODEL_SCALE);
        } else {
            await loadSTL(url, scene, controls, MODEL_SCALE);
        }
        
        // Model loaded, now create canvas and remove spinner
        container.innerHTML = "";
        
        const canvas = document.createElement("canvas");
        canvas.width = PREVIEW_SIZE;
        canvas.height = PREVIEW_SIZE;
        canvas.style.width = PREVIEW_SIZE + "px";
        canvas.style.height = PREVIEW_SIZE + "px";
        container.appendChild(canvas);
        
        const renderer = rendererPool.getRenderer();
        let animationId;
        let isActive = true;
        
        function animate() {
            if (!isActive) {
                if (animationId) cancelAnimationFrame(animationId);
                return;
            }
            animationId = requestAnimationFrame(animate);
            controls.update();
            renderer.setSize(PREVIEW_SIZE, PREVIEW_SIZE);
            renderer.render(scene, camera);
            
            const ctx = canvas.getContext("2d");
            ctx.drawImage(renderer.domElement, 0, 0);
        }
        
        container._stopAnimation = () => { isActive = false; };
        animate();
    } catch (error) {
        console.error("Failed to load 3D preview:", error);
        container.innerHTML = "<div style=\"display:flex;align-items:center;justify-content:center;height:100%;color:#f093fb;\">Load failed</div>";
    }
}

export async function loadCardPreview(fileId, container) {
    rendererPool.cleanup(container);
    container.innerHTML = "";
    
    await waitForThree();
    
    const canvas = document.createElement("canvas");
    canvas.width = container.clientWidth;
    canvas.height = container.clientHeight;
    canvas.style.width = "100%";
    canvas.style.height = "100%";
    container.appendChild(canvas);
    
    const scene = createScene();
    const camera = new THREE.PerspectiveCamera(50, container.clientWidth / container.clientHeight, 0.1, 1000);
    camera.position.set(CARD_CAMERA_DISTANCE, CARD_CAMERA_DISTANCE, CARD_CAMERA_DISTANCE);
    camera.lookAt(0, 0, 0);
    
    const lights = createLights();
    lights.forEach(light => scene.add(light));
    
    try {
        const fileInfo = await fetchFile(fileId);
        const ext = getFileExtension(fileInfo.filename);
        const url = getFileDownloadUrl(fileId);
        
        let object;
        if (ext === "3mf") {
            object = await load3MF(url, scene, null, MODEL_SCALE);
        } else if (ext === "obj") {
            object = await loadOBJ(url, scene, null, MODEL_SCALE);
        } else {
            object = await loadSTL(url, scene, null, MODEL_SCALE);
        }
        
        const renderer = rendererPool.getRenderer();
        let animationId;
        let isActive = true;
        
        function animate() {
            if (!isActive) {
                if (animationId) cancelAnimationFrame(animationId);
                return;
            }
            animationId = requestAnimationFrame(animate);
            if (object) object.rotation.z += 0.005;
            renderer.setSize(canvas.width, canvas.height);
            renderer.render(scene, camera);
            
            const ctx = canvas.getContext("2d");
            ctx.drawImage(renderer.domElement, 0, 0);
        }
        
        container._stopAnimation = () => { isActive = false; };
        animate();
    } catch (error) {
        console.error("Failed to load card preview:", error);
    }
}

export function loadImagePreview(fileId, container) {
    const img = document.createElement("img");
    img.src = getFileDownloadUrl(fileId);
    img.style.maxWidth = "100%";
    img.style.maxHeight = "100%";
    img.style.objectFit = "contain";
    container.appendChild(img);
}

export function loadCardImagePreview(fileId, container) {
    container.innerHTML = "";
    const img = document.createElement("img");
    img.src = getFileDownloadUrl(fileId);
    img.style.width = "100%";
    img.style.height = "100%";
    img.style.objectFit = "cover";
    container.appendChild(img);
}

async function load3MF(url, scene, controls, targetSize) {
    const loader = new ThreeMFLoader();
    return new Promise((resolve, reject) => {
        loader.load(url, (object) => {
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
            object.rotation.x = -Math.PI / 2;
            
            const heightAboveGrid = (size.z * scale) / 2;
            object.position.set(0, heightAboveGrid, 0);
            
            scene.add(object);
            
            if (controls) {
                controls.target.set(0, heightAboveGrid, 0);
                controls.update();
            }
            
            resolve(object);
        }, undefined, reject);
    });
}

async function loadOBJ(url, scene, controls, targetSize) {
    const loader = new OBJLoader();
    return new Promise((resolve, reject) => {
        loader.load(url, (object) => {
            const box = new THREE.Box3().setFromObject(object);
            const center = new THREE.Vector3();
            box.getCenter(center);
            
            object.traverse((child) => {
                if (child.isMesh && child.geometry) {
                    child.geometry.translate(-center.x, -center.y, -center.z);
                    child.material = new THREE.MeshPhongMaterial({ color: 0xcccccc });
                }
            });
            
            box.setFromObject(object);
            const size = new THREE.Vector3();
            box.getSize(size);
            const maxDim = Math.max(size.x, size.y, size.z);
            const scale = targetSize / maxDim;
            object.scale.setScalar(scale);
            
            object.rotation.x = -Math.PI / 2;
            scene.add(object);
            if (controls) controls.update();
            resolve(object);
        }, undefined, reject);
    });
}

async function loadSTL(url, scene, controls, targetSize) {
    const loader = new STLLoader();
    return new Promise((resolve, reject) => {
        loader.load(url, (geometry) => {
            geometry.computeBoundingBox();
            const center = new THREE.Vector3();
            geometry.boundingBox.getCenter(center);
            geometry.translate(-center.x, -center.y, -center.z);
            
            const size = new THREE.Vector3();
            geometry.boundingBox.getSize(size);
            const maxDim = Math.max(size.x, size.y, size.z);
            const scale = targetSize / maxDim;
            
            const material = new THREE.MeshPhongMaterial({ color: 0xcccccc });
            const mesh = new THREE.Mesh(geometry, material);
            mesh.scale.setScalar(scale);
            mesh.rotation.x = -Math.PI / 2;
            mesh.position.y = (size.z * scale) / 2;
            
            scene.add(mesh);
            
            if (controls) {
                controls.target.set(0, (size.z * scale) / 2, 0);
                controls.update();
            }
            
            resolve(mesh);
        }, undefined, reject);
    });
}
