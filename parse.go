package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"
)

var rxTag = regexp.MustCompile(`(?i)//gistsnip:(\w+):(\w+)`)
var rxLine = regexp.MustCompile(`(?im)^\s*//gistsnip:(\w+):(\w+)\s*$\n`)

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
		snipPath := SnippetPath(filename, content.Name)

		snip, multiple := gist.Snippets[snipPath]
		if multiple {
			snip.Content += "\n\n" + content.Content
			continue
		}

		gist.Snippets[snipPath] = &Snippet{
			Name:    content.Name,
			Line:    content.Line,
			File:    filepath.ToSlash(filename),
			Path:    filepath.ToSlash(snipPath),
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
	Line    int
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

			for tail := head + 1; tail < len(tags); tail++ {
				if tags[tail].Action == "end" && strings.EqualFold(tags[tail].Value, start.Value) {
					end = tags[tail]
					tags = append(tags[:tail], tags[tail+1:]...)
					break
				}
			}

			text := string(content[start.End:end.Start])

			// remove any nested gistsnips
			text = rxLine.ReplaceAllString(text, "")

			text = strings.TrimLeft(text, "\n\r")
			text = strings.TrimRight(text, " \n\r\t")

			text = Dedent(text, '\t')
			text = Dedent(text, ' ')

			snippets = append(snippets, SnippetContent{
				Name:    start.Value,
				Line:    bytes.Count(content[:start.End], []byte{'\n'}) + 1,
				Content: text,
			})
		}
	}

	return snippets
}

func Dedent(text string, delimiter rune) string {
	minIndent := 1000
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := 0
		for _, r := range line {
			if r == delimiter {
				indent++
			} else {
				break
			}
		}

		if indent < minIndent {
			minIndent = indent
		}
	}

	rxDedent := regexp.MustCompile(`(?m)^` + strings.Repeat(string(delimiter), minIndent))
	text = rxDedent.ReplaceAllString(text, "")

	return text
}
