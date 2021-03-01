package archives

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	//Parentfolder is the name of the cdda directory
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

//CreateBackup create a zip archive of the cataclysm directory, with the flag all set to false it will only backup Saves, Sound, Config and Font.
func CreateBackup(name string, all bool) error {
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
			fmt.Println("Folder 'save' not found... skipping")
		}

		err = filepath.Walk(name+"/config", walker)
		if err != nil {
			fmt.Println("Folder 'config' not found... skipping")
		}

		err = filepath.Walk(name+"/sound", walker)
		if err != nil {
			fmt.Println("Folder 'sound' not found... skipping")
		}

		err = filepath.Walk(name+"/font", walker)
		if err != nil {
			fmt.Println("Folder 'font' not found... skipping")
		}

	}
	return nil
}

//ExtractBackup extract the cdda backup.
func ExtractBackup(name string) error {
	reader, err := zip.OpenReader(name)
	defer reader.Close()
	if err != nil {
		return err
	}

	for _, file := range reader.File {

		err = os.MkdirAll("./"+filepath.Dir(file.Name), 0755)
		if err != nil {
			return err
		}

		f, err := os.Create("./" + file.Name)
		if err != nil {
			return err
		}
		defer f.Close()

		fd, err := file.Open()
		if err != nil {
			return err
		}
		defer fd.Close()

		_, err = io.Copy(f, fd)
		if err != nil {
			return err
		}
	}
	return nil
}

//CheckFolder check if the folder exists
func CheckFolder(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
