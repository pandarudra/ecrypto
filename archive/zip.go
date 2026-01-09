package archive

import (
	"archive/zip"
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// ZipFolder compresses a folder into a ZIP archive (bytes).
func ZipFolder(root string) ([]byte, error) {
    buf := &bytes.Buffer{}
    zw := zip.NewWriter(buf)
    defer zw.Close()

    root = filepath.Clean(root)
    err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if d.IsDir() {
            return nil
        }

        rel, err := filepath.Rel(root, path)
        if err != nil {
            return err
        }

        info, err := d.Info()
        if err != nil {
            return err
        }

        hdr := &zip.FileHeader{
            Name:   filepath.ToSlash(rel),
            Method: zip.Deflate,
        }
        hdr.SetModTime(info.ModTime())

        w, err := zw.CreateHeader(hdr)
        if err != nil {
            return err
        }

        f, err := os.Open(path)
        if err != nil {
            return err
        }
        defer f.Close()

        _, err = io.Copy(w, f)
        return err
    })
    if err != nil {
        return nil, err
    }

    if err := zw.Close(); err != nil {
        return nil, err
    }

    return buf.Bytes(), nil
}

// UnzipTo extracts a ZIP archive (bytes) to a target folder.
func UnzipTo(outDir string, zipBytes []byte) error {
    readerAt := bytes.NewReader(zipBytes)
    zr, err := zip.NewReader(readerAt, int64(len(zipBytes)))
    if err != nil {
        return err
    }

    for _, f := range zr.File {
        destPath := filepath.Join(outDir, filepath.FromSlash(f.Name))

        if f.FileInfo().IsDir() {
            if err := os.MkdirAll(destPath, 0o755); err != nil {
                return err
            }
            continue
        }

        if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
            return err
        }

        rc, err := f.Open()
        if err != nil {
            return err
        }

        tmp := destPath + ".tmp"
        df, err := os.Create(tmp)
        if err != nil {
            rc.Close()
            return err
        }

        if _, err := io.Copy(df, rc); err != nil {
            rc.Close()
            df.Close()
            return err
        }

        rc.Close()
        if err := df.Close(); err != nil {
            return err
        }

        if err := os.Rename(tmp, destPath); err != nil {
            return err
        }
    }

    return nil
}