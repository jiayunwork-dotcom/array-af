package web

import (
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

type Assets struct {
	WebFS    fs.FS
	Examples map[string][]byte
}

type Server struct {
	webFS    fs.FS
	examples map[string][]byte
	mux      *http.ServeMux
}

func NewServer(assets Assets) *Server {
	s := &Server{
		webFS:    assets.WebFS,
		examples: assets.Examples,
		mux:      http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) readAsset(name string) ([]byte, error) {
	if s.webFS == nil {
		return nil, fmt.Errorf("no web assets configured")
	}
	name = path.Clean(name)
	if data, err := fs.ReadFile(s.webFS, name); err == nil {
		return data, nil
	}
	return fs.ReadFile(s.webFS, "web/"+name)
}

func (s *Server) routes() {
	s.mux.HandleFunc("/api/af", s.requirePOST(s.handleAF))
	s.mux.HandleFunc("/api/scan", s.requirePOST(s.handleScan))
	s.mux.HandleFunc("/api/examples", s.handleExamples)
	s.mux.HandleFunc("/", s.handleIndex)
	s.mux.HandleFunc("/static/", s.handleStatic)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFound(w, "no such page: "+r.URL.Path)
		return
	}
	data, err := s.readAsset("index.html")
	if err != nil {
		internalError(w, "web asset index.html missing: "+err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/static/")
	if name == "" || strings.Contains(name, "..") {
		notFound(w, "missing or invalid static asset name")
		return
	}
	data, err := s.readAsset(name)
	if err != nil {
		notFound(w, "static asset not found: "+name)
		return
	}
	w.Header().Set("Content-Type", contentType(name))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func contentType(name string) string {
	switch {
	case strings.HasSuffix(name, ".css"):
		return "text/css; charset=utf-8"
	case strings.HasSuffix(name, ".js"):
		return "text/javascript; charset=utf-8"
	case strings.HasSuffix(name, ".html"):
		return "text/html; charset=utf-8"
	case strings.HasSuffix(name, ".json"):
		return "application/json; charset=utf-8"
	default:
		return "text/plain; charset=utf-8"
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	recoverer(logger(s.mux)).ServeHTTP(w, r)
}
