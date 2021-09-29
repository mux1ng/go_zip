package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	excludeExt = []string{".jpg", ".png"}
)

func main() {
	defaultOutput := fmt.Sprintf("%v.zip", time.Now().Unix())
	srcPath := flag.String("i", "", "the source file")
	dstPath := flag.String("o", defaultOutput, "the output file")
	flag.Parse()
	goZip(*srcPath, *dstPath)
}

func goZip(srcPath string, dstPath string) {
	os.RemoveAll(dstPath)
	zipFile, _ := os.Create(dstPath)
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	filepath.Walk(srcPath, func(path string, info fs.FileInfo, _ error) error {
		srcPath, _ = filepath.Abs(srcPath)
		path, _ = filepath.Abs(path)

		for _, ext := range excludeExt {
			if filepath.Ext(path) == ext {
				return nil
			}
		}

		header, _ := zip.FileInfoHeader(info)
		header.Name = filepath.Base(srcPath) + strings.TrimPrefix(path, srcPath)
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, _ := archive.CreateHeader(header)

		if !info.IsDir() {
			file, _ := os.Open(path)
			defer file.Close()
			io.Copy(writer, file)
		}
		return nil
	})
}
