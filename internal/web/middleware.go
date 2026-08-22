package web

import (
	"log"
	"net/http"
	"time"
)

// requirePOST 只放行 POST，其余返回 405 错误体。
func (s *Server) requirePOST(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed,
				"method "+r.Method+" not allowed on "+r.URL.Path+" (use POST)")
			return
		}
		next(w, r)
	}
}

// logger 记录每次请求的方法、路径与耗时。
func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s (%s)", r.Method, r.URL.Path, time.Since(start))
	})
}

// recoverer 捕获 handler panic 并返回 500 错误体，
// 避免连接被粗暴切断后前端只看到空响应。
func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic serving %s %s: %v", r.Method, r.URL.Path, rec)
				writeError(w, http.StatusInternalServerError,
					"internal error while serving request")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
