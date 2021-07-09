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
	Short: "Download a new version of C:BN",
	Long: `The download command by default downloads the latest experimental
of C:BN, then it extracts it and moves old data from the previous version
to the new one.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		cursesSwitch, _ := cmd.Flags().GetBool("ncurses")
		versionFlag, _ := cmd.Flags().GetString("version")
		build := web.Build{}

		if versionFlag != "" {
			build.Version = versionFlag
			if cursesSwitch {
				build.Graphic = "curses"
			} else {
				build.Graphic = "tiles"
			}
		} else {
			build, err = web.LastBuild(cursesSwitch)
			if err != nil {
				log.Fatal(err)
			}
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

		fmt.Println("Extracting C:BN...")
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

	downloadCmd.Flags().BoolP("ncurses", "n", false, "Ncurses version")
	downloadCmd.Flags().StringP("version", "v", "", "Experimental version string")
}
