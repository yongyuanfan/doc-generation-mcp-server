package apiserver

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/yong/doc-generation-mcp-server/internal/common"
	"github.com/yong/doc-generation-mcp-server/internal/config"
	"github.com/yong/doc-generation-mcp-server/internal/model"
	docsvc "github.com/yong/doc-generation-mcp-server/internal/service/document"
)

type Handler struct {
	service *docsvc.Service
	config  config.Config
}

func NewHandler(cfg config.Config, service *docsvc.Service) http.Handler {
	h := &Handler{service: service, config: cfg}
	mux := http.NewServeMux()
	mux.HandleFunc("/documents/generate", h.handleGenerate)
	mux.HandleFunc("/documents/render-template", h.handleRenderTemplate)
	mux.HandleFunc("/documents/files/", h.handleDownload)
	mux.HandleFunc("/templates", h.handleTemplates)
	mux.HandleFunc("/capabilities", h.handleCapabilities)
	return withMaxBodyBytes(mux, cfg.DocxMaxRequestBodyBytes)
}

func Healthz(w http.ResponseWriter, _ *http.Request) {
	common.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req model.GenerateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteError(w, http.StatusBadRequest, "invalid json body")
		return
	}
	result, err := h.service.Generate(r.Context(), req)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	common.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) handleRenderTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var req model.RenderTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteError(w, http.StatusBadRequest, "invalid json body")
		return
	}
	result, err := h.service.RenderTemplate(r.Context(), req)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	common.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) handleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		common.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	name := strings.TrimPrefix(r.URL.Path, "/documents/files/")
	path, err := h.service.DownloadPath(name)
	if err != nil {
		common.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	w.Header().Set("Content-Disposition", `attachment; filename="`+name+`"`)
	http.ServeFile(w, r, path)
}

func (h *Handler) handleCapabilities(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	common.WriteJSON(w, http.StatusOK, h.service.Capabilities())
}

func (h *Handler) handleTemplates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	result, err := h.service.ListTemplates()
	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	common.WriteJSON(w, http.StatusOK, result)
}

func withMaxBodyBytes(next http.Handler, limit int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.Body = http.MaxBytesReader(w, r.Body, limit)
		}
		next.ServeHTTP(w, r)
	})
}
