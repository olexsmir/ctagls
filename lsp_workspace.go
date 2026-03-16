package main

import (
	"encoding/json"
	"fmt"
)

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

func (s *LspServer) reindex() error {
	return nil
}
