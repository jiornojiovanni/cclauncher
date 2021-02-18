package download

import (
	"cclauncher/web"
	"fmt"
	"io"
	"net/http"
	"os"
)

//GetBuild download the specified version
func GetBuild(build web.Build) (string, error) {
	resp, err := http.Get("https://github.com/CleverRaven/Cataclysm-DDA/releases/download/cdda-jenkins-b" + fmt.Sprint(build.Version) + "/cataclysmdda-0.E-Linux_x64-" + build.Graphic + "-b" + fmt.Sprint(build.Version) + ".tar.gz")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return "", fmt.Errorf("build not found on Github")
	}

	filename := "cataclysm-" + build.Graphic + "-" + fmt.Sprint(build.Version) + ".tar.gz"
	err = createFile(resp, filename)
	if err != nil {
		return "", err
	}

	return filename, err
}

func createFile(resp *http.Response, filename string) error {

	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fd.Close()

	bytes, err := io.Copy(fd, resp.Body)
	if err != nil {
		return err
	} else if bytes == 0 {
		return error(fmt.Errorf("nothing was written on disk, error"))
	}

	return nil
}
