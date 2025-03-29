package handler

import (
	"net/http"

	"github.com/Thoustick/SlugKiller/internal/service"
	"github.com/Thoustick/SlugKiller/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.URLService
	logger  logger.Logger
}

func NewHandler(s service.URLService, l logger.Logger) *Handler {
	return &Handler{
		service: s,
		logger:  l,
	}
}

func (h *Handler) Start(addr string) error {
	r := gin.Default()
	h.RegisterRoutes(r)
	h.logger.Info("starting server", map[string]interface{}{"addr": addr})
	return r.Run(addr)
}

// func (h *Handler) ResolveURL(c *gin.Context) { /* ... */ }

func (h *Handler) ShortenURL(c *gin.Context) {
	h.logger.Info("Handling shorten request", map[string]interface{}{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"ip":     c.ClientIP(),
	})
	ctx := c.Request.Context()
	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid shorten request", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	slug, err := h.service.Shorten(ctx, req.URL)
	if err != nil {
		h.logger.Error("Failed to shorten URL", err, map[string]interface{}{
			"url": req.URL,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to shorten URL"})
		return
	}

	h.logger.Info("URL shortened successfully", map[string]interface{}{
		"url":  req.URL,
		"slug": slug,
	})

	c.JSON(http.StatusOK, ShortenResponse{Slug: slug})
}

func (h *Handler) ResolveURL(c *gin.Context) {
	slug := c.Param("slug")
	h.logger.Info("Handling resolve request", map[string]interface{}{
		"slug":   slug,
		"path":   c.Request.URL.Path,
		"ip":     c.ClientIP(),
		"method": c.Request.Method,
	})

	ctx := c.Request.Context()
	originalURL, err := h.service.Resolve(ctx, slug)
	if err != nil {
		h.logger.Error("Failed to resolve URL", err, map[string]interface{}{
			"slug": slug,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve URL"})
		return
	}

	if originalURL == "" {
		h.logger.Warn("Slug not found", map[string]interface{}{
			"slug": slug,
		})
		c.JSON(http.StatusNotFound, gin.H{"error": "Slug not found"})
		return
	}

	h.logger.Info("Redirecting to original URL", map[string]interface{}{
		"slug": slug,
		"url":  originalURL,
	})

	c.Redirect(http.StatusMovedPermanently, originalURL)
}
