package main

import (
	"bufio"
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var ErrCTagsNotInstalled = errors.New("ctags not found in $PATH")

type CTags struct {
	bin  string
	tags string
}

func NewCTags(ctagsBin, tagsFile string) *CTags {
	// TODO: set tagsFile as abs path to the file
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

type TagEntry struct {
	Name    string
	Path    string
	Pattern string
	Kind    string
	Line    int
	Scope   string
}

func (c *CTags) Parse() ([]TagEntry, error) {
	file, err := os.Open(c.tags)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	kindMap := newKindMap()
	entries := make([]TagEntry, 0, 1024)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if strings.HasPrefix(trimmed, "!") {
			c.parseKindDescription(line, kindMap)
			continue
		}

		entry, ok := c.parseEntry(line, kindMap)
		if ok {
			entries = append(entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func (c *CTags) parseKindDescription(line string, kindMap *kindMap) {
	if !strings.HasPrefix(line, "!_TAG_KIND_DESCRIPTION") {
		return
	}

	fields := strings.Split(line, "\t")
	if len(fields) < 2 {
		return
	}

	language := strings.TrimPrefix(fields[0], "!_TAG_KIND_DESCRIPTION")
	if after, ok := strings.CutPrefix(language, "!"); ok {
		language = after
	} else {
		language = ""
	}

	parts := strings.SplitN(fields[1], ",", 2)
	if len(parts) != 2 {
		return
	}

	letter := parts[0]
	kind := parts[1]
	if letter == "" || kind == "" {
		return
	}

	kindMap.add(language, letter, kind)
}

func (c *CTags) parseEntry(line string, kindMap *kindMap) (TagEntry, bool) {
	fields := strings.Split(line, "\t")
	if len(fields) < 3 {
		return TagEntry{}, false
	}

	entry := TagEntry{
		Name:    fields[0],
		Path:    fields[1], // TODO: make it relative to the tags file
		Pattern: strings.TrimSuffix(fields[2], ";\""),
	}

	scopeKindSet := false
	kindField, language := "", ""

	nextFieldIndex := 3
	if len(fields) > 3 && !strings.Contains(fields[3], ":") {
		kindField = fields[3]
		nextFieldIndex = 4
	}

	for _, field := range fields[nextFieldIndex:] {
		if field == "" {
			continue
		}
		key, value, ok := strings.Cut(field, ":")
		if !ok {
			continue
		}

		switch key {
		case "line":
			if lineNum, err := strconv.Atoi(value); err == nil {
				entry.Line = lineNum
			}
		case "language":
			language = value
		case "kind":
			kindField = value
		case "scope":
			entry.Scope = value
		case "scopeKind":
			scopeKindSet = true
		default:
			if entry.Scope == "" && !scopeKindSet && kindMap.isKindName(key) {
				entry.Scope = value
				scopeKindSet = true
			}
		}
	}

	if entry.Line == 0 {
		if lineNum, err := strconv.Atoi(entry.Pattern); err == nil {
			entry.Line = lineNum
		}
	}

	if kindField != "" {
		kindField = c.resolveKind(kindField, language, kindMap)
		entry.Kind = kindField
	}

	return entry, true
}

func (c *CTags) resolveKind(kindField, language string, kindMap *kindMap) string {
	if len(kindField) != 1 {
		return kindField
	}

	if mapped, ok := kindMap.resolve(language, kindField); ok {
		return mapped
	}
	return kindField
}

type kindMap struct {
	byLanguage map[string]map[string]string
	any        map[string]string
	kindNames  map[string]bool
}

func newKindMap() *kindMap {
	return &kindMap{
		byLanguage: make(map[string]map[string]string),
		any:        make(map[string]string),
		kindNames:  make(map[string]bool),
	}
}

func (k *kindMap) isKindName(kind string) bool { return k.kindNames[kind] }

func (k *kindMap) add(language, letter, kind string) {
	if language == "" {
		language = "default"
	}
	if _, ok := k.byLanguage[language]; !ok {
		k.byLanguage[language] = make(map[string]string)
	}
	k.byLanguage[language][letter] = kind
	if _, ok := k.any[letter]; !ok {
		k.any[letter] = kind
	}
	k.kindNames[kind] = true
}

func (k *kindMap) resolve(language, letter string) (string, bool) {
	if language != "" {
		if byLang, ok := k.byLanguage[language]; ok {
			if kind, ok := byLang[letter]; ok {
				return kind, true
			}
		}
	}
	if kind, ok := k.any[letter]; ok {
		return kind, true
	}
	return "", false
}
