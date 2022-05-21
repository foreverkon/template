package cmd

import (
	"sort"
	"template/request"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// search these repos
	urls = []string{
		"https://github.com/github/gitignore/tree/main/Global",
		"https://github.com/github/gitignore",
		"https://github.com/foreverkon/template-files",
	}

	// coressponding dowload prefix
	dowload_prefix = []string{
		"https://raw.githubusercontent.com/github/gitignore/main/Global/",
		"https://raw.githubusercontent.com/github/gitignore/main/",
		"https://raw.githubusercontent.com/foreverkon/template-files/main/",
	}

	// ignore_files
	ignore_files = []string{
		"README.md",
		"LICENSE",
		".gitignore",
	}
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "update templates.json",
	Aliases: []string{"u"},
	Run: func(cmd *cobra.Command, args []string) {
		names := make([]string, 0, 1024)
		for index, url := range urls {
			green.Print("fetching ", blue(url), "...")
			resp, err := request.Get(url, nil, []byte("http://127.0.0.1:7890"))
			if err != nil {
				red.Println(err.Error())
				return
			}
			defer func() { _ = resp.Body.Close() }()

			doc, _ := goquery.NewDocumentFromReader(resp.Body)

			doc.Find("div[role=row]").Each(func(i int, s *goquery.Selection) {
				if label, _ := s.Find("div[role] svg").Attr("aria-label"); label == "File" {
					name, exist := s.Find("div[role=rowheader] span a[title]").Attr("title")
					if in(name, ignore_files) {
						return
					}
					if exist {
						names = append(names, name)
						viper.Set("templates."+name, dowload_prefix[index]+name)
					}
				}
			})
			green.Println("done")
		}
		viper.Set("names", names)
		viper.Set("update", time.Now().Unix())
		viper.WriteConfig()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func in(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	if index < len(str_array) && str_array[index] == target {
		return true
	}
	return false
}
