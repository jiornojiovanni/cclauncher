package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
)

type githubResponse []struct {
	Name   string `json:"name"`
	Assets []struct {
		URL                string    `json:"url"`
		ID                 int       `json:"id"`
		NodeID             string    `json:"node_id"`
		Name               string    `json:"name"`
		Label              string    `json:"label"`
		ContentType        string    `json:"content_type"`
		State              string    `json:"state"`
		Size               int       `json:"size"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		BrowserDownloadURL string    `json:"browser_download_url"`
	} `json:"assets"`
}

//Build contains the version and the graphics (tiles or curses) of a CDDA build
type Build struct {
	Version string
	Graphic string
}

//LastBuild return a Build struct containing the version of the last successful build, with it's graphics type.
func LastBuild(curses bool) (Build, error) {
	build := Build{}

	if curses {
		build.Graphic = "curses"
	} else {
		build.Graphic = "tiles"
	}

	resp, err := http.Get("https://api.github.com/repos/cataclysmbnteam/Cataclysm-BN/releases?per_page=1")
	if err != nil {
		return Build{}, err
	}
	defer resp.Body.Close()

	gitResp, err := decodeResponse(resp)
	if err != nil {
		return Build{}, err
	}

	build.Version = strings.Split((*gitResp)[0].Name, " ")[3]
	return build, nil
}

func decodeResponse(resp *http.Response) (*githubResponse, error) {
	response := &githubResponse{}

	err := json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

//GetBuild download the specified version and return the name of the zip if successful.
func GetBuild(build Build) (string, error) {
	resp, err := http.Get("https://github.com/cataclysmbnteam/Cataclysm-BN/releases/download/cbn-experimental-" + build.Version + "/cbn-linux-" + build.Graphic + "-x64-" + build.Version + ".tar.gz")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return "", fmt.Errorf("this build version does not exists at the moment. Try later or with a different version")
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
