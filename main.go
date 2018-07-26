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
	description = flag.String("description", "", "gist description")

	gistsnip    = flag.String("gistsnip", ".gistsnip", "gistsnip info file")
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

	oldGist, err := LoadGist(*gistsnip)
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

	newGist.Description = *description
	if newGist.Description == "" {
		newGist.Description = oldGist.Description
	}

	pretty.Println(oldGist)
	pretty.Println(newGist)

	pretty.Println(newGist.ChangedSnippets(oldGist))

	err = SaveGist(*gistsnip, newGist)
	if err != nil {
		log.Fatal(err)
	}

	return

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
	gist.Description = github.String(newGist.Description)
	gist.Public = github.Bool(false)
	gist.Files = map[github.GistFilename]github.GistFile{}

	for _, snippet := range newGist.Snippets {
		// todo better
		gist.Files[github.GistFilename(snippet.Path)] = github.GistFile{
			Content: github.String(snippet.Content),
		}
	}

	result, _, err := client.Gists.Create(gist)
	if err != nil {
		log.Fatal(err)
	}
	pretty.Println(result)
}

//gistsnip:end:main
