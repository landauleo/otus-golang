package main

import (
	"flag"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

// go run . -from="testdata/input.txt" -to="out.txt" -limit=100 -offset=0
func main() {
	flag.Parse() //сначала парсим флаги
	Copy(from, to, offset, limit)
}
