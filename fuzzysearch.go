package symspell

import (
	"container/heap"
	"fmt"
	"time"
)

// FuzzySearch ...
type FuzzySearch struct {
	editdistance int
	dict         map[string][]*Word
}

func push(bs heap.Interface, words ...*Word) {
	for _, word := range words {
		heap.Push(bs, word)
	}
}

// Search ...
func (fs *FuzzySearch) Search(query string, n, ed int) (v []*Word) {
	defer tick(time.Now(), "searchn "+query)
	var bs byspell
	heap.Init(&bs)
	if candis, ok := fs.dict[query]; ok {
		push(&bs, candis...)
	}
	q := []string{query}
	for i := 0; i < ed; i++ {
		q = edit1(q)
		for _, spell := range q {
			if candis, ok := fs.dict[spell]; ok {
				push(&bs, candis...)
			}
		}
	}

	var prev *Word
	for len(v) < n && bs.Len() > 0 {
		last := heap.Pop(&bs).(*Word)
		if prev == nil || prev.Spell != last.Spell {
			v = append(v, last)
			prev = last
		}
	}
	return
}

// AddWord ...
func (fs *FuzzySearch) AddWord(wd string, snippet interface{}, edmax int) {
	word := &Word{wd, snippet}
	fs.dict[wd] = append(fs.dict[wd], word)

	q := []string{wd}
	for i := 0; i < edmax; i++ {
		q = edit1(q)
		for _, spell := range q {
			fs.dict[spell] = append(fs.dict[spell], word)
		}
	}
}

func edit1(strs []string) (v []string) {
	for _, q := range strs {
		q := []rune(q)
		for i := 0; i < len(q); i++ {
			x := remove(q, i)
			v = append(v, string(x))
		}
	}
	return
}

// only delete one, no transposes, replaces and inserts
func edits(q []rune, ed, maxed int) (v []string) {
	v = append(v, string(q))
	ed++

	for i := 0; i < len(q); i++ {
		x := remove(q, i)
		v = append(v, string(x))
		if ed < maxed {
			v = append(v, edits(x, ed, maxed)...)
		}
	}
	return
}

// MakeFuzzySearch ...
func MakeFuzzySearch(ed int) *FuzzySearch {
	return &FuzzySearch{ed, make(map[string][]*Word)}
}

func remove(runes []rune, i int) []rune {
	var v = append([]rune{}, runes[:i]...)
	return append(v, runes[i+1:]...)
}

func tick(begin time.Time, prompt string) {
	fmt.Println(prompt, time.Now().Sub(begin))
}
