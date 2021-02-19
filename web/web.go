package web

import (
	"encoding/json"
	"fmt"
	"net/http"
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
func LastBuild(graphics string) (Build, error) {
	build := Build{}

	build.Graphic = graphics

	resp, err := http.Get("https://ci.narc.ro/job/Cataclysm-Matrix/Graphics=" + build.Graphic + ",Platform=Linux_x64/lastSuccessfulBuild/api/json?pretty=true")
	if err != nil {
		return Build{}, err
	}
	defer resp.Body.Close()

	json, err := decodeResponse(resp)
	if err != nil {
		return Build{}, err
	}

	//Unnecessary check, but better safe than sorry
	if json.Result == "SUCCESS" {
		build.Version = json.Number
		return build, nil
	}

	return Build{}, nil
}

//CheckBuild check if the version/graphic combo was compiled succesfully.
func CheckBuild(build Build) (bool, error) {
	resp, err := http.Get("https://ci.narc.ro/job/Cataclysm-Matrix/Graphics=" + build.Graphic + ",Platform=Linux_x64/" + fmt.Sprint(build.Version) + "/api/json?pretty=true")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	json, err := decodeResponse(resp)
	if err != nil {
		return false, err
	}

	//Here the check is useful
	if json.Result == "SUCCESS" {
		return true, nil
	}

	return false, nil

}

func decodeResponse(resp *http.Response) (*jsonResponse, error) {
	var response jsonResponse

	err := json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
