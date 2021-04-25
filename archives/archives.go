package archives

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	//ParentFolder is the name of the cdda directory
	ParentFolder string = "cataclysm"
	TempFolder   string = "previous_version"
	Suffix       string = ".zip"
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

	var folder string
	reader := tar.NewReader(archive)
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(ParentFolder, header.Name)
		var newPath string

		info := header.FileInfo()

		/* Here we find the name of the folder inside the downloaded archive (eg. cataclysm-0.E).
		 * This name change at every stable release, so we need to find it at runtime.
		 */
		if filepath.Dir(path) == ParentFolder {
			folder = filepath.Base(path)
			continue
		} else {
			//We remove that name so the top folder it's always a generic "cataclysm".
			newPath, err = filepath.Rel(ParentFolder+"/"+folder, path)
			if err != nil {
				return err
			}

		}

		if info.IsDir() {
			if err = os.MkdirAll(ParentFolder+"/"+newPath, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(ParentFolder+"/"+newPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		_, err = io.Copy(file, reader)
		if err != nil {
			return err
		}
		file.Close()
	}

	return nil
}

//CreateBackup create a zip archive of the cataclysm directory.
func CreateBackup(folder string) error {
	name := ParentFolder + " " + time.Now().Format(time.Stamp) + Suffix
	file, err := os.Create(name)
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

	err = filepath.Walk(folder, walker)
	if err != nil {
		return err
	}

	return nil
}

//ExtractBackup extract the cdda backup.
func ExtractBackup(name string) error {
	reader, err := zip.OpenReader(name)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {

		err = os.MkdirAll("./"+filepath.Dir(file.Name), 0755)
		if err != nil {
			return err
		}

		f, err := os.Create("./" + file.Name)
		if err != nil {
			return err
		}

		fd, err := file.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(f, fd)
		if err != nil {
			return err
		}

		err = f.Close()
		if err != nil {
			return err
		}

		err = fd.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
