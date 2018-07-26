package main

import (
	"bufio"
	"encoding/json"
	"os"
)

type Gist struct {
	GistID  string `json:"omitempty"`
	GistURL string `json:"omitempty"`

	Description string
	Files       map[string]*File
}

type File struct {
	Path     string
	Snippets map[string]*Snippet
}

type Snippet struct {
	GistID  string `json:"omitempty"`
	GistURL string `json:"omitempty"`

	Name    string
	Content string
}

func NewGist() *Gist {
	return &Gist{Files: make(map[string]*File)}
}

func NewFile(path string) *File {
	return &File{Path: path, Snippets: make(map[string]*Snippet)}
}

func LoadGist(name string) (*Gist, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gist := NewGist()
	err = json.NewDecoder(bufio.NewReader(f)).Decode(f)
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

func (gist *Gist) EqualContent(next *Gist) bool {
	// TODO
	return false
}
