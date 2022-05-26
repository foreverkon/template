package cmd

import (
	"errors"
	"io"
	"os"
	"sort"
	"template/request"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/spf13/viper"
)

func search(name string) fuzzy.Ranks {
	names := viper.GetStringSlice("names")
	ranks := fuzzy.RankFind(name, names)

	if len(ranks) == 0 {
		red.Println("no template found")
		return nil
	}

	sort.Sort(ranks)
	green.Println("Found these files:")
	for i := len(ranks) - 1; i >= 0; i-- {
		magenta.Printf("[%d] %s\n", i+1, names[ranks[i].OriginalIndex])
	}

	return ranks
}

func downloadToFile(url, path string) error {
	resp, err := request.Get(url, nil, []byte("http://127.0.0.1:7890"))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return errors.New("IO Err")
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return errors.New("IO Err")
	}
	return nil
}
