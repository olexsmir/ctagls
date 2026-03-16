package main

import (
	"encoding/json"
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

	uri, err := normalizeFileURI(params.TextDocument.URI)
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

	uri, err := normalizeFileURI(params.TextDocument.Text)
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

	uri, err := normalizeFileURI(params.TextDocument.URI)
	if err != nil {
		return
	}

	s.mu.Lock()
	delete(s.fileCache, uri)
	s.mu.Unlock()
}

type TextDocumentDidSaveParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

func (s *LspServer) handleTextDocumentDidSave(req RPCRequest) {
	var params TextDocumentDidSaveParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return
	}

	uri, err := normalizeFileURI(params.TextDocument.URI)
	if err != nil {
		return
	}

	_ = uri

	// TODO: rescan for root markers
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

	uri, err := normalizeFileURI(params.TextDocument.URI)
	if err != nil {
		s.invalidParams(req.ID, err)
		return
	}

	_ = uri
}

func (s *LspServer) handleTextDocumentDocumentSymbol(req RPCRequest) {}
