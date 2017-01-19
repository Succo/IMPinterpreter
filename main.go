package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	var in *bufio.Reader
	args := flag.Args()
	if len(args) == 0 {
		in = bufio.NewReader(os.Stdin)
	} else if len(args) == 1 {
		source := args[0]
		file, err := os.Open(source)
		if err != nil {
			panic(err.Error())
		}
		in = bufio.NewReader(file)
	} else {
		fmt.Println("IMPinterpreter <file> to interpret a file")
		fmt.Println("if no file is provided, reads from stdin")
		return
	}

	l := NewLexer(in)
	p := l.Scan()
	i := NewInterpreter(p.parseInst())
	i.execute()
}
