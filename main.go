package main

import (
	"cclauncher/archives"
	"cclauncher/download"
	"cclauncher/web"
	"flag"
	"fmt"
	"log"
	"strconv"
)

func main() {
	var err error
	var version string
	var graphics string
	var downloadOnly bool
	flag.StringVar(&version, "v", "latest", "Specify the version of CDDA.")
	flag.StringVar(&graphics, "g", "t", "Specify the graphic version of CDDA, t for tiles and c for curses.")
	flag.BoolVar(&downloadOnly, "d", false, "Set this flag to true to only download the build")
	flag.Parse()

	var g string
	if graphics == "c" {
		g = "Curses"
	} else {
		g = "Tiles"
	}

	var build web.Build

	if version != "latest" {
		v, err := strconv.Atoi(version)
		if err != nil {
			log.Fatal("There was an error while parsing the version number.")
		}

		build = web.Build{Version: v, Graphic: g}

		res, err := web.CheckBuild(build)
		if err != nil {
			log.Fatal(err)
		}

		if !res {
			log.Fatal("This version is unavailable, try another.")
		}
	} else {

		build, err = web.LastBuild(graphics)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Last build: ", build.Version)
	}

	fmt.Println("Trying to download...")
	filename, err := download.GetBuild(build)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Build successfully downloaded.")

	if !downloadOnly {
		err = archives.Extract(filename)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Files extracted.")
	}
}
