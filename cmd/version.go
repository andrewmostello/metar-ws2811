package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/andrewmostello/metar-ws2811/version"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `The version of this particular executable.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", version.Version)
		fmt.Printf("Built:   %s\n", version.BuildDate)
		fmt.Printf("Commit:  %s\n", version.GitCommit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
