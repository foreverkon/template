package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var downloadCmd = &cobra.Command{
	Use: "download <name>",
	Example: `
template download report
template d report
template d --name report
template d --name report --output report.tex`,
	Short:   "download a template file",
	Aliases: []string{"d"},
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.PersistentFlags().GetString("name")
		output, _ := cmd.LocalFlags().GetString("output")
		if name == "" {
			if len(args) == 0 {
				red.Println("no template name specified")
				return
			} else if len(args) == 1 {
				name = args[0]
			} else if len(args) == 2 {
				name = args[0]
				output = args[1]
			} else {
				red.Println("too many arguments")
				return
			}
		}
		download(name, output)
	},
}

func init() {
	// add subcommand
	rootCmd.AddCommand(downloadCmd)

	// define local flags
	downloadCmd.Flags().StringP("output", "o", "", "name of the output file")
}

func download(name, output string) {
	names := viper.GetStringSlice("names")
	ranks := search(name)
	if ranks == nil {
		return
	}

	if len(ranks) > 1 {
		var num int
		for {
			yellow.Print("Select a number: ")
			fmt.Scanf("%d", &num)
			if num > 0 && num <= len(ranks) {
				name = names[ranks[num-1].OriginalIndex]
				break
			} else {
				red.Println("Invalid number")
			}
		}
	} else {
		name = names[ranks[0].OriginalIndex]
	}

	if output == "" {
		output = name
	}

	url := viper.GetString("templates." + name)
	green.Printf("Downloading %s to %s ...\n", blue(name), blue(output))
	err := downloadToFile(url, output)
	if err != nil {
		red.Println(err.Error())
		return
	}
	green.Println("Downloaded successfully")
}
