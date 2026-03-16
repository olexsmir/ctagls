package main

import "encoding/json"

func (s *LspServer) handleWorkspaceSymbol(req RPCRequest) {}

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
	s.settings.EnsureDefaults()
}
