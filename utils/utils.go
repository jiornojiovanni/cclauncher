package utils

import (
	"cclauncher/archives"
	"io/ioutil"
	"os"
)

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

//GetCustomFolders compares the clean install of cdda to the player one and detects
func getCustomFolders(path string, previousPath string) ([]string, error) {
	folder, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var officialSet = make([]string, 0)
	for _, f := range folder {
		if f.IsDir() {
			officialSet = append(officialSet, f.Name())
		}
	}

	previousFolder, err := ioutil.ReadDir(previousPath)
	if err != nil {
		return nil, err
	}

	var previousSet = make([]string, 0)
	for _, f := range previousFolder {
		if f.IsDir() {
			previousSet = append(previousSet, f.Name())
		}
	}

	for i := 0; i < len(previousSet); i++ {
		for j := 0; j < len(officialSet); j++ {
			if previousSet[i] == officialSet[j] {
				previousSet = append(previousSet[:i], previousSet[i+1:]...)
			}
		}
	}

	return previousSet, nil
}

func restoreFolders(folder string, content []string) error {
	for _, value := range content {
		err := os.Rename(archives.TempFolder+"/"+folder+"/"+value, archives.ParentFolder+"/"+folder+"/"+value)
		if err != nil {
			return err
		}
	}

	return nil
}

//RestoreCustomContent restore user downloaded custom content.
func RestoreCustomContent(curses bool) error {
	var content []string
	if curses {
		content = []string{"data/sound", "data/mods"}
	} else {
		content = []string{"gfx", "data/sound", "data/mods"}
	}

	for _, folder := range content {
		exists, err := CheckFolder(archives.TempFolder + "/" + folder)
		if err != nil {
			return err
		}

		if exists {
			stuffs, err := getCustomFolders(archives.ParentFolder+"/"+folder, archives.TempFolder+"/"+folder)
			if err != nil {
				return err
			}
			if len(stuffs) != 0 {
				err = restoreFolders(folder, stuffs)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

//RestoreData moves restore user generated data.
func RestoreData() error {
	content := [...]string{"save", "config"}
	for _, folder := range content {
		exists, err := CheckFolder(archives.ParentFolder + "/" + folder)
		if err != nil {
			return err
		}

		if exists {
			dir, err := ioutil.ReadDir(archives.ParentFolder + "/" + folder)
			if len(dir) != 0 {
				if err != nil {
					return err
				}

				set := make([]string, len(dir))
				for _, f := range dir {
					set = append(set, f.Name())
				}

				err = restoreFolders(folder, set)
				if err != nil {
					return err
				}
			}

		}
	}
	return nil
}
