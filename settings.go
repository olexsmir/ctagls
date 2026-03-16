package main

type Settings struct {
	CTags    string
	TagsFile string
}

func (s *Settings) EnsureDefaults() {
	if s.CTags == "" {
		s.CTags = "ctags"
	}
	if s.TagsFile == "" {
		s.TagsFile = "tags"
	}
}
