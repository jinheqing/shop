package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"tea-system/internal/repository"
	"tea-system/internal/service"
)

// SlowPresetHandler — 慢直播独立接口组（独立于 live_rooms）
type SlowPresetHandler struct {
	svc *service.SlowPresetService
}

func NewSlowPresetHandler(svc *service.SlowPresetService) *SlowPresetHandler {
	return &SlowPresetHandler{svc: svc}
}

type slowPresetCreateReq struct {
	Name        string `json:"name" binding:"required"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

type slowPresetUpdateReq struct {
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

// Create — POST /slow-presets
func (h *SlowPresetHandler) Create(c *gin.Context) {
	var req slowPresetCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	room, err := h.svc.Create(c.Request.Context(), req.Name, req.Location, req.Description)
	if err != nil {
		log.Error().Err(err).Msg("slow_preset: create failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "create failed"})
		return
	}
	c.JSON(http.StatusCreated, room)
}

// List — GET /slow-presets
func (h *SlowPresetHandler) List(c *gin.Context) {
	enabledOnly := c.Query("enabled") == "true"
	items, err := h.svc.List(c.Request.Context(), enabledOnly)
	if err != nil {
		log.Error().Err(err).Msg("slow_preset: list failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "list failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "total": len(items)})
}

// PublicList — GET /public/slow-presets（公开路由）
func (h *SlowPresetHandler) PublicList(c *gin.Context) {
	items, err := h.svc.PublicList(c.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("slow_preset: public list failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

// Update — PUT /slow-presets/:id
func (h *SlowPresetHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}
	var req slowPresetUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	room, err := h.svc.Update(c.Request.Context(), id, req.Name, req.Location, req.Description)
	if err != nil {
		if errors.Is(err, repository.ErrLiveRoomNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "slow preset not found"})
			return
		}
		log.Error().Err(err).Msg("slow_preset: update failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "update failed"})
		return
	}
	c.JSON(http.StatusOK, room)
}

// Delete — DELETE /slow-presets/:id
func (h *SlowPresetHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrLiveRoomNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "slow preset not found"})
			return
		}
		log.Error().Err(err).Msg("slow_preset: delete failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "deleted"})
}

// Enable — POST /slow-presets/:id/enable
func (h *SlowPresetHandler) Enable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}
	room, err := h.svc.Enable(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrLiveRoomNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "slow preset not found"})
			return
		}
		log.Error().Err(err).Msg("slow_preset: enable failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "enable failed"})
		return
	}
	c.JSON(http.StatusOK, room)
}

// Disable — POST /slow-presets/:id/disable
func (h *SlowPresetHandler) Disable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}
	room, err := h.svc.Disable(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrLiveRoomNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "slow preset not found"})
			return
		}
		log.Error().Err(err).Msg("slow_preset: disable failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "disable failed"})
		return
	}
	c.JSON(http.StatusOK, room)
}
