package cmd

import (
	"cclauncher/utils"
	"cclauncher/web"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// changelogCmd represents the changelog command
var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Show the changelog",
	Long: `The changelog command retrieve the changelog of the latest CDDA experimental version
or of a specific version (with the flag -v).
You can get a json version with --json.`,
	Run: func(cmd *cobra.Command, args []string) {
		version, _ := cmd.Flags().GetInt("version")
		var err error
		var build web.Build

		if version != -1 {
			build.Version = version
			exists, err := web.CheckBuild(build)
			if err != nil {
				log.Fatal(err)
			}

			if !exists {
				log.Fatal("This build version does not exist at the moment. Try later or with a different version.")
			}
		} else {
			build, err = web.LastBuild(false)
			if err != nil {
				log.Fatal(err)
			}
		}

		commits, err := web.GetChangelog(build)
		if err != nil {
			log.Fatal(err)
		}

		jsonVersion, _ := cmd.Flags().GetBool("json")
		if jsonVersion {
			fmt.Println(commits)
		} else {
			utils.PrintChangelog(commits)
		}
	},
}

func init() {
	rootCmd.AddCommand(changelogCmd)
	changelogCmd.Flags().Bool("json", false, "Json version of changelog")
	changelogCmd.Flags().IntP("version", "v", -1, "Experimental version number")
}
