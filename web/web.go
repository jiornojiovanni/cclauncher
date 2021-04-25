package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cheggaaa/pb/v3"
)

type jsonResponse struct {
	Number    int    `json:"number"`
	Result    string `json:"result"`
	ChangeSet struct {
		Item []CommitData `json:"items"`
	} `json:"changeSet"`
}

//CommitData represent a single commit.
type CommitData struct {
	CommitID string `json:"commitId"`
	Date     string `json:"date"`
	Msg      string `json:"msg"`
}

//Build contains the version and the graphics (tiles or curses) of a CDDA build
type Build struct {
	Version int
	Graphic string
}

//LastBuild return a Build struct containing the version of the last successful build, with it's graphics type.
func LastBuild(curses bool) (Build, error) {
	build := Build{}

	if curses {
		build.Graphic = "Curses"
	} else {
		build.Graphic = "Tiles"
	}

	resp, err := http.Get("https://ci.narc.ro/job/Cataclysm-Matrix/Graphics=" + build.Graphic + ",Platform=Linux_x64/lastSuccessfulBuild/api/json?pretty=true")
	if err != nil {
		return Build{}, err
	}
	defer resp.Body.Close()

	jsonResp, err := decodeResponse(resp)
	if err != nil {
		return Build{}, err
	}

	//Unnecessary check, but better safe than sorry
	if jsonResp.Result == "SUCCESS" {
		build.Version = jsonResp.Number
		return build, nil
	}
	return Build{}, nil
}

//CheckBuild check if the version/graphic combo was compiled successfully.
func CheckBuild(build Build) (bool, error) {
	resp, err := http.Get("https://ci.narc.ro/job/Cataclysm-Matrix/Graphics=" + build.Graphic + ",Platform=Linux_x64/" + fmt.Sprint(build.Version) + "/api/json?pretty=true")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	jsonResponse, err := decodeResponse(resp)
	if err != nil {
		return false, err
	}

	//Here the check is useful
	if jsonResponse.Result == "SUCCESS" {
		return true, nil
	}

	return false, nil
}

//GetChangelog returns an array containing the various commits of the build.
func GetChangelog(build Build) ([]CommitData, error) {
	resp, err := http.Get("https://ci.narc.ro/job/Cataclysm-Matrix/Graphics=" + build.Graphic + ",Platform=Linux_x64/" + fmt.Sprint(build.Version) + "/api/json?pretty=true")
	if err != nil {
		return []CommitData{}, err
	}
	defer resp.Body.Close()

	jsonResponse, err := decodeResponse(resp)
	if err != nil {
		return []CommitData{}, err
	}

	return jsonResponse.ChangeSet.Item, nil

}

func decodeResponse(resp *http.Response) (*jsonResponse, error) {
	var response jsonResponse

	err := json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

//GetBuild download the specified version and return the name of the zip if successful.
func GetBuild(build Build) (string, error) {
	resp, err := http.Get("https://github.com/CleverRaven/Cataclysm-DDA/releases/download/cdda-jenkins-b" + fmt.Sprint(build.Version) + "/cataclysmdda-0.E-Linux_x64-" + build.Graphic + "-b" + fmt.Sprint(build.Version) + ".tar.gz")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return "", fmt.Errorf("build not found on Github")
	}

	filename := "cataclysm-" + build.Graphic + "-" + fmt.Sprint(build.Version) + ".tar.gz"
	err = downloadFile(resp, filename)
	if err != nil {
		return "", err
	}

	return filename, err
}

func downloadFile(resp *http.Response, filename string) error {

	bar := pb.New(int(resp.ContentLength))
	bar.Set(pb.Bytes, true)

	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fd.Close()

	proxy := bar.NewProxyReader(resp.Body)

	bar.Start()
	bytes, err := io.Copy(fd, proxy)
	if err != nil {
		return err
	} else if bytes == 0 {
		return fmt.Errorf("nothing was written on disk, error")
	}

	bar.Finish()

	return nil
}
