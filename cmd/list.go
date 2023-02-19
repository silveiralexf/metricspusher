package cmd

import (
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Prints list of metrics currently supported",
		Run: func(cmd *cobra.Command, args []string) {
			printMetricList()
		},
	}
)

func printMetricList() {
	metrics := map[string]string{
		"gaugetimestamp": "Records the timestamp of the last successful execution and result of a given job",
	}

	//headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	//columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Name", "Description")
	//tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for name, description := range metrics {
		tbl.AddRow(name, description)
	}

	tbl.Print()

}

func init() {
	rootCmd.AddCommand(listCmd)
}
