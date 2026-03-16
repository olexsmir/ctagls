package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type InitializeParams struct {
	RootURI          string            `json:"rootUri"`
	RootPath         string            `json:"rootPath"`
	WorkspaceFolders []WorkspaceFolder `json:"workspaceFolders"`
}

type WorkspaceFolder struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type TextDocumentSyncOptions struct {
	Change    int  `json:"change"`
	OpenClose bool `json:"openClose"`
	Save      bool `json:"save"`
}

type ServerCapabilities struct {
	TextDocumentSync        TextDocumentSyncOptions `json:"textDocumentSync"`
	Workspace               WorkspaceCapabilities   `json:"workspace"`
	DefinitionProvider      bool                    `json:"definitionProvider,omitempty"`
	WorkspaceSymbolProvider bool                    `json:"workspaceSymbolProvider,omitempty"`
	DocumentSymbolProvider  bool                    `json:"documentSymbolProvider,omitempty"`
}

type WorkspaceCapabilities struct {
	Configuration          bool                             `json:"configuration,omitempty"`
	DidChangeConfiguration DidChangeConfigurationCapability `json:"didChangeConfiguration"`
}

type DidChangeConfigurationCapability struct {
	DynamicRegistration bool `json:"dynamicRegistration"`
}

type InitializeResponse struct {
	Info         ServerInfo         `json:"serverInfo"`
	Capabilities ServerCapabilities `json:"capabilities"`
}

func (s *LspServer) handleInitialize(req RPCRequest) {
	var params InitializeParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.sendError(req.ID, InvalidParamsCode, "Invalid params", nil)
		return
	}

	rootURI, err := s.selectRootURI(params)
	if err != nil {
		s.sendError(req.ID, InvalidParamsCode, "Invalid paraams", err.Error())
		return
	}

	if rootURI != "" {
		if err := s.setRootURI(rootURI); err != nil {
			s.sendError(req.ID, InvalidParamsCode, "Invalid paraams", err.Error())
			return
		}
	}

	s.sendResult(req.ID, InitializeResponse{
		Info: ServerInfo{
			Name:    name,
			Version: version,
		},
		Capabilities: ServerCapabilities{
			TextDocumentSync: TextDocumentSyncOptions{
				Change:    1, //  TextDocumentSyncKindFull // TODO: improve me
				OpenClose: true,
				Save:      true,
			},
			WorkspaceSymbolProvider: true,
			DefinitionProvider:      true,
			DocumentSymbolProvider:  true,
			Workspace: WorkspaceCapabilities{
				Configuration: true,
				DidChangeConfiguration: DidChangeConfigurationCapability{
					DynamicRegistration: true,
				},
			},
		},
	})
	s.initialized = true
}

func (s *LspServer) handleInitialized(_ RPCRequest) {
	s.setupSettingsAndReindexTags(&Settings{})
}

func (s *LspServer) handleShutdown(req RPCRequest) {
	// TODO: check the proper way of handing it
	s.sendResult(req.ID, nil)
}

func (s *LspServer) handleExit(_ RPCRequest) {
	// TODO: check the proper way of handing it
	os.Exit(0)
}

func (s *LspServer) selectRootURI(params InitializeParams) (string, error) {
	if len(params.WorkspaceFolders) > 0 {
		// TODO: support multiple workspaces
		return params.WorkspaceFolders[0].URI, nil
	}

	if params.RootURI != "" {
		return params.RootURI, nil
	}

	if params.RootPath != "" {
		cleanPath := filepath.Clean(params.RootPath)
		absPath, err := filepath.Abs(cleanPath)
		if err != nil {
			return "", err
		}
		return pathToFileURI(absPath), nil
	}

	return "", nil
}
