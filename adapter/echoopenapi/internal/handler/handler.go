package handler

import (
	"github.com/labstack/echo/v4"
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

func (h *Handler) Docs(c echo.Context) error {
	ui := h.handler.Docs()
	ui.ServeHTTP(c.Response(), c.Request())
	return nil
}

func (h *Handler) DocsFile(c echo.Context) error {
	ui := h.handler.DocsFile()
	ui.ServeHTTP(c.Response(), c.Request())
	return nil
}
