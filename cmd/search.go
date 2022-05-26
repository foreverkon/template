package cmd

import (
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search from local templates.json",
	Example: `
template search report
template s report
template s --name report
template s -n report`,
	Aliases: []string{"s"},
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.PersistentFlags().GetString("name")
		if name == "" {
			if len(args) == 0 {
				red.Println("no template name specified")
				return
			} else if len(args) == 1 {
				name = args[0]
			} else {
				red.Println("too many arguments")
				return
			}
		}

		search(name)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
