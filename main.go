package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var (
	description = flag.String("description", "{{.RepositoryLink}} {{.SourceLink}}", "gist description template")

	gistsnip    = flag.String("gistsnip", ".gistsnip", "gistsnip info file")
	githubToken = flag.String("github", os.Getenv("GISTSNIP_TOKEN"), "github authentication token")
)

func main() {
	ctx := context.Background()
	flag.Parse()

	paths := flag.Args()
	if len(paths) == 0 {
		paths = []string{"."}
	}

	descriptionTemplate, err := template.New("").Parse(*description)
	if err != nil {
		log.Fatal(err)
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

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *githubToken})
	httpClient := oauth2.NewClient(ctx, tokenSource)

	client := github.NewClient(httpClient)

	currentUser, _, err := client.Users.Get(ctx, "")
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

		repositoryLink, sourceLink, err := GithubLinkToFile(snippet.File, snippet.Line)
		if err != nil {
			log.Println("Failed to create github link", err)
		}
		
		var description strings.Builder
		err = descriptionTemplate.Execute(&description, map[string]interface{}{
			"RepositoryLink": repositoryLink,
			"SourceLink": sourceLink,
		})
		snippet.Description = description.String()

		oldSnippet, exists := oldGist.Snippets[gistName]
		if exists && oldSnippet.EqualContent(snippet) && oldSnippet.GistID != "" {
			log.Println("Skipping ", gistName)
			continue
		}

		if exists {
			log.Println("Updating ", gistName)
		} else {
			log.Println("Creating ", gistName)
		}

		gist := &github.Gist{}
		gist.Owner = currentUser
		gist.Description = github.String(snippet.Description)
		gist.Public = github.Bool(false)
		gist.Files = map[github.GistFilename]github.GistFile{}

		sanitizedPath := strings.Replace(snippet.Path, "/", "\\", -1)
		gist.Files[github.GistFilename(sanitizedPath)] = github.GistFile{
			Content: github.String(snippet.Content),
		}

		if oldSnippet, ok := oldGist.Snippets[gistName]; ok {
			if oldSnippet.GistID != "" {
				_, _, err := client.Gists.Edit(ctx, oldSnippet.GistID, gist)
				if err != nil {
					log.Fatal(err)
				}
				continue
			}
		}

		result, _, err := client.Gists.Create(ctx, gist)
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
