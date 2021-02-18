package main

import (
	"cclauncher/archives"
	"cclauncher/download"
	"cclauncher/web"
	"flag"
	"fmt"
	"log"
)

func main() {
	var version string
	var graphics string
	var downloadOnly bool
	flag.StringVar(&version, "v", "latest", "Specify the version of CDDA.")
	flag.StringVar(&graphics, "g", "t", "Specify the graphic version of CDDA, t for tiles and c for curses.")
	flag.BoolVar(&downloadOnly, "d", false, "Set this flag to true to only download the build")
	flag.Parse()

	if version == "latest" {
		if graphics != "t" && graphics != "c" {
			log.Fatal("Unrecognized graphics option. Use [t]iles or [c]urses")
		}

		build, err := web.LastBuild(graphics)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Last build: ", build.Version)
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
	} else {
		fmt.Println("Sorry, not implemented yet.")
	}

}
