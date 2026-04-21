package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zero-life/server/internal/dto/request"
	"github.com/zero-life/server/internal/pkg/errcode"
	"github.com/zero-life/server/internal/service"
)

type ImportController struct {
	importService *service.ImportService
}

func NewImportController(importService *service.ImportService) *ImportController {
	return &ImportController{importService: importService}
}

func (ctrl *ImportController) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrImportFileInvalid)
		return
	}
	defer file.Close()

	result, err := ctrl.importService.Upload(header.Filename, file)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *ImportController) Parse(c *gin.Context) {
	userID := c.GetUint64("user_id")
	_ = userID

	var req request.ImportParseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.importService.Parse(req.FileID, req.Mapping)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *ImportController) Execute(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.ImportExecuteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.importService.Execute(userID, req.FileID, req.Mapping)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}
