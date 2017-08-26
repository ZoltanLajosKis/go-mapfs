# go-mapfs

[![Build Status](https://travis-ci.org/ZoltanLajosKis/go-mapfs.svg?branch=master)](https://travis-ci.org/ZoltanLajosKis/go-mapfs)
[![Go Report Card](https://goreportcard.com/badge/github.com/ZoltanLajosKis/go-mapfs)](https://goreportcard.com/report/github.com/ZoltanLajosKis/go-mapfs)
[![Coverage Status](https://coveralls.io/repos/github/ZoltanLajosKis/go-mapfs/badge.svg?branch=master)](https://coveralls.io/github/ZoltanLajosKis/go-mapfs?branch=master)
[![GoDoc](https://godoc.org/github.com/ZoltanLajosKis/go-mapfs?status.svg)](https://godoc.org/github.com/ZoltanLajosKis/go-mapfs)

A map-based [vfs.FileSystem][vfsfs] implementation that stores both the
contents and the modification times for files.


## Usage
Each input file is represented by the `File` structure that stores the contents
and the modification time of the file.
```go
type File struct {
  Data    []byte
  ModTime time.Time
}
```

These `File`s are collected in a `Files` map that maps each file to the file's
path.
```go
type Files map[string]*File
```

The `New` function creates the file system from the map.
```go
func New(Files) (vfs.FileSystem, error)
```


### Example
```go
package main

import (
  "time"

  "github.com/ZoltanLajosKis/go-mapfs"
)

func main() {
  files := mapfs.Files{
    "test/hello.txt": {
      []byte("Hello."),
      time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC),
    },
    // additional files ...
  }

  fs, err := mapfs.New(files)
}
```


[vfsfs]: https://godoc.org/golang.org/x/tools/godoc/vfs#FileSystem
