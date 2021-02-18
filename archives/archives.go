package archives

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

//Extract a tar.gz archive to the root of the folder
func Extract(name string) error {
	gz, err := os.Open(name)
	if err != nil {
		return err
	}
	defer gz.Close()

	archive, err := gzip.NewReader(gz)
	if err != nil {
		return err
	}
	defer archive.Close()

	tar := tar.NewReader(archive)
	for {
		header, err := tar.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join("./", header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tar)
		if err != nil {
			return err
		}
	}
	return nil
}
