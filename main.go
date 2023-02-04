package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"sync/atomic"
)

// Sizes
const (
	KB int64 = 1000
	MB       = 1000 * KB
	GB       = 1000 * MB
	TB       = 1000 * GB
)

// Base conv
const base int64 = 10

// Traverse a given directory recursively and compare file contents to identify duplicate
func traverseDirAndProcessDuplicates(hashes, duplicates map[string]string, dupeSize *int64, entries []os.FileInfo, directory string) {
	for _, entry := range entries {
		fullpath := (path.Join(directory, entry.Name()))

		//
		if !entry.Mode().IsDir() && !entry.Mode().IsRegular() {
			continue
		}

		if entry.IsDir() {
			dirFiles, err := ioutil.ReadDir(fullpath)
			if err != nil {
				panic(err)
			}
			traverseDirAndProcessDuplicates(hashes, duplicates, dupeSize, dirFiles, fullpath)
			continue
		}
		file := ReadFileContent(fullpath)
		h := CalculateHash(file)
		StoreDuplicates(hashes, h, duplicates, fullpath, dupeSize, entry)
	}
}

func ReadFileContent(fullpath string) []byte {
	file, err := ioutil.ReadFile(fullpath)
	if err != nil {
		panic(err)
	}
	return file
}

func CalculateHash(file []byte) string {
	hash := sha1.New()
	if _, err := hash.Write(file); err != nil {
		panic(err)
	}
	hashSum := hash.Sum(nil)
	hashString := fmt.Sprintf("%x", hashSum)
	return hashString
}

func StoreDuplicates(hashes map[string]string, hashString string, duplicates map[string]string, fullpath string, dupeSize *int64, entry fs.FileInfo) {
	// store duplicate file path along with its size
	if hashEntry, ok := hashes[hashString]; ok {
		duplicates[hashEntry] = fullpath
		atomic.AddInt64(dupeSize, entry.Size())
	} else {
		hashes[hashString] = fullpath
	}
}

func toReadableSize(nbytes int64) string {
	var readableSz string
	switch {
	case nbytes > TB:
		readableSz = ConvertByteToSize(nbytes, TB) + " TB"
	case nbytes > GB:
		readableSz = ConvertByteToSize(nbytes, GB) + " GB"
	case nbytes > MB:
		readableSz = ConvertByteToSize(nbytes, MB) + " MB"
	case nbytes > KB:
		readableSz = ConvertByteToSize(nbytes, KB) + " KB"
	default:
		readableSz = ConvertByteToSize(nbytes, 1) + " B"
	}
	return readableSz
}

func ConvertByteToSize(n int64, size int64) string {
	return strconv.FormatInt(n/size, int(base))
}

func main() {
	var err error
	dir := flag.String("path", "", "the path to traverse searching for duplicates")
	flag.Parse()

	if *dir == "" {
		*dir, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}

	hashes := map[string]string{}
	duplicates := map[string]string{}
	var dupeSize int64

	entries, err := ioutil.ReadDir(*dir)
	if err != nil {
		panic(err)
	}

	traverseDirAndProcessDuplicates(hashes, duplicates, &dupeSize, entries, *dir)

	fmt.Println("DUPLICATES")

	fmt.Println("TOTAL FILES:", len(hashes))
	fmt.Println("DUPLICATES:", len(duplicates))
	fmt.Println("TOTAL DUPLICATE SIZE:", toReadableSize(dupeSize))
}

// running into problems of not being able to open directories inside .app folders
