package main

import (
	"clean-code-workshop/constants"
	"crypto/sha1"
	"errors"
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

type duplicateFileInfo struct {
	hashes     map[string]string
	duplicates map[string]string
	dupeSize   *int64
}

// Traverse a given directory recursively and compare file contents to identify duplicate
func traverseDirAndProcessDuplicates(dmap duplicateFileInfo, entries []os.FileInfo, directory string) error {
	for _, entry := range entries {
		fullpath := (path.Join(directory, entry.Name()))

		// Ignore directory or file
		if !entry.Mode().IsDir() && !entry.Mode().IsRegular() {
			continue
		}

		if entry.IsDir() {
			dirFiles, err := ioutil.ReadDir(fullpath)
			if err != nil {
				return errors.New(constants.FAILEDTOREADDIR)
			}
			traverseDirAndProcessDuplicates(dmap, dirFiles, fullpath)
			continue
		}
		file, err := ReadFileContent(fullpath)
		if err != nil {
			return err
		}
		hash, err := CalculateHash(file)
		if err != nil {
			fmt.Printf(constants.FAILEDTOCALCHASH)
		}
		StoreDuplicates(dmap, hash, fullpath, entry)
	}
	return nil
}

func ReadFileContent(fullpath string) ([]byte, error) {
	file, err := ioutil.ReadFile(fullpath)
	if err != nil {
		return nil, errors.New(constants.FAILEDTOREADFILE)
	}
	// defer

	return file, nil
}

func CalculateHash(file []byte) (string, error) {
	hash := sha1.New()
	if _, err := hash.Write(file); err != nil {
		return "", errors.New(constants.FAILEDTOWRITEHASH)
	}
	hashSum := hash.Sum(nil)
	hashString := fmt.Sprintf("%x", hashSum)
	return hashString, nil
}

func StoreDuplicates(dmap duplicateFileInfo, hashString string, fullpath string, entry fs.FileInfo) {
	// store duplicate file path along with its size
	if hashEntry, ok := dmap.hashes[hashString]; ok {
		dmap.duplicates[hashEntry] = fullpath
		atomic.AddInt64(dmap.dupeSize, entry.Size())
	} else {
		dmap.hashes[hashString] = fullpath
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
			fmt.Printf(constants.FAILEDTOGETCURRDIR)
		}
	}

	entries, err := ioutil.ReadDir(*dir)
	if err != nil {
		fmt.Printf(constants.FAILEDTOREADDIR)
	}
	d := duplicateFileInfo{}
	err = traverseDirAndProcessDuplicates(d, entries, *dir)
	if err != nil {
		fmt.Printf("Unexpected error: %s", err.Error())
	}
	fmt.Println("DUPLICATES")

	fmt.Println("TOTAL FILES:", len(d.hashes))
	fmt.Println("DUPLICATES:", len(d.duplicates))
	fmt.Println("TOTAL DUPLICATE SIZE:", toReadableSize(*d.dupeSize))
}

// running into problems of not being able to open directories inside .app folders
