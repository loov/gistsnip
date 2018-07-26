package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

var rxTag = regexp.MustCompile(`(?i)//gistsnip:(\w+):(\w+)`)

func GistFromGlobs(globs []string) (*Gist, error) {
	gist := NewGist()
	for _, glob := range globs {
		if err := gist.IncludeGlob(glob); err != nil {
			return gist, err
		}
	}
	return gist, nil
}

func (gist *Gist) IncludeGlob(glob string) error {
	matches, err := filepath.Glob(glob)
	if err != nil {
		return err
	}
	for _, match := range matches {
		if match == "." {
			gist.IncludeDir(".")
		} else {
			if strings.HasPrefix(filepath.Base(match), ".") {
				continue
			}

			stat, err := os.Lstat(match)
			if err != nil {
				return err
			}
			if stat.IsDir() {
				gist.IncludeDir(match)
			} else {
				gist.IncludeFile(match)
			}
		}
	}
	return nil
}

func (gist *Gist) IncludeDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || path == "." || path == ".." {
			return err
		}

		// ignore hidden
		rel, _ := filepath.Rel(dir, path)
		if strings.HasPrefix(rel, ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			return nil
		}

		err = gist.IncludeFile(path)
		if err != nil {
			return err
		}
		return nil
	})
}

func (gist *Gist) IncludeFile(filename string) error {
	path := filepath.ToSlash(filename)
	file, exists := gist.Files[path]
	if !exists {
		file = NewFile(path)
		defer func() {
			if len(file.Snippets) > 0 {
				gist.Files[path] = file
			}
		}()
	}

	stat, _ := os.Stat(filename)
	if stat.Size() > 1<<20 {
		// larger than megabyte, probably not a source file
		return nil
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if !utf8.Valid(content) {
		// probably a binary file
		return nil
	}

	tags := ParseTags(content)
	contents := ParseSnippetContent(content, tags)

	for _, content := range contents {
		snip, multiple := file.Snippets[content.Name]
		if multiple {
			snip.Content += "\n\n" + content.Content
			continue
		}

		file.Snippets[content.Name] = &Snippet{
			Name:    content.Name,
			Content: content.Content,
		}
	}

	return nil
}

type Tag struct {
	Start  int
	End    int
	Action string
	Value  string
}

func ParseTags(content []byte) []Tag {
	tagsIndices := rxTag.FindAllSubmatchIndex(content, -1)

	tags := []Tag{}
	for _, ix := range tagsIndices {
		tags = append(tags, Tag{
			Start:  ix[0],
			End:    ix[1],
			Action: strings.ToLower(string(content[ix[2]:ix[3]])),
			Value:  strings.ToLower(string(content[ix[4]:ix[5]])),
		})
	}

	return tags
}

type SnippetContent struct {
	Name    string
	Content string
}

func ParseSnippetContent(content []byte, initialTags []Tag) []SnippetContent {
	snippets := []SnippetContent{}

	tags := append(initialTags[:0], initialTags...)

	for head := 0; head < len(tags); head++ {
		switch tags[head].Action {
		case "start":
			start := tags[head]
			end := Tag{Start: len(content), End: len(content), Action: "end", Value: start.Value}

			for tail := len(tags) - 1; head < tail; tail-- {
				if tags[tail].Action == "end" && strings.EqualFold(tags[tail].Value, start.Value) {
					end = tags[tail]
					tags = append(tags[:tail], tags[tail+1:]...)
					break
				}
			}

			snippets = append(snippets, SnippetContent{
				Name:    start.Value,
				Content: strings.Trim(string(content[start.End:end.Start]), "\n\r"),
			})
		}
	}

	return snippets
}