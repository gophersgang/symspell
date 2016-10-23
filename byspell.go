package symspell

import (
	"bufio"
	"os"
	"time"
)

// Word ...
type Word struct {
	Spell   string
	Snippet interface{}
}

type byspell []*Word

func (bs byspell) Less(i, j int) bool {
	return bs[i].Spell < bs[j].Spell
}
func (bs byspell) Swap(i, j int) {
	bs[i], bs[j] = bs[j], bs[i]
}
func (bs byspell) Len() int {
	return len(bs)
}
func (bs *byspell) Push(x interface{}) {
	*bs = append(*bs, x.(*Word))
}
func (bs *byspell) Pop() interface{} {
	old, n := *bs, len(*bs)
	x := old[n-1]
	*bs = old[0 : n-1]
	return x
}

// EachLine ...
func EachLine(fp string, proc func(string)) {
	defer tick(time.Now(), "read corpus")
	if file, err := os.Open(fp); err == nil {
		defer file.Close() // error return not checked
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			proc(scanner.Text())
		}
	}
}
