package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"template/request"

	"github.com/lithammer/fuzzysearch/fuzzy"
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
		name, _ := cmd.LocalFlags().GetString("name")
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
	downloadCmd.Flags().StringP("name", "n", "", "name of the template to download")
	downloadCmd.Flags().StringP("output", "o", "", "name of the output file")
}

func download(name, output string) {
	names := viper.GetStringSlice("names")
	nosuffix_names := make([]string, len(names))
	for i := range names {
		tmp := names[i]
		// tmp = strings.TrimSuffix(tmp, ".tex")
		// tmp = strings.TrimSuffix(tmp, ".gitignore")
		// tmp = strings.TrimSuffix(tmp, ".go")
		// tmp = strings.TrimSuffix(tmp, ".sh")
		// tmp = strings.TrimSuffix(tmp, ".py")
		nosuffix_names[i] = tmp
	}

	ranks := fuzzy.RankFind(name, nosuffix_names)
	sort.Sort(ranks)

	if len(ranks) == 0 {
		red.Println("no template found")
		return
	}
	if len(ranks) > 1 {
		green.Println("Found multiple files:")
		for i := len(ranks) - 1; i >= 0; i-- {
			magenta.Printf("[%d] %s\n", i+1, names[ranks[i].OriginalIndex])
		}
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

	err := downloadFile(name, output)
	if err != nil {
		red.Println(err.Error())
		return
	}
}

func downloadFile(name, path string) error {
	if path == "" {
		path = name
	}
	green.Printf("Downloading %s to %s ...\n", blue(name), blue(path))

	// if exists(path) {
	// 	return errors.New("file already exists")
	// }

	url := viper.GetString("templates." + name)
	resp, err := request.Get(url, nil, []byte("http://127.0.0.1:7890"))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// file, err := os.Create(path)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return errors.New("IO Err")
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return errors.New("IO Err")
	}
	green.Println("Downloaded successfully")
	return nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
