// Package web 提供 array-af 的 HTTP 控制台：求解 API 与静态页面。
package web

import (
	"encoding/json"
	"net/http"
)

// errorBody 是统一的错误响应体，前端与 curl 都能直接读到失败原因。
type errorBody struct {
	Error string `json:"error"`
}

// writeError 写出 JSON 错误体。
func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorBody{Error: msg})
}

// writeJSON 写出成功 JSON。
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// badRequest 写 400 错误体。
func badRequest(w http.ResponseWriter, msg string) {
	writeError(w, http.StatusBadRequest, msg)
}

// notFound 写 404 错误体。
func notFound(w http.ResponseWriter, msg string) {
	writeError(w, http.StatusNotFound, msg)
}

// internalError 写 500 错误体。
func internalError(w http.ResponseWriter, msg string) {
	writeError(w, http.StatusInternalServerError, msg)
}

// decodeJSON 解析请求体并限制大小。
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	body := http.MaxBytesReader(w, r.Body, 1<<20) // 1 MiB
	dec := json.NewDecoder(body)
	if err := dec.Decode(dst); err != nil {
		badRequest(w, "invalid JSON request body: "+err.Error())
		return false
	}
	return true
}
