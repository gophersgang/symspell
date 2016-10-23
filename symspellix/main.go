package main

import (
	"flag"
	"fmt"
	"strings"

	ss "github.com/hearts.zhang/symspell"
)

func main() {
	flag.Parse()
	config.fs = ss.MakeFuzzySearch(config.edindex)
	ss.EachLine(config.corpus, func(line string) {
		words, snippet := readWord(line)
		for _, word := range words {
			config.fs.AddWord(word, snippet, config.edindex)
		}
	})
	words := config.fs.Search(config.q, config.n, config.edsearch)
	for _, word := range words {
		fmt.Println(word.Spell, word.Snippet.(*media))
	}
}

type media struct {
	name string
	id   string
}

func readWord(line string) ([]string, interface{}) {
	fields := strings.Split(line, ";")
	m := &media{id: fields[0], name: fields[1]}
	pinyin := fields[5]
	pinyins := strings.Fields(pinyin)
	full := strings.Join(pinyins, "")
	var firsts []rune
	for _, pinyin := range pinyins {
		firsts = append(firsts, []rune(pinyin)[0])
	}
	return []string{full, string(firsts)}, m
}

func init() {
	flag.StringVar(&config.corpus, "corpus", "media-resites-v3.tsv", "")
	flag.StringVar(&config.q, "q", "xiaoaobang", "")
	flag.IntVar(&config.edindex, "ed", 1, "")
	flag.IntVar(&config.edsearch, "edsearch", 1, "")
	flag.IntVar(&config.n, "n", 3, "")
}

var config struct {
	n        int
	edindex  int
	edsearch int
	corpus   string
	q        string
	fs       *ss.FuzzySearch
}
