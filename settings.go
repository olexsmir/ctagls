package main

import "path/filepath"

type Settings struct {
	CTags    string
	TagsFile string
}

func (s *Settings) EnsureDefaults(rootPath string) {
	if s.CTags == "" {
		s.CTags = "ctags"
	}
	if s.TagsFile == "" {
		s.TagsFile = "tags"
	}
	if !filepath.IsAbs(s.TagsFile) {
		s.TagsFile = filepath.Join(rootPath, s.TagsFile)
	}
}
