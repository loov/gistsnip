package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/google/go-github/github"
	"github.com/kr/pretty"
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

	client := github.NewClient(httpClient)

	currentUser, _, err := client.Users.Get("")
	if err != nil {
		log.Fatal(err)
	}

	pretty.Println(currentUser)

	gist := &github.Gist{}
	gist.Owner = currentUser
	gist.Description = github.String("current-directory")
	gist.Public = github.Bool(false)
	gist.Files = map[github.GistFilename]github.GistFile{}

	for _, file := range newGist.Files {
		for _, snippet := range file.Snippets {
			// todo better
			name := file.Path + "." + snippet.Name + ".go"
			gist.Files[github.GistFilename(name)] = github.GistFile{
				Content: github.String(snippet.Content),
			}
		}
	}

	result, _, err := client.Gists.Create(gist)
	if err != nil {
		log.Fatal(err)
	}
	pretty.Println(result)
}

//gistsnip:end:main
