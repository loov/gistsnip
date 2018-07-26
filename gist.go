package main

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
)

type Gist struct {
	GistID  string `json:"omitempty"`
	GistURL string `json:"omitempty"`

	Description string
	Snippets    map[string]*Snippet
}

type Snippet struct {
	GistID  string `json:"omitempty"`
	GistURL string `json:"omitempty"`

	Path    string
	Name    string
	Content string
}

func SnippetPath(file, snippetName string) string {
	cfile := filepath.ToSlash(file)
	ext := filepath.Ext(cfile)
	root := file[:len(cfile)-len(ext)]
	return root + "#" + snippetName + ext
}

func NewGist() *Gist {
	return &Gist{Snippets: make(map[string]*Snippet)}
}

func LoadGist(name string) (*Gist, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gist := NewGist()
	err = json.NewDecoder(bufio.NewReader(f)).Decode(gist)
	return gist, err
}

func SaveGist(name string, gist *Gist) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	wr := bufio.NewWriter(f)
	defer wr.Flush()

	enc := json.NewEncoder(wr)
	enc.SetIndent("", "    ")

	err = enc.Encode(gist)
	return err
}

func (gist *Gist) EqualContent(old *Gist) bool {
	if len(gist.Snippets) != len(old.Snippets) {
		return false
	}

	return len(gist.ChangedSnippets(old)) == 0
}

func (gist *Gist) ChangedSnippets(old *Gist) []*Snippet {
	changed := []*Snippet{}
	for newSnipName, newSnippet := range gist.Snippets {
		oldSnippet, found := old.Snippets[newSnipName]
		if !found {
			changed = append(changed, newSnippet)
			continue
		}

		if newSnippet.Content != oldSnippet.Content {
			changed = append(changed, newSnippet)
			continue
		}
	}
	return changed
}
