// Single shared renderer for all previews
class RendererPool {
    constructor() {
        this.sharedRenderer = null;
    }

    getRenderer() {
        if (!this.sharedRenderer) {
            this.sharedRenderer = new THREE.WebGLRenderer({ antialias: true });
        }
        return this.sharedRenderer;
    }

    cleanup(container) {
        // Just stop animation, keep renderer
        if (container._stopAnimation) {
            container._stopAnimation();
        }
    }

    cleanupAll() {
        // Dispose shared renderer only when really needed
        if (this.sharedRenderer) {
            this.sharedRenderer.dispose();
            this.sharedRenderer = null;
        }
    }
}

export const rendererPool = new RendererPool();
