package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/kr/pretty"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
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

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *githubToken})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)

	client := githubv4.NewClient(httpClient)

	var query struct {
		Viewer struct {
			Login     githubv4.String
			CreatedAt githubv4.DateTime
		}
	}

	err = client.Query(context.Background(), &query, nil)
	if err != nil {
		log.Fatal(err)
	}
	pretty.Println(query)
}

//gistsnip:end:main
