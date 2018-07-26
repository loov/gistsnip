package main

import (
	"bufio"
	"encoding/json"
	"os"
)

type Gist struct {
	GistID string
	Desc   string
	Files  map[string]*File
}

type File struct {
	Path     string
	Snippets map[string]*Snippet
}

type Snippet struct {
	Name    string
	GistID  string
	Content string
}

func NewGist() *Gist {
	return &Gist{Desc: "", Files: make(map[string]*File)}
}

func NewFile(path string) *File {
	return &File{Path: path, Snippets: make(map[string]*Snippet)}
}

func LoadConfig(name string) (*Gist, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gist := NewGist()
	err = json.NewDecoder(bufio.NewReader(f)).Decode(f)
	return gist, err
}
