package main

import (
	"bufio"
	"container/heap"
	"flag"
	"fmt"
	"hash"
	"os"
	"strings"
	"time"

	xh "github.com/OneOfOne/xxhash"
)

func main() {
	flag.Parse()
	eachLine(config.corpus, func(line string) {
		word := readWord(line)
		//config.dict[word.spell] = append(config.dict[word.spell], word)

		for _, spell := range edits([]rune(word.spell), 0) {
			config.dict[spell] = append(config.dict[spell], word)
		}
	})
	words := search(config.q)
	for _, word := range words {
		fmt.Println(word)
	}
}
func readWord(line string) *word {
	pinyin := ssplit(line, ";")[5]
	pinyins := strings.Fields(pinyin)
	full := strings.Join(pinyins, "")
	return &word{full, nil}
}
func ssplit(str, sep string) []string {
	return strings.FieldsFunc(str, func(r rune) bool {
		return strings.ContainsRune(sep, r)
	})
}
func eachLine(fp string, proc func(string)) {
	if file, err := os.Open(fp); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			proc(scanner.Text())
		}
	}
}

// only delete one, no transposes, replaces and inserts
func edits(q []rune, ed int) (v []string) {
	v = append(v, string(q))
	ed++

	for i := 0; i < len(q); i++ {
		x := remove(q, i)
		v = append(v, string(x))
		if ed < config.editdistance {
			v = append(v, edits(x, ed)...)
		}
	}
	return
}
func remove(runes []rune, i int) []rune {
	var v = append([]rune{}, runes[:i]...)
	return append(v, runes[i+1:]...)
}
func search(q string) (v []*word) {
	var bs byspell
	for _, spell := range edits([]rune(q), 0) {
		if candis, ok := config.dict[spell]; ok {
			bs = append(bs, candis...)
		}
	}

	heap.Init(&bs)
	var prev *word
	for len(v) < config.n && bs.Len() > 0 {
		last := heap.Pop(&bs).(*word)
		if prev == nil || prev.spell != last.spell {
			v = append(v, last)
			prev = last
		}
	}
	return
}

func init() {
	flag.StringVar(&config.corpus, "corpus", "media-resites-v3.tsv", "")
	flag.StringVar(&config.q, "q", "xiaoaobang", "")
	flag.IntVar(&config.editdistance, "ed", 1, "")
	flag.IntVar(&config.n, "n", 3, "")
	config.dict = make(map[string][]*word)
	config.hash = xh.NewS64(uint64(time.Now().UnixNano()))
}

var config struct {
	editdistance int
	n            int
	corpus       string
	q            string
	hash         hash.Hash64
	dict         map[string][]*word
}

type word struct {
	spell   string
	snippet interface{}
}

type byspell []*word

func (bs byspell) Less(i, j int) bool {
	return bs[i].spell < bs[j].spell
}
func (bs byspell) Swap(i, j int) {
	bs[i], bs[j] = bs[j], bs[i]
}
func (bs byspell) Len() int {
	return len(bs)
}
func (bs *byspell) Push(x interface{}) {
	*bs = append(*bs, x.(*word))
}
func (bs *byspell) Pop() interface{} {
	old, n := *bs, len(*bs)
	x := old[n-1]
	*bs = old[0 : n-1]
	return x
}
