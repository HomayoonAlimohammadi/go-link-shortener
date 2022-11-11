package linkshortener

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Link Shortener",
	Long:  `All software has versions. This is Link Shortener's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Link Shortener v1.0.0")
	},
}
