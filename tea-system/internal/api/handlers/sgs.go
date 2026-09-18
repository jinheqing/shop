package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"tea-system/internal/models"
	"tea-system/internal/repository"
)

// ============================================================
// SgsReportHandler — SGS 检测报告管理 + 公开查询
// ============================================================

type SgsReportHandler struct {
	repo *repository.InvoiceRepo
}

func NewSgsReportHandler(repo *repository.InvoiceRepo) *SgsReportHandler {
	return &SgsReportHandler{repo: repo}
}

// SgsReportCreateRequest — POST /sgs-reports
type SgsReportCreateRequest struct {
	ReportNo  string          `json:"report_no" binding:"required"`
	BatchNo   string          `json:"batch_no" binding:"required"`
	TeaType   string          `json:"tea_type" binding:"required"`
	TestDate  string          `json:"test_date" binding:"required"`
	IssueDate string          `json:"issue_date" binding:"required"`
	PdfURL    string          `json:"pdf_url" binding:"required"`
	TestItems json.RawMessage `json:"test_items,omitempty"`
}

// SgsReportUpdateRequest — PUT /sgs-reports/:id
type SgsReportUpdateRequest struct {
	PdfURL *string `json:"pdf_url,omitempty"`
}

// ==================== Staff only ====================

// Create — POST /sgs-reports
func (h *SgsReportHandler) Create(c *gin.Context) {
	var req SgsReportCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	testDate, err := time.Parse("2006-01-02", req.TestDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid test_date (YYYY-MM-DD)"})
		return
	}
	issueDate, err := time.Parse("2006-01-02", req.IssueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid issue_date (YYYY-MM-DD)"})
		return
	}

	var testItems models.JSONMap
	if len(req.TestItems) > 0 {
		if err := json.Unmarshal(req.TestItems, &testItems); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid test_items json"})
			return
		}
	}

	s := &models.SgsReport{
		ReportNo:  req.ReportNo,
		BatchNo:   req.BatchNo,
		TeaType:   req.TeaType,
		TestDate:  testDate,
		IssueDate: issueDate,
		PdfURL:    req.PdfURL,
		TestItems: testItems,
	}

	if err := h.repo.CreateSgsReport(c.Request.Context(), s); err != nil {
		log.Error().Err(err).Msg("sgs: create failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "create failed"})
		return
	}
	c.JSON(http.StatusCreated, s)
}

// List — GET /sgs-reports (staff)
func (h *SgsReportHandler) List(c *gin.Context) {
	teaType := c.Query("tea_type")
	batchNo := c.Query("batch_no")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	list, total, err := h.repo.ListSgsReports(c.Request.Context(), repository.SgsListFilter{
		TeaType: teaType,
		BatchNo: batchNo,
		Page:    page,
		Size:    size,
	})
	if err != nil {
		log.Error().Err(err).Msg("sgs: list failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "list failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": list,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// GetByID — GET /sgs-reports/:id (staff)
func (h *SgsReportHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	s, err := h.repo.GetSgsReportByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrSgsReportNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "sgs report not found"})
			return
		}
		log.Error().Err(err).Msg("sgs: get failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
		return
	}

	c.JSON(http.StatusOK, s)
}

// Update — PUT /sgs-reports/:id (staff)
func (h *SgsReportHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	var req SgsReportUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	patch := map[string]interface{}{}
	if v := req.PdfURL; v != nil {
		patch["pdf_url"] = *v
	}

	if len(patch) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "no fields to update"})
		return
	}

	if err := h.repo.UpdateSgsReport(c.Request.Context(), id, patch); err != nil {
		if errors.Is(err, repository.ErrSgsReportNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "sgs report not found"})
			return
		}
		log.Error().Err(err).Msg("sgs: update failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "update failed"})
		return
	}

	s, _ := h.repo.GetSgsReportByID(c.Request.Context(), id)
	c.JSON(http.StatusOK, s)
}

// Delete — DELETE /sgs-reports/:id (staff)
func (h *SgsReportHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid id"})
		return
	}

	if err := h.repo.DeleteSgsReport(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrSgsReportNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "sgs report not found"})
			return
		}
		log.Error().Err(err).Msg("sgs: delete failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "delete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "deleted"})
}

// ==================== 公开（无 JWT） ====================

// PublicList — GET /public/sgs-reports (公开查询，按 report_no 查或分页列表)
// 支持 ?report_no=xxx (精确查) 或 ?tea_type=xxx (过滤)
func (h *SgsReportHandler) PublicList(c *gin.Context) {
	reportNo := c.Query("report_no")
	if reportNo != "" {
		// 精确查单份（公开"通过凭证查询"场景）
		s, err := h.repo.GetSgsReportByReportNo(c.Request.Context(), reportNo)
		if err != nil {
			if errors.Is(err, repository.ErrSgsReportNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "sgs report not found"})
				return
			}
			log.Error().Err(err).Msg("sgs: public get by report_no failed")
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "internal error"})
			return
		}
		c.JSON(http.StatusOK, s)
		return
	}

	// 或者：分页列表（公开浏览）
	teaType := c.Query("tea_type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	list, total, err := h.repo.ListSgsReports(c.Request.Context(), repository.SgsListFilter{
		TeaType: teaType,
		Page:    page,
		Size:    size,
	})
	if err != nil {
		log.Error().Err(err).Msg("sgs: public list failed")
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "list failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": list,
		"total": total,
		"page":  page,
		"size":  size,
	})
}
