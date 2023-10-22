package main

import (
	"flag"
	"fmt"
	"os"
	"log"
)

const usage = `Usage of sodium:
  -i, --init verbose output
  -h, --help prints help information 
`

func main() {
	shortInit := flag.Bool("i", false, "Init aliases from current working directory")
	longInit := flag.Bool("init", false, "Init aliases from current working directory")

	flag.Usage = func() {fmt.Print(usage)}
	flag.Parse()

	init := *shortInit || *longInit

	fmt.Printf("init: %v\n", init)

	if(init) {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal("Could not get current working dir")
		}

		fmt.Println(cwd)
	}
}