# 3D Library - Quick Start Guide

## Start the Application

```bash
# Check if running
systemctl status 3d-library

# Start if not running
systemctl start 3d-library

# View logs
journalctl -u 3d-library -f
```

**Access:** http://192.168.3.26:3000

---

## Stop the Application

```bash
systemctl stop 3d-library
```

---

## Restart After Changes

```bash
cd /root/3d-library
systemctl restart 3d-library
journalctl -u 3d-library -f
```

---

## Project Location

- **Host:** Proxmox (skynet)
- **Container:** LXC 104 (dev-course)
- **IP:** 192.168.3.26
- **Path:** `/root/3d-library`
- **Port:** 3000
- **GitHub:** https://github.com/billsbdb3/Go3D

---

## Key Features

✅ 3D Preview (STL, OBJ, 3MF)
✅ ZIP extraction with directory preservation
✅ Slicer integration (PrusaSlicer, Bambu Studio, OrcaSlicer, Cura)
✅ Full REST API (26 endpoints)
✅ Background job processing
✅ Professional web UI

---

## Database

```bash
# Access database
PGPASSWORD=dev123 psql -h localhost -U library3d -d library3d

# Check data
SELECT COUNT(*) FROM models;
SELECT COUNT(*) FROM model_files;
```

**Credentials:**
- Database: `library3d`
- User: `library3d`
- Password: `dev123`

---

## Git Operations

```bash
cd /root/3d-library
git status
git add -A
git commit -m "Your message"
git push origin main
```

---

## 3D Viewer Details

- **THREE.js:** Latest version (r170+)
- **Formats:** STL, OBJ, 3MF
- **Features:** Interactive controls, auto-centering, grid floor
- **Material:** Light gray (0xcccccc)
- **Cache:** `?v=latest7`

---

## If Something Breaks

1. Check logs: `journalctl -u 3d-library -f`
2. Check if Postgres running: `systemctl status postgresql`
3. Check if Redis running: `systemctl status redis`
4. Restart: `systemctl restart 3d-library`

---

**That's it! Service runs automatically on boot.**
