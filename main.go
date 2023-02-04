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

const (
	TB = 1000 * 1000 * 1000 * 1000
	GB = 1000 * 1000 * 1000
	MB = 1000 * 1000
	KB = 1000
)

func traverseDir(hashes, duplicates map[string]string, dupeSize *int64, entries []os.FileInfo, directory string) {
	for _, entry := range entries {
		fullpath := (path.Join(directory, entry.Name()))

		if !entry.Mode().IsDir() && !entry.Mode().IsRegular() {
			continue
		}

		if entry.IsDir() {
			dirFiles, err := ioutil.ReadDir(fullpath)
			if err != nil {
				panic(err)
			}
			traverseDir(hashes, duplicates, dupeSize, dirFiles, fullpath)
			continue
		}
		file, err := ioutil.ReadFile(fullpath)
		if err != nil {
			panic(err)
		}
		h := CalculateHash(file)
		StoreDuplicates(hashes, h, duplicates, fullpath, dupeSize, entry)
	}
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
	if hashEntry, ok := hashes[hashString]; ok {
		duplicates[hashEntry] = fullpath
		atomic.AddInt64(dupeSize, entry.Size())
	} else {
		hashes[hashString] = fullpath
	}
}

func toReadableSize(nbytes int64) string {
	if nbytes > TB {
		return strconv.FormatInt(nbytes/(TB), 10) + " TB"
	}
	if nbytes > GB {
		return strconv.FormatInt(nbytes/(GB), 10) + " GB"
	}
	if nbytes > MB {
		return strconv.FormatInt(nbytes/(MB), 10) + " MB"
	}
	if nbytes > KB {
		return strconv.FormatInt(nbytes/KB, 10) + " KB"
	}
	return strconv.FormatInt(nbytes, 10) + " B"
}

func main() {
	var err error
	dir := flag.String("path", "~/duplicates_files_directory", "the path to traverse searching for duplicates")
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

	traverseDir(hashes, duplicates, &dupeSize, entries, *dir)

	fmt.Println("DUPLICATES")

	fmt.Println("TOTAL FILES:", len(hashes))
	fmt.Println("DUPLICATES:", len(duplicates))
	fmt.Println("TOTAL DUPLICATE SIZE:", toReadableSize(dupeSize))
}

// running into problems of not being able to open directories inside .app folders
