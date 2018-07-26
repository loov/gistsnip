package main

import (
	"bufio"
	"encoding/json"
	"os"
)

type Config struct {
	Desc  string
	Files map[string]*File
}

type File struct {
	Path     string
	Snippets map[string]Snippet
}

type Snippet struct {
	Name    string
	Github  string
	Content string
}

func NewConfig() *Config {
	return &Config{Desc: "", Files: make(map[string]*File)}
}

func LoadConfig(name string) (*Config, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	config := NewConfig()
	err = json.NewDecoder(bufio.NewReader(f)).Decode(f)
	return config, err
}
