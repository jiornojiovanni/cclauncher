package archives

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

const (
	Parentfolder string = "cataclysm"
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

	var topfolder string
	tar := tar.NewReader(archive)
	for {
		header, err := tar.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(Parentfolder, header.Name)
		var newpath string

		info := header.FileInfo()

		/* Here we find the name of the folder inside the downloaded archive (eg. cataclysmdda-0.E).
		 * This name change at every stable release, so we need to find it at runtime.
		 */
		if filepath.Dir(path) == Parentfolder {
			topfolder = filepath.Base(path)
			continue
		} else {
			//We remove that name so the top folder it's always a generic "cataclysm".
			newpath, err = filepath.Rel(Parentfolder+"/"+topfolder, path)
			if err != nil {
				return err
			}

		}

		if info.IsDir() {
			if err = os.MkdirAll(Parentfolder+"/"+newpath, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(Parentfolder+"/"+newpath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
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

//Compress create a zip archive, with the flag all set to false it will only backup Saves, Sound, Config and Font.
func Compress(name string, all bool) error {
	file, err := os.Create(name + ".zip")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := zip.NewWriter(file)
	defer writer.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		f, err := writer.Create(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}

	if all {
		err = filepath.Walk(name, walker)
		if err != nil {
			return err
		}
	} else {
		err = filepath.Walk(name+"/save", walker)
		if err != nil {
			return err
		}

		err = filepath.Walk(name+"/config", walker)
		if err != nil {
			return err
		}

		err = filepath.Walk(name+"/sound", walker)
		if err != nil {
			return err
		}

		err = filepath.Walk(name+"/font", walker)
		if err != nil {
			return err
		}
	}
	return err
}
