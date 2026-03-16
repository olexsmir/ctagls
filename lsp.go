package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const (
	ParseErrorCode     = -32700
	InvalidRequestCode = -32600
	MethodNotFoundCode = -32601
	InvalidParamsCode  = -32602
	InternalErrorCode  = -32603

	MessageTypeWarning = 2
)

type RPCRequest struct {
	RPC    string           `json:"jsonrpc"`
	ID     *json.RawMessage `json:"id,omitempty"`
	Method string           `json:"method"`
	Params json.RawMessage  `json:"params,omitempty"`
}

func (r RPCRequest) isNotification() bool { return r.ID == nil }

type RPCNotification struct {
	RPC    string `json:"jsonrpc"`
	Method string `json:"method"`
	Params any    `json:"params,omitempty"`
}

type ShowMessageParams struct {
	Type    int    `json:"type"`
	Message string `json:"message"`
}

type RPCSuccessResponse struct {
	RPC    string           `json:"jsonrpc"`
	ID     *json.RawMessage `json:"id"`
	Result any              `json:"result"`
}

type RPCErrorResponse struct {
	RPC   string           `json:"jsonrpc"`
	ID    *json.RawMessage `json:"id"`
	Error *RPCError        `json:"error"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type LspServer struct {
	initialized bool
	rootURI     string

	mu        sync.RWMutex
	fileCache map[string][]string

	ctags    *CTags
	settings *Settings

	input  io.Reader
	output io.Writer
}

func NewLspServer(in io.Reader, out io.Writer) *LspServer {
	return &LspServer{
		input:     in,
		output:    out,
		fileCache: make(map[string][]string),
	}
}

func (s *LspServer) Handle() error {
	reader := bufio.NewReader(s.input)
	for {
		req, err := s.readMsg(reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			s.sendError(nil, InvalidRequestCode, "Malformed request", err.Error())
			continue
		}
		go s.handleRequest(req)
	}
}

func (s *LspServer) handleRequest(req RPCRequest) {
	b, _ := json.Marshal(req)
	slog.Debug("got request", "req", b)

	if !s.initialized && req.Method != "initialize" && req.Method != "shutdown" && req.Method != "exit" {
		s.sendError(req.ID, InvalidParamsCode, "Server not initialized.", "Received request before successful initialization")
		return
	}

	switch req.Method {
	case "initialize":
		s.handleInitialize(req)
	case "initialized":
		// here so we wound't fail
	case "shutdown":
		s.handleShutdown(req)
	case "exit":
		s.handleExit(req)
	case "textDocument/didOpen":
		s.handleTextDocumentDidOpen(req)
	case "textDocument/didChange":
		s.handleTextDocumentDidChange(req)
	case "textDocument/didClose":
		s.handleTextDocumentDidClose(req)
	case "textDocument/didSave":
		s.handleTextDocumentDidSave(req)
	case "textDocument/definition":
		s.handleTextDocumentDefinition(req)
	case "textDocument/documentSymbol":
		s.handleTextDocumentDocumentSymbol(req)
	case "textDocument/codeAction":
		s.handleTextDocumentCodeAction(req)
	case "workspace/executeCommand":
		s.handleWorkspaceExecuteCommand(req)
	case "workspace/didChangeConfiguration":
		s.handleWorkspaceDidChangeConfiguration(req)
	case "workspace/symbol":
		s.handleWorkspaceSymbol(req)
	default:
		s.sendError(req.ID, MethodNotFoundCode,
			fmt.Sprintf("Method not found: %s", req.Method), nil)
	}
}

func (s *LspServer) readMsg(r *bufio.Reader) (RPCRequest, error) {
	contentLength := 0
	for {
		line, err := r.ReadString('\r')
		if err != nil {
			return RPCRequest{}, fmt.Errorf("error reading header: %w", err)
		}

		b, err := r.ReadByte()
		if err != nil {
			return RPCRequest{}, fmt.Errorf("error reading header: %w", err)
		}
		if b != '\n' {
			return RPCRequest{}, fmt.Errorf("line endings must be \\r\\n")
		}
		if line == "\r" {
			break
		}
		if after, ok := strings.CutPrefix(line, "Content-Length:"); ok {
			clStr := strings.TrimSpace(after)
			cl, err := strconv.Atoi(clStr)
			if err != nil {
				return RPCRequest{}, fmt.Errorf("invalid Content-Length: %v", err)
			}
			contentLength = cl
		}
	}

	body := make([]byte, contentLength)
	_, err := io.ReadFull(r, body)
	if err != nil {
		return RPCRequest{}, fmt.Errorf("error reading body: %w", err)
	}

	var req RPCRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return RPCRequest{}, fmt.Errorf("invalid JSON-RPC request: %v", err)
	}
	if isInvalidID(req.ID) {
		return RPCRequest{}, fmt.Errorf("id must be a string or integer")
	}
	return req, nil
}

func (s *LspServer) isRootSet() bool { return s.rootURI != "" }
func (s *LspServer) setRootURI(rootURI string) error {
	s.rootURI = rootURI
	// TODO: resolve tags files in the root
	return nil
}

func (s *LspServer) sendResponse(resp any) {
	body, _ := json.Marshal(resp)
	fmt.Fprintf(s.output, "Content-Length: %d\r\n\r\n%s",
		len(body), string(body))
}

func (s *LspServer) sendResult(id *json.RawMessage, result any) {
	b, _ := json.Marshal(result)
	slog.Debug("send result", "id", id, "result", string(b))

	s.sendResponse(RPCSuccessResponse{
		RPC:    "2.0",
		ID:     id,
		Result: result,
	})
}

func (s *LspServer) sendMessage(err error) {
	slog.Debug("send message", "err", err)
	s.sendResponse(RPCNotification{
		RPC:    "2.0",
		Method: "window/showMessage",
		Params: ShowMessageParams{
			Type:    MessageTypeWarning,
			Message: msg,
		},
	})
}

func (s *LspServer) sendError(id *json.RawMessage, code int, msg string, data any) {
	b, _ := json.Marshal(data)
	slog.Debug("send error: ", "id", id, "code", code, "msg", msg, "data", b)

	s.sendResponse(RPCErrorResponse{
		RPC: "2.0",
		ID:  id,
		Error: &RPCError{
			Code:    code,
			Message: msg,
			Data:    data,
		},
	})
}

func (s *LspServer) invalidParams(id *json.RawMessage, err error) {
	s.sendError(id, InvalidParamsCode, "Invalid params", err.Error())
}

func isInvalidID(id *json.RawMessage) bool {
	if id == nil {
		return false
	}

	var s string
	if json.Unmarshal(*id, &s) == nil {
		return false
	}

	var n int64
	return json.Unmarshal(*id, &n) != nil
}

func normalizeFileURI(uri string) (string, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return "", fmt.Errorf("failed to parse URI %q: %w", uri, err)
	}
	if parsed.Scheme != "file" {
		return "", fmt.Errorf("expected file:// URI: %q", uri)
	}
	if parsed.Path == "" {
		return "", fmt.Errorf("empty file URI")
	}

	path := filepath.Clean(filepath.FromSlash(parsed.Path))
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path %q: %w", path, err)
	}
	return pathToFileURI(absPath), nil
}

func pathToFileURI(path string) string {
	slashPath := filepath.ToSlash(path)
	return (&url.URL{Scheme: "file", Path: slashPath}).String()
}
