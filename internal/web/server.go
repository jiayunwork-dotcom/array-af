package web

import (
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// Assets 是 Web 服务所需的静态资源。
type Assets struct {
	// WebFS 是嵌入的 web/ 目录（embed.FS）。
	WebFS fs.FS
	// Examples 是内置算例，键为算例名。
	Examples map[string][]byte
}

// Server 是 array-af 的 HTTP 控制台。
type Server struct {
	webFS    fs.FS
	examples map[string][]byte
	mux      *http.ServeMux
}

// NewServer 构造控制台服务器。
func NewServer(assets Assets) *Server {
	s := &Server{
		webFS:    assets.WebFS,
		examples: assets.Examples,
		mux:      http.NewServeMux(),
	}
	s.routes()
	return s
}

// readAsset 读取静态资源。embed 的 web/ 目录带 "web" 前缀，
// 测试用的 MapFS 不带，这里两种路径都尝试。
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

// routes 注册全部路由。method 约束在 handler 内检查，
// 以兼容 Go 1.21 的 ServeMux（不支持 "POST /path" 写法）。
func (s *Server) routes() {
	s.mux.HandleFunc("/api/af", s.requirePOST(s.handleAF))
	s.mux.HandleFunc("/api/scan", s.requirePOST(s.handleScan))
	s.mux.HandleFunc("/api/examples", s.handleExamples)
	s.mux.HandleFunc("/", s.handleIndex)
	s.mux.HandleFunc("/static/", s.handleStatic)
}

// handleIndex 返回控制台页面。
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

// handleStatic 提供 web/ 下的静态资源。
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

// contentType 按扩展名返回静态资源类型。
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

// ServeHTTP 实现 http.Handler，走中间件链。
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	recoverer(logger(s.mux)).ServeHTTP(w, r)
}
