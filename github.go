package main

import (
	"errors"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// one of:
// https://github.com/loov/watchrun.git
// git@github.com:loov/watchrun.git
var rxGithub = regexp.MustCompile(`^(?:https://github.com/|git@github.com:)(.*)\.git$`)

func GithubLinkToFile(path string, line int) (repositoryLink, sourceLink string, err error) {
	dir := filepath.Dir(path)

	remoteurl, err := gitexec(dir, "remote", "get-url", "origin")
	if err != nil {
		return "", "", err
	}

	filename, err := gitexec(dir, "ls-files", "--full-name", filepath.Base(path))
	if err != nil {
		return "", "", err
	}

	hash, err := gitexec(dir, "rev-parse", "HEAD")
	if err != nil {
		return "", "", err
	}

	// remoteurl = git@github.com:loov/watchrun.git
	// filename  = 00_yolo/main.go
	// hash      = 872c42f3b01ebc324fb2e0ac1a51e7b5539c7b50
	// result: https://github.com/egonelbre/db-demo/blob/872c42f3b01ebc324fb2e0ac1a51e7b5539c7b50/00_yolo/main.go#L11

	matches := rxGithub.FindStringSubmatch(remoteurl)
	if len(matches) == 0 {
		return "", "", errors.New("not a github repository")
	}

	repository := matches[1]
	return "https://github.com/" + repository, "https://github.com/" + repository + "/blob/" + hash + "/" + filename + "#L" + strconv.Itoa(line), nil
}

func gitexec(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	result, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(result)), err
}
