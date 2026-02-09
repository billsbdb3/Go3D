# 3D Library - Documentation Index

## Quick Start
- **START_HERE.md** - Quick commands to start/stop the application
- **README.md** - Main documentation with setup and usage

## Feature Documentation
- **PREVIEW_UPDATE.md** - Preview selection feature details
- **API_COMPLETE.md** - Complete API endpoint documentation

## Development
- **REFACTORING_COMPLETE.md** - Summary of best practices refactoring
- **REFACTORING.md** - Detailed refactoring notes
- **BEST_PRACTICES.md** - Best practices checklist and comparison

## Status
- **FINAL.md** - Complete build summary (pre-refactoring)
- **READY.md** - Production readiness checklist (pre-refactoring)
- **NEXT_STEPS.md** - Future enhancement ideas

## Archived
- **README-old.md** - Original README (backup)

## Project Structure

### Frontend (Modular)
```
web/static/js/
├── app.js                    # Entry point
├── three-loader.js           # THREE.js loader
└── modules/
    ├── config.js             # Configuration
    ├── api.js                # API client
    ├── renderer-pool.js      # WebGL management
    ├── three-utils.js        # THREE.js utilities
    ├── model-viewer.js       # 3D rendering
    └── ui.js                 # UI components
```

### Backend (Go)
```
cmd/
├── web/main.go              # HTTP server
└── worker/main.go           # Background worker

internal/
├── config/                  # Configuration
├── database/                # Database
├── handlers/                # HTTP handlers
├── jobs/                    # Background jobs
├── models/                  # Data models
└── scanner/                 # File scanner
```

## Key Features

- 3D preview (STL, OBJ, 3MF)
- Image preview (PNG, JPG)
- ZIP upload with extraction
- Slicer integration (4 slicers)
- Background job processing
- Auto-preview selection
- REST API (26 endpoints)

## Access

- **URL:** http://192.168.3.26:3000
- **GitHub:** https://github.com/billsbdb3/Go3D

## Recent Changes

### Best Practices Refactoring (Latest)
- Modular frontend architecture
- WebGL context pooling
- Centralized configuration
- Comprehensive documentation
- Improved error handling

See **REFACTORING_COMPLETE.md** for details.
