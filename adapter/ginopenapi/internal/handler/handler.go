package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/module/specui"
)

type Handler struct {
	handler *specui.Handler
}

func NewHandler(gen spec.Generator) *Handler {
	return &Handler{
		handler: specui.NewHandler(gen),
	}
}

func (h *Handler) DocsPath() string {
	return h.handler.DocsPath()
}

func (h *Handler) DocsFilePath() string {
	return h.handler.DocsFilePath()
}

func (h *Handler) Docs(c *gin.Context) {
	ui := h.handler.Docs()
	ui.ServeHTTP(c.Writer, c.Request)
}

func (h *Handler) DocsFile(c *gin.Context) {
	ui := h.handler.DocsFile()
	ui.ServeHTTP(c.Writer, c.Request)
}
