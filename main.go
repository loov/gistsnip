package main

import (
	"flag"
	"log"
	"os"

	"github.com/kr/pretty"
)

var (
	configName = flag.String("config", ".gistsnip", "default configuration file")
)

//gistsnip:start:main
func main() {
	flag.Parse()

	paths := flag.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}

	oldConfig, err := LoadConfig(*configName)
	if err == os.ErrNotExist {
		oldConfig = &Config{}
		err = nil
	}
	if err != nil {
		log.Fatal(err)
	}

	newConfig, err := ConfigFromGlobs(paths)
	if err != nil {
		log.Fatal(err)
	}

	pretty.Print(oldConfig)
	pretty.Print(newConfig)
}

//gistsnip:end:main
