package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Print("IMPinterpreter <source.imp>")
		return
	}
	source := args[0]
	file, err := os.Open(source)
	if err != nil {
		panic(err.Error())
	}
	i := NewInterpreter(bufio.NewReader(file))
	i.Scan()
	fmt.Println(i.scanned)
}
