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
	var dontShowChangelog bool
	flag.StringVar(&version, "v", "latest", "Specify the version of CDDA.")
	flag.StringVar(&g, "g", "t", "Specify the graphic version of CDDA, t for tiles and c for curses.")
	flag.BoolVar(&downloadOnly, "d", false, "Set this flag to true to only download the build")
	flag.BoolVar(&dontShowChangelog, "no-c", false, "Set this flag to true to don't show the changelog")
	flag.Parse()

	isInstalled, err := archives.CheckFolder("./" + archives.Parentfolder)
	if err != nil {
		log.Fatal(err)
	}

	if isInstalled && !downloadOnly {
		fmt.Println("Older CDDA version found, backupping data...")
		err = archives.CreateBackup("./"+archives.Parentfolder, false)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Backup saved.")

		err = os.RemoveAll("./" + archives.Parentfolder)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Older CDDA version deleted.")
	}

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

	commits, err := web.GetChangelog(build)

	if !dontShowChangelog {
		fmt.Println("Changelog: ")
		fmt.Printf("\033[32m# %-26s%-80sCommit\n", "Date", "Message")
		for i := 0; i < len(commits); i++ {
			fmt.Print("\033[31m", i+1)
			fmt.Print(" \033[33m", commits[i].Date)
			fmt.Printf("\033[37m %-80s", commits[i].Msg)
			fmt.Print("\033[34m", commits[i].CommitID)
			fmt.Print("\033[0m")
			fmt.Print("\n")
		}
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

		isBackupped, err := archives.CheckFolder(archives.Parentfolder + ".zip")
		if err != nil {
			log.Fatal(err)
		}

		if isBackupped {
			err = archives.ExtractBackup(archives.Parentfolder + ".zip")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Backup extracted.")
		}
	}
}
