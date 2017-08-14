package mapfs

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"golang.org/x/tools/godoc/vfs"
)

func TestOpenRoot(t *testing.T) {
	files := make(Files)
	files["foo/bar/three.txt"] = &File{[]byte("a"), time.Now()}
	files["foo/bar.txt"] = &File{[]byte("b"), time.Now()}
	files["top.txt"] = &File{[]byte("c"), time.Now()}
	files["other-top.txt"] = &File{[]byte("d"), time.Now()}
	fs, _ := New(files)

	testRead(t, fs, "/foo/bar/three.txt", []byte("a"))
	testRead(t, fs, "foo/bar/three.txt", []byte("a"))
	testRead(t, fs, "foo/bar.txt", []byte("b"))
	testRead(t, fs, "top.txt", []byte("c"))
	testRead(t, fs, "/top.txt", []byte("c"))
	testRead(t, fs, "other-top.txt", []byte("d"))
	testRead(t, fs, "/other-top.txt", []byte("d"))

	p := "/xxxx"
	_, err := fs.Open(p)
	if !os.IsNotExist(err) {
		t.Errorf("Read(%q) = %v; want os.IsNotExist error", p, err)
	}
}

func testRead(t *testing.T, fs vfs.FileSystem, p string, data []byte) {
	r, err := fs.Open(p)
	if err != nil {
		t.Errorf("Open(%q) = %v", p, err)
		return
	}
	defer r.Close()

	fdata, err := ioutil.ReadAll(r)
	if err != nil {
		t.Error(err)
		return
	}

	assertEqual(t, string(data), string(fdata))
}

func TestReadDir(t *testing.T) {
	ts := []time.Time{time.Unix(1200000000, 0), time.Unix(1300000000, 0),
		time.Unix(1400000000, 0), time.Unix(1500000000, 0)}

	files := make(Files)
	files["foo/bar/three.txt"] = &File{[]byte("333"), ts[0]}
	files["foo/bar.txt"] = &File{[]byte("22"), ts[1]}
	files["top.txt"] = &File{[]byte("top.txt file"), ts[2]}
	files["other-top.txt"] = &File{[]byte("other-top.txt file"), ts[3]}
	fs, _ := New(files)

	dir1 := "/"
	fis1, err := fs.ReadDir(dir1)
	if err != nil {
		t.Errorf("ReadDir(%q) = %v", dir1, err)
		return
	}

	assertEqual(t, len(fis1), 3)
	assertFI(t, fis1[0], "foo", 2, ts[1], true)
	assertFI(t, fis1[1], "other-top.txt", len("other-top.txt file"), ts[3], false)
	assertFI(t, fis1[2], "top.txt", len("top.txt file"), ts[2], false)

	dir2 := "/foo"
	fis2, err := fs.ReadDir(dir2)
	if err != nil {
		t.Errorf("ReadDir(%q) = %v", dir2, err)
		return
	}

	assertEqual(t, len(fis2), 2)
	assertFI(t, fis2[0], "bar", 1, ts[0], true)
	assertFI(t, fis2[1], "bar.txt", 2, ts[1], false)

	dir3 := "/foo/"
	fis3, err := fs.ReadDir(dir3)
	if err != nil {
		t.Errorf("ReadDir(%q) = %v", dir3, err)
		return
	}

	assertEqual(t, len(fis3), 2)
	assertFI(t, fis3[0], "bar", 1, ts[0], true)
	assertFI(t, fis3[1], "bar.txt", 2, ts[1], false)

	dir4 := "/foo/bar"
	fis4, err := fs.ReadDir(dir4)
	if err != nil {
		t.Errorf("ReadDir(%q) = %v", dir4, err)
		return
	}

	assertEqual(t, len(fis4), 1)
	assertFI(t, fis4[0], "three.txt", 3, ts[0], false)

	dir5 := "/xxxx"
	_, err = fs.ReadDir(dir5)
	if !os.IsNotExist(err) {
		t.Errorf("ReadDir (%q) = %v; want os.IsNotExist error", dir5, err)
	}
}

func assertFI(t *testing.T, fi os.FileInfo, name string, size int, modTime time.Time, dir bool) {
	assertEqual(t, fi.Name(), name)
	assertEqual(t, fi.Size(), int64(size))
	assertEqual(t, fi.ModTime(), modTime)
	assertEqual(t, fi.IsDir(), dir)
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	t.Fatal(fmt.Sprintf("%v != %v", a, b))
}
