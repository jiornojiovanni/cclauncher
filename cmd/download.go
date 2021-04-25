package cmd

import (
	"cclauncher/archives"
	"cclauncher/utils"
	"cclauncher/web"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a new version of CDDA",
	Long: `The download command by default downloads the latest experimental
of CDDA, then it extracts it and moves old data from the previous version
to the new one.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		changelogSwitch, _ := cmd.Flags().GetBool("changelog")
		cursesSwitch, _ := cmd.Flags().GetBool("ncurses")
		versionFlag, _ := cmd.Flags().GetInt("version")
		build := web.Build{}

		if versionFlag != -1 {
			build.Version = versionFlag
			if cursesSwitch {
				build.Graphic = "Curses"
			} else {
				build.Graphic = "Tiles"
			}
			exists, err := web.CheckBuild(build)
			if err != nil {
				log.Fatal(err)
			}

			if !exists {
				log.Fatal("This build version does not exists at the moment. Try later or with a different version.")
			}

		} else {
			build, err = web.LastBuild(cursesSwitch)
			if err != nil {
				log.Fatal(err)
			}
		}

		if changelogSwitch {
			commits, err := web.GetChangelog(build)
			if err != nil {
				log.Fatal(err)
			}

			utils.PrintChangelog(commits)
		}

		fmt.Println("Downloading version", build.Version)
		zip, err := web.GetBuild(build)
		if err != nil {
			log.Fatal(err)
		}

		previousExists, err := utils.CheckFolder(archives.ParentFolder)
		if err != nil {
			log.Fatal(err)
		}

		if previousExists {
			fmt.Println("Backupping old folder...")
			err = archives.CreateBackup(archives.ParentFolder)
			if err != nil {
				log.Fatal(err)
			}

			err = os.Rename(archives.ParentFolder, archives.TempFolder)
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println("Extracting CDDA...")
		err = archives.Extract(zip)
		if err != nil {
			log.Fatal(err)
		}

		if previousExists {
			fmt.Println("Restoring data...")
			err = utils.RestoreCustomContent(cursesSwitch)
			if err != nil {
				log.Fatal(err)
			}

			err = utils.RestoreData()
			if err != nil {
				log.Fatal(err)
			}

			err = os.RemoveAll(archives.TempFolder)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = os.Remove(zip)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Done!")
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().BoolP("changelog", "c", false, "Display changelog")
	downloadCmd.Flags().BoolP("ncurses", "n", false, "Ncurses version")
	downloadCmd.Flags().IntP("version", "v", -1, "Experimental version number")
}
