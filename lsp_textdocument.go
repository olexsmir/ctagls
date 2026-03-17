package main

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
)

type TextDocumentDidOpenParams struct {
	TextDocument TextDocument `json:"textDocument"`
}

type TextDocument struct {
	URI  string `json:"uri"`
	Text string `json:"text"`
}

func (s *LspServer) handleTextDocumentDidOpen(req RPCRequest) {
	if !s.isRootSet() {
		return
	}

	var params TextDocumentDidOpenParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return
	}

	uri, err := uriToPath(params.TextDocument.URI)
	if err != nil {
		return
	}

	s.mu.Lock()
	s.fileCache[uri] = strings.Split(params.TextDocument.Text, "\n")
	s.mu.Unlock()
}

type TextDocumentDidChangeParams struct {
	TextDocument   TextDocument                     `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type TextDocumentContentChangeEvent struct {
	Text string `json:"text"`
}

func (s *LspServer) handleTextDocumentDidChange(req RPCRequest) {
	var params TextDocumentDidChangeParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return
	}

	uri, err := uriToPath(params.TextDocument.URI)
	if err != nil {
		return
	}

	// FIXME: is it right ???
	if len(params.ContentChanges) > 0 {
		s.mu.Lock()
		s.fileCache[uri] = strings.Split(params.ContentChanges[0].Text, "\n")
		s.mu.Unlock()
	}
}

type TextDocumentDidCloseParams struct {
	TextDocument TextDocument `json:"textDocument"`
}

func (s *LspServer) handleTextDocumentDidClose(req RPCRequest) {
	var params TextDocumentDidCloseParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return
	}

	uri, err := uriToPath(params.TextDocument.URI)
	if err != nil {
		return
	}

	s.mu.Lock()
	delete(s.fileCache, uri)
	s.mu.Unlock()
}

type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type TextDocumentDefinitionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

type Location struct {
	URI   string `json:"uri"`
	Range Range  `json:"range"`
}

func (s *LspServer) handleTextDocumentDefinition(req RPCRequest) {
	var params TextDocumentDefinitionParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.invalidParams(req.ID, err)
		return
	}

	uri, err := uriToPath(params.TextDocument.URI)
	if err != nil {
		s.invalidParams(req.ID, err)
		return
	}

	symbol, err := s.wordAtPosition(uri, params.Position)
	if err != nil {
		s.sendResult(req.ID, nil)
		return
	}

	var locations []Location

	s.mu.Lock()
	// TODO: i bet there's some algo i should use
	for _, entry := range s.tagEntries {
		if entry.Name == symbol {
			content, err := s.getFileContent(entry.Path)
			if err != nil {
				s.sendMessage("failed to load content for: " + entry.Path)
				continue
			}

			symbolRange := s.findSymbolRangeInFile(content, entry.Name, entry.Line)
			locations = append(locations, Location{
				URI:   pathToFileURI(entry.Path),
				Range: symbolRange,
			})
		}
	}
	s.mu.Unlock()

	if len(locations) == 0 {
		s.sendResult(req.ID, nil)
	} else if len(locations) == 1 {
		s.sendResult(req.ID, locations[0])
	} else {
		s.sendResult(req.ID, locations)
	}
}

type TextDocumentDocumentSymbol struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

func (s *LspServer) handleTextDocumentDocumentSymbol(req RPCRequest) {
	var params TextDocumentDocumentSymbol
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.invalidParams(req.ID, err)
		return
	}

	uri, err := uriToPath(params.TextDocument.URI)
	if err != nil {
		s.invalidParams(req.ID, err)
		return
	}

	var symbols []SymbolInformation
	s.mu.RLock()

	var documentEntries []TagEntry
	for _, entry := range s.tagEntries {
		if entry.Path != uri {
			continue
		}
		documentEntries = append(documentEntries, entry)
	}

	sort.SliceStable(documentEntries, func(i, j int) bool {
		return documentEntries[i].Line < documentEntries[j].Line
	})

	for _, entry := range documentEntries {
		kind := s.ctags.LspSymbolKind(entry.Kind)
		// if kind == 0 {
		// 	continue
		// }

		content, err := s.getFileContent(entry.Path)
		if err != nil {
			s.sendMessage("failed to load content for: " + entry.Path)
			continue
		}

		symbolRange := s.findSymbolRangeInFile(content, entry.Name, entry.Line)
		symbols = append(symbols, SymbolInformation{
			Name:          entry.Name,
			Kind:          kind,
			ContainerName: entry.Scope,
			Location: Location{
				URI:   pathToFileURI(entry.Path),
				Range: symbolRange,
			},
		})
	}

	s.mu.RUnlock()

	s.sendResult(req.ID, symbols)
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type TextDocumentCodeActionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                  `json:"range"`
	Context      CodeActionContext      `json:"context"`
}

type CodeActionContext struct {
	Diagnostics []Diagnostic `json:"diagnostics"`
}

type Diagnostic struct {
	Range    Range  `json:"range"`
	Message  string `json:"message"`
	Severity *int   `json:"severity,omitempty"`
	Source   string `json:"source,omitempty"`
}

type Command struct {
	Title     string `json:"title"`
	Command   string `json:"command"`
	Arguments []any  `json:"arguments,omitempty"`
}

type CodeAction struct {
	Title   string   `json:"title"`
	Kind    string   `json:"kind"`
	Command *Command `json:"command"`
}

const CtagLSReindexAction = "ctagls.reindex"

func (s *LspServer) handleTextDocumentCodeAction(req RPCRequest) {
	var params TextDocumentCodeActionParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.invalidParams(req.ID, err)
		return
	}

	s.sendResult(req.ID, []CodeAction{
		{
			Title: "Re-index tags",
			Kind:  "source",
			Command: &Command{
				Title:   "Re-index tags",
				Command: CtagLSReindexAction,
			},
		},
	})
}

func (s *LspServer) findSymbolRangeInFile(lines []string, symbolName string, lineNumber int) Range {
	lineIdx := lineNumber - 1
	if lineIdx < 0 || lineIdx >= len(lines) {
		return Range{
			Start: Position{Line: lineIdx, Character: 0},
			End:   Position{Line: lineIdx, Character: 0},
		}
	}

	lineContent := lines[lineIdx]
	startChar := strings.Index(lineContent, symbolName)
	if startChar == -1 {
		return Range{
			Start: Position{Line: lineIdx, Character: 0},
			End:   Position{Line: lineIdx, Character: len([]rune(lineContent))},
		}
	}

	endChar := startChar + len([]rune(symbolName))

	return Range{
		Start: Position{Line: lineIdx, Character: startChar},
		End:   Position{Line: lineIdx, Character: endChar},
	}
}

// getFileContent ... PLEASE LOCK THE MUTEXT FOR IT
// gets file content from cache, or fallbacks to reading form the disk
func (s *LspServer) getFileContent(fpath string) ([]string, error) {
	content, ok := s.fileCache[fpath]
	if ok {
		return content, nil
	}

	contentBytes, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	return strings.Split(
		string(contentBytes),
		"\n",
	), nil
}
