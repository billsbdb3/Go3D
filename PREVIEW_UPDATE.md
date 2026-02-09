# Preview Selection Feature

## Changes Made

### Database
- Added `preview_file_id` column to `models` table
- References `model_files(id)` with ON DELETE SET NULL

### Backend (Go)
1. **models.go** - Added `PreviewFileID *int64` field to Model struct
2. **model.go** - Added `SetPreview()` handler method
3. **main.go** - Added route: `POST /api/models/{id}/preview`

### Frontend (JavaScript)
1. **viewModelFiles()** - Updated to:
   - Display PNG/JPG images in viewport (similar to 3D models)
   - Show "Set as Preview" button for 3D models and images
   - Mark current preview with ★ indicator

2. **loadModelPreview()** - Updated to:
   - Check model.preview_file_id first
   - Fall back to first 3D/image file if no preview set
   - Support both 3D and image previews in cards

3. **New functions:**
   - `loadImagePreview()` - Display images in detail view
   - `loadCardImagePreview()` - Display images in card view
   - `setPreview()` - API call to set preview file

### Usage
1. Upload ZIP with STL/3MF and PNG files
2. View model details
3. Click "Set as Preview" on any 3D model or image file
4. Preview will show on model cards in library view

### API
```bash
# Set preview
POST /api/models/{id}/preview
Content-Type: application/json
{"file_id": 123}

# Clear preview
POST /api/models/{id}/preview
Content-Type: application/json
{"file_id": null}
```

### Cache Version
Updated to `?v=latest8` in index.html

## Auto-Preview Selection (Added)

### Backend Changes
- **upload.go**: Added `setDefaultPreview()` function, called after file upload
- **jobs.go**: Added `setDefaultPreview()` function, called after scanning each model

### Logic
1. First, look for image files (PNG, JPG, JPEG)
2. If no images, look for 3D model files (STL, OBJ, 3MF)
3. Set the first matching file as preview_file_id

### Applied To
- New uploads (ZIP or individual files)
- Library scans (background job)
- Existing models (via SQL script)

All models now automatically get a preview without manual selection.
