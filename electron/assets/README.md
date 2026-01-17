# Icon Placeholder

Place your application icons here:

- **icon.png** (512x512 PNG) - For Linux AppImage
- **icon.ico** (256x256 ICO) - For Windows installer
- **icon.icns** (512x512 ICNS) - For macOS .app

## Quick Icon Generation

You can use online tools or ImageMagick:

### From PNG to ICO (Windows)

```bash
magick icon.png -define icon:auto-resize=256,128,64,48,32,16 icon.ico
```

### From PNG to ICNS (macOS)

```bash
mkdir icon.iconset
sips -z 16 16 icon.png --out icon.iconset/icon_16x16.png
sips -z 32 32 icon.png --out icon.iconset/icon_16x16@2x.png
# ... (more sizes)
iconutil -c icns icon.iconset
```

## Temporary Placeholder

For now, the app will use the default Electron icon. Replace these files before publishing.
