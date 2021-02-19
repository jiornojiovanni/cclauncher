package main

import (
	"cclauncher/archives"
	"cclauncher/download"
	"cclauncher/web"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	var err error
	var version string
	var g string
	var downloadOnly bool
	flag.StringVar(&version, "v", "latest", "Specify the version of CDDA.")
	flag.StringVar(&g, "g", "t", "Specify the graphic version of CDDA, t for tiles and c for curses.")
	flag.BoolVar(&downloadOnly, "d", false, "Set this flag to true to only download the build")
	flag.Parse()

	var graphics string
	if g == "c" {
		graphics = "Curses"
	} else {
		graphics = "Tiles"
	}

	var build web.Build

	if version != "latest" {
		v, err := strconv.Atoi(version)
		if err != nil {
			log.Fatal("There was an error while parsing the version number.")
		}

		build = web.Build{Version: v, Graphic: graphics}

		res, err := web.CheckBuild(build)
		if err != nil || !res {
			fmt.Println(("This version is unavailable, try another."))
			os.Exit(0)
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
