package mapfs

import (
	"bytes"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	vfs "golang.org/x/tools/godoc/vfs"
)

// File represents a file's contents and modification time.
type File struct {
	Data    []byte
	ModTime time.Time
}

// Files is the collection of files to init a MapFS file system with. It maps
// the file's full path to its File descriptor. The map only contains files,
// directories will be created implicitly.
type Files map[string]*File

type fsEntry struct {
	name    string
	data    []byte
	entries fsEntries
	modTime time.Time
}

type fsEntries []*fsEntry

// MapFS is a vfs.FileSystem implementation.
type MapFS map[string]*fsEntry

// New creates a new MapFS instance from the input files.
func New(files Files) (vfs.FileSystem, error) {
	fs := make(MapFS)

	for path, file := range files {
		if err := fs.addFile(path, file); err != nil {
			return nil, err
		}
	}

	return fs, nil
}

func fixpath(p string) string {
	p = strings.TrimSuffix(p, "/")
	if !strings.HasPrefix(p, "/") {
		p = strings.Join([]string{"/", p}, "")
	}

	return p
}

func (fs MapFS) addFile(p string, f *File) error {
	p = fixpath(p)

	entry := &fsEntry{path.Base(p), f.Data, nil, f.ModTime}

	_, exists := fs[p]
	if exists {
		return os.ErrExist
	}

	fs[p] = entry

	elems := strings.Split(p, "/")
	for i := len(elems) - 1; i >= 1; i-- {
		ps := strings.Join(elems[0:i], "/")
		if ps == "" {
			ps = "/"
		}

		e, ok := fs[ps]
		if ok {
			if e.data != nil {
				return os.ErrExist
			}
			e.entries = append(e.entries, entry)
			sort.Sort(e.entries)
			if e.modTime.Before(f.ModTime) {
				e.modTime = f.ModTime
			}
			break
		}

		entry = &fsEntry{path.Base(ps), nil, []*fsEntry{entry}, f.ModTime}
		fs[ps] = entry
	}

	return nil
}

func (ents fsEntries) Len() int {
	return len(ents)
}

func (ents fsEntries) Less(i int, j int) bool {
	return strings.Compare(ents[i].name, ents[j].name) < 1
}

func (ents fsEntries) Swap(i int, j int) {
	ents[i], ents[j] = ents[j], ents[i]
}

func (fs MapFS) String() string {
	return "mapfs"
}

// Open implements vfs.Opener.
func (fs MapFS) Open(p string) (vfs.ReadSeekCloser, error) {
	e, ok := fs[fixpath(p)]
	if !ok {
		return nil, os.ErrNotExist
	}
	return nopCloser{bytes.NewReader(e.data)}, nil
}

// Lstat returns the fileinfo of a file or link.
func (fs MapFS) Lstat(p string) (os.FileInfo, error) {
	e, ok := fs[fixpath(p)]
	if ok {
		return e, nil
	}
	return nil, os.ErrNotExist
}

// Stat returns the fileinfo of a file.
func (fs MapFS) Stat(p string) (os.FileInfo, error) {
	return fs.Lstat(p)
}

// ReadDir reads the directory named by path and returns a list of sorted directory entries.
func (fs MapFS) ReadDir(p string) ([]os.FileInfo, error) {
	e, ok := fs[fixpath(p)]
	if !ok {
		return nil, os.ErrNotExist
	}

	fis := make([]os.FileInfo, len(e.entries))
	for i, e := range e.entries {
		fis[i] = os.FileInfo(e)
	}
	return fis, nil
}

func (e *fsEntry) IsDir() bool {
	return e.data == nil
}

func (e *fsEntry) ModTime() time.Time {
	return e.modTime
}

func (e *fsEntry) Mode() os.FileMode {
	if e.IsDir() {
		return 0755 | os.ModeDir
	}
	return 0444
}

func (e *fsEntry) Name() string {
	return e.name
}

func (e *fsEntry) Size() int64 {
	if e.IsDir() {
		return int64(len(e.entries))
	}
	return int64(len(e.data))
}

func (e *fsEntry) Sys() interface{} {
	return nil
}

type nopCloser struct {
	io.ReadSeeker
}

func (nc nopCloser) Close() error { return nil }
