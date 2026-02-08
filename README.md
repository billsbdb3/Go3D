# 3D Model Library - Go Edition

High-performance 3D model library manager. Built for speed.

## Why Go?

- **10-50x faster** than Ruby/Rails
- **10x less memory** usage
- **Instant startup** (<100ms vs 5-10s)
- **True concurrency** for file scanning
- **Single binary** deployment

## Stack

- Go 1.23
- PostgreSQL
- HTMX + Alpine.js
- THREE.js for 3D rendering
- Asynq for background jobs

## Quick Start

```bash
# Install dependencies
make deps

# Set up database
createdb 3d_library
# Run migrations (TODO: add goose)

# Run dev server
make dev
```

Server runs on http://localhost:3000

## Project Structure

```
cmd/
  web/          - Web server
  worker/       - Background job processor
  cli/          - CLI tools

internal/
  models/       - Database models
  handlers/     - HTTP handlers
  services/     - Business logic
  jobs/         - Background jobs
  storage/      - File storage abstraction
  scanner/      - Fast file system scanner
  analyzer/     - 3D file analysis

web/
  templates/    - HTML templates
  static/       - CSS, JS, images
  components/   - Reusable UI components

migrations/     - SQL migrations
```

## Features

### Core
- [x] Fast file scanning (10k+ files/sec)
- [x] Database models
- [ ] 3D model preview
- [ ] Collections & tags
- [ ] Search
- [ ] Background jobs

### Storage
- [ ] Local filesystem
- [ ] S3-compatible storage
- [ ] Multi-storage support

### Analysis
- [ ] File type detection
- [ ] Duplicate detection (SHA256)
- [ ] 3D geometry analysis
- [ ] Thumbnail generation

## Performance Targets

- Scan 10,000 files: <1 second
- Memory usage: <50MB base
- Startup time: <100ms
- API response: <10ms (p95)

## Development

```bash
# Run tests
make test

# Build binaries
make build

# Clean
make clean
```

## TODO

- [ ] Add goose for migrations
- [ ] Implement handlers
- [ ] Add HTMX templates
- [ ] Background job system
- [ ] S3 storage
- [ ] 3D file analysis
- [ ] API endpoints
- [ ] Authentication
