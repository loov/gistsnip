package main

import (
	"flag"
	"log"
	"os"

	"github.com/kr/pretty"
)

var (
	configName  = flag.String("config", ".gistsnip", "default configuration file")
	githubToken = flag.String("github", os.Getenv("GISTSNIP_TOKEN"), "github authentication token")
)

//gistsnip:start:main
func main() {
	//gistsnip:start:parse
	flag.Parse()
	//gistsnip:end:parse

	paths := flag.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}

	oldGist, err := LoadConfig(*configName)
	if os.IsNotExist(err) {
		oldGist = NewGist()
		err = nil
	}
	if err != nil {
		log.Fatal(err)
	}

	newGist, err := GistFromGlobs(paths)
	if err != nil {
		log.Fatal(err)
	}

	pretty.Println(oldGist)
	pretty.Println(newGist)
}

//gistsnip:end:main
