package main

import (
	"errors"
	"os/exec"
	"strings"
)

var ErrCTagsNotInstalled = errors.New("ctags not found in $PATH")

type CTags struct {
	bin  string
	tags string
}

func NewCTags(ctagsBin, tagsFile string) *CTags {
	return &CTags{
		bin:  ctagsBin,
		tags: tagsFile,
	}
}

func (c *CTags) CheckIfInstalled() error {
	out, err := exec.Command(c.bin, "--version").Output()
	if err != nil || !strings.Contains(string(out), "Universal Ctags") {
		return ErrCTagsNotInstalled
	}
	return nil
}

type TagEntry struct{}

func (c *CTags) Parse() ([]TagEntry, error) {
	return nil, nil
}

func (c *CTags) run() {
}
