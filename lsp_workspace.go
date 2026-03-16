package main

import (
	"encoding/json"
	"fmt"
)

type SymbolInformation struct {
	Name          string   `json:"name"`
	Kind          int      `json:"kind"`
	Location      Location `json:"location"`
	ContainerName string   `json:"containerName,omitempty"`
}

type WorkspaceSymbolParams struct {
	Query string `json:"query"`
}

func (s *LspServer) handleWorkspaceSymbol(req RPCRequest) {
	var params WorkspaceSymbolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.invalidParams(req.ID, err)
		return
	}

	var symbols []SymbolInformation
	s.mu.RLock()
	for _, entry := range s.tagEntries {
		if params.Query != "" && entry.Name != params.Query {
			continue
		}

		kind := s.ctags.LspSymbolKind(entry.Kind)
		if kind == 0 {
			continue
		}

		content, err := s.getFileContent(entry.Path)
		if err != nil {
			s.sendMessage("failed to load content for: " + entry.Path)
			continue
		}

		symbolRange := s.findSymbolRangeInFile(content, entry.Name, entry.Line)
		symbols = append(symbols, SymbolInformation{
			Name:          entry.Kind,
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

type WorkspaceDidChangeConfiguration struct {
	Settings ServerSettings `json:"settings"`
}

type ServerSettings struct {
	CTags    string `json:"ctags"`
	TagsFile string `json:"tagsFile"`
}

func (s *LspServer) handleWorkspaceDidChangeConfiguration(req RPCRequest) {
	// TODO: handle wrong config types passed

	var params WorkspaceDidChangeConfiguration
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return
	}

	s.settings = &Settings{
		CTags:    params.Settings.CTags,
		TagsFile: params.Settings.TagsFile,
	}
	s.settings.EnsureDefaults(s.root)
	s.setupCtags()

	if err := s.reindex(); err != nil {
		s.internalError(req.ID, err)
		return
	}
}

type WorkspaceExecuteCommandParams struct {
	Command   string `json:"command"`
	Arguments []any  `json:"arguments,omitempty"`
}

func (s *LspServer) handleWorkspaceExecuteCommand(req RPCRequest) {
	var params WorkspaceExecuteCommandParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.invalidParams(req.ID, err)
		return
	}

	switch params.Command {
	case CtagLSReindexAction:
		if err := s.regenerateTags(); err != nil {
			s.sendError(req.ID, InternalErrorCode, "ctag regenerate failed", err.Error())
			return
		}
		if err := s.reindex(); err != nil {
			s.sendError(req.ID, InternalErrorCode, "reindex failed", err.Error())
			return
		}

		s.sendResult(req.ID, struct{}{})

	default:
		s.sendError(req.ID, MethodNotFoundCode, fmt.Sprintf("Unknown command: %s", params.Command), nil)
		return
	}
}

func (s *LspServer) setupCtags() {
	ctags := NewCTags(s.settings.CTags, s.settings.TagsFile)
	s.ctags = ctags
}

// reindex, reindexes the tags file
func (s *LspServer) reindex() error {
	tags, err := s.ctags.Parse()
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.tagEntries = tags
	s.mu.Unlock()

	return nil
}

func (s *LspServer) regenerateTags() error {
	return s.ctags.Regenerate()
}
