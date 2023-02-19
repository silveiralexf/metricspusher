package cmd

import "github.com/spf13/cobra"

var (
	releaseNumber string = "local"
	versionCmd           = &cobra.Command{
		Use:   "version",
		Short: "Prints version information",
		Run: func(cmd *cobra.Command, args []string) {
			print(releaseNumber + "\n")
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
