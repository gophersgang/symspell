package main

import (
	"flag"
	"fmt"

	ss "github.com/hearts.zhang/symspell"
)

var config struct {
	a  string
	b  string
	ed int
}

func init() {
	flag.StringVar(&config.a, "a", "", "")
	flag.StringVar(&config.b, "b", "", "")
	flag.IntVar(&config.ed, "ed", 1, "")
}
func main() {
	flag.Parse()

	x := ss.EditMatch(config.a, config.b, config.ed)
	fmt.Println(x)
}
