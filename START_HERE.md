# 3D Library - Quick Start Guide

## Start the Application

```bash
# SSH into LXC 104
pct enter 104

# Navigate to project
cd /root/3d-library

# Start both services (web + worker)
./start-all.sh
```

**Access:** http://192.168.3.26:3000

---

## Stop the Application

```bash
# Kill all processes
pkill -9 -f '3d-library'
pkill -9 -f 'go run'
```

---

## Check Status

```bash
# View server log
tail -f /tmp/server.log

# View worker log
tail -f /tmp/worker.log

# Check if running
ps aux | grep -E 'go run|3d-library'
```

---

## Rebuild After Changes

```bash
cd /root/3d-library
export PATH=/usr/local/go/bin:$PATH

# Stop old processes
pkill -9 -f 'go run'

# Restart
./start-all.sh
```

---

## Project Location

- **Host:** Proxmox (skynet)
- **Container:** LXC 104 (dev-course)
- **IP:** 192.168.3.26
- **Path:** `/root/3d-library`
- **Port:** 3000

---

## Key Files

- `cmd/web/main.go` - Web server
- `cmd/worker/main.go` - Background worker
- `internal/handlers/` - API handlers
- `web/static/` - UI files
- `.env` - Configuration
- `FINAL.md` - Complete documentation

---

## Database

```bash
# Access database
sudo -u postgres psql library3d

# Check data
SELECT COUNT(*) FROM models;
SELECT COUNT(*) FROM libraries;
```

---

## Quick Test

```bash
# Create library
curl -X POST http://192.168.3.26:3000/api/libraries \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","path":"/tmp/test","storage":"local"}'

# List libraries
curl http://192.168.3.26:3000/api/libraries
```

---

## If Something Breaks

1. Check logs: `tail -f /tmp/server.log`
2. Check if Postgres running: `systemctl status postgresql`
3. Check if Redis running: `systemctl status redis`
4. Restart everything: `./start-all.sh`

---

## Next Session Goals

- [ ] Add 3D model preview (THREE.js)
- [ ] File download endpoints
- [ ] Thumbnail generation
- [ ] Docker deployment
- [ ] Authentication

---

**That's it! Just run `./start-all.sh` and you're good to go.**
