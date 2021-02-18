package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type jsonResponse struct {
	Number int
	Result string
}

//Build contains the version and the graphics (tiles or curses) of a CDDA build
type Build struct {
	Version int
	Graphic string
}

//LastBuild return a Build struct containing the version of the last successful build, with it's graphics type.
func LastBuild(graphics string) (Build, error) {
	build := Build{}
	if graphics == "c" {
		build.Graphic = "Curses"
	} else {
		build.Graphic = "Tiles"
	}
	resp, err := http.Get("https://ci.narc.ro/job/Cataclysm-Matrix/Graphics=" + build.Graphic + ",Platform=Linux_x64/lastBuild/api/json?pretty=true")
	if err != nil {
		return Build{}, err
	}
	defer resp.Body.Close()

	var response jsonResponse

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return Build{}, err
	}

	if response.Result == "SUCCESS" {
		build.Version = response.Number
		return build, nil
	}

	for build.Version = response.Number; build.Version > 0; build.Version-- {
		resp, err := http.Get("https://ci.narc.ro/job/Cataclysm-Matrix/Graphics=" + build.Graphic + ",Platform=Linux_x64/" + fmt.Sprint(build.Version) + "/api/json?pretty=true")
		if err != nil {
			return Build{}, err
		}
		defer resp.Body.Close()

		if response.Result == "SUCCESS" {
			return build, nil
		}
	}
	return build, nil
}
