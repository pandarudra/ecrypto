package archive

// Manifest placeholder for future chunked encryption support.
// Currently unused; kept for plan alignment.
type Manifest struct {
    Version int
    Files   []FileEntry
}

type FileEntry struct {
    Name  string
    Size  int64
    Mtime int64
}