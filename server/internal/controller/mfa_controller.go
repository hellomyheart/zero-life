package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

type MFAController struct {
	mfaService *service.MFAService
}

func NewMFAController(mfaService *service.MFAService) *MFAController {
	return &MFAController{mfaService: mfaService}
}

func (ctrl *MFAController) Enable(c *gin.Context) {
	userID := c.GetUint64("user_id")

	result, err := ctrl.mfaService.Enable(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *MFAController) Confirm(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.MFAConfirmReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	if err := ctrl.mfaService.Confirm(userID, req.Code); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

func (ctrl *MFAController) Disable(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.MFADisableReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	if err := ctrl.mfaService.Disable(userID, req.Code); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

func (ctrl *MFAController) Verify(c *gin.Context) {
	var req request.MFAVerifyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.mfaService.VerifyMFA(req.MFAToken, req.Code)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *MFAController) GetBackupCodes(c *gin.Context) {
	userID := c.GetUint64("user_id")

	result, err := ctrl.mfaService.GetBackupCodes(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *MFAController) RegenerateBackupCodes(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.MFACodeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.mfaService.RegenerateBackupCodes(userID, req.Code)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}
