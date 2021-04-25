package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cclauncher",
	Short: "Cataclysm: DDA Linux Launcher",
	Long: `cclauncher is a Linux launcher for Cataclysm: DDA,
capable of downloading new versions, backups and restoring of
tilesets, sfx, mods, saves etc.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
}
