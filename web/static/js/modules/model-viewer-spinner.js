export async function loadDetailPreview(fileId, container) {
    rendererPool.cleanup(container);
    container.innerHTML = `
        <div style="display:flex;flex-direction:column;align-items:center;justify-content:center;height:100%;color:#999;">
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
    
    container.innerHTML = "";
    
    const scene = createScene();
    const camera = new THREE.PerspectiveCamera(50, 1, 0.1, 1000);
    camera.position.set(DETAIL_CAMERA_DISTANCE, DETAIL_CAMERA_DISTANCE, DETAIL_CAMERA_DISTANCE);
    
    const renderer = new THREE.WebGLRenderer({ antialias: true });
    renderer.setSize(PREVIEW_SIZE, PREVIEW_SIZE);
    container.appendChild(renderer.domElement);
    rendererPool.register(container, renderer);
    
    const controls = new OrbitControls(camera, renderer.domElement);
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
        
        let animationId;
        function animate() {
            if (!container._renderer) {
                if (animationId) cancelAnimationFrame(animationId);
                return;
            }
            animationId = requestAnimationFrame(animate);
            controls.update();
            renderer.render(scene, camera);
        }
        animate();
    } catch (error) {
        console.error("Failed to load 3D preview:", error);
        container.innerHTML = "<div style=\"display:flex;align-items:center;justify-content:center;height:100%;color:#f093fb;\">Load failed</div>";
    }
}
