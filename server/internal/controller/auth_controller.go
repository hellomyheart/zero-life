package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellomyheart/zero-life/server/internal/dto/request"
	"github.com/hellomyheart/zero-life/server/internal/pkg/errcode"
	"github.com/hellomyheart/zero-life/server/internal/service"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var req request.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.authService.Register(&req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req request.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.authService.Login(&req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	var req request.RefreshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *AuthController) ForgotPassword(c *gin.Context) {
	var req request.ForgotPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	if err := ctrl.authService.ForgotPassword(req.Email); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

func (ctrl *AuthController) ResetPassword(c *gin.Context) {
	var req request.ResetPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	if err := ctrl.authService.ResetPassword(req.Token, req.Password); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

func (ctrl *AuthController) GetProfile(c *gin.Context) {
	userID := c.GetUint64("user_id")

	result, err := ctrl.authService.GetProfile(userID)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *AuthController) UpdateProfile(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	result, err := ctrl.authService.UpdateProfile(userID, &req)
	if err != nil {
		handleError(c, err)
		return
	}

	Success(c, result)
}

func (ctrl *AuthController) ChangePassword(c *gin.Context) {
	userID := c.GetUint64("user_id")

	var req request.ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, errcode.ErrBadRequest)
		return
	}

	if err := ctrl.authService.ChangePassword(userID, &req); err != nil {
		handleError(c, err)
		return
	}

	Success(c, nil)
}

func handleError(c *gin.Context, err error) {
	if e, ok := err.(*errcode.Error); ok {
		status := http.StatusInternalServerError
		switch e.Code {
		case 400:
			status = http.StatusBadRequest
		case 401:
			status = http.StatusUnauthorized
		case 403:
			status = http.StatusForbidden
		case 404:
			status = http.StatusNotFound
		case 429:
			status = http.StatusTooManyRequests
		default:
			if e.Code >= 10000 {
				status = http.StatusBadRequest
			}
		}
		Error(c, status, e)
		return
	}
	Error(c, http.StatusInternalServerError, errcode.ErrInternal)
}
