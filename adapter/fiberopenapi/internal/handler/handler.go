package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
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

func (h *Handler) Docs(c *fiber.Ctx) error {
	ui := h.handler.Docs()
	return adaptor.HTTPHandler(ui)(c)
}

func (h *Handler) DocsFile(c *fiber.Ctx) error {
	ui := h.handler.DocsFile()
	return adaptor.HTTPHandler(ui)(c)
}
