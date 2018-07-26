package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	description = flag.String("description", "", "gist description")

	gistsnip    = flag.String("gistsnip", ".gistsnip", "gistsnip info file")
	githubToken = flag.String("github", os.Getenv("GISTSNIP_TOKEN"), "github authentication token")
)

func main() {
	flag.Parse()

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

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *githubToken})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)

	client := github.NewClient(httpClient)

	currentUser, _, err := client.Users.Get("")
	if err != nil {
		log.Fatal(err)
	}

	for gistName, snippet := range newGist.Snippets {
		oldSnippet, exists := oldGist.Snippets[gistName]
		if exists {
			snippet.GistID = oldSnippet.GistID
			snippet.GistURL = oldSnippet.GistURL
		}
	}

	snippets := []*Snippet{}
	for gistName, snippet := range newGist.Snippets {
		snippets = append(snippets, snippet)

		oldSnippet, exists := oldGist.Snippets[gistName]
		if exists && oldSnippet.EqualContent(snippet) && oldSnippet.GistID != "" {
			continue
		}

		description := newGist.Description

		if link, err := GithubLinkToFile(snippet.File, snippet.Line); err == nil {
			if description != "" {
				description += "\n"
			}
			description += link
		}

		gist := &github.Gist{}
		gist.Owner = currentUser
		gist.Description = github.String(description)
		gist.Public = github.Bool(false)
		gist.Files = map[github.GistFilename]github.GistFile{}

		gist.Files[github.GistFilename(snippet.Path)] = github.GistFile{
			Content: github.String(snippet.Content),
		}

		if oldSnippet, ok := oldGist.Snippets[gistName]; ok {
			if oldSnippet.GistID != "" {
				_, _, err := client.Gists.Edit(oldSnippet.GistID, gist)
				if err != nil {
					log.Fatal(err)
				}
				continue
			}
		}

		result, _, err := client.Gists.Create(gist)
		if err != nil {
			log.Fatal(err)
		}

		snippet.GistID = *result.ID
		snippet.GistURL = *result.HTMLURL
	}

	err = SaveGist(*gistsnip, newGist)
	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(snippets, func(i, k int) bool {
		return snippets[i].Path < snippets[k].Path
	})

	for _, snippet := range snippets {
		fmt.Println(snippet.Path, snippet.GistURL)
	}
}
